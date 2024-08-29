package middleware

import (
	"net/http"

	"github.com/CustodiaJS/custodiajs-core/global/types"
)

// Wird verwendet um Die Domain einer Anfrage zu überprüfen
func ValidateRequestDomain(domains []string) MiddlewareFunction {
	return func(core types.CoreInterface, w http.ResponseWriter, r *http.Request) *types.SpecificError {
		return nil
	}
}
