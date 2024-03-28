package webservice

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"vnh1/types"
	"vnh1/utils"
)

type RPCFunctionParameter struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type RPCFunctionCall struct {
	FunctionName string                 `json:"name"`
	Parms        []RPCFunctionParameter `json:"parms"`
}

func (o *Webservice) vmRPCHandler(w http.ResponseWriter, r *http.Request) {
	// Es wird eine neue Process Log Session erzeugt
	proc := utils.NewProcLogSession()

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
	proc.LogPrint("RPC: &[%s]: searching function '%s'\n", vmid, data.FunctionName)
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
	proc.LogPrint("RPC: &[%s]: convert function '%s' parameters\n", vmid, foundFunction.GetName())
	extractedValues := make([]*types.FunctionParameterCapsle, 0)
	for x := range foundFunction.GetParmTypes() {
		// Es wird geprüft ob es sich bei dem Angefordeten Parameter um einen zulässigen Parameter handelt
		if foundFunction.GetParmTypes()[x] != data.Parms[x].Type {
			http.Error(w, "Invalid parmtype XY", http.StatusBadRequest)
			return
		}

		// Es wird versucht den Datentypen umzuwandeln
		switch data.Parms[x].Type {
		case "boolean":
			// Der Datentyp wird umgewandelt
			converted, ok := data.Parms[x].Value.(bool)
			if !ok {
				http.Error(w, "Datatype converting error: bool", http.StatusBadRequest)
				return
			}

			// Der Eintrag wird erzeugt
			newEntry := &types.FunctionParameterCapsle{Value: converted, CapsleType: "bool"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case "number":
			// Der Datentyp wird umgewandelt
			converted, ok := data.Parms[x].Value.(uint64)
			if !ok {
				http.Error(w, "Datatype converting error: number", http.StatusBadRequest)
				return
			}

			// Der Eintrag wird erzeugt
			newEntry := &types.FunctionParameterCapsle{Value: converted, CapsleType: "number"}

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
			newEntry := &types.FunctionParameterCapsle{Value: converted, CapsleType: "string"}

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
			newEntry := &types.FunctionParameterCapsle{Value: converted, CapsleType: "array"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case "object":
			// Der Eintrag wird erzeugt
			newEntry := &types.FunctionParameterCapsle{Value: data.Parms[x].Value, CapsleType: "object"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case "bytes":
			// Der Datentyp wird umgewandelt
			converted, ok := data.Parms[x].Value.(string)
			if !ok {
				http.Error(w, "Datatype converting error: string", http.StatusBadRequest)
				return
			}

			// Die Daten werden mittels Base64 Dekodiert
			decoded, err := base64.StdEncoding.DecodeString(converted)
			if err != nil {
				http.Error(w, "Die VM wurde nicht gefunden", http.StatusBadRequest)
				return
			}

			// Der Eintrag wird erzeugt
			newEntry := &types.FunctionParameterCapsle{Value: decoded, CapsleType: "bytearray"}

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
			newEntry := &types.FunctionParameterCapsle{Value: timeObj, CapsleType: "timestamp"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		}
	}

	// Die Funktion wird aufgerufen
	proc.LogPrint("RPC: &[%s]: call functions '%s'\n", vmid, foundFunction.GetName())
	result, err := foundFunction.EnterFunctionCall(extractedValues...)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Calling error", http.StatusBadRequest)
		return
	}
	proc.LogPrintSuccs("RPC: &[%s]: function '%s' call, done\n", vmid, foundFunction.GetName())

	_, _ = foundedVM, result
}
