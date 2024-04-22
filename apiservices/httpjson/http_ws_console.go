package httpjson

import (
	"context"
	"log"
	"net/http"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func (o *HttpApiService) handleConsoleStreamWebsocket(w http.ResponseWriter, r *http.Request) {
	// Es wird geprüft ob es sich um die POST Methode handelt
	request, err := validateWSRequestAndGetRequestData(r)
	if err != nil {
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

	// Die URL wird gelesen
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("WebSocket Accept Error:", err)
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
