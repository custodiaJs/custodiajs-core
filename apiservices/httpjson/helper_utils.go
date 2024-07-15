package httpjson

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/static/errormsgs"
	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/utils"
	"github.com/btcsuite/btcutil/base58"

	"github.com/fxamacker/cbor/v2"
)

// Wandelt ein Paket in ein Bytearry um
func convertResponseCapsleToByteSlice(rc *ResponseCapsle, requestContentType types.HttpRequestContentType) ([]byte, *types.SpecificError) {
	switch requestContentType {
	case static.HTTP_CONTENT_JSON:
		b, err := json.Marshal(rc)
		if err != nil {
			return nil, errormsgs.HTTP_API_SERVICE_RESPONSE_WRITING_ENCODING_ERROR("convertResponseCapsleToByteSlice", requestContentType, err)
		}
		return b, nil
	case static.HTTP_CONTENT_CBOR:
		b, err := cbor.Marshal(rc)
		if err != nil {
			return nil, errormsgs.HTTP_API_SERVICE_RESPONSE_WRITING_ENCODING_ERROR("convertResponseCapsleToByteSlice", requestContentType, err)
		}
		return b, nil
	default:
		return nil, nil
	}
}

// Wird verwendet um einen Request durchzuführen
func closeHTTPRequest(r *http.Request, coreWebSession types.WebRequestBasedRPCSessionInterface) {
	r.Body.Close()                                // Schließt den HTTP Boddy
	coreWebSession.CloseBecauseFunctionReturned() // Schließt die Gesamte Sitzung
}

// Erzeugt einen SHA3-256 Hash aus einem Response Capsle
func computeSHA3_256HashFromResponseCapsle(rc *ResponseCapsle, requestContentType types.HttpRequestContentType) (string, *types.SpecificError) {
	byted, err := convertResponseCapsleToByteSlice(rc, requestContentType)
	if err != nil {
		err.AddCallerFunctionToHistory("computeSHA3_256HashFromResponseCapsle")
		return "", err
	}

	hexvalued := hex.EncodeToString(byted)
	hash := utils.HashOfString(hexvalued)

	return hash, nil
}

// Funktion zum Abrufen der lokalen IP-Adresse des Servers
func getLocalIPFromRequest(r *http.Request) (string, *types.SpecificError) {
	// Ermittle den Listener vom Context des Requests
	addr := r.Context().Value(http.LocalAddrContextKey).(net.Addr)
	// Konvertiere zu net.TCPAddr, um die IP-Adresse zu erhalten
	tcpAddr := addr.(*net.TCPAddr)
	return tcpAddr.IP.String(), nil
}

