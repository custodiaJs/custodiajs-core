package httpjson

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/static/errormsgs"
	"github.com/CustodiaJS/custodiajs-core/utils/procslog"
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
	// Es wird eine neue Webbasierte RPC Sitzung im Core Registriert
	coreWebSession, csError := o.core.GetCoreSessionManagmentUnit().NewWebRequestBasedRPCSession(r)
	if csError != nil {
		// Der Fehler wird zurückgesendet
		csError.AddCallerFunctionToHistory("HttpApiService->httpCallFunction")
		coreWebSession.GetProcLogSession().LogPrint("HTTP-RPC", fmt.Sprintf("%s\n", csError.GetGoProcessErrorMessage()))
		responseFrame := &ResponseCapsle{Error: csError.GetRemoteApiOrRpcErrorMessage()}

		// Es wird geprüft ob die Verbindung getrennt wurde
		if coreWebSession.IsConnected() {
			// Es wird versucht den Fehler zurückzusenden
			rwerr := responseWrite(static.HTTP_CONTENT_JSON, w, responseFrame)
			if rwerr != nil {
				rwerr.AddCallerFunctionToHistory("HttpApiService->vmInfo")
				coreWebSession.GetProcLogSession().LogPrint("HTTP-RPC", fmt.Sprintf("%s\n", rwerr.GetGoProcessErrorMessage()))
				coreWebSession.SignalsThatAnErrorHasOccurredWhenTheErrorIsSent(getResponseCapsleSize(responseFrame, static.HTTP_CONTENT_JSON), rwerr)
				return
			}
		}

		// Rückgabe
		return
	}

	// Setze den Content-Type der Antwort auf application/json.
	w.Header().Set("Content-Type", "application/json")

	// Es wird geprüft ob es sich um eine Zulässige GET Anfrage handelt
	request, err := validateGETRequestAndGetRequestData(r, o.core)
	if err != nil {
		// Set the 'Allow' header to indicate that only POST is allowed
		w.Header().Set("Allow", "POST")

		// Send the HTTP status code 405 Method Not Allowed
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)

		// Der Vorgang wird beendet
		return
	}

	// Es wird geprüft ob es sich um eine bekannte VM handelt
	foundedVM, foundVm, err := o.core.GetScriptContainerVMByID(request.VmNameOrID)
	if err != nil {
		http.Error(w, "internal error", http.StatusBadRequest)
		return
	}
	if !foundVm {
		http.Error(w, "vm not found", http.StatusBadRequest)
		return
	}

	// Es wird geprüft ob es sich um eine WebRequest aus einem Webbrowser handelt,
	// wenn ja wird ermittelt ob es sich um eine Zulässige Quelle handelt,
	// wenn es sich nicht um eine zulässige Quelle handelt, wird der Vorgang abgebrochen.
	if request.XRequestedWith != nil {
		if request.XRequestedWith != EMPTY_X_REQUEST_WITH {
			if !foundedVM.IsAllowedXRequested(request.XRequestedWith) {
				responeError := errormsgs.HTTP_REQUEST_NOT_AUTHORIZED_X_SOURCE("HttpApiService->httpCallFunction", request.XRequestedWith)
				coreWebSession.GetProcLogSession().LogPrint("HTTP-RPC", fmt.Sprintf("%s\n", responeError.GetGoProcessErrorMessage()))
				responseFrame := &ResponseCapsle{Error: responeError.GetRemoteApiOrRpcErrorMessage()}

				// Es wird geprüft ob die Verbindung getrennt wurde
				if coreWebSession.IsConnected() {
					// Es wird versucht den Fehler zurückzusenden
					rwerr := responseWrite(static.HTTP_CONTENT_JSON, w, responseFrame)
					if rwerr != nil {
						rwerr.AddCallerFunctionToHistory("HttpApiService->vmInfo")
						coreWebSession.GetProcLogSession().LogPrint("HTTP-RPC", fmt.Sprintf("%s\n", rwerr.GetGoProcessErrorMessage()))
						coreWebSession.SignalsThatAnErrorHasOccurredWhenTheErrorIsSent(getResponseCapsleSize(responseFrame, static.HTTP_CONTENT_JSON), rwerr)
						return
					}
				}

				// Rückgabe
				return
			}
		}
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
	procslog.LogPrint(fmt.Sprintf("HTTP-API: retrive vm '%s' informations\n", strings.ToUpper(request.VmNameOrID)))
}
