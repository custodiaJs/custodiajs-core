package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/CustodiaJS/custodiajs-core/context"
	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/static/errormsgs"
	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/utils"
	"github.com/CustodiaJS/custodiajs-core/utils/procslog"
)

// Nimmt Anfragen für "/" entgegen
func (o *HttpApiService) httpIndex(w http.ResponseWriter, r *http.Request) {
	// Die Core Session wird aus dem Context extrahiert,
	// sollte kein Context vorhanden sein, wird die Verbindung abgebrochen
	coreSession, ok := r.Context().Value(static.CORE_SESSION_CONTEXT_KEY).(*context.HttpContext)
	if !ok {
		// Es muss ein neuer Process Log erzeugt werden um den Fehler auszugeben
		tempLogProc := procslog.NewProcLogSession()

		// Der Fehler wird zurückgegeben
		BuildErrorHttpRequestResponseAndWrite("", errormsgs.HTTP_API_CORE_CONTEXT_EXTRACTION_ERROR("httpIndex"), nil, tempLogProc, w)

		// Der Vorgang wird beendet
		return
	}

	// Es wird ein neuer ProcLog erzeugt
	procLog := coreSession.GetChildProcessLog("ServerOverview")
	procLog.Log("New server overview request [%s]", "")

	// Die Sitzung wird geschlossen
	defer coreSession.Close()

	// Es werden alle Script Container extrahiert
	procLog.Debug("Try to call all available VMs...")
	scriptContainers := o.core.GetAllActiveScriptContainerIDs(procLog)
	ucscontainers := []string{}
	for _, item := range scriptContainers {
		ucscontainers = append(ucscontainers, strings.ToUpper(item))
	}

	// Erstelle ein Response-Objekt mit deiner Nachricht.
	response := types.Response{
		Version:          utils.FormatNumberWithDots(int(static.C_VESION)),
		ScriptContainers: ucscontainers,
	}

	// Die Größe das Paketes wird ermittelt, der Hash des Paketes wird berechnet
	procLog.Debug("Get size of response package")
	responseSize, grErr := GetResponseSize(response, coreSession.GetContentType())
	if grErr != nil {
		// Der Fehler wird übermittelt
		BuildErrorHttpRequestResponseAndWrite("httpIndex", grErr, coreSession, procLog, w)

		// Der Vorgang wird beendet
		return
	}

	// Es wird geprüft ob mindestens 1 Eintrag vorhanden ist
	if responseSize < 0 {
		// Es wird ein neuer Fehler erzeugt
		responeError := errormsgs.HTTP_API_SERVICE_INTERNAL_ERROR_BY_EMITTING_CAPSLE_SIZE("httpIndex")

		// Der Fehler wird übermittelt
		BuildErrorHttpRequestResponseAndWrite("", responeError, coreSession, procLog, w)

		// Der Vorgang wird beendet
		return
	}

	// Es wird ein Hash aus dem Response Capsle erzeugt
	procLog.Debug("Compute sha3 hash of response package")
	responsePackageHash, packageComputingError := ComputeSHA3_256HashFromResponseCapsle(response, coreSession.GetContentType())
	if packageComputingError != nil {
		// Der Fehler wird übermittelt
		BuildErrorHttpRequestResponseAndWrite("httpIndex", packageComputingError, coreSession, procLog, w)

		// Der Vorgang wird beendet
		return
	}

	// Die Antwort wird an den Client zurückgesendet
	procLog.Debug("Write response to client")
	if writingErr := HttpResponseWrite(coreSession.GetContentType(), w, response); writingErr != nil {
		// Die Aktuelle Funktion wird der Historie hinzugefügt
		writingErr.AddCallerFunctionToHistory("httpIndex")

		// Es wird Signalisiert dass die Daten nicht übermittelt werden konnten
		coreSession.SignalTheResponseCouldNotBeSent(responseSize, writingErr)

		// Der Vorgang wird geschlossen
		return
	}

	// Es wird Signalisiert dass die Daten erfolgreich übermittelt wurden
	coreSession.SignalTheResponseWasTransmittedSuccessfully(responseSize, responsePackageHash)

	// Es wird Signalisiert das die Sitzung beendet wurde
	coreSession.Done()
}