// validateRequestAndGetRequestData überprüft einen HTTP-Request auf Gültigkeit und gibt ein RequestData-Objekt zurück.
// Es validiert die HTTP-Methode, die TLS-Verbindung, den Content-Type, den Query-Parameter 'id' und erstellt ein RequestData-Objekt mit den extrahierten Daten.
func validateRequestAndGetRequestData(methode string, r *http.Request, core types.CoreInterface) (*RequestData, *types.SpecificError) {
	// Es wird geprüft ob es sich um die POST Methode handelt
	if r.Method != methode {
		return nil, errormsgs.HTTP_API_SERVICE_INVALID_METHODE("validateRequestAndGetRequestData", r.Method, "POST")
	}

	// Es wird geprüft ob eine TLS Verbindung vorhanden ist
	if r.TLS == nil {
		return nil, errormsgs.HTTP_API_SERVICE_HAS_NO_TLS_ENCRYPTION("validateRequestAndGetRequestData", r.RemoteAddr)
	}

	// Der Content Typ wird geprüft
	var contentType types.HttpRequestContentType
	switch ctype := r.Header.Get("content-type"); ctype {
	case "application/json":
		contentType = static.HTTP_CONTENT_JSON
	case "application/cbor":
		contentType = static.HTTP_CONTENT_CBOR
	default:
		return nil, errormsgs.HTTP_API_SERVICE_INVALID_CODEC("validateRequestAndGetRequestData", ctype)
	}

	// Es wird geprüft ob der VM Query angegeben wurde
	queryParams := r.URL.Query()

	// Prüfe, ob mehr als ein Query-Parameter vorhanden ist oder der Parameter 'name' nicht existiert oder mehr als einen Wert hat.
	if len(queryParams) != 1 {
		return nil, errormsgs.HTTP_API_SERVICE_INVALID_QUERY_PARM_SIZE("validateRequestAndGetRequestData")
	}

	// Es wird versucht die VMId oder den VM-Namen zu extrahieren
	var vlaueType types.RPCRequestVMIdentificationMethode
	var finalValue string
	if vmId, ok := queryParams["id"]; ok {
		// Es wird geprüft ob genau 1 Eintrag vorhanden ist
		if len(vmId) != 1 {
			return nil, errormsgs.HTTP_API_SERVICE_INVALID_QUERY_PARMS("validateRequestAndGetRequestData")
		}

		// Es wird geprüft ob der Eintrag 64 Zeichen lang ist
		if len(vmId[1]) != 64 {
			return nil, errormsgs.HTTP_API_SERVICE_INVALID_QUERY_PARM_ID_SIZE("validateRequestAndGetRequestData")
		}

		// Dert wert wird Dekodiert und wieder Kodiert
		// Es wird geprüft ob es sich um einen Hexwert handelt
		decodedId, err := hex.DecodeString(vmId[0])
		if err != nil {
			return nil, errormsgs.HTTP_API_SERVICE_INVALID_QUERY_PARM_HEX_DECODING("validateRequestAndGetRequestData")
		}

		// Der Wert sowie der Typ werden geschrieben
		finalValue, vlaueType = hex.EncodeToString(decodedId), static.RPC_REQUEST_METHODE_VM_IDENT_ID
	} else if vmName, ok := queryParams["name"]; ok {
		// Es wird geprüft ob genau 1 Eintrag vorhanden ist
		if len(vmId) != 1 {
			return nil, errormsgs.HTTP_API_SERVICE_INVALID_QUERY_PARMS("validateRequestAndGetRequestData")
		}

		// Der Name wird getrimmt
		trimmedName := strings.TrimSpace(vmName[0])

		// Es wird geprüft ob es sich um einen Zulässigen Namen handelt
		if !utils.ValidateVMName(trimmedName) {
			return nil, errormsgs.HTTP_API_SERVICE_INVALID_QUERY_PARM_VMNAME("validateRequestAndGetRequestData")
		}

		// Der Wert sowie der Typ werden geschrieben
		finalValue, vlaueType = trimmedName, static.RPC_REQUEST_METHODE_VM_IDENT_NAME
	} else {
		return nil, errormsgs.HTTP_API_SERVICE_INVALID_QUERY_PARM_HAS_NOT_ID_AND_NAME("validateRequestAndGetRequestData")
	}

	// Es wird versucht die Lokale IP-Adresse einzulesen
	localip, err := getLocalIPFromRequest(r)
	if err != nil {
		// Die Aktuelle Funktion wird der Historie hinzugefügt
		err.AddCallerFunctionToHistory("validateRequestAndGetRequestData")

		// Der Fehler wird zurückgegeben
		return nil, err
	}

	// Es wird eine Verified Core IP Address (LRSAP) aus der Remote sowie aus der Lokalen IP-Addresse zu erzeugen
	LRSAP, cerr := core.ConvertLagacyIPAddressToLRSAP(r.RemoteAddr, localip)
	if cerr != nil {
		// Die Aktuelle Funktion wird der Historie hinzugefügt
		cerr.AddCallerFunctionToHistory("validateRequestAndGetRequestData")

		// Der Fehler wird zurückgegeben
		return nil, cerr
	}

	// Es wird geprüft ob ein XRequestedWith vorhanden ist,
	// sollte dieser Vorhanden sein, wird er versucht einzulesen
	var xreqwith *types.XRequestedWithData
	var xerr *types.SpecificError
	if r.Header.Get("X-Requested-With") != "" {
		xreqwith, xerr = readXRequestedWithData(r.Header.Get("X-Requested-With"))
	} else {
		xreqwith, xerr = EMPTY_X_REQUEST_WITH, nil
	}

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if xerr != nil {
		xerr.AddCallerFunctionToHistory("validateRequestAndGetRequestData")
		return nil, xerr
	}

	// Es wird geprüft ob der Referer vorhanden ist, wenn ja wird dieser ausgelesen
	var refererURL *url.URL
	if refr := r.Header.Get("Referer"); refr != "" {
		// Es wird versucht die Referer URL einzulesen
		var terr error
		refererURL, terr = url.Parse(refr)
		if terr != nil {
			return nil, errormsgs.HTTP_API_SERVICE_REFERER_READING_ERROR("validateRequestAndGetRequestData", refr)
		}
	}

	// Es wird versucht den Origin einzulesen
	var originURL *url.URL
	if orign := r.Header.Get("Origin"); orign != "" {
		// Es wird versucht die Origin Adresse einzulesen
		var terr error
		originURL, terr = url.Parse(orign)
		if terr != nil {
			return nil, errormsgs.HTTP_API_SERVICE_ORIGIN_READING_ERROR("validateRequestAndGetRequestData", orign)
		}
	}

	// Das Rückgabe Objekt wird erstellt
	returnObj := &RequestData{
		VmIdentificationMethode: vlaueType,
		TransportProtocol:       static.HTTP_JSON,
		ContentType:             contentType,
		Cookies:                 r.Cookies(),
		XRequestedWith:          xreqwith,
		Referer:                 refererURL,
		Origin:                  originURL,
		VmNameOrID:              finalValue,
		Source:                  LRSAP,
		TLS:                     r.TLS,
	}

	// Die VM ID wird zurückgegeben
	return returnObj, nil
}

