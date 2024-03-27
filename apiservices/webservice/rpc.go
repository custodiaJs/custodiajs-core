package webservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"vnh1/static"
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
	// Es wird geprüft ob es sich um die POST Methode handelt
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
	foundedVM, err := o.core.GetScriptContainerVMByID(vmid)
	if err != nil {
		http.Error(w, "Die VM wurde nicht gefunden", http.StatusBadRequest)
		return
	}

	// Es wird versucht den Datensatz einzulesen
	var data RPCFunctionCall
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Es wird versucht die Passende Funktion zu ermitteln
	var foundFunction static.SharedFunctionInterface
	for _, item := range foundedVM.GetLocalShareddFunctions() {
		if item.GetName() == data.FunctionName {
			foundFunction = item
			break
		}
	}
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
		http.Error(w, "Die VM wurde nicht gefunden", http.StatusBadRequest)
		return
	}

	// Es wird ermitelt ob die Datentypen korrekt sind
	if len(foundFunction.GetParmTypes()) != len(data.Parms) {
		http.Error(w, "Die VM wurde nicht gefunden", http.StatusBadRequest)
		return
	}

	// Die Einzelnen Parameter werden geprüft und abgearbeitet
	extractedValues := make([]interface{}, 0)
	for x := range foundFunction.GetParmTypes() {
		// Es wird geprüft ob es sich bei dem Angefordeten Parameter um einen zulässigen Parameter handelt
		if foundFunction.GetParmTypes()[x] != data.Parms[x].Type {
			http.Error(w, "Die VM wurde nicht gefunden", http.StatusBadRequest)
			return
		}

		// Es wird versucht den Datentypen umzuwandeln
		switch data.Parms[x].Type {
		case "boolean":
			converted, ok := data.Parms[x].Value.(bool)
			if !ok {
				http.Error(w, "Die VM wurde nicht gefunden", http.StatusBadRequest)
				return
			}
			extractedValues = append(extractedValues, converted)
		case "number":
			converted, ok := data.Parms[x].Value.(uint64)
			if !ok {
				http.Error(w, "Die VM wurde nicht gefunden", http.StatusBadRequest)
				return
			}
			extractedValues = append(extractedValues, converted)
		case "string":
			converted, ok := data.Parms[x].Value.(string)
			if !ok {
				http.Error(w, "Die VM wurde nicht gefunden", http.StatusBadRequest)
				return
			}
			extractedValues = append(extractedValues, converted)
		case "array":
			converted, ok := data.Parms[x].Value.([]interface{})
			if !ok {
				http.Error(w, "Die VM wurde nicht gefunden", http.StatusBadRequest)
				return
			}
			extractedValues = append(extractedValues, converted)
		case "object":
			extractedValues = append(extractedValues, data.Parms[x].Value)
		}
	}

	// Die Funktion wird aufgerufen
	result, err := foundFunction.EnterFunctionCall(extractedValues...)
	if err != nil {
		http.Error(w, "Die VM wurde nicht gefunden", http.StatusBadRequest)
		return
	}

	_, _ = foundedVM, result
}
