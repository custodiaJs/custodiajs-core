package middlewares

import (
	"net/http"
	"strings"

	"github.com/CustodiaJS/custodiajs-core/context"
	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/static/errormsgs"
	"github.com/CustodiaJS/custodiajs-core/types"
)

func ParseAndPassVmRpcUrlFunctionSignature(core types.CoreInterface, w http.ResponseWriter, r *http.Request) *types.SpecificError {
	// Die Core Session wird aus dem Context extrahiert,
	// sollte kein Context vorhanden sein, wird die Verbindung abgebrochen
	coreSession, ok := r.Context().Value(static.CORE_SESSION_CONTEXT_KEY).(*context.HttpContext)
	if !ok {
		return nil
	}

	// URL-Parameter parsen
	queryParams := r.URL.Query()

	// Extrahiere die VM-ID
	vm := queryParams.Get("vm")
	if vm == "" {
		return nil
	}

	// Extrahiere den Funktionsnamen
	functionName := queryParams.Get("function")
	if functionName == "" {
		return nil
	}

	// Extrahiere die Parameter (kann leer sein)
	params := strings.Split(queryParams.Get("parms"), ",")
	if len(params) == 1 && params[0] == "" {
		params = []string{}
	}

	// Extrahiere den Rückgabewert
	returnType := queryParams.Get("return")

	// Die Funktionssignatur wird erstellt
	searchedFunctionSignature := &types.FunctionSignature{
		VMID:         strings.ToLower(string(vm)),
		FunctionName: functionName,
		Params:       params,
		ReturnType:   returnType,
	}

	// Die Funktionssignatur wird hinzugefügt
	coreSession.AddSearchedFunctionSignature(searchedFunctionSignature)

	// Es ist kein Fehler aufgetreten
	return nil
}

func ValidateRPCRequest(core types.CoreInterface, w http.ResponseWriter, r *http.Request) *types.SpecificError {
	// Die Core Session wird aus dem Context extrahiert,
	// sollte kein Context vorhanden sein, wird die Verbindung abgebrochen
	coreSession, ok := r.Context().Value(static.CORE_SESSION_CONTEXT_KEY).(*context.HttpContext)
	if !ok {
		return nil
	}

	// Es wird versucht die Function Signature anzurufen
	functionSignature := coreSession.GetSearchedFunctionSignature()

	// Es wird versucht die VM zu finden
	vmInstance, foundVM, vmSearchError := core.GetScriptContainerVMByID(functionSignature.VMID)
	if vmSearchError != nil {
		return vmSearchError.AddCallerFunctionToHistory("ValidateRPCRequest")
	}

	// Sollte die VM nicht gefunden werden können, wird der Vorgang abgebrochen
	if !foundVM {
		return errormsgs.HTTP_API_SERVICE_REQUEST_VM_NOT_FOUND("HttpApiService->httpCallFunction", functionSignature.VMID, static.RPC_REQUEST_METHODE_VM_IDENT_ID)
	}

	// Es wird ermittelt ob die VM Instanz ausgeführt wird
	if vmInstance.GetState() != static.Running {
		return nil
	}

	// Es wird versucht die Funktion anhand ihrer Signatur zu ermitteln
	// procslog.ProcFormatConsoleText(coreSession.GetProcLogSession(), "HTTP-PRC", types.DETERMINE_THE_FUNCTION, vmInstance.GetVMName(), functionSignature.FunctionName)
	foundFunction, hasFound, gsfbserr := vmInstance.GetSharedFunctionBySignature(static.LOCAL, functionSignature)
	if gsfbserr != nil {
		return gsfbserr.AddCallerFunctionToHistory("ValidateRPCRequest")
	}

	// Sollte keine Passende Funktion gefunden werden, wird der Vorgang abgebrochen
	if !hasFound {
		return errormsgs.HTTP_API_SERVICE_REQUEST_VM_FUNCTION_NOT_FOUND_ERROR("HttpApiService->httpCallFunction", functionSignature.VMName, static.RPC_REQUEST_METHODE_VM_IDENT_ID, functionSignature)
	}

	// Die Anzahl der Übertragenene Parameter muss mit den Vorhandenen Parametern übereinstimmen
	if len(foundFunction.GetParmTypes()) != len(functionSignature.Params) {
		return nil
	}

	// Die Einzelnen Parameter werden abgearbeitet, es wird geprüft ob die Datentypen passen
	for parmHight, parmType := range foundFunction.GetParmTypes() {
		if parmType != functionSignature.Params[parmHight] {
			return nil
		}
	}

	// Es ist kein Fehler aufgetreten
	return nil
}
