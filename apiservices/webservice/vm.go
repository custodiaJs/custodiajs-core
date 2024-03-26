package webservice

import (
	"encoding/json"
	"log"
	"net/http"
	"vnh1/static"
)

type SharedFunction struct {
	Name      string   `json:"name"`
	ParmTypes []string `json:"parmtypes"`
}

type SharedFunctions struct {
	Public []SharedFunction `json:"public"`
	Local  []SharedFunction `json:"local"`
}

type vmInfoResponse struct {
	Name            string          `json:"name"`
	Hash            string          `json:"hash"`
	Modules         []string        `json:"modules"`
	State           string          `json:"state"`
	SharedFunctions SharedFunctions `json:"sharedfunctions"`
}

func (o *Webservice) vmInfo(w http.ResponseWriter, r *http.Request) {
	// Setze den Content-Type der Antwort auf application/json.
	w.Header().Set("Content-Type", "application/json")

	// Es wird geprüft ob es sich um eine Zulässige GET Anfrage handelt
	vmId, invalidRequest := validateGETRequestAndGetVMId(w, r)
	if !invalidRequest {
		// Set the 'Allow' header to indicate that only POST is allowed
		w.Header().Set("Allow", "POST")

		// Send the HTTP status code 405 Method Not Allowed
		http.Error(w, "405 Method Not Allowed: Only POST method is allowed", http.StatusMethodNotAllowed)

		// Der Vorgang wird beendet
		return
	}

	// Es wird geprüft ob es sich um eine bekannte VM handelt
	foundedVM, err := o.core.GetScriptContainerVMByID(vmId)
	if err != nil {
		http.Error(w, "Die VM wurde nicht gefunden", http.StatusBadRequest)
		return
	}

	// Die Lokalen Funktionen welche geteilt wurden, werden extrahiert
	localSharedFunctions := make([]SharedFunction, 0)
	for _, item := range foundedVM.GetLocalShareddFunctions() {
		newobj := SharedFunction{Name: item.GetName(), ParmTypes: item.GetParmTypes()}
		localSharedFunctions = append(localSharedFunctions, newobj)
	}

	// Die Öffentlichen Funktionen welche geteilt wurden, werden extrahiert
	publicSharedFunctions := make([]SharedFunction, 0)
	for _, item := range foundedVM.GetPublicShareddFunctions() {
		newobj := SharedFunction{Name: item.GetName(), ParmTypes: item.GetParmTypes()}
		localSharedFunctions = append(localSharedFunctions, newobj)
	}

	// Der Status wird eingelesen
	var stateStrValue string
	switch foundedVM.GetState() {
	case static.Closed:
		stateStrValue = "CLOSED"
	case static.Running:
		stateStrValue = "RUNNING"
	case static.Starting:
		stateStrValue = "STARTING"
	default:
		stateStrValue = "UNKOWN"
	}

	// Erstelle ein Response-Objekt mit deiner Nachricht.
	response := vmInfoResponse{
		Name:    foundedVM.GetVMName(),
		Hash:    foundedVM.GetFingerprint(),
		Modules: foundedVM.GetVMModuleNames(),
		SharedFunctions: SharedFunctions{
			Public: publicSharedFunctions,
			Local:  localSharedFunctions,
		},
		State: stateStrValue,
	}

	// Schreibe die JSON-Daten in den ResponseWriter.
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Fehler beim Senden der JSON-Antwort: %v", err)
		http.Error(w, "Ein interner Fehler ist aufgetreten", http.StatusInternalServerError)
	}
}
