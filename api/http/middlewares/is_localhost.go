package middlewares

import (
	"net/http"

	"github.com/CustodiaJS/custodiajs-core/global/types"
)

// Wird verwendet um Localhost Anfragen zu Validieren
func IsLocalhostPOSTRequest(core types.CoreInterface, w http.ResponseWriter, r *http.Request) *types.SpecificError {
	return nil
}