// Nimmt Anfragen für "/vm" entgegen
func (o *HttpApiService) httpVmInfo(w http.ResponseWriter, r *http.Request) {
	// Die Core Session wird aus dem Context extrahiert,
	// sollte kein Context vorhanden sein, wird die Verbindung abgebrochen
	coreSession, ok := r.Context().Value(static.CORE_SESSION_CONTEXT_KEY).(*context.HttpContext)
	if !ok {
		// Es muss ein neuer Process Log erzeugt werden um den Fehler auszugeben
		tempLogProc := procslog.NewProcLogSession()

		// Der Fehler wird zurückgegeben
		BuildErrorHttpRequestResponseAndWrite("", errormsgs.HTTP_API_CORE_CONTEXT_EXTRACTION_ERROR("httpIndex"), nil, tempLogProc, w)

		// Der Vorgang wird beendet
		return
	}

	// Die Sitzung wird geschlossen
	defer coreSession.Close()

	// Setze den Content-Type der Antwort auf application/json.
	w.Header().Set("Content-Type", "application/json")

	// Erstelle ein Response-Objekt mit deiner Nachricht.
	response := types.VmInfoResponse{
		//Name:            foundedVM.GetVMName(),
		//Id:              string(foundedVM.GetFingerprint()),
		//SharedFunctions: sharedFunctions,
		//State:           stateStrValue,
	}

	// Schreibe die JSON-Daten in den ResponseWriter.
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Fehler beim Senden der JSON-Antwort: %v", err)
		http.Error(w, "Ein interner Fehler ist aufgetreten", http.StatusInternalServerError)
	}
}

