package httpjson

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"vnh1/static"
	"vnh1/utils/procslog"
)

func (o *HttpApiService) indexHandler(w http.ResponseWriter, r *http.Request) {
	// Es werden alle Script Container extrahiert
	scriptContainers := o.core.GetAllActiveScriptContainerIDs()
	ucscontainers := []string{}
	for _, item := range scriptContainers {
		ucscontainers = append(ucscontainers, strings.ToUpper(item))
	}

	// Erstelle ein Response-Objekt mit deiner Nachricht.
	response := Response{Version: uint32(static.C_VESION), ScriptContainers: ucscontainers}

	// Setze den Content-Type der Antwort auf application/json.
	w.Header().Set("Content-Type", "application/json")

	// Log
	procslog.LogPrint("HTTP-API: retrive host informations\n")

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
	foundedVM, foundVm, err := o.core.GetScriptContainerVMByID(request.VmId)
	if err != nil {
		http.Error(w, "internal error", http.StatusBadRequest)
		return
	}
	if !foundVm {
		http.Error(w, "vm not found", http.StatusBadRequest)
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
	sharedFunctions := make([]SharedFunction, 0)
	for _, item := range foundedVM.GetAllSharedFunctions() {
		newobj := SharedFunction{Name: item.GetName(), ParmTypes: item.GetParmTypes(), ReturnDatatype: item.GetReturnDatatype()}
		sharedFunctions = append(sharedFunctions, newobj)
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
	case static.StillWait:
		stateStrValue = "STILL_WAIT"
	default:
		stateStrValue = "UNKOWN"
	}

	// Erstelle ein Response-Objekt mit deiner Nachricht.
	response := vmInfoResponse{
		Name:            foundedVM.GetVMName(),
		Id:              string(foundedVM.GetFingerprint()),
		SharedFunctions: sharedFunctions,
		State:           stateStrValue,
	}

	// Schreibe die JSON-Daten in den ResponseWriter.
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Fehler beim Senden der JSON-Antwort: %v", err)
		http.Error(w, "Ein interner Fehler ist aufgetreten", http.StatusInternalServerError)
	}

	// Log
	procslog.LogPrint(fmt.Sprintf("HTTP-API: retrive vm '%s' informations\n", strings.ToUpper(request.VmId)))
}