// validatePOSTRequestAndGetRequestData ist eine Hilfsfunktion, die validateRequestAndGetRequestData für POST-Requests verwendet.
func validatePOSTRequestAndGetRequestData(r *http.Request, core types.CoreInterface) (*RequestData, *types.SpecificError) {
	return validateRequestAndGetRequestData("POST", r, core)
}

// validateGETRequestAndGetRequestData ist eine Hilfsfunktion, die validateRequestAndGetRequestData für GET-Requests verwendet.
func validateGETRequestAndGetRequestData(r *http.Request, core types.CoreInterface) (*RequestData, *types.SpecificError) {
	return validateRequestAndGetRequestData("GET", r, core)
}

// validateWSRequestAndGetRequestData ist eine Hilfsfunktion, die validateRequestAndGetRequestData für GET-Requests verwendet.
func validateWSRequestAndGetRequestData(r *http.Request, core types.CoreInterface) (*RequestData, *types.SpecificError) {
	return validateRequestAndGetRequestData("GET", r, core)
}

// Ließt den X-Requested-With wert aus
func readXRequestedWithData(xRequestedWith string) (*types.XRequestedWithData, *types.SpecificError) {
	switch strings.ToLower(xRequestedWith) {
	case "xmlhttprequest":
		return &types.XRequestedWithData{IsKnown: true, Value: xRequestedWith}, nil
	case "fetch":
		return &types.XRequestedWithData{IsKnown: true, Value: xRequestedWith}, nil
	case "mobileapp":
		return &types.XRequestedWithData{IsKnown: true, Value: xRequestedWith}, nil
	case "reactnative":
		return &types.XRequestedWithData{IsKnown: true, Value: xRequestedWith}, nil
	default:
		return &types.XRequestedWithData{IsKnown: false, Value: xRequestedWith}, nil
	}
}

// Versucht den Inhalt des HTTP Bodys auszulesen
func tryExtractHttpRpcBody(requestContentType types.HttpRequestContentType, body io.ReadCloser) (*RPCFunctionCall, *types.SpecificError) {
	var data *RPCFunctionCall
	switch requestContentType {
	case static.HTTP_CONTENT_CBOR:
		body, err := io.ReadAll(body)
		if err != nil {
			return nil, errormsgs.HTTP_REQUEST_READING_ERROR("extractHttpRpcBody", err.Error())
		}
		if err := cbor.Unmarshal(body, &data); err != nil {
			return nil, errormsgs.HTTP_REQUEST_READING_CBOR_ERROR("extractHttpRpcBody", err.Error())
		}
	case static.HTTP_CONTENT_JSON:
		if err := json.NewDecoder(body).Decode(&data); err != nil {
			return nil, errormsgs.HTTP_REQUEST_READING_JSON_ERROR("extractHttpRpcBody", err.Error())
		}
	default:
		return nil, errormsgs.HTTP_REQUEST_READING_UNKOWN_ERROR("extractHttpRpcBody")
	}
	return data, nil
}

