package httpjson

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/static/errormsgs"
	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/utils"
	"github.com/CustodiaJS/custodiajs-core/utils/procslog"
)

func (o *HttpApiService) httpCallFunction(w http.ResponseWriter, r *http.Request) {
	// Es wird eine neue Webbasierte RPC Sitzung im Core Registriert
	coreWebSession, stmerr := o.core.GetCoreSessionManagmentUnit().NewWebRequestBasedRPCSession(r)
	if stmerr != nil {
		// Die Aktuelle Funktion wird der Historie hinzugefügt
		stmerr.AddCallerFunctionToHistory("HttpApiService->httpCallFunction")

		// LOG
		coreWebSession.GetProcLogSession().LogPrint("HTTP-RPC", fmt.Sprintf("%s\n", stmerr.GetGoProcessErrorMessage()))

		// Das Response Frame wird erzeugt
		responseFrame := &ResponseCapsle{Error: stmerr.GetRemoteApiOrRpcErrorMessage()}

		// Es wird geprüft ob die Verbindung getrennt wurde
		if coreWebSession.IsConnected() {
			// Es wird versucht den Fehler zurückzusenden
			if err := responseWrite(static.HTTP_CONTENT_JSON, w, responseFrame); err != nil {
				// Die Aktuelle Funktion wird der Historie hinzugefügt
				err.AddCallerFunctionToHistory("HttpApiService->httpCallFunction")

				// LOG
				coreWebSession.GetProcLogSession().LogPrint("HTTP-RPC", fmt.Sprintf("%s\n", err.GetGoProcessErrorMessage()))

				// Der HTTP Boddy wird geschlossen
				r.Body.Close()

				// Rückgabe
				return
			}
		}

		// Der HTTP Boddy wird geschlossen
		r.Body.Close()

		// Rückgabe
		return
	}

	// Diese Funktion wird ausgeführt sobald ein Return durchgeführt wird
	defer closeHTTPRequest(r, coreWebSession)

	// validateRequestAndGetRequestData überprüft einen HTTP-Request auf Gültigkeit und gibt ein RequestData-Objekt zurück.
	// Es validiert die HTTP-Methode auf Übereinstimmung mit der angegebenen Methode (POST), die TLS-Verbindung,
	// den Content-Type (JSON oder CBOR), und den Query-Parameter 'id' auf Existenz und korrekten Hexadezimalwert.
	// Bei erfolgreicher Validierung wird ein RequestData-Objekt mit den extrahierten Daten zurückgegeben.
	procslog.ProcFormatConsoleText(coreWebSession.GetProcLogSession(), "HTTP-RPC", types.VALIDATE_INCOMMING_REMOTE_FUNCTION_CALL_REQUEST_FROM, r.RemoteAddr) // Logausgabe
	request, vpragrderr := validatePOSTRequestAndGetRequestData(r, o.core)                                                                                   //  Validierung der Anfrage
	if vpragrderr != nil {
		// Der Fehler wird übermittelt
		buildErrorHTTPRequestResponseAndWrite("HttpApiService->httpCallFunction", vpragrderr, coreWebSession, static.HTTP_CONTENT_JSON, w)

		// Rückgabe
		return
	}

	// Es wird versucht die Passende VM zu ermitteln
	var vmSearchError *types.SpecificError
	var foundedVM types.VmInterface
	var foundVM bool
	switch request.VmIdentificationMethode {
	case static.RPC_REQUEST_METHODE_VM_IDENT_ID:
		// Log
		procslog.ProcFormatConsoleText(coreWebSession.GetProcLogSession(), "HTTP-RPC", types.DETERMINE_THE_SCRIPT_CONTAINER, strings.ToUpper(request.VmNameOrID))

		// Es wird versucht die Passende VM anhand der ID zu ermitteln
		foundedVM, foundVM, vmSearchError = o.core.GetScriptContainerVMByID(request.VmNameOrID)
	case static.RPC_REQUEST_METHODE_VM_IDENT_NAME:
		// Log
		procslog.ProcFormatConsoleText(coreWebSession.GetProcLogSession(), "HTTP-RPC", types.DETERMINE_THE_SCRIPT_CONTAINER, strings.ToUpper(request.VmNameOrID))

		// Es wird versucht die Passende VM anhand des Namen zu ermitteln
		foundedVM, foundVM, vmSearchError = o.core.GetScriptContainerByVMName(request.VmNameOrID)
	default:
		// Es wird ein Fehler zurückgegeben: Es wurde keine gültige Methode angegeben
		vmSearchError = errormsgs.HTTP_API_SERVICE_REQUEST_HAS_UNKOWN_VM_IDENT_METHODE("")
	}

	// Es wird geprüft ob ein Fehler aufgetreten ist beim Ermitteln der Passenden VM
	if vmSearchError != nil {
		// Der Fehler wird übermittelt
		buildErrorHTTPRequestResponseAndWrite("HttpApiService->httpCallFunction", vmSearchError, coreWebSession, request.ContentType, w)

		// Rückgabe
		return
	}

	// Es wird geprüft ob eine VM Gefunden wurde
	if !foundVM {
		// Es wird ein neuer Fehler erzeugt
		responeError := errormsgs.HTTP_API_SERVICE_REQUEST_VM_NOT_FOUND("HttpApiService->httpCallFunction", request.VmNameOrID, request.VmIdentificationMethode)

		// Der Fehler wird übermittelt
		buildErrorHTTPRequestResponseAndWrite("", responeError, coreWebSession, request.ContentType, w)

		// Rückgabe
		return
	}

	// Wenn es sich nicht um eine Lokale Adresse handelt,
	// wird geprüft ob die Quelle berechtigt ist diesen Knotenpunkt zu verwenden.
	if !o.core.LRSAPSourceIsAllowed(request.Source) {
		// Es wird ein neuer Fehler erzeugt
		responeError := errormsgs.HTTP_REQUEST_NOT_ALLOWED_SOURCE_IP("HttpApiService->httpCallFunction")

		// Der Fehler wird übermittelt
		buildErrorHTTPRequestResponseAndWrite("", responeError, coreWebSession, request.ContentType, w)

		// Rückgabe
		return
	}

	// Es wird geprüft ob es sich um eine WebRequest aus einem Webbrowser handelt,
	// wenn ja wird ermittelt ob es sich um eine Zulässige Quelle handelt,
	// wenn es sich nicht um eine zulässige Quelle handelt, wird der Vorgang abgebrochen.
	if request.XRequestedWith != nil {
		if request.XRequestedWith != EMPTY_X_REQUEST_WITH {
			if !foundedVM.IsAllowedXRequested(request.XRequestedWith) {
				// Es wird ein neuer Fehler erzeugt
				responeError := errormsgs.HTTP_REQUEST_NOT_AUTHORIZED_X_SOURCE("HttpApiService->httpCallFunction", request.XRequestedWith)

				// Der Fehler wird übermittelt
				buildErrorHTTPRequestResponseAndWrite("", responeError, coreWebSession, request.ContentType, w)

				// Rückgabe
				return
			}
		}
	}

	// Es wird versucht die Kerndaten aus dem Request einztulesen
	data, searchedFunctionSignature, ttcrfcerr := tryToReadCompleteFunctionCallFromRequest(foundedVM.GetKId(), request.ContentType, r.Body)
	if ttcrfcerr != nil {
		// Der Fehler wird übermittelt
		buildErrorHTTPRequestResponseAndWrite("HttpApiService->httpCallFunction", ttcrfcerr, coreWebSession, request.ContentType, w)

		// Rückgabe
		return
	}

	// Es wird versucht die Funktion anhand ihrer Signatur zu ermitteln
	procslog.ProcFormatConsoleText(coreWebSession.GetProcLogSession(), "HTTP-PRC", types.DETERMINE_THE_FUNCTION, foundedVM.GetVMName(), data.FunctionName)
	foundFunction, hasFound, gsfbserr := foundedVM.GetSharedFunctionBySignature(static.LOCAL, searchedFunctionSignature)
	if gsfbserr != nil {
		// Der Fehler wird übermittelt
		buildErrorHTTPRequestResponseAndWrite("HttpApiService->httpCallFunction", gsfbserr, coreWebSession, request.ContentType, w)

		// Rückgabe
		return
	}

	// Sollte keine Passende Funktion gefunden werden, wird der Vorgang abgebrochen
	if !hasFound {
		// Es wird ein neuer Fehler erzeugt
		responeError := errormsgs.HTTP_API_SERVICE_REQUEST_VM_FUNCTION_NOT_FOUND_ERROR("HttpApiService->httpCallFunction", request.VmNameOrID, request.VmIdentificationMethode, searchedFunctionSignature)

		// Der Fehler wird übermittelt
		buildErrorHTTPRequestResponseAndWrite("", responeError, coreWebSession, request.ContentType, w)

		// Rückgabe
		return
	}

	// Die Einzelnen Parameter werden geprüft und abgearbeitet
	extractedValues, trfperr := tryReadFunctionParameter(data, foundFunction)
	if trfperr != nil {
		// Der Fehler wird übermittelt
		buildErrorHTTPRequestResponseAndWrite("HttpApiService->httpCallFunction", trfperr, coreWebSession, request.ContentType, w)

		// Rückgabe
		return
	}

	// Das HTTP Request Objekt wird erstellt
	requestHttpObject := &types.HttpRpcRequest{
		IsConnected:      coreWebSession.IsConnected,
		ContentLength:    r.ContentLength,
		PostForm:         r.PostForm,
		Header:           r.Header,
		Host:             r.Host,
		Form:             r.Form,
		Proto:            r.Proto,
		RemoteAddr:       r.RemoteAddr,
		RequestURI:       r.RequestURI,
		TLS:              r.TLS,
		TransferEncoding: r.TransferEncoding,
		URL:              r.URL,
		Cookies:          r.Cookies(),
		UserAgent:        r.UserAgent(),
	}

	// Diese Funktion nimmt die Antwort entgegen
	resolveFunction := func(response *types.FunctionCallReturn) error {
		// Es wird geprüft ob das Response Null ist, wenn ja wird ein Panic ausgelöst
		if response == nil {
			panic("http rpc function call response is null, critical error")
		}

		// Es wird geprüft ob der Vorgang bereits abgeschlossen wurde
		if coreWebSession.GetReturnChan().IsClosed() && coreWebSession.IsConnected() {
			return utils.MakeHttpConnectionIsClosedError()
		} else if coreWebSession.IsConnected() {
			return utils.MakeAlreadyAnsweredRPCRequestError()
		}

		// Die Antwort wird geschrieben
		coreWebSession.GetReturnChan().WriteAndClose(response)

		// Es ist kein Fehler aufgetreten
		return nil
	}

	// Das Request Objekt wird erzeugt
	requestObject := &types.RpcRequest{
		Parms:         extractedValues,
		RpcRequest:    request,
		ProcessLog:    coreWebSession.GetProcLogSession(),
		RequestType:   static.HTTP_REQUEST,
		HttpRequest:   requestHttpObject,
		WriteResponse: resolveFunction,
	}

	// Definiert die Lambda Funktion welche verwendet wird um das Script Panic Sicher auszuführen,
	// durch das Verwenden diese Funktion, soll verhindert werden dass ein Ausführungs Golang Panic dazu führt, dass das gesamte Programm beendet wird
	safeBoxFunctionCall := func(requestObject *types.RpcRequest) (err *types.SpecificError) {
		// Fängt Panics ab
		defer func() {
			if r := recover(); r != nil {
				err = errormsgs.HTTP_REQUEST_FUNCTION_CALL_PANIC("HttpApiService->httpCallFunction", fmt.Sprintf("%s", r))
			}
		}()

		// Ruft die Funktion innerhalb der V8Go Engine auf
		goerr := foundFunction.EnterFunctionCall(requestObject)
		if goerr != nil {
			goerr.AddCallerFunctionToHistory("HttpApiService->httpCallFunction")
			return goerr
		}

		// Gibt nichts zurück
		err = nil
		return
	}

	// Die Funktion wird aufgerufen
	if sfbError := safeBoxFunctionCall(requestObject); sfbError != nil {
		// Der Fehler wird übermittelt
		buildErrorHTTPRequestResponseAndWrite("HttpApiService->httpCallFunction", sfbError, coreWebSession, request.ContentType, w)

		// Rückgabe
		return
	}

	// Es wird auf das Ergebniss gewartet
	result, ok := coreWebSession.GetReturnChan().Read()
	if result == nil || !ok {
		// Es wird ein neuer Fehler erzeugt
		responeError := errormsgs.HTTP_API_SERVICE_REQUEST_INTERNAL_ERROR_BY_READING_RETURN_CHAN("HttpApiService->httpCallFunction")

		// Der Fehler wird übermittelt
		buildErrorHTTPRequestResponseAndWrite("", responeError, coreWebSession, request.ContentType, w)

		// Der Vorgang wird beendet
		return
	}

	// Die Antwort wird gebaut
	var responseData *ResponseCapsle
	if result.State == "ok" {
		// Erstellt ein Array welches alle Antwortdaten enthält
		dt := make([]*RPCResponseData, 0)

		// Die Antwortdaten werden abgearbeitet
		for _, item := range result.Return {
			dt = append(dt, &RPCResponseData{DType: item.Type, Value: item.Value})
		}

		// Das Response Capsle wird erzeugt
		responseData = &ResponseCapsle{Data: dt}
	} else if result.State == "failed" {
		// Die Response Capsle wird erzeugt
		responseData = &ResponseCapsle{Error: result.Error}
	} else if result.State == "exception" {
		// Die Response Capsle wird erzeugt
		responseData = &ResponseCapsle{Error: result.Error}
	} else {
		// Die Response Capsle wird erzeugt
		responseData = &ResponseCapsle{Error: "unkown return state"}
	}

	// Die Größe das Paketes wird ermittelt, der Hash des Paketes wird berechnet
	responseSize := getResponseCapsleSize(responseData, request.ContentType)
	if responseSize < 0 {
		// Es wird ein neuer Fehler erzeugt
		responeError := errormsgs.HTTP_API_SERVICE_INTERNAL_ERROR_BY_EMITTING_CAPSLE_SIZE("HttpApiService->httpCallFunction")

		// Der Fehler wird übermittelt
		buildErrorHTTPRequestResponseAndWrite("", responeError, coreWebSession, request.ContentType, w)

		// Der Vorgang wird beendet
		return
	}

	// Es wird ein Hash aus dem Response Capsle erzeugt
	responsePackageHash, packageComputingError := computeSHA3_256HashFromResponseCapsle(responseData, request.ContentType)
	if packageComputingError != nil {
		// Der Fehler wird übermittelt
		buildErrorHTTPRequestResponseAndWrite("HttpApiService->httpCallFunction", packageComputingError, coreWebSession, request.ContentType, w)

		// Der Vorgang wird beendet
		return
	}

	// Es wird versucht die Daten zurückzusenden
	if frwrerr := responseWrite(request.ContentType, w, responseData); frwrerr != nil {
		// Die Aktuelle Funktion wird der Historie hinzugefügt
		frwrerr.AddCallerFunctionToHistory("HttpApiService->httpCallFunction")

		// Es wird Signalisiert dass die Daten nicht übermittelt werden konnten
		coreWebSession.SignalTheResponseCouldNotBeSent(responseSize, frwrerr)
		//result.Reject()

		// Log
		coreWebSession.GetProcLogSession().LogPrint("HTTP-RPC: &[%s]: call function response sending '%s' error\n\t%s\n", foundedVM.GetVMName(), foundFunction.GetName(), frwrerr)

		// Der Vorgang wird geschlossen
		return
	}

	// Es wird Signalisiert dass die Daten erfolgreich übermittelt wurden
	coreWebSession.SignalTheResponseWasTransmittedSuccessfully(responseSize, responsePackageHash)
	//result.Resolve()

	// Log
	procslog.ProcFormatConsoleText(coreWebSession.GetProcLogSession(), "HTTP-PRC", types.RPC_CALL_DONE_RESPONSE, fmt.Sprintf("%d", getResponseCapsleSize(responseData, request.ContentType)))

	// Signalisiert dass die Funktion erfolgreich durchgelaufen ist und nun als nächstes der Return ansteht
	coreWebSession.Done()
}
