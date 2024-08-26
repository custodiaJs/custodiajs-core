package http

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

// Erzeugt einen SHA3-256 Hash aus einem Response Capsle
func ComputeSHA3_256HashFromResponseCapsle(rc interface{}, requestContentType types.HttpRequestContentType) (string, *types.SpecificError) {
	b, err := cbor.Marshal(rc)
	if err != nil {
		return "", errormsgs.HTTP_API_SERVICE_RESPONSE_WRITING_ENCODING_ERROR("ComputeSHA3_256HashFromResponseCapsle", requestContentType, err)
	}

	hexvalued := hex.EncodeToString(b)
	hash := utils.HashOfString(hexvalued)

	return hash, nil
}

// Versucht den Inhalt des HTTP Bodys auszulesen
func TryExtractHttpRpcBody(requestContentType types.HttpRequestContentType, body io.ReadCloser) (*types.RPCFunctionCall, *types.SpecificError) {
	var data *types.RPCFunctionCall
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

// Versucht die Mittels RPC Übermitteltens Funktionsparameter einzulesen
func TryReadFunctionParameter(data *types.RPCFunctionCall, foundFunction types.SharedFunctionInterface) ([]*types.FunctionParameterCapsle, *types.SpecificError) {
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
func TryToReadCompleteFunctionCallFromRequest(_ types.KernelID, requestContentType types.HttpRequestContentType, body io.ReadCloser) (*types.RPCFunctionCall, *types.SpecificError) {
	// Es wird versucht den Body einzulesen
	data, err := TryExtractHttpRpcBody(requestContentType, body)
	if err != nil {
		err.AddCallerFunctionToHistory("tryToReadCompleteFunctionCallFromRequest")
		return nil, err
	}

	// Es wird geprüft ob es sich um einen Zulässigen Funktionsnamen handelt
	// Es wird geprüft, ob der Funktionsname korrekt ist
	if !utils.ValidateFunctionName(data.FunctionName) {
		responeError := errormsgs.HTTP_REQUEST_INVALID_BODY_DATA_CALLED_FUNCTION_NAME("tryToReadCompleteFunctionCallFromRequest", data.FunctionName)
		return nil, responeError
	}

	// Die Daten werden zurückgegeben
	return data, nil
}

// Gibt die Größe der Antwort welche versendet werden soll zurück
func GetResponseSize(rc interface{}, requestContentType types.HttpRequestContentType) (int, *types.SpecificError) {
	switch requestContentType {
	case static.HTTP_CONTENT_JSON:
		b, err := json.Marshal(rc)
		if err != nil {
			return -1, errormsgs.HTTP_API_SERVICE_RESPONSE_WRITING_ENCODING_ERROR("getResponseSize", requestContentType, err)
		}
		return len(b), nil
	case static.HTTP_CONTENT_CBOR:
		b, err := cbor.Marshal(rc)
		if err != nil {
			return -1, errormsgs.HTTP_API_SERVICE_RESPONSE_WRITING_ENCODING_ERROR("convertResponseCapsleToByteSlice", requestContentType, err)
		}
		return len(b), nil
	default:
		return -1, nil
	}
}

// Schreibt eine Allgemein gültige Antwort
func HttpResponseWrite(contentType types.HttpRequestContentType, w http.ResponseWriter, response interface{}) *types.SpecificError {
	// Der Passende Typ wird festgelegt
	switch contentType {
	case static.HTTP_CONTENT_CBOR:
		w.Header().Set("Content-Type", "application/cbor")
	case static.HTTP_CONTENT_JSON:
		w.Header().Set("Content-Type", "application/json")
	default:
		return errormsgs.HTTP_API_SERVICE_RESPONSE_WRITING_UNKOWN_ENCODING_ERROR("responseWrite")
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

// Schreibt eine Antwort
func HttpRpcResponseWrite(contentType types.HttpRequestContentType, w http.ResponseWriter, rpcd *types.HttpRpcResponseCapsle) *types.SpecificError {
	// Die Antwort wird gebaut
	var response *types.RPCResponse
	if rpcd.Data != nil {
		w.WriteHeader(http.StatusOK)
		response = &types.RPCResponse{Result: "ok", Error: nil, Data: rpcd.Data}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		response = &types.RPCResponse{Result: "failed", Error: &rpcd.Error, Data: nil}
	}

	// Die Daten werden final geschrieben
	if wrerr := HttpResponseWrite(contentType, w, response); wrerr != nil {
		// Die Aktuelle Funktion wird der Historie hinzugefügt
		wrerr.AddCallerFunctionToHistory("httpRpcResponseWrite")

		// Der Fehler wird zurückgegeben
		return wrerr
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

// Schreibt einen Fehler
func BuildErrorRpcHttpRequestResponseAndWrite(funcname string, spefe *types.SpecificError, coreWebSession types.CoreHttpContextInterface, contentType types.HttpRequestContentType, logproc types.ProcessLogSessionInterface, w http.ResponseWriter) {
	// Sollte 'funcname' nicht "" sein, wird der Name in der Historie hinzugefügt
	if funcname != "" {
		// Die Aktuelle Funktion wird der Historie hinzugefügt
		spefe.AddCallerFunctionToHistory(funcname)
	}

	// Das Response Frame wird erzeugt
	responseFrame := &types.HttpRpcResponseCapsle{Error: spefe.GetRemoteApiOrRpcErrorMessage()}

	// LOG
	logproc.Log(fmt.Sprintf("%s\n", spefe.GetGoProcessErrorMessage()))

	// Es wird geprüft ob die Verbindung getrennt wurde,
	// wenn nicht wird der Fehler an den Client zurückgesendet
	if coreWebSession.IsConnected() {
		// Es wird versucht den Fehler zu übertragen
		if rwerr := HttpRpcResponseWrite(static.HTTP_CONTENT_JSON, w, responseFrame); rwerr != nil {
			// Die Aktuelle Funktion wird in der History hinzugefügt
			rwerr.AddCallerFunctionToHistory(funcname)

			// Es wird Signalisiert dass der Vorgang erfolgreich war
			size, _ := GetResponseSize(responseFrame, contentType)
			coreWebSession.SignalsThatAnErrorHasOccurredWhenTheErrorIsSent(size, rwerr)

			// Die Funktion wird beendet
			return
		}

		// Es wird Signalisiert das der Vorgang aufgrund eines Fehler abgebrochen wurde und der Fehler erfolgreich übermittelt wurde

		// Es wird Signalisiert dass der Fehler erfolgreich übertragen wurde
		size, _ := GetResponseSize(responseFrame, contentType)
		coreWebSession.SignalThatTheErrorWasSuccessfullyTransmitted(size)
	} else {
		// Es wird Signalisiert das der Vorgang aufgrund eines Fehler abgebrochen wurde und der Fehler nicht übertragen werden konnte

		// Es wird Signalisiert das der Fehler nicht übertragen werden konnte,
		// da die Verbindung getrennt wurde
		size, _ := GetResponseSize(responseFrame, static.HTTP_CONTENT_JSON)
		coreWebSession.SignalTheErrorSignalCouldNotBeTransmittedTheConnectionWasLost(size, spefe)
	}
}

// Schreibt einen Fehler
func BuildErrorHttpRequestResponseAndWrite(funcname string, spefe *types.SpecificError, coreWebSession types.CoreHttpContextInterface, logproc types.ProcessLogSessionInterface, w http.ResponseWriter) {
	// Sollte 'funcname' nicht "" sein, wird der Name in der Historie hinzugefügt
	if funcname != "" {
		// Die Aktuelle Funktion wird der Historie hinzugefügt
		spefe.AddCallerFunctionToHistory(funcname)
	}

	// Das Response Frame wird erzeugt
	responseFrame := &types.HttpRpcResponseCapsle{Error: spefe.GetRemoteApiOrRpcErrorMessage()}

	// LOG
	logproc.Debug(fmt.Sprintf("%s\n", spefe.GetGoProcessErrorMessage()))

	// Es wird geprüft ob die Verbindung getrennt wurde,
	// wenn nicht wird der Fehler an den Client zurückgesendet
	if coreWebSession.IsConnected() {
		// Das Paket wird gesendet
		bytedMsg := []byte(spefe.RemoteApiOrRpcError.ResponseMessageText)
		_, err := w.Write(bytedMsg)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				//return errormsgs.HTTP_API_SERVICE_RESPONSE_WRITING_CONNECTION_CLOSED_ERROR("responseWrite", len(bytedResponse))
			} else {
				//return errormsgs.HTTP_API_SERVICE_RESPONSE_WRITING_UNKOWN_ERROR("responseWrite", err)
			}
		}

		// Es wird Signalisiert dass der Fehler erfolgreich übertragen wurde
		coreWebSession.SignalThatTheErrorWasSuccessfullyTransmitted(len(bytedMsg))
	} else {
		// Es wird Signalisiert das der Vorgang aufgrund eines Fehler abgebrochen wurde und der Fehler nicht übertragen werden konnte

		// Es wird Signalisiert das der Fehler nicht übertragen werden konnte,
		// da die Verbindung getrennt wurde
		size, _ := GetResponseSize(responseFrame, static.HTTP_CONTENT_JSON)
		coreWebSession.SignalTheErrorSignalCouldNotBeTransmittedTheConnectionWasLost(size, spefe)
	}
}
