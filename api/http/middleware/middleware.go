package middleware

import (
	"net/http"

	"github.com/CustodiaJS/custodiajs-core/global/types"
)

func RequestMiddleware(nextHandlerFunction httpRequest, midlewareHandlers MiddlewareFunctionList, core types.CoreInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Es wird eine neue Sitzung im Core Registriert
		for _, item := range midlewareHandlers {
			if checkError := item(core, w, r); checkError != nil {

				return
			}
		}
		// Weiterleitung an den nächsten Handler
		next := http.HandlerFunc(nextHandlerFunction)
		next.ServeHTTP(w, r)
	})
}

func GlobalMiddleware(next http.Handler, globalHandlers MiddlewareFunctionList, core types.CoreInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Es wird eine neue Sitzung im Core Registriert
		for _, item := range globalHandlers {
			if checkError := item(core, w, r); checkError != nil {
				// Der Vorgang wird abgebrochen, der Fehler wird zurückgegeben
				return
			}
		}

		// Nächster Handler
		next.ServeHTTP(w, r)
	})
}
