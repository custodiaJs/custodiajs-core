package webservice

import (
	"encoding/json"
	"log"
	"net/http"
)

// Definiere eine Struktur, die das Format deiner JSON-Antwort repräsentiert.
type Response struct {
	Version          uint32   `json:"version"`
	RemoteConsole    bool     `json:"remoteconsole"`
	ScriptContainers []string `json:"scriptcontainers"`
}

func (o *Webservice) indexHandler(w http.ResponseWriter, r *http.Request) {
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
