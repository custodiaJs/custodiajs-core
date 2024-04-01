package httpapi

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
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
	proc.LogPrint("RPC: new incomming request from '%s'\n", r.RemoteAddr)

	// Es wird geprüft ob es sich um die POST Methode handelt
	proc.LogPrint("RPC: validate incomming rpc request from '%s'\n", r.RemoteAddr)
	vmid, isValidateRequest := validatePOSTRequestAndGetVMId(w, r)
	if !isValidateRequest {
		// Set the 'Allow' header to indicate that only POST is allowed
		w.Header().Set("Allow", "POST")

		// Send the HTTP status code 405 Method Not Allowed
		http.Error(w, "405 Method Not Allowed: Only POST method is allowed", http.StatusMethodNotAllowed)

		// Der Vorgang wird beendet
		return
	}

	// Es wird geprüft ob es sich um eine bekannte VM handelt
	proc.LogPrint("RPC: searching script container '%s'\n", vmid)
	foundedVM, err := o.core.GetScriptContainerVMByID(vmid)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Die VM wurde nicht gefunden", http.StatusBadRequest)
		return
	}

	// Es wird versucht den Datensatz einzulesen
	var data RPCFunctionCall
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Der Body wird geschlossen sobald der Vorgang beendet wurde
	defer r.Body.Close()

	// Es wird versucht die Passende Funktion zu ermitteln
	proc.LogPrint("RPC: &[%s]: searching function '%s'\n", foundedVM.GetVMName(), data.FunctionName)
	var foundFunction types.SharedFunctionInterface
	for _, item := range foundedVM.GetLocalShareddFunctions() {
		if item.GetName() == data.FunctionName {
			foundFunction = item
			break
		}
	}

	// Es wird geprüft ob eine Funktion gefunden wurde
	if foundFunction == nil {
		for _, item := range foundedVM.GetPublicShareddFunctions() {
			if item.GetName() == data.FunctionName {
				foundFunction = item
				break
			}
		}
	}

	// Sollte keine Passende Funktion gefunden werden, wird der Vorgang abgebrochen
	if foundFunction == nil {
		http.Error(w, "Unkown function", http.StatusBadRequest)
		return
	}

	// Es wird ermitelt ob die Datentypen korrekt sind
	if len(foundFunction.GetParmTypes()) != len(data.Parms) {
		http.Error(w, "Invalid total parameters", http.StatusBadRequest)
		return
	}

	// Die Einzelnen Parameter werden geprüft und abgearbeitet
	proc.LogPrint("RPC: &[%s]: convert function '%s' parameters\n", foundedVM.GetVMName(), foundFunction.GetName())
	extractedValues := make([]types.FunctionParameterBundleInterface, 0)
	for x := range foundFunction.GetParmTypes() {
		// Es wird geprüft ob es sich bei dem Angefordeten Parameter um einen zulässigen Parameter handelt
		if foundFunction.GetParmTypes()[x] != data.Parms[x].Type {
			http.Error(w, "Invalid parmtype XY", http.StatusBadRequest)
			return
		}

		// Es wird versucht den Datentypen umzuwandeln
		switch data.Parms[x].Type {
		case "boolean":
			// Es wird geprüft ob es sich um ein Boolean handelt
			if reflect.TypeOf(data.Parms[x].Value).Kind() != reflect.Bool {
				http.Error(w, "Datatype converting error: bool", http.StatusBadRequest)
				return
			}

			// Der Datentyp wird umgewandelt
			converted, ok := data.Parms[x].Value.(bool)
			if !ok {
				http.Error(w, "Datatype converting error: bool", http.StatusBadRequest)
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
					fmt.Println(data.Parms[x].Value)
					fmt.Println(reflect.TypeOf(data.Parms[x].Value))
					http.Error(w, "Datatype converting error: number", http.StatusBadRequest)
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
				http.Error(w, "Datatype converting error: string", http.StatusBadRequest)
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
				http.Error(w, "Datatype converting error: array", http.StatusBadRequest)
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
				http.Error(w, "Datatype converting error: to enocoded string", http.StatusBadRequest)
				return
			}

			// Es wird geprüft ob der String aus 2 teilen besteht, der este Teil gibt an welches Codec verwendet wird,
			// der Zweite teil enthält die eigentlichen Daten
			splitedValue := strings.Split("://", converted)
			if len(splitedValue) != 2 {
				http.Error(w, "Datatype converting error: invalid byte string coded", http.StatusBadRequest)
				return
			}

			// Es wird geprüft ob es sich um ein zulässiges Codec handelt
			var decodedDataSlice []byte
			switch strings.ToLower(splitedValue[0]) {
			case "base64":
				decodedDataSlice, err = base64.StdEncoding.DecodeString(converted)
				if err != nil {
					http.Error(w, "Die VM wurde nicht gefunden", http.StatusBadRequest)
					return
				}
			case "base32":
				decodedDataSlice, err = base32.StdEncoding.DecodeString(converted)
				if err != nil {
					http.Error(w, "Die VM wurde nicht gefunden", http.StatusBadRequest)
					return
				}
			case "hex":
				decodedDataSlice, err = hex.DecodeString(converted)
				if err != nil {
					http.Error(w, "Die VM wurde nicht gefunden", http.StatusBadRequest)
					return
				}
			case "base58":
				decodedDataSlice = base58.Decode(converted)
			default:
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
	proc.LogPrint("RPC: &[%s]: call function '%s'\n", foundedVM.GetVMName(), foundFunction.GetName())
	result, err := foundFunction.EnterFunctionCall(nil, &RpcRequest{parms: extractedValues})
	if err != nil {
		proc.LogPrint("RPC: &[%s]: call function '%s' error\n\t%s\n", foundedVM.GetVMName(), foundFunction.GetName(), err)
		http.Error(w, "Calling error", http.StatusBadRequest)
		return
	}
	proc.LogPrintSuccs("RPC: &[%s]: function '%s' call, done\n", foundedVM.GetVMName(), foundFunction.GetName())

	// Die Antwortdaten werden Extrahiert
	var responseData RPCResponseData
	if result == nil {
		responseData = RPCResponseData{DType: "null", Value: nil}
	} else if result.ExportType() == goja.Undefined().ExportType() && result.Export() == nil {
		responseData = RPCResponseData{DType: "undefined", Value: nil}
	} else {
		switch result.ExportType().Kind() {
		case reflect.Bool:
			responseData = RPCResponseData{DType: "boolean", Value: result.ToBoolean()}
		case reflect.Int64:
			responseData = RPCResponseData{DType: "number", Value: result.ToInteger()}
		case reflect.Float64:
			responseData = RPCResponseData{DType: "number", Value: result.ToFloat()}
		case reflect.String:
			responseData = RPCResponseData{DType: "string", Value: result.String()}
		case reflect.Slice:
			slicedObject, isConverted := result.Export().([]interface{})
			if !isConverted {
				http.Error(w, "invalid object datatype, slice", http.StatusBadRequest)
				return
			}
			responseData = RPCResponseData{DType: "array", Value: slicedObject}
		case reflect.Map:
			mapObjected, isConverted := result.Export().(map[string]interface{})
			if !isConverted {
				http.Error(w, "invalid object datatype, object", http.StatusBadRequest)
				return
			}
			responseData = RPCResponseData{DType: "object", Value: mapObjected}
		case reflect.Func:
			http.Error(w, "function return not allowed in rpc post request", http.StatusBadRequest)
			return
		default:
			fmt.Println(result.ExportType().Kind())
			http.Error(w, "Calling error a", http.StatusBadRequest)
			return
		}
	}

	// Die Antwort wird erzeugt
	response := &RPCResponse{Result: "success", Data: responseData}
	bytedResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Calling error", http.StatusBadRequest)
		return
	}

	// Die Daten werden zurückgesendet
	proc.LogPrint("RPC: &[%s]: sending function '%s' call response\n", foundedVM.GetVMName(), data.FunctionName)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytedResponse)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
	proc.LogPrint("RPC: &[%s]: done\n", foundedVM.GetVMName())
}
