package webservice

import (
	"context"
	"log"
	"net/http"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func (o *Webservice) handleConsoleStreamWebsocket(w http.ResponseWriter, r *http.Request) {
	// Es wird geprüft ob es sich um die POST Methode handelt
	vmid, isValidateRequest := validateWSRequestAndGetVMId(w, r)
	if !isValidateRequest {
		// Send the HTTP status code 405 Method Not Allowed
		http.Error(w, "405 Method Not Allowed: Only WS method is allowed", http.StatusMethodNotAllowed)

		// Der Vorgang wird beendet
		return
	}

	// Es wird geprüft ob es sich um eine bekannte VM handelt
	foundedVM, err := o.core.GetScriptContainerVMByID(vmid)
	if err != nil {
		http.Error(w, "Die VM wurde nicht gefunden", http.StatusBadRequest)
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
	c.Close(websocket.StatusNormalClosure, "")
}