// Extrhiert die Datentypen als Strings einzelner RPC Parameter
func extractParmDatatypeStrings(rpffp []RPCFunctionParameter) ([]string, *types.SpecificError) {
	// Die Datentypen der Parameter werden ausgeslesen
	dataTypeParms := make([]string, 0)
	for _, item := range rpffp {
		dataTypeParms = append(dataTypeParms, item.Type)
	}
	return dataTypeParms, nil
}

// Versucht die Mittels RPC Übermitteltens Funktionsparameter einzulesen
func tryReadFunctionParameter(data *RPCFunctionCall, foundFunction types.SharedFunctionInterface) ([]*types.FunctionParameterCapsle, *types.SpecificError) {
	// Es wird ermitelt ob die Datentypen korrekt sind
	if len(foundFunction.GetParmTypes()) != len(data.Parms) {
		return nil, errormsgs.HTTP_REQUEST_INVALID_PARAMETER_SLICE_SIZE("tryReadFunctionParameter", len(data.Parms), len(foundFunction.GetParmTypes()))
	}

	// Die Einzelnen Parameter werden geprüft und abgearbeitet
	invalidDatatypes := []*types.RPCParmeterReadingError{}
	extractedValues := make([]*types.FunctionParameterCapsle, 0)
	for x := range foundFunction.GetParmTypes() {
		// Es wird geprüft ob es sich bei dem Aktuellen Parameter um den Angeforderten Parametertypen an Stelle x handelt
		if foundFunction.GetParmTypes()[x] != data.Parms[x].Type {
			invalidDatatypes = append(invalidDatatypes, &types.RPCParmeterReadingError{Pos: x, Has: data.Parms[x].Type, Need: foundFunction.GetParmTypes()[x], SpeficMsg: "datatype dissmatch"})
			continue
		}

		// Es wird versucht den Datentypen umzuwandeln
		var goerr error
		switch data.Parms[x].Type {
		case "boolean":
			// Es wird geprüft ob es sich um ein Boolean handelt
			if reflect.TypeOf(data.Parms[x].Value).Kind() != reflect.Bool {
				invalidDatatypes = append(invalidDatatypes, &types.RPCParmeterReadingError{Pos: x, Has: data.Parms[x].Type, Need: foundFunction.GetParmTypes()[x], SpeficMsg: "datatype not same"})
				continue
			}

			// Der Datentyp wird umgewandelt
			converted, ok := data.Parms[x].Value.(bool)
			if !ok {
				invalidDatatypes = append(invalidDatatypes, &types.RPCParmeterReadingError{Pos: x, Has: data.Parms[x].Type, Need: foundFunction.GetParmTypes()[x], SpeficMsg: "data reading error"})
				continue
			}

			// Sollte ein Ungültiger Datentyp vorhanden sein, wird das Element übersprungen
			if len(invalidDatatypes) > 0 {
				continue
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
					invalidDatatypes = append(invalidDatatypes, &types.RPCParmeterReadingError{Pos: x, Has: data.Parms[x].Type, Need: foundFunction.GetParmTypes()[x], SpeficMsg: "datatype not same"})
					continue
				}

				// Der Eintrag wird erzeugt
				newEntry := &types.FunctionParameterCapsle{Value: onvertedfloat, CType: "number"}

				// Die Daten werden hinzugefügt
				extractedValues = append(extractedValues, newEntry)
				break
			}

			// Sollte ein Ungültiger Datentyp vorhanden sein, wird das Element übersprungen
			if len(invalidDatatypes) > 0 {
				continue
			}

			// Der Eintrag wird erzeugt
			newEntry := &types.FunctionParameterCapsle{Value: converted, CType: "number"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case "string":
			// Es wird geprüft ob es sich um ein String handelt
			if reflect.TypeOf(data.Parms[x].Value).Kind() != reflect.String {
				invalidDatatypes = append(invalidDatatypes, &types.RPCParmeterReadingError{Pos: x, Has: data.Parms[x].Type, Need: foundFunction.GetParmTypes()[x], SpeficMsg: "datatype dissmatch"})
				continue
			}

			// Der Datentyp wird umgewandelt
			converted, ok := data.Parms[x].Value.(string)
			if !ok {
				invalidDatatypes = append(invalidDatatypes, &types.RPCParmeterReadingError{Pos: x, Has: data.Parms[x].Type, Need: foundFunction.GetParmTypes()[x], SpeficMsg: "data reading error"})
				continue
			}

			// Sollte ein Ungültiger Datentyp vorhanden sein, wird das Element übersprungen
			if len(invalidDatatypes) > 0 {
				continue
			}

			// Der Eintrag wird erzeugt
			newEntry := &types.FunctionParameterCapsle{Value: converted, CType: "string"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case "array":
			// Es wird geprüft ob es sich um ein Slice/Array handelt
			if reflect.TypeOf(data.Parms[x].Value).Kind() != reflect.Slice {
				invalidDatatypes = append(invalidDatatypes, &types.RPCParmeterReadingError{Pos: x, Has: data.Parms[x].Type, Need: foundFunction.GetParmTypes()[x], SpeficMsg: "datatype dissmatch"})
				continue
			}

			// Der Datentyp wird umgewandelt
			converted, ok := data.Parms[x].Value.([]interface{})
			if !ok {
				invalidDatatypes = append(invalidDatatypes, &types.RPCParmeterReadingError{Pos: x, Has: data.Parms[x].Type, Need: foundFunction.GetParmTypes()[x], SpeficMsg: "data reading error"})
				continue
			}

			// Sollte ein Ungültiger Datentyp vorhanden sein, wird das Element übersprungen
			if len(invalidDatatypes) > 0 {
				continue
			}

			// Der Eintrag wird erzeugt
			newEntry := &types.FunctionParameterCapsle{Value: converted, CType: "array"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case "object":
			// Es wird geprüft ob es sich um ein Objekt handelt
			if reflect.TypeOf(data.Parms[x].Value).Kind() != reflect.Struct && reflect.TypeOf(data.Parms[x].Value).Kind() != reflect.Map {
				invalidDatatypes = append(invalidDatatypes, &types.RPCParmeterReadingError{Pos: x, Has: data.Parms[x].Type, Need: foundFunction.GetParmTypes()[x], SpeficMsg: "datatype dissmatch"})
				continue
			}

			// Sollte ein Ungültiger Datentyp vorhanden sein, wird das Element übersprungen
			if len(invalidDatatypes) > 0 {
				continue
			}

			// Der Eintrag wird erzeugt
			newEntry := &types.FunctionParameterCapsle{Value: data.Parms[x].Value, CType: "object"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case "bytes":
			// Es wird geprüft ob es sich um ein Byte Slice handet handelt
			if !utils.IsByteSlice(data.Parms[x].Value) {
				invalidDatatypes = append(invalidDatatypes, &types.RPCParmeterReadingError{Pos: x, Has: data.Parms[x].Type, Need: foundFunction.GetParmTypes()[x], SpeficMsg: "datatype dissmatch"})
				continue
			}

			// Der Datentyp wird umgewandelt
			converted, ok := data.Parms[x].Value.(string)
			if !ok {
				invalidDatatypes = append(invalidDatatypes, &types.RPCParmeterReadingError{Pos: x, Has: data.Parms[x].Type, Need: foundFunction.GetParmTypes()[x], SpeficMsg: "data reading error"})
				continue
			}

			// Es wird geprüft ob der String aus 2 teilen besteht, der este Teil gibt an welches Codec verwendet wird,
			// der Zweite teil enthält die eigentlichen Daten
			splitedValue := strings.Split("://", converted)
			if len(splitedValue) != 2 {
				invalidDatatypes = append(invalidDatatypes, &types.RPCParmeterReadingError{Pos: x, Has: data.Parms[x].Type, Need: foundFunction.GetParmTypes()[x], SpeficMsg: "invalid byte format"})
				continue
			}

			// Es wird geprüft ob es sich um ein zulässiges Codec handelt
			var decodedDataSlice []byte
			switch strings.ToLower(splitedValue[0]) {
			case "base64":
				decodedDataSlice, goerr = base64.StdEncoding.DecodeString(converted)
				if goerr != nil {
					invalidDatatypes = append(invalidDatatypes, &types.RPCParmeterReadingError{Pos: x, Has: data.Parms[x].Type, Need: foundFunction.GetParmTypes()[x], SpeficMsg: "base64 decoding error"})
					continue
				}
			case "base32":
				decodedDataSlice, goerr = base32.StdEncoding.DecodeString(converted)
				if goerr != nil {
					invalidDatatypes = append(invalidDatatypes, &types.RPCParmeterReadingError{Pos: x, Has: data.Parms[x].Type, Need: foundFunction.GetParmTypes()[x], SpeficMsg: "base64 decoding error"})
					continue
				}
			case "hex":
				decodedDataSlice, goerr = hex.DecodeString(converted)
				if goerr != nil {
					invalidDatatypes = append(invalidDatatypes, &types.RPCParmeterReadingError{Pos: x, Has: data.Parms[x].Type, Need: foundFunction.GetParmTypes()[x], SpeficMsg: "base64 decoding error"})
					continue
				}
			case "base58":
				decodedDataSlice = base58.Decode(converted)
			default:
				invalidDatatypes = append(invalidDatatypes, &types.RPCParmeterReadingError{Pos: x, Has: data.Parms[x].Type, Need: foundFunction.GetParmTypes()[x], SpeficMsg: "unkown encoding"})
				continue
			}

			// Sollte ein Ungültiger Datentyp vorhanden sein, wird das Element übersprungen
			if len(invalidDatatypes) > 0 {
				continue
			}

			// Der Eintrag wird erzeugt
			newEntry := &types.FunctionParameterCapsle{Value: decodedDataSlice, CType: "bytearray"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case "timestamp":
			// Der Datentyp wird umgewandelt
			converted, ok := data.Parms[x].Value.(int64)
			if !ok {
				invalidDatatypes = append(invalidDatatypes, &types.RPCParmeterReadingError{Pos: x, Has: data.Parms[x].Type, Need: foundFunction.GetParmTypes()[x], SpeficMsg: "data reading error"})
				continue
			}

			// Umwandlung von Unix-Zeit in time.Time
			timeObj := time.Unix(converted, 0)

			// Sollte ein Ungültiger Datentyp vorhanden sein, wird das Element übersprungen
			if len(invalidDatatypes) > 0 {
				continue
			}

			// Der Eintrag wird erzeugt
			newEntry := &types.FunctionParameterCapsle{Value: timeObj, CType: "timestamp"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		}
	}

	// Es wird geprüft ob Invalid Datatypes vorhanden sind
	if len(invalidDatatypes) > 0 {
		return nil, errormsgs.HTTP_REQUEST_INVALID_RPC_FUNCTION_DATATYPES("tryReadFunctionParameter", invalidDatatypes)
	}

	// Die Extrahierten Parameter werden zurückgegeben
	return extractedValues, nil
}

// Versucht den Gesamten Request einzulesen
func tryToReadCompleteFunctionCallFromRequest(kid types.KernelID, requestContentType types.HttpRequestContentType, body io.ReadCloser) (*RPCFunctionCall, *types.FunctionSignature, *types.SpecificError) {
	// Es wird versucht den Body einzulesen
	data, err := tryExtractHttpRpcBody(requestContentType, body)
	if err != nil {
		err.AddCallerFunctionToHistory("tryToReadCompleteFunctionCallFromRequest")
		return nil, nil, err
	}

	// Es wird geprüft ob es sich um einen Zulässigen Funktionsnamen handelt
	// Es wird geprüft, ob der Funktionsname korrekt ist
	if !utils.ValidateFunctionName(data.FunctionName) {
		responeError := errormsgs.HTTP_REQUEST_INVALID_BODY_DATA_CALLED_FUNCTION_NAME("tryToReadCompleteFunctionCallFromRequest", data.FunctionName)
		return nil, nil, responeError
	}

	// Die Datentypen der Parameter werden ausgeslesen
	dataTypeParms, err := extractParmDatatypeStrings(data.Parms)
	if err != nil {
		err.AddCallerFunctionToHistory("tryToReadCompleteFunctionCallFromRequest")
		return nil, nil, err
	}

	// Aus dem Request wird eine Funktionssignatur erzeugt
	searchedFunctionSignature := &types.FunctionSignature{
		VMID:         strings.ToLower(string(kid)),
		FunctionName: data.FunctionName,
		Params:       dataTypeParms,
		ReturnType:   data.ReturnDataType,
	}

	// Die Daten werden zurückgegeben
	return data, searchedFunctionSignature, nil
}

// Gibt die Größe der Antwort welche versendet werden soll zurück
func getResponseCapsleSize(rc *ResponseCapsle, requestContentType types.HttpRequestContentType) int {
	// Das Response Capsle Paket wird in ein Byteslice umgewandelt
	bslice, err := convertResponseCapsleToByteSlice(rc, requestContentType)
	if err != nil {
		return -1
	}

	// Die Größe wird zurückgegeben
	return len(bslice)
}

// Schreibt eine Antwort
func responseWrite(contentType types.HttpRequestContentType, w http.ResponseWriter, rpcd *ResponseCapsle) *types.SpecificError {
	// Der Passende Typ wird festgelegt
	switch contentType {
	case static.HTTP_CONTENT_CBOR:
		w.Header().Set("Content-Type", "application/cbor")
	case static.HTTP_CONTENT_JSON:
		w.Header().Set("Content-Type", "application/json")
	default:
		return errormsgs.HTTP_API_SERVICE_RESPONSE_WRITING_UNKOWN_ENCODING_ERROR("responseWrite")
	}

	// Die Antwort wird gebaut
	var response *RPCResponse
	if rpcd.Data != nil {
		w.WriteHeader(http.StatusOK)
		response = &RPCResponse{Result: "ok", Error: nil, Data: rpcd.Data}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		response = &RPCResponse{Result: "failed", Error: &rpcd.Error, Data: nil}
	}

	// Das Antwortpaket wird ferigestellt
	var bytedResponse []byte
	var err error
	switch contentType {
	case static.HTTP_CONTENT_CBOR:
		bytedResponse, err = json.Marshal(response)
	default:
		bytedResponse, err = json.Marshal(response)
	}

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if err != nil {
		return errormsgs.HTTP_API_SERVICE_RESPONSE_WRITING_ENCODING_ERROR("responseWrite", contentType, err)
	}

	// Das Paket wird gesendet
	_, err = w.Write(bytedResponse)
	if err != nil {
		if errors.Is(err, net.ErrClosed) {
			return errormsgs.HTTP_API_SERVICE_RESPONSE_WRITING_CONNECTION_CLOSED_ERROR("responseWrite", len(bytedResponse))
		} else {
			return errormsgs.HTTP_API_SERVICE_RESPONSE_WRITING_UNKOWN_ERROR("responseWrite", err)
		}
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

// Schreibt einen Fehler
func buildErrorHTTPRequestResponseAndWrite(funcname string, spefe *types.SpecificError, coreWebSession types.WebRequestBasedRPCSessionInterface, contentType types.HttpRequestContentType, w http.ResponseWriter) {
	// Sollte 'funcname' nicht "" sein, wird der Name in der Historie hinzugefügt
	if funcname != "" {
		// Die Aktuelle Funktion wird der Historie hinzugefügt
		spefe.AddCallerFunctionToHistory(funcname)
	}

	// Das Response Frame wird erzeugt
	responseFrame := &ResponseCapsle{Error: spefe.GetRemoteApiOrRpcErrorMessage()}

	// LOG
	coreWebSession.GetProcLogSession().LogPrint("HTTP-RPC", fmt.Sprintf("%s\n", spefe.GetGoProcessErrorMessage()))

	// Es wird geprüft ob die Verbindung getrennt wurde,
	// wenn nicht wird der Fehler an den Client zurückgesendet
	if coreWebSession.IsConnected() {
		// Es wird versucht den Fehler zu übertragen
		if rwerr := responseWrite(static.HTTP_CONTENT_JSON, w, responseFrame); rwerr != nil {
			// Die Aktuelle Funktion wird in der History hinzugefügt
			rwerr.AddCallerFunctionToHistory(funcname)

			// Es wird Signalisiert dass der Vorgang erfolgreich war
			coreWebSession.SignalsThatAnErrorHasOccurredWhenTheErrorIsSent(getResponseCapsleSize(responseFrame, contentType), rwerr)

			// Die Funktion wird beendet
			return
		}

		// Es wird Signalisiert dass der Fehler erfolgreich übertragen wurde
		coreWebSession.SignalThatTheErrorWasSuccessfullyTransmitted(getResponseCapsleSize(responseFrame, contentType))
	} else {
		// Es wird Signalisiert das der Fehler nicht übertragen werden konnte,
		// da die Verbindung getrennt wurde
		coreWebSession.SignalTheErrorSignalCouldNotBeTransmittedTheConnectionWasLost(getResponseCapsleSize(responseFrame, static.HTTP_CONTENT_JSON), spefe)
	}
}
