package middlewares

import (
	"net/http"
	"net/url"

	"github.com/CustodiaJS/custodiajs-core/context"
	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/static/errormsgs"
	"github.com/CustodiaJS/custodiajs-core/types"
)

func ProxyOrBrowserRequestValidation(core types.CoreInterface, w http.ResponseWriter, r *http.Request) *types.SpecificError {
	// Die Core Session wird aus dem Context extrahiert,
	// sollte kein Context vorhanden sein, wird die Verbindung abgebrochen
	coreSession, ok := r.Context().Value(static.CORE_SESSION_CONTEXT_KEY).(*context.HttpContext)
	if !ok {
		return nil
	}

	// Prüfe, ob ein XRequestedWith vorhanden ist, sollte dieser Vorhanden sein, wird er versucht einzulesen
	if r.Header.Get("X-Requested-With") != "" {
		// Der X-Request With header wird eingelesen
		xreqwith, xerr := readXRequestedWithData(r.Header.Get("X-Requested-With"))
		if xerr != nil {
			return nil
		}

		// Der Header wird in der Sitzung abgespeichert
		coreSession.SetXRequestedWith(xreqwith)
	}

	// Prüfe, ob der Referer vorhanden ist, wenn ja wird dieser eingelesen
	if refr := r.Header.Get("Referer"); refr != "" {
		// Es wird versucht die Referer URL einzulesen
		refererURL, terr := url.Parse(refr)
		if terr != nil {
			return errormsgs.HTTP_API_SERVICE_REFERER_READING_ERROR("ProxyOrBrowserRequestValidation", refr)
		}

		// Der Refrer wird übergeben
		coreSession.SetReferer(refererURL)
	}

	// Prüfe, ob der Origin vorhanden ist, wenn ja wird dieser eingelesen
	if orign := r.Header.Get("Origin"); orign != "" {
		// Es wird versucht die Origin Adresse einzulesen
		originURL, terr := url.Parse(orign)
		if terr != nil {
			return errormsgs.HTTP_API_SERVICE_ORIGIN_READING_ERROR("ProxyOrBrowserRequestValidation", orign)
		}

		// Die Origin URL wird an die Core Sitzung übergeben
		coreSession.SetOrigin(originURL)
	}

	// validateRequestAndGetRequestData überprüft einen HTTP-Request auf Gültigkeit und gibt ein RequestData-Objekt zurück.
	// Es validiert die HTTP-Methode auf Übereinstimmung mit der angegebenen Methode (POST), die TLS-Verbindung,
	// den Content-Type (JSON oder CBOR), und den Query-Parameter 'id' auf Existenz und korrekten Hexadezimalwert.
	/* Bei erfolgreicher Validierung wird ein RequestData-Objekt mit den extrahierten Daten zurückgegeben.
	procslog.ProcFormatConsoleText(coreSession.GetProcLogSession(), "HTTP-RPC", types.VALIDATE_INCOMMING_REMOTE_FUNCTION_CALL_REQUEST_FROM, r.RemoteAddr) // Logausgabe
	request, vpragrderr := validatePOSTRequestAndGetRequestData(r, coreSession.GetCore())                                                                 //  Validierung der Anfrage
	if vpragrderr != nil {
		// Die Aktuelle Funktion wird dem Aufrufenden Fehler hinzugefügt
		vpragrderr.AddCallerFunctionToHistory("ValidateRPCRequest")

		// Rückgabe des Fehlers
		return vpragrderr
	}*/

	return nil
}