// Nimmt Anfragen für "/rpc" entgegen
func (o *HttpApiService) httpRPC(w http.ResponseWriter, r *http.Request) {
	// Die Core Session wird aus dem Context extrahiert,
	// sollte kein Context vorhanden sein, wird die Verbindung abgebrochen
	coreWebSession, ok := r.Context().Value(static.CORE_SESSION_CONTEXT_KEY).(*context.HttpContext)
	if !ok {
		// Es muss ein neuer Process Log erzeugt werden um den Fehler auszugeben
		tempLogProc := procslog.NewProcLogSession()

		// Der Fehler wird zurückgegeben
		BuildErrorRpcHttpRequestResponseAndWrite("", errormsgs.HTTP_API_CORE_CONTEXT_EXTRACTION_ERROR("httpRPC"), nil, static.HTTP_CONTENT_JSON, tempLogProc, w)

		// Der Vorgang wird beendet
		return
	}

	// Es wird ein neuer ProcLog erzeugt
	procLog := coreWebSession.GetChildProcessLog("httpRpc")

	// Die Sitzung wird geschlossen
	defer coreWebSession.Close()

	// Die gesuchte Funktionssignatur wird extrahiert
	functionSignature := coreWebSession.GetSearchedFunctionSignature()

	// Es wird versucht die VM abzurufen
	vmInstance, foundVM, vmSearchError := o.core.GetScriptContainerVMByID(functionSignature.VMID)
	if vmSearchError != nil {
		// Der Fehler wird zurückgegeben
		BuildErrorRpcHttpRequestResponseAndWrite("httpRPC", vmSearchError, coreWebSession, coreWebSession.GetContentType(), procLog, w)

		// Der Vorgang wird beendet
		return
	}

	// Sollte die VM nicht gefunden werden können, wird der Vorgang abgebrochen
	if !foundVM {
		// Der Fehler wird zurückgegeben
		BuildErrorRpcHttpRequestResponseAndWrite(
			"",
			errormsgs.HTTP_API_SERVICE_REQUEST_VM_NOT_FOUND("httpRPC", functionSignature.VMID, static.RPC_REQUEST_METHODE_VM_IDENT_ID),
			coreWebSession,
			coreWebSession.GetContentType(),
			procLog,
			w,
		)

		// Der Vorgang wird beendet
		return
	}

	// Es wird ermittelt ob die VM Instanz ausgeführt wird
	if vmInstance.GetState() != static.Running {
		// Der Fehler wird zurückgegeben
		BuildErrorRpcHttpRequestResponseAndWrite("", errormsgs.HTTP_API_RPC_VM_NOT_RUNNING("httpRPC"), coreWebSession, coreWebSession.GetContentType(), procLog, w)

		// Der Vorgang wird beendet
		return
	}

	// Es wird versucht die Funktion anhand ihrer Signatur zu ermitteln
	// procslog.ProcFormatConsoleText(coreSession.GetProcLogSession(), "HTTP-PRC", types.DETERMINE_THE_FUNCTION, vmInstance.GetVMName(), functionSignature.FunctionName)
	foundFunction, hasFound, gsfbserr := vmInstance.GetSharedFunctionBySignature(static.LOCAL, functionSignature)
	if gsfbserr != nil {
		// Der Fehler wird zurückgegeben
		BuildErrorRpcHttpRequestResponseAndWrite("httpRPC", gsfbserr, coreWebSession, coreWebSession.GetContentType(), procLog, w)

		// Der Vorgang wird beendet
		return
	}

	// Sollte keine Passende Funktion gefunden werden, wird der Vorgang abgebrochen
	if !hasFound {
		// Der Fehler wird zurückgegeben
		BuildErrorRpcHttpRequestResponseAndWrite(
			"",
			errormsgs.HTTP_API_SERVICE_REQUEST_VM_FUNCTION_NOT_FOUND_ERROR("httpRPC", functionSignature.VMName, static.RPC_REQUEST_METHODE_VM_IDENT_ID, functionSignature),
			coreWebSession,
			coreWebSession.GetContentType(),
			procLog,
			w,
		)

		// Der Vorgang wird beendet
		return
	}

	// Es wird versucht die Parameter aus dem Request auszulesen, dazu gehören der Datensatz, sowie der Datentyp der Parameter
	data, ttcrfcerr := TryToReadCompleteFunctionCallFromRequest(vmInstance.GetKId(), coreWebSession.GetContentType(), r.Body)
	if ttcrfcerr != nil {
		// Der Fehler wird zurückgegeben
		BuildErrorRpcHttpRequestResponseAndWrite("httpRPC", ttcrfcerr, coreWebSession, coreWebSession.GetContentType(), procLog, w)

		// Der Vorgang wird beendet
		return
	}

	// Die Einzelnen Parameter werden geprüft und abgearbeitet
	extractedValues, trfperr := TryReadFunctionParameter(data, foundFunction)
	if trfperr != nil {
		// Der Fehler wird übermittelt
		BuildErrorRpcHttpRequestResponseAndWrite("httpRPC", trfperr, coreWebSession, coreWebSession.GetContentType(), procLog, w)

		// Rückgabe
		return
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
		Parms: extractedValues,
		Request: &types.HttpContext{
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
		},
		ProcessLog:    procLog,
		RequestType:   static.HTTP_REQUEST,
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
		BuildErrorRpcHttpRequestResponseAndWrite("HttpApiService->httpCallFunction", sfbError, coreWebSession, coreWebSession.GetContentType(), procLog, w)

		// Rückgabe
		return
	}

	// Es wird auf das Ergebniss gewartet
	result, ok := coreWebSession.GetReturnChan().Read()
	if result == nil || !ok {
		// Es wird ein neuer Fehler erzeugt
		responeError := errormsgs.HTTP_API_SERVICE_REQUEST_INTERNAL_ERROR_BY_READING_RETURN_CHAN("HttpApiService->httpCallFunction")

		// Der Fehler wird übermittelt
		BuildErrorRpcHttpRequestResponseAndWrite("", responeError, coreWebSession, coreWebSession.GetContentType(), procLog, w)

		// Der Vorgang wird beendet
		return
	}

	// Die Antwort wird gebaut
	var responseData *types.HttpRpcResponseCapsle
	if result.State == "ok" {
		// Erstellt ein Array welches alle Antwortdaten enthält
		dt := make([]*types.RPCResponseData, 0)

		// Die Antwortdaten werden abgearbeitet
		for _, item := range result.Return {
			dt = append(dt, &types.RPCResponseData{DType: item.Type, Value: item.Value})
		}

		// Das Response Capsle wird erzeugt
		responseData = &types.HttpRpcResponseCapsle{Data: dt}
	} else if result.State == "failed" {
		// Die Response Capsle wird erzeugt
		responseData = &types.HttpRpcResponseCapsle{Error: result.Error}
	} else if result.State == "exception" {
		// Die Response Capsle wird erzeugt
		responseData = &types.HttpRpcResponseCapsle{Error: result.Error}
	} else {
		// Die Response Capsle wird erzeugt
		responseData = &types.HttpRpcResponseCapsle{Error: "unkown return state"}
	}

	// Die Größe das Paketes wird ermittelt, der Hash des Paketes wird berechnet
	responseSize, _ := GetResponseSize(responseData, coreWebSession.GetContentType())
	if responseSize < 0 {
		// Es wird ein neuer Fehler erzeugt
		responeError := errormsgs.HTTP_API_SERVICE_INTERNAL_ERROR_BY_EMITTING_CAPSLE_SIZE("HttpApiService->httpCallFunction")

		// Der Fehler wird übermittelt
		BuildErrorRpcHttpRequestResponseAndWrite("", responeError, coreWebSession, coreWebSession.GetContentType(), procLog, w)

		// Der Vorgang wird beendet
		return
	}

	// Es wird ein Hash aus dem Response Capsle erzeugt
	responsePackageHash, packageComputingError := ComputeSHA3_256HashFromResponseCapsle(responseData, coreWebSession.GetContentType())
	if packageComputingError != nil {
		// Der Fehler wird übermittelt
		BuildErrorRpcHttpRequestResponseAndWrite("HttpApiService->httpCallFunction", packageComputingError, coreWebSession, coreWebSession.GetContentType(), procLog, w)

		// Der Vorgang wird beendet
		return
	}

	// Es wird versucht die Daten zurückzusenden
	if frwrerr := HttpRpcResponseWrite(coreWebSession.GetContentType(), w, responseData); frwrerr != nil {
		// Die Aktuelle Funktion wird der Historie hinzugefügt
		frwrerr.AddCallerFunctionToHistory("HttpApiService->httpCallFunction")

		// Es wird Signalisiert dass die Daten nicht übermittelt werden konnten
		coreWebSession.SignalTheResponseCouldNotBeSent(responseSize, frwrerr)
		//result.Reject()

		// Log
		procLog.Log("%s call function response sending '%s' error\n\t%s\n", vmInstance.GetVMName(), foundFunction.GetName(), frwrerr)

		// Der Vorgang wird geschlossen
		return
	}

	// Es wird Signalisiert dass die Daten erfolgreich übermittelt wurden
	coreWebSession.SignalTheResponseWasTransmittedSuccessfully(responseSize, responsePackageHash)
	//result.Resolve()

	// Log
	size, _ := GetResponseSize(responseData, coreWebSession.GetContentType())
	procslog.ProcFormatConsoleText(procLog, "HTTP-PRC", types.RPC_CALL_DONE_RESPONSE, fmt.Sprintf("%d", size))

	// Signalisiert dass die Funktion erfolgreich durchgelaufen ist und nun als nächstes der Return ansteht
	coreWebSession.Done()
}
