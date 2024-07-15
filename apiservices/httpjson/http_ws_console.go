package httpjson

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/static/errormsgs"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func (o *HttpApiService) handleConsoleStreamWebsocket(w http.ResponseWriter, r *http.Request) {
	// Der Body wird geschlossen sobald der Vorgang beendet wurde
	defer r.Body.Close()

	// Es wird eine neue Webbasierte RPC Sitzung im Core Registriert
	coreWebSession, csError := o.core.GetCoreSessionManagmentUnit().NewWebRequestBasedRPCSession(r)
	if csError != nil {
		// Der Fehler wird zurückgesendet
		csError.AddCallerFunctionToHistory("HttpApiService->httpCallFunction")
		coreWebSession.GetProcLogSession().LogPrint("HTTP-RPC", fmt.Sprintf("%s\n", csError.GetGoProcessErrorMessage()))

		// Es wird geprüft ob die Verbindung getrennt wurde
		if coreWebSession.IsConnected() {
			// Es wird versucht den Fehler zurückzusenden
			if err := responseWrite(static.HTTP_CONTENT_JSON, w, &ResponseCapsle{Error: csError.GetRemoteApiOrRpcErrorMessage()}); err != nil {
				err.AddCallerFunctionToHistory("HttpApiService->httpCallFunction")
				coreWebSession.GetProcLogSession().LogPrint("HTTP-RPC", fmt.Sprintf("%s\n", err.GetGoProcessErrorMessage()))
				return
			}
		}

		// Rückgabe
		return
	}

	// Es wird geprüft ob es sich um die POST Methode handelt
	request, err := validateWSRequestAndGetRequestData(r, o.core)
	if err != nil {
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

				// Es wird geprüft ob die Verbindung getrennt wurde
				if coreWebSession.IsConnected() {
					// Es wird versucht den Fehler zurückzusenden
					if err := responseWrite(static.HTTP_CONTENT_JSON, w, &ResponseCapsle{Error: responeError.GetRemoteApiOrRpcErrorMessage()}); err != nil {
						err.AddCallerFunctionToHistory("HttpApiService->httpCallFunction")
						coreWebSession.GetProcLogSession().LogPrint("HTTP-RPC", fmt.Sprintf("%s\n", err.GetGoProcessErrorMessage()))
						return
					}
				}

				// Rückgabe
				return
			}
		}
	}

	// Die URL wird gelesen
	c, gerr := websocket.Accept(w, r, nil)
	if gerr != nil {
		log.Println("WebSocket Accept Error:", gerr)
		return
	}

	// Wird ausgeführt wenn die Funktion fertig ist
	defer c.Close(websocket.StatusInternalError, "Der interne Serverfehler ist aufgetreten")

	// Es wird auf eintreffende Daten gewaret
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Es wird ein neuer Console Watcher erstellt
	consoleWatcher := foundedVM.GetConsoleOutputWatcher()

	// Die Schleife wird solange ausgeführt bis die Verbindung geschlossen wurde
	for {
		// Es wird auf eine neue Ausgabe aus dem Watcher gewartet
		watcherOutput := consoleWatcher.Read()

		// Die Ausgabe wird in den Stream geschrieben
		if err := wsjson.Write(ctx, c, watcherOutput); err != nil {
			log.Println("Senden fehlgeschlagen:", err)
			break
		}
	}

	// Die Console wurde geschlossen
	c.Close(websocket.StatusNormalClosure, "")
}
