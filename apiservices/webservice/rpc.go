package webservice

import (
	"net/http"
)

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
	_ = foundedVM

}
