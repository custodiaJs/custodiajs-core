package httpapi

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"
	"vnh1/types"
	"vnh1/utils"

	"github.com/btcsuite/btcutil/base58"
	"github.com/dop251/goja"
)

func (o *HttpApiService) httpRPCHandler(w http.ResponseWriter, r *http.Request) {
	// Es wird eine neue Process Log Session erzeugt
	proc := utils.NewProcLogSession()

	// Es wird geprüft ob es sich um die POST Methode handelt
	proc.LogPrint("RPC: validate incomming remote function call request from '%s'\n", r.RemoteAddr)
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
	proc.LogPrint("RPC: determine the script container '%s'\n", strings.ToUpper(request.VmId))
	foundedVM, err := o.core.GetScriptContainerVMByID(request.VmId)
	if err != nil {
		proc.LogPrint("RPC: determine the script container '%s' failed, unkown script vm container\n", strings.ToUpper(request.VmId))
		errorResponse(request.ContentType, w, err.Error())
		proc.LogPrint("RPC: failed\n")
		return
	}

	// Es wird geprüft ob es sich um eine WebRequest aus einem Webbrowser handelt,
	// wenn ja wird ermittelt ob es sich um eine Zulässige Quelle handelt,
	// wenn es sich nicht um eine zulässige Quelle handelt, wird der Vorgang abgebrochen.
	requestHttpSource := getRefererOrXRequestedWith(request)
	if hasRefererOrXRequestedWith(request) && !foundedVM.ValidateRPCRequestSource(requestHttpSource) {
		proc.LogPrint("RPC: process aborted, not allowed request websource '%s'\n", getRefererOrXRequestedWith(request))
		errorResponse(request.ContentType, w, "not allowed request from webresource")
		return
	}

	// Es wird versucht den Body einzulesen
	data, err := extractHttpRpcBody(request.ContentType, r.Body)
	if err != nil {
		errorResponse(request.ContentType, w, "invalid body data")
		proc.LogPrint("RPC: failed, invalid request\n")
		return
	}

	// Es wird geprüft ob es sich um einen Zulässigen Funktionsnamen handelt
	// Es wird geprüft, ob der Funktionsname korrekt ist
	if !utils.ValidateFunctionName(data.FunctionName) {
		errorResponse(request.ContentType, w, "invalid function name")
		proc.LogPrint("RPC: failed, invalid function name\n")
		return
	}

	// Der Body wird geschlossen sobald der Vorgang beendet wurde
	defer r.Body.Close()

	// Es wird versucht die Passende Funktion zu ermitteln
	proc.LogPrint("RPC: &[%s]: determine the function '%s'\n", foundedVM.GetVMName(), data.FunctionName)
	var foundFunction types.SharedFunctionInterface
	for _, item := range foundedVM.GetLocalSharedFunctions() {
		if item.GetName() == data.FunctionName {
			foundFunction = item
			break
		}
	}

	// Es wird geprüft ob eine Funktion gefunden wurde, wenn nicht werden die Öffentlichen Funktionen durchsucht
	if foundFunction == nil {
		for _, item := range foundedVM.GetPublicSharedFunctions() {
			if item.GetName() == data.FunctionName {
				foundFunction = item
				break
			}
		}
	}

	// Sollte keine Passende Funktion gefunden werden, wird der Vorgang abgebrochen
	if foundFunction == nil {
		proc.LogPrint("RPC: &[%s]: determine the function '%s' failed, unkown function\n", foundedVM.GetVMName(), data.FunctionName)
		errorResponse(request.ContentType, w, "function not found")
		proc.LogPrint("RPC: failed\n")
		return
	}

	// Es wird ermitelt ob die Datentypen korrekt sind
	if len(foundFunction.GetParmTypes()) != len(data.Parms) {
		errorResponse(request.ContentType, w, "the number of parameters required does not match the number of parameters submitted")
		return
	}

	// Die Einzelnen Parameter werden geprüft und abgearbeitet
	proc.LogPrint("RPC: &[%s@%s]: convert %d function parameters\n", foundFunction.GetName(), foundedVM.GetVMName(), len(foundFunction.GetParmTypes()))
	extractedValues := make([]types.FunctionParameterBundleInterface, 0)
	for x := range foundFunction.GetParmTypes() {
		// Es wird geprüft ob es sich bei dem Angefordeten Parameter um einen zulässigen Parameter handelt
		if foundFunction.GetParmTypes()[x] != data.Parms[x].Type {
			errorResponse(request.ContentType, w, fmt.Sprintf("the data type of parameter %d does not match the required data type", x))
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
			newEntry := &FunctionParameterCapsle{Value: converted, CType: "bool"}

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
				newEntry := &FunctionParameterCapsle{Value: onvertedfloat, CType: "number"}

				// Die Daten werden hinzugefügt
				extractedValues = append(extractedValues, newEntry)
				break
			}

			// Der Eintrag wird erzeugt
			newEntry := &FunctionParameterCapsle{Value: converted, CType: "number"}

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
			newEntry := &FunctionParameterCapsle{Value: converted, CType: "string"}

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
			newEntry := &FunctionParameterCapsle{Value: converted, CType: "array"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case "object":
			// Der Eintrag wird erzeugt
			newEntry := &FunctionParameterCapsle{Value: data.Parms[x].Value, CType: "object"}

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
			newEntry := &FunctionParameterCapsle{Value: decodedDataSlice, CType: "bytearray"}

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
			newEntry := &FunctionParameterCapsle{Value: timeObj, CType: "timestamp"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		}
	}

	// Die Funktion wird aufgerufen
	proc.LogPrint("RPC: &[%s@%s]: call function\n", foundFunction.GetName(), foundedVM.GetVMName())
	result, err := foundFunction.EnterFunctionCall(request, &RpcRequest{parms: extractedValues})
	if err != nil {
		proc.LogPrint("RPC: &[%s]: call function '%s' error\n\t%s\n", foundedVM.GetVMName(), foundFunction.GetName(), err)
		errorResponse(request.ContentType, w, "an error occurred when calling the function, error: "+err.Error())
		return
	}

	// Die Antwortdaten werden Extrahiert
	var responseData *RPCResponseData
	if result == nil {
		responseData = &RPCResponseData{DType: "null", Value: nil}
		proc.LogPrintSuccs("RPC: &[%s@%s]: function call, return null\n", foundFunction.GetName(), foundedVM.GetVMName())
	} else if result.ExportType() == goja.Undefined().ExportType() && result.Export() == nil {
		responseData = &RPCResponseData{DType: "undefined", Value: nil}
		proc.LogPrintSuccs("RPC: &[%s@%s]: function call, return undefined\n", foundFunction.GetName(), foundedVM.GetVMName())
	} else {
		switch result.ExportType().Kind() {
		case reflect.Bool:
			responseData = &RPCResponseData{DType: "boolean", Value: result.ToBoolean()}
			proc.LogPrintSuccs("RPC: &[%s@%s]: function call, return boolean\n", foundFunction.GetName(), foundedVM.GetVMName())
		case reflect.Int64:
			responseData = &RPCResponseData{DType: "number", Value: result.ToInteger()}
			proc.LogPrintSuccs("RPC: &[%s@%s]: function call, return int64\n", foundFunction.GetName(), foundedVM.GetVMName())
		case reflect.Float64:
			responseData = &RPCResponseData{DType: "number", Value: result.ToFloat()}
			proc.LogPrintSuccs("RPC: &[%s@%s]: function call, return float64\n", foundFunction.GetName(), foundedVM.GetVMName())
		case reflect.String:
			responseData = &RPCResponseData{DType: "string", Value: result.String()}
			proc.LogPrintSuccs("RPC: &[%s@%s]: function call, return string\n", foundFunction.GetName(), foundedVM.GetVMName())
		case reflect.Slice:
			slicedObject, isConverted := result.Export().([]interface{})
			if !isConverted {
				http.Error(w, "invalid object datatype, slice", http.StatusBadRequest)
				return
			}
			responseData = &RPCResponseData{DType: "array", Value: slicedObject}
			proc.LogPrintSuccs("RPC: &[%s@%s]: function call, return slice\n", foundFunction.GetName(), foundedVM.GetVMName())
		case reflect.Map:
			mapObjected, isConverted := result.Export().(map[string]interface{})
			if !isConverted {
				http.Error(w, "invalid object datatype, object", http.StatusBadRequest)
				return
			}
			responseData = &RPCResponseData{DType: "object", Value: mapObjected}
			proc.LogPrintSuccs("RPC: &[%s@%s]: function call, return map\n", foundFunction.GetName(), foundedVM.GetVMName())
		case reflect.Func:
			proc.LogPrintSuccs("RPC: &[%s@%s]: function call, return function\n", foundFunction.GetName(), foundedVM.GetVMName())
			errorResponse(request.ContentType, w, "function return not allowed in web remote function call request")
			return
		default:
			errorResponse(request.ContentType, w, fmt.Sprintf("fhe function returned a data type (%s) which is not supported, the function was executed without errors", result.ExportType().Kind().String()))
			return
		}
	}

	// Die Daten werden zurückgesendet
	proc.LogPrint("RPC: &[%s@%s]: sending function call response\n", foundFunction.GetName(), foundedVM.GetVMName())
	responsSize, err := finalResponse(request.ContentType, w, responseData)
	if err != nil {
		proc.LogPrint("RPC: &[%s]: call function response sending '%s' error\n\t%s\n", foundedVM.GetVMName(), foundFunction.GetName(), err)
	}
	proc.LogPrint(fmt.Sprintf("RPC: done, response size = %d bytes\n", responsSize))
}
