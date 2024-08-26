package middlewares

import (
	"net/http"

	"github.com/CustodiaJS/custodiajs-core/context"
	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/static/errormsgs"
	"github.com/CustodiaJS/custodiajs-core/types"
)

func ForceTLS(core types.CoreInterface, w http.ResponseWriter, r *http.Request) *types.SpecificError {
	// Die Core Session wird aus dem Context extrahiert,
	// sollte kein Context vorhanden sein, wird die Verbindung abgebrochen
	coreSession, ok := r.Context().Value(static.CORE_SESSION_CONTEXT_KEY).(*context.HttpContext)
	if !ok {
		return errormsgs.HTTP_API_CORE_CONTEXT_EXTRACTION_ERROR("ForceTLS")
	}

	// DEBUG
	coreSession.GetChildProcessLog("Middleware::ForceTLS")

	// Es wird geprüft ob es sich um eine TLS Verbindung handelt,
	// sollte es sich nicht um eine TLS Verbindung handeln, wird der Vorgang abgebrochen
	if r.TLS == nil {
		// Es wird der Fehler zurückgegeben dass
		return errormsgs.HTTP_API_SERVICE_HAS_NO_TLS_ENCRYPTION("ForceTLS", r.RemoteAddr)
	}

	// Das Certificate wird zwischengespeichert
	coreSession.SetTLSCertificate(r.TLS.PeerCertificates)

	// Der Vorgang wurde Ohne Fehler beeendet
	return nil
}
