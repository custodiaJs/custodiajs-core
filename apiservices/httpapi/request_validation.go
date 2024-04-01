package httpapi

import "net/http"

func validateRequestAndGetVMID(methode string, w http.ResponseWriter, r *http.Request) (string, bool) {
	// Es wird geprüft ob es sich um die POST Methode handelt
	if r.Method != methode {
		// Set the 'Allow' header to indicate that only POST is allowed
		w.Header().Set("Allow", "POST")

		// Der Vorgang wird beendet
		return "", false
	}

	// Es wird geprüft ob der VM Query angegeben wurde
	queryParams := r.URL.Query()

	// Prüfe, ob mehr als ein Query-Parameter vorhanden ist oder der Parameter 'name' nicht existiert oder mehr als einen Wert hat.
	if len(queryParams) != 1 {
		http.Error(w, "Anfrage muss genau einen Query-Parameter 'name' mit genau einem Wert enthalten", http.StatusBadRequest)
		return "", false
	}

	// Prüfen, ob 'name' existiert und genau einen Wert hat.
	value, ok := queryParams["id"]
	if !ok || len(queryParams["id"]) != 1 {
		http.Error(w, "Anfrage muss genau einen Query-Parameter 'name' mit genau einem Wert enthalten", http.StatusBadRequest)
		return "", false
	}

	// Die ID wird geprüft
	if len(value) != 1 {
		http.Error(w, "Bad Request: Der Wert des 'vm'-Parameters muss genau 64 Zeichen lang sein", http.StatusBadRequest)
		return "", false
	}
	if len(value[0]) != 64 {
		http.Error(w, "Bad Request: Der Wert des 'vm'-Parameters muss genau 64 Zeichen lang sein", http.StatusBadRequest)
		return "", false
	}

	// Die VM ID wird zurückgegeben
	return value[0], true
}

func validatePOSTRequestAndGetVMId(w http.ResponseWriter, r *http.Request) (string, bool) {
	return validateRequestAndGetVMID("POST", w, r)
}

func validateGETRequestAndGetVMId(w http.ResponseWriter, r *http.Request) (string, bool) {
	return validateRequestAndGetVMID("GET", w, r)
}

func validateWSRequestAndGetVMId(w http.ResponseWriter, r *http.Request) (string, bool) {
	return validateRequestAndGetVMID("GET", w, r)
}

func isRequestFromIframe(r *http.Request) bool {
	referer := r.Header.Get("Referer")
	return referer != "" && referer != r.URL.String()
}

func isRequestFromJS(r *http.Request) bool {
	// Überprüfe den X-Requested-With Header, typischerweise gesetzt für AJAX-Anfragen
	requestedWith := r.Header.Get("X-Requested-With")
	return requestedWith == "XMLHttpRequest"
}
