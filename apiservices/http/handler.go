package http

import (
	"context"
	"net/http"

	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/static/errormsgs"
	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/utils/procslog"
)

// Wird ausgeführt um eine neue HTTP Sitzung zu erstellen und zu schließen
func (o *HttpApiService) newSessionHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sollte es sich um das Faveicon handeln, wird der Vorgang abgebrochen
		if r.URL.Path == "/favicon.ico" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Es wird eine neue Sitzum im Core erzeugt
		coreSession, err := o.core.GetCoreSessionManagmentUnit().NewHTTPBasesSession(r)
		if err != nil {
			// Es muss ein neuer Process Log erzeugt werden um den Fehler auszugeben
			tempLogProc := procslog.NewProcLogSession()

			// Der Fehler wird zurückgegeben
			BuildErrorRpcHttpRequestResponseAndWrite("newSessionHandler", err, nil, static.HTTP_CONTENT_JSON, tempLogProc, w)

			// Der Fehler wird zurückgegeben
			return
		}

		// Es wird ein neuer ProcLog erzeugt
		procLog := coreSession.GetProcessLog()

		// Log
		procLog.Log("New incomming http request to '%s", r.URL)

		// Beispielhafte Verwendung des method-Werts
		skipContentTypeCheck := false
		switch types.HTTP_METHOD(r.Method) {
		case static.GET:
			// Es wird geprüft ob es sich um eine Websocket verbindung handelt
			if r.Header.Get("Upgrade") == "websocket" && r.Header.Get("Connection") == "Upgrade" {
				coreSession.SetMethod(static.WEBSOCKET)
				skipContentTypeCheck = true
				break
			}

			// Es handelt sich um eine normales GET Abfrage
			coreSession.SetMethod(static.GET)
		case static.POST:
			coreSession.SetMethod(static.POST)
		case static.PUT:
			coreSession.SetMethod(static.PUT)
		case static.DELETE:
			coreSession.SetMethod(static.DELETE)
		case static.PATCH:
			coreSession.SetMethod(static.PATCH)
		case static.HEAD:
			coreSession.SetMethod(static.HEAD)
		case static.OPTIONS:
			coreSession.SetMethod(static.OPTIONS)
		case static.CONNECT:
			coreSession.SetMethod(static.CONNECT)
		case static.TRACE:
			coreSession.SetMethod(static.TRACE)
		default:
			// Die Fehlermeldung wird zurückgegeben
			BuildErrorRpcHttpRequestResponseAndWrite("", errormsgs.HTTP_API_SERVICE_UNKOWN_METHODE("newSessionHandler", r.Method), nil, static.HTTP_CONTENT_JSON, procLog, w)

			// Der Vorgang wird abgebrochen
			return
		}

		// Wenn skipContentTypeCheck 'true' ist,
		// dann wird der Content Typ wird geprüft, sollte es sich nicht um:
		// - application/json oder um application/cbor handeln,
		// wird der Vorgang mit einer Fehlermeldung abgebrochen.
		// Sollte es sich um einen zulässigen Content Typen handeln,
		// wird dieser an die Session übergeben und zwischengespeichert.
		if !skipContentTypeCheck {
			switch ctype := r.Header.Get("content-type"); ctype {
			case "application/json":
				coreSession.SetContentType(static.HTTP_CONTENT_JSON)
			case "application/cbor":
				coreSession.SetContentType(static.HTTP_CONTENT_CBOR)
			case "":
				coreSession.SetContentType(static.HTTP_CONTENT_JSON)
			default:
				// Der Fehler wird geschrieben
				unsupoortedContentType := errormsgs.HTTP_API_SERVICE_INVALID_CODEC("newSessionHandler", ctype)
				BuildErrorRpcHttpRequestResponseAndWrite("", unsupoortedContentType, nil, static.HTTP_CONTENT_JSON, procLog, w)

				// Der Vorgang wird
				return
			}
		}

		// Füge einen neuen Eintrag in den Kontext der Anfrage hinzu
		ctx := context.WithValue(r.Context(), static.CORE_SESSION_CONTEXT_KEY, coreSession)
		r = r.WithContext(ctx)

		// Falls keiner der obigen Fälle zutrifft, rufe den nächsten Handler in der Kette auf
		next.ServeHTTP(w, r)

	})
}
