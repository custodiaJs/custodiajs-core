package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
	"vnh1/types"
)

func (o *HttpApiService) indexHandler(w http.ResponseWriter, r *http.Request) {
	// Es werden alle Script Container extrahiert
	scriptContainer := o.core.GetAllActiveScriptContainerIDs()

	// Erstelle ein Response-Objekt mit deiner Nachricht.
	response := Response{Version: 1000000000, ScriptContainers: scriptContainer}

	// Setze den Content-Type der Antwort auf application/json.
	w.Header().Set("Content-Type", "application/json")

	// Schreibe die JSON-Daten in den ResponseWriter.
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Fehler beim Senden der JSON-Antwort: %v", err)
		http.Error(w, "Ein interner Fehler ist aufgetreten", http.StatusInternalServerError)
	}
}

func (o *HttpApiService) vmInfo(w http.ResponseWriter, r *http.Request) {
	// Setze den Content-Type der Antwort auf application/json.
	w.Header().Set("Content-Type", "application/json")

	// Es wird geprüft ob es sich um eine Zulässige GET Anfrage handelt
	request, err := validateGETRequestAndGetRequestData(r)
	if err != nil {
		// Set the 'Allow' header to indicate that only POST is allowed
		w.Header().Set("Allow", "POST")

		// Send the HTTP status code 405 Method Not Allowed
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)

		// Der Vorgang wird beendet
		return
	}

	// Es wird geprüft ob es sich um eine bekannte VM handelt
	foundedVM, err := o.core.GetScriptContainerVMByID(request.VmId)
	if err != nil {
		http.Error(w, "Die VM wurde nicht gefunden", http.StatusBadRequest)
		return
	}

	// Es wird geprüft ob es sich um eine WebRequest aus einem Webbrowser handelt,
	// wenn ja wird ermittelt ob es sich um eine Zulässige Quelle handelt
	requestHttpSource := getRefererOrXRequestedWith(request)
	if hasRefererOrXRequestedWith(request) && !foundedVM.ValidateRPCRequestSource(requestHttpSource) {
		// Der Vorgang wird abgebrochen, es handelt sich nicht nicht um eine zulässige Quelle
		http.Error(w, "Unzulässige Quelle", http.StatusBadRequest)
		return
	}

	// Die Lokalen Funktionen welche geteilt wurden, werden extrahiert
	localSharedFunctions := make([]SharedFunction, 0)
	for _, item := range foundedVM.GetLocalSharedFunctions() {
		newobj := SharedFunction{Name: item.GetName(), ParmTypes: item.GetParmTypes()}
		localSharedFunctions = append(localSharedFunctions, newobj)
	}

	// Die Öffentlichen Funktionen welche geteilt wurden, werden extrahiert
	publicSharedFunctions := make([]SharedFunction, 0)
	for _, item := range foundedVM.GetPublicSharedFunctions() {
		newobj := SharedFunction{Name: item.GetName(), ParmTypes: item.GetParmTypes()}
		localSharedFunctions = append(localSharedFunctions, newobj)
	}

	// Der Status wird eingelesen
	var stateStrValue string
	switch foundedVM.GetState() {
	case types.Closed:
		stateStrValue = "CLOSED"
	case types.Running:
		stateStrValue = "RUNNING"
	case types.Starting:
		stateStrValue = "STARTING"
	case types.StillWait:
		stateStrValue = "STILL_WAIT"
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
