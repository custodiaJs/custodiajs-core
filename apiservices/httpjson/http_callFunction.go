package httpjson

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/CustodiaJS/custodiajs-core/saftychan"
	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/utils"
	"github.com/CustodiaJS/custodiajs-core/utils/grsbool"
	"github.com/CustodiaJS/custodiajs-core/utils/procslog"

	"github.com/btcsuite/btcutil/base58"
)

func (o *HttpApiService) httpCallFunction(w http.ResponseWriter, r *http.Request) {
	// Es wird eine neue Process Log Session erzeugt
	proc := procslog.NewProcLogSession()

	// Holen Sie sich den Kontext der Anfrage
	ctx := r.Context()

	// Gibt an ob die Verbindung getrennt wurde
	isConnected := grsbool.NewGrsbool(true)

	// Der ResultChan wird erezugt
	saftyResponseChan := saftychan.NewFunctionCallReturnChan()

	// Starte eine Go-Routine, um die Verbindung zu überwachen
	go func() {
		// Es wird darauf gewartet dass die Verbindung geschlossen wird
		<-ctx.Done()

		// Es wird Signalisiert dass die Verbindung geschlossen wurde
		isConnected.Set(false)

		// Das SaftyChan wird geschlossen
		saftyResponseChan.Close()
	}()

	// Es wird geprüft ob es sich um die POST Methode handelt
	procslog.ProcFormatConsoleText(proc, "HTTP-RPC", types.VALIDATE_INCOMMING_REMOTE_FUNCTION_CALL_REQUEST_FROM, r.RemoteAddr)
	request, err := validatePOSTRequestAndGetRequestData(r)
	if err != nil {
		// Set the 'Allow' header to indicate that only POST is allowed
		w.Header().Set("Allow", "POST")

		// Send the HTTP status code 405 Method Not Allowed
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)

		// Der Vorgang wird beendet
		return
	}

	// Es wird geprüft ob es sich um eine bekannte VM handelt
	procslog.ProcFormatConsoleText(proc, "HTTP-RPC", types.DETERMINE_THE_SCRIPT_CONTAINER, strings.ToUpper(request.VmId))
	foundedVM, foundVM, err := o.core.GetScriptContainerVMByID(request.VmId)
	if err != nil {
		proc.LogPrint("HTTP-RPC", "failed\n")
		errorResponse(request.ContentType, w, "internal error")
		return
	}
	if !foundVM {
		proc.LogPrint("HTTP-RPC", "failed\n")
		errorResponse(request.ContentType, w, "not found")
		return
	}

	// Es wird geprüft ob es sich um eine WebRequest aus einem Webbrowser handelt,
	// wenn ja wird ermittelt ob es sich um eine Zulässige Quelle handelt,
	// wenn es sich nicht um eine zulässige Quelle handelt, wird der Vorgang abgebrochen.
	requestHttpSource := getRefererOrXRequestedWith(request)
	if hasRefererOrXRequestedWith(request) && !foundedVM.ValidateRPCRequestSource(requestHttpSource) {
		proc.LogPrint("HTTP-RPC: process aborted, not allowed request websource '%s'\n", getRefererOrXRequestedWith(request))
		errorResponse(request.ContentType, w, "not allowed request from webresource")
		return
	}

	// Es wird versucht den Body einzulesen
	data, err := extractHttpRpcBody(request.ContentType, r.Body)
	if err != nil {
		errorResponse(request.ContentType, w, "invalid body data")
		proc.LogPrint("HTTP-RPC", "failed, invalid request\n")
		return
	}

	// Es wird geprüft ob es sich um einen Zulässigen Funktionsnamen handelt
	// Es wird geprüft, ob der Funktionsname korrekt ist
	if !utils.ValidateFunctionName(data.FunctionName) {
		errorResponse(request.ContentType, w, "invalid function name")
		proc.LogPrint("HTTP-RPC", "failed, invalid function name\n")
		return
	}

	// Der Body wird geschlossen sobald der Vorgang beendet wurde
	defer r.Body.Close()

	// Die Datentypen der Parameter werden ausgeslesen
	dataTypeParms := make([]string, 0)
	for _, item := range data.Parms {
		dataTypeParms = append(dataTypeParms, item.Type)
	}

	// Aus dem Request wird eine Funktionssignatur erzeugt
	searchedFunctionSignature := &types.FunctionSignature{
		VMID:         strings.ToLower(request.VmId),
		FunctionName: data.FunctionName,
		Params:       dataTypeParms,
		ReturnType:   data.ReturnDataType,
	}

	// Es wird versucht die Passende Funktion zu ermitteln
	procslog.ProcFormatConsoleText(proc, "HTTP-PRC", types.DETERMINE_THE_FUNCTION, foundedVM.GetVMName(), data.FunctionName)
	foundFunction, hasFound, err := foundedVM.GetSharedFunctionBySignature(static.LOCAL, searchedFunctionSignature)
	if err != nil {
		errorResponse(request.ContentType, w, "internal error")
		proc.LogPrint("HTTP-RPC", "failed, invalid function name\n")
		return
	}

	// Sollte keine Passende Funktion gefunden werden, wird der Vorgang abgebrochen
	if !hasFound {
		proc.LogPrint("HTTP-RPC: &[%s]: determine the function '%s' failed, unkown function\n", foundedVM.GetVMName(), data.FunctionName)
		errorResponse(request.ContentType, w, "function not found")
		proc.LogPrint("HTTP-RPC", "failed\n")
		return
	}

	// Es wird ermitelt ob die Datentypen korrekt sind
	if len(foundFunction.GetParmTypes()) != len(data.Parms) {
		proc.LogPrint("HTTP-RPC: &[%s]: the number of parameters required does not match the number of parameters submitted\n", foundedVM.GetVMName())
		errorResponse(request.ContentType, w, "the number of parameters required does not match the number of parameters submitted")
		proc.LogPrint("HTTP-RPC", "failed\n")
		return
	}

	// Die Einzelnen Parameter werden geprüft und abgearbeitet
	extractedValues := make([]*types.FunctionParameterCapsle, 0)
	for x := range foundFunction.GetParmTypes() {
		// Es wird geprüft ob es sich bei dem Angefordeten Parameter um einen zulässigen Parameter handelt
		if foundFunction.GetParmTypes()[x] != data.Parms[x].Type {
			proc.LogPrint("HTTP-RPC: &[%s@%s]: the data type of parameter %d does not match the required data type\n", foundFunction.GetName(), foundedVM.GetVMName(), x)
			errorResponse(request.ContentType, w, fmt.Sprintf("the data type of parameter %d does not match the required data type", x))
			proc.LogPrint("HTTP-RPC", "failed\n")
			return
		}

		// Es wird versucht den Datentypen umzuwandeln
		switch data.Parms[x].Type {
		case "boolean":
			// Es wird geprüft ob es sich um ein Boolean handelt
			if reflect.TypeOf(data.Parms[x].Value).Kind() != reflect.Bool {
				errorResponse(request.ContentType, w, fmt.Sprintf("error reading parameter %d, wrong data type", x))
				return
			}

			// Der Datentyp wird umgewandelt
			converted, ok := data.Parms[x].Value.(bool)
			if !ok {
				errorResponse(request.ContentType, w, fmt.Sprintf("error reading parameter %d, wrong data type", x))
				return
			}

			// Der Eintrag wird erzeugt
			newEntry := &types.FunctionParameterCapsle{Value: converted, CType: "bool"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case "number":
			// Der Datentyp wird umgewandelt
			converted, ok := data.Parms[x].Value.(int64)
			if !ok {
				// Es wird geprüft ob es sich um ein Float handelt
				onvertedfloat, ok := data.Parms[x].Value.(float64)
				if !ok {
					errorResponse(request.ContentType, w, fmt.Sprintf("error reading parameter %d, wrong data type", x))
					return
				}

				// Der Eintrag wird erzeugt
				newEntry := &types.FunctionParameterCapsle{Value: onvertedfloat, CType: "number"}

				// Die Daten werden hinzugefügt
				extractedValues = append(extractedValues, newEntry)
				break
			}

			// Der Eintrag wird erzeugt
			newEntry := &types.FunctionParameterCapsle{Value: converted, CType: "number"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case "string":
			// Der Datentyp wird umgewandelt
			converted, ok := data.Parms[x].Value.(string)
			if !ok {
				errorResponse(request.ContentType, w, fmt.Sprintf("error reading parameter %d, wrong data type", x))
				return
			}

			// Der Eintrag wird erzeugt
			newEntry := &types.FunctionParameterCapsle{Value: converted, CType: "string"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case "array":
			// Der Datentyp wird umgewandelt
			converted, ok := data.Parms[x].Value.([]interface{})
			if !ok {
				errorResponse(request.ContentType, w, fmt.Sprintf("error reading parameter %d, wrong data type", x))
				return
			}

			// Der Eintrag wird erzeugt
			newEntry := &types.FunctionParameterCapsle{Value: converted, CType: "array"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case "object":
			// Der Eintrag wird erzeugt
			newEntry := &types.FunctionParameterCapsle{Value: data.Parms[x].Value, CType: "object"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case "bytes":
			// Der Datentyp wird umgewandelt
			converted, ok := data.Parms[x].Value.(string)
			if !ok {
				errorResponse(request.ContentType, w, fmt.Sprintf("error reading parameter %d, wrong data type", x))
				return
			}

			// Es wird geprüft ob der String aus 2 teilen besteht, der este Teil gibt an welches Codec verwendet wird,
			// der Zweite teil enthält die eigentlichen Daten
			splitedValue := strings.Split("://", converted)
			if len(splitedValue) != 2 {
				errorResponse(request.ContentType, w, fmt.Sprintf("error reading parameter %d, wrong data type", x))
				return
			}

			// Es wird geprüft ob es sich um ein zulässiges Codec handelt
			var decodedDataSlice []byte
			switch strings.ToLower(splitedValue[0]) {
			case "base64":
				decodedDataSlice, err = base64.StdEncoding.DecodeString(converted)
				if err != nil {
					errorResponse(request.ContentType, w, fmt.Sprintf("error reading parameter %d, wrong data type", x))
					return
				}
			case "base32":
				decodedDataSlice, err = base32.StdEncoding.DecodeString(converted)
				if err != nil {
					errorResponse(request.ContentType, w, fmt.Sprintf("error reading parameter %d, wrong data type", x))
					return
				}
			case "hex":
				decodedDataSlice, err = hex.DecodeString(converted)
				if err != nil {
					errorResponse(request.ContentType, w, fmt.Sprintf("error reading parameter %d, wrong data type", x))
					return
				}
			case "base58":
				decodedDataSlice = base58.Decode(converted)
			default:
				errorResponse(request.ContentType, w, fmt.Sprintf("error reading parameter %d, wrong data type", x))
				return
			}

			// Der Eintrag wird erzeugt
			newEntry := &types.FunctionParameterCapsle{Value: decodedDataSlice, CType: "bytearray"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case "timestamp":
			// Der Datentyp wird umgewandelt
			converted, ok := data.Parms[x].Value.(int64)
			if !ok {
				http.Error(w, "Datatype converting error: timestamp", http.StatusBadRequest)
				return
			}

			// Umwandlung von Unix-Zeit in time.Time
			timeObj := time.Unix(converted, 0)

			// Der Eintrag wird erzeugt
			newEntry := &types.FunctionParameterCapsle{Value: timeObj, CType: "timestamp"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		}
	}

	// Das HTTP Request Objekt wird erstellt
	requestHttpObject := &types.HttpRpcRequest{
		IsConnected:      isConnected,
		ContentLength:    r.ContentLength,
		PostForm:         r.PostForm,
		Header:           r.Header,
		Host:             r.Host,
		Form:             r.Form,
		Proto:            r.Proto,
		RemoteAddr:       r.RemoteAddr,
		RequestURI:       r.RequestURI,
		TLS:              r.TLS,
		TransferEncoding: r.TransferEncoding,
		URL:              r.URL,
		Cookies:          r.Cookies(),
		UserAgent:        r.UserAgent(),
	}

	// Diese Funktion nimmt die Antwort entgegen
	resolveFunction := func(response *types.FunctionCallReturn) error {
		// Es wird geprüft ob das Response Null ist, wenn ja wird ein Panic ausgelöst
		if response == nil {
			panic("http rpc function call response is null, critical error")
		}

		// Es wird geprüft ob der Vorgang bereits abgeschlossen wurde
		if saftyResponseChan.IsClosed() && isConnected.Bool() {
			return utils.MakeHttpConnectionIsClosedError()
		} else if !isConnected.Bool() {
			return utils.MakeAlreadyAnsweredRPCRequestError()
		}

		// Die Antwort wird geschrieben
		saftyResponseChan.WriteAndClose(response)

		// Es ist kein Fehler aufgetreten
		return nil
	}

	// Das Request Objekt wird erzeugt
	requestObject := &types.RpcRequest{
		Parms:       extractedValues,
		RpcRequest:  request,
		ProcessLog:  proc,
		RequestType: static.HTTP_REQUEST,
		HttpRequest: requestHttpObject,
		Resolve:     resolveFunction,
	}

	// Die Funktion wird aufgerufen
	err = foundFunction.EnterFunctionCall(requestObject)
	if err != nil {
		proc.LogPrint("HTTP-RPC: &[%s]: call function '%s' error\n\t%s\n", foundedVM.GetVMName(), foundFunction.GetName(), err)
		errorResponse(request.ContentType, w, "an error occurred when calling the function, error: "+err.Error())
		return
	}

	// Es wird auf das Ergebniss gewartet
	result, ok := saftyResponseChan.Read()
	if !ok || result == nil {
		// Es wird geprüft ob die Verbindung aufgebaut ist
		if !isConnected.Bool() {
			// Log
			proc.LogPrint("HTTP-RPC", "aborted, connection closed\n")

			// Rückgabe
			return
		}

		// Es handelt sich um einen unbekannten Fehler, Panic
		panic("unkown internal http rpc calling error")
	}

	// Die Antwort wird gebaut
	var responseData *ResponseCapsle
	if result.State == "ok" {
		dt := make([]*RPCResponseData, 0)
		for _, item := range result.Return {
			dt = append(dt, &RPCResponseData{DType: item.Type, Value: item.Value})
		}
		responseData = &ResponseCapsle{Data: dt}
	} else if result.State == "failed" {
		responseData = &ResponseCapsle{Error: result.Error}
	} else if result.State == "exception" {
		responseData = &ResponseCapsle{Error: result.Error}
	} else {
		responseData = &ResponseCapsle{Error: "unkown return state"}
	}

	// Die Daten werden zurückgesendet
	responsSize, err := responseWrite(request.ContentType, w, responseData)
	if err != nil {
		// Log
		proc.LogPrint("HTTP-RPC: &[%s]: call function response sending '%s' error\n\t%s\n", foundedVM.GetVMName(), foundFunction.GetName(), err)

		// Es wird Signalisiert dass die Daten nicht übermittelt werden konnten
		result.Reject()

		// Der Vorgang wird geschlossen
		return
	}

	// Log
	procslog.ProcFormatConsoleText(proc, "HTTP-PRC", types.RPC_CALL_DONE_RESPONSE, fmt.Sprintf("%d", responsSize))

	// Es wird Signalisiert dass die Daten erfolgreich übermittelt wurden
	result.Resolve()
}
