package rpcjsprocessor

import (
	"time"

	"github.com/CustodiaJS/custodiajs-core/global/static/errormsgs"
	"github.com/CustodiaJS/custodiajs-core/global/types"
	"github.com/CustodiaJS/custodiajs-core/global/utils"
	"github.com/CustodiaJS/custodiajs-core/vm/context/eventloop"

	v8 "rogchap.com/v8go"
)

// callInKernelEventLoopCheck überprüft den Status eines Promises in der Kernel-Eventschleife.
// Bei einem Pending-Promise plant es die nächste Überprüfung ohne aktives Warten.
// Bei einem Rejected-Promise führt es einen Microtask-Checkpoint durch.
func callInKernelEventLoopCheck(o *SharedFunction, ctx *v8.Context, prom *v8.Promise, request *SharedFunctionRequestContext, req *types.RpcRequest) *types.SpecificError {
	// Der Stauts des Objektes wird ermittelt
	switch prom.State() {
	case v8.Pending:
		// Planen Sie die nächste Überprüfung, ohne aktives Warten zu verwenden
		go func() {
			// Es wird 1ne Milisekunde gewartet
			time.Sleep(1 * time.Millisecond)

			// Es wird eine neue Kernel Funktion erzeugt
			eventloopFunction := eventloop.NewKernelEventLoopFunctionOperation(func(ctx *v8.Context, klopr types.KernelEventLoopContextInterface) {
				callInKernelEventLoopCheck(o, ctx, prom, request, req)
			})

			// Es wird ein neues Event zum Kernel hinzugefügt
			o.kernel.AddToEventLoop(eventloopFunction)
		}()
	case v8.Rejected:
		// PerformMicrotaskCheckpoint runs the default MicrotaskQueue until empty. This is used to make progress on Promises.
		ctx.PerformMicrotaskCheckpoint()
	}

	// Keine Rückgabe
	return nil
}

// functionCallInEventloopFinall führt den abschließenden Schritt eines Funktionsaufrufs durch.
// Es fügt einen neuen Eintrag zur Eventschleife hinzu, prüft den Promise-Status und behandelt etwaige Fehler.
// Bei Erfolg wird das Ergebnis der Operation signalisiert.
func functionCallInEventloopFinall(o *SharedFunction, request *SharedFunctionRequestContext, req *types.RpcRequest, prom *v8.Promise) *types.SpecificError {
	// Die Eventloop Funktion wird erzeugt
	eventloopFunction := eventloop.NewKernelEventLoopFunctionOperation(func(ctx *v8.Context, klopr types.KernelEventLoopContextInterface) {
		err := callInKernelEventLoopCheck(o, ctx, prom, request, req)
		if err != nil {
			// Der Fehler wird zurückgegeben
			klopr.SetError(err)
		}

		// Signalisiert dass der Vorgang erfolgreich war
		klopr.SetResult(nil)
	})

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if err := o.kernel.AddToEventLoop(eventloopFunction); err != nil {
		return err.AddCallerFunctionToHistory("functionCallInEventloopFinall")
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

// functionCallInEventloopPromiseOperation verarbeitet das Ergebnis eines Funktionsaufrufs, der ein Promise zurückgibt.
// Es prüft, ob die Verbindung noch besteht, behandelt das Promise und führt die finalen Schritte des Funktionsaufrufs durch.
// Bei Erfolg wird das Ergebnis der Operation signalisiert.
func functionCallInEventloopPromiseOperation(o *SharedFunction, request *SharedFunctionRequestContext, req *types.RpcRequest, result *v8.Value) *types.SpecificError {
	// Die Eventloop Funktion wird erzeugt
	eventloopFunction := eventloop.NewKernelEventLoopFunctionOperation(func(ctx *v8.Context, klopr types.KernelEventLoopContextInterface) {
		// Es wird geprüft ob es sich um ein Promises handelt
		if !result.IsPromise() {
			panic("isnr promise")
		}

		// Das Promises Objekt wird erzeugt
		prom, err := result.AsPromise()
		if err != nil {
			panic(err)
		}

		// Wird ausgeführt wenn der Funktionsaufruf durchgeführt wurde
		funcFinal := func(info *v8.FunctionCallbackInfo) *v8.Value {
			request.functionCallFinal()
			return v8.Undefined(info.Context().Isolate())
		}

		// Wird ausgeführt wenn ein Throw durchgeführt wurde
		throwProm := func(info *v8.FunctionCallbackInfo) *v8.Value {
			request.functionCallException(info.Args()[0].String())
			return v8.Undefined(info.Context().Isolate())
		}

		// Die Then und Catch funktionen werden festgelegt
		prom = prom.Then(funcFinal, throwProm)
		prom = prom.Catch(throwProm)

		// Der 5. Schritt des Funktionsaufrufes wird durchgeführt
		if err := functionCallInEventloopFinall(o, request, req, prom); err != nil {
			return
		}

		// Signalisiert dass der Vorgang erfolgreich war
		klopr.SetResult(nil)
	})

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if err := o.kernel.AddToEventLoop(eventloopFunction); err != nil {
		return err.AddCallerFunctionToHistory("functionCallInEventloopPromiseOperation")
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

// functionCallInEventloop führt den vorbereiteten Funktionsaufruf innerhalb der Eventschleife aus.
// Es prüft, ob die Verbindung noch besteht, führt die Funktion aus und behandelt das Ergebnis.
// Bei Erfolg wird das Ergebnis der Operation signalisiert.
func functionCallInEventloop(o *SharedFunction, request *SharedFunctionRequestContext, req *types.RpcRequest, proxFunction *v8.Function, proxArguments []v8.Valuer) *types.SpecificError {
	// Die Eventloop Funktion wird erzeugt
	eventloopFunction := eventloop.NewKernelEventLoopFunctionOperation(func(ctx *v8.Context, klopr types.KernelEventLoopContextInterface) {
		// Die Funktion wird ausgeführt
		result, err := proxFunction.Call(v8.Undefined(ctx.Isolate()), proxArguments...)
		if err != nil {
			panic(err)
		}

		// Der 4. Schritt des Funktionsaufrufes wird durchgeführt
		if err := functionCallInEventloopPromiseOperation(o, request, req, result); err != nil {
			return
		}

		// Signalisiert dass der Vorgang erfolgreich war
		klopr.SetResult(nil)
	})

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if err := o.kernel.AddToEventLoop(eventloopFunction); err != nil {
		return err.AddCallerFunctionToHistory("functionCallInEventloop")
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

// functionCallInEventloopProxyObjectPrepare bereitet den Proxy-Objekt-Funktionsaufruf innerhalb der Eventschleife vor.
// Es erstellt die finalen Argumente, setzt den JavaScript-Code für den Proxy-Wrap,
// führt die Funktion in der Eventschleife aus und behandelt mögliche Fehler.
// Bei Erfolg wird das Ergebnis der Operation signalisiert.
func functionCallInEventloopProxyObjectPrepare(o *SharedFunction, request *SharedFunctionRequestContext, req *types.RpcRequest, requestObj *v8.Object, convertedValues []v8.Valuer) *types.SpecificError {
	// Die Eventloop Funktion wird erzeugt
	eventloopFunction := eventloop.NewKernelEventLoopFunctionOperation(func(ctx *v8.Context, klopr types.KernelEventLoopContextInterface) {
		// Die Finalen Argumente werden erstellt
		finalArguments := make([]v8.Valuer, 0)
		finalArguments = append(finalArguments, requestObj)
		finalArguments = append(finalArguments, convertedValues...)

		// Der Code für die Proxy Shield Funktion wird ersteltl
		procxyFunction, err := ctx.RunScript(testJsProxySource, "rpc_function_call_proxy_shield.js")
		if err != nil {
			return
		}

		// Es wird geprüft ob es sich um eine Funktion handelt,
		// wenn ja wird die Funktion extrahiert
		proxFunction, err := procxyFunction.AsFunction()
		if err != nil {
			return
		}

		// Das Proxy Objekt wird erzeugt
		proxyObject, err := v8makeProxyForRPCCall(ctx, request)
		if err != nil {
			return
		}

		// Die Argumente für den Proxy werden erstellt
		proxArguments := []v8.Valuer{o.v8Function, proxyObject}
		proxArguments = append(proxArguments, finalArguments...)

		// Der 3. Schritt des Funktionsaufrufes wird durchgeführt
		if err := functionCallInEventloop(o, request, req, proxFunction, proxArguments); err != nil {
			return
		}

		// Signalisiert dass der Vorgang erfolgreich war
		klopr.SetResult(nil)
	})

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if err := o.kernel.AddToEventLoop(eventloopFunction); err != nil {
		return err.AddCallerFunctionToHistory("functionCallInEventloopProxyObjectPrepare")
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

// functionCallInEventloopInit initialisiert einen Funktionsaufruf innerhalb der Eventschleife.
// Es prüft, ob die Verbindung besteht, wandelt die Parameter um, erstellt ein Request-Objekt,
// und führt die vorbereitenden Schritte des Funktionsaufrufs durch.
// Die Funktion wird zur Eventschleife hinzugefügt und das Ergebnis des Aufrufs wird verarbeitet.
func functionCallInEventloopInit(o *SharedFunction, request *SharedFunctionRequestContext, req *types.RpcRequest) *types.SpecificError {
	// Die Eventloop Funktion wird erzeugt
	eventloopFunction := eventloop.NewKernelEventLoopFunctionOperation(func(ctx *v8.Context, klopr types.KernelEventLoopContextInterface) {
		// Die Parameter werden umgewandelt
		convertedValues, err := convertRequestParametersToV8Parameters(ctx.Isolate(), o.signature.Params, req.Parms)
		if err != nil {
			return
		}

		// Das Request Objekt wird erstellt
		requestObj, err := v8makeSharedFunctionObject(ctx, request, req)
		if err != nil {
			return
		}

		// Der 2. Schritt des Funktionsaufrufes wird durchgeführt
		if err := functionCallInEventloopProxyObjectPrepare(o, request, req, requestObj, convertedValues); err != nil {
			return
		}

		// Signalisiert dass der Vorgang erfolgreich war
		klopr.SetResult(nil)
	})

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if err := o.kernel.AddToEventLoop(eventloopFunction); err != nil {
		return err.AddCallerFunctionToHistory("functionCallInEventloopInit")
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

// convertRequestParametersToV8Parameters wandelt die RPC-Argumente in V8-Argumente für den aktuellen Kontext um.
// Es überprüft die Datentypen und konvertiert sie in die entsprechenden V8-Typen.
// Bei einem Fehler wird eine entsprechende Fehlermeldung zurückgegeben.
func convertRequestParametersToV8Parameters(iso *v8.Isolate, parmTypes []string, reqparms []*types.FunctionParameterCapsle) ([]v8.Valuer, *types.SpecificError) {
	// Es wird versucht die Paraemter in den Richtigen v8 Datentypen umzuwandeln
	convertedValues := make([]v8.Valuer, 0)
	for hight, item := range reqparms {
		// Es wird geprüft ob der Datentyp gewünscht ist
		if item.CType != parmTypes[hight] {
			return nil, errormsgs.KMDOULE_RPC_SHARED_FUNCTION_FUNCTION_REQUEST_CONVERTING_PARM_INVALID_DTYPE_ERROR("convertRequestParametersToV8Parameters", hight, item.CType, parmTypes[hight])
		}

		// Es wird geprüft ob es sich um einen Zulässigen Datentypen handelt
		switch item.CType {
		case "boolean":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, errormsgs.KMDOULE_RPC_SHARED_FUNCTION_REQUEST_CONVERTING_ERROR("convertRequestParametersToV8Parameters", hight, err, "boolean")
			}
			convertedValues = append(convertedValues, val)
		case "number":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, errormsgs.KMDOULE_RPC_SHARED_FUNCTION_REQUEST_CONVERTING_ERROR("convertRequestParametersToV8Parameters", hight, err, "number")
			}
			convertedValues = append(convertedValues, val)
		case "string":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, errormsgs.KMDOULE_RPC_SHARED_FUNCTION_REQUEST_CONVERTING_ERROR("convertRequestParametersToV8Parameters", hight, err, "string")
			}
			convertedValues = append(convertedValues, val)
		case "array":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, errormsgs.KMDOULE_RPC_SHARED_FUNCTION_REQUEST_CONVERTING_ERROR("convertRequestParametersToV8Parameters", hight, err, "array")
			}
			convertedValues = append(convertedValues, val)
		case "object":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, errormsgs.KMDOULE_RPC_SHARED_FUNCTION_REQUEST_CONVERTING_ERROR("convertRequestParametersToV8Parameters", hight, err, "object")
			}
			convertedValues = append(convertedValues, val)
		case "bytes":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, errormsgs.KMDOULE_RPC_SHARED_FUNCTION_REQUEST_CONVERTING_ERROR("convertRequestParametersToV8Parameters", hight, err, "bytes")
			}
			convertedValues = append(convertedValues, val)
		default:
			return nil, errormsgs.KMDOULE_RPC_SHARED_FUNCTION_REQUEST_UNKOWN_DATATYPE("convertRequestParametersToV8Parameters", hight)
		}
	}

	// Rückgabe ohne Fehler
	return convertedValues, nil
}

// Überprüft ob ein SharedFunctionRequestContext korrekt aufgebaut ist
func validateSharedFunctionRequestContextObject(o *SharedFunctionRequestContext) bool {
	// Sollte die SharedFunctionRequestContext "o" NULL sein, wird ein False zurückgegeben
	if o == nil {
		return false
	}

	// Die Einzelnen Elemente werden geprüft
	if o.kernel == nil {
		return false
	}
	if o._rprequest == nil {
		return false
	}
	if o._returnDataType == "" {
		return false
	}
	if o._destroyed == nil {
		return false
	}

	// Es handelt sich um ein zulässiges Objekt
	return true
}

// Die Funktion wird erstellt
func v8makeSharedFunctionObject(context *v8.Context, request *SharedFunctionRequestContext, rrpcrequest *types.RpcRequest) (*v8.Object, *types.SpecificError) {
	// Das Requestobjekt wird ersellt
	objTemplate := v8.NewObjectTemplate(context.Isolate())

	// Die Resolve Funktion wird festgelegt
	if err := objTemplate.Set("Resolve", v8.NewFunctionTemplate(context.Isolate(), request.resolveFunctionCallbackV8)); err != nil {
		return nil, utils.V8ObjectWritingError()
	}

	// Die Reject Funktion wird festgelegt
	if err := objTemplate.Set("Reject", v8.NewFunctionTemplate(context.Isolate(), request.rejectFunctionCallbackV8)); err != nil {
		return nil, utils.V8ObjectWritingError()
	}

	// Das Objekt wird erzeugt
	obj, err := objTemplate.NewInstance(context)
	if err != nil {
		return nil, utils.V8ObjectWritingError()
	}

	// Wird von V8 Verwendet um zu ermitteln ob die Verbindung mit der Anfragendenseite noch besteht
	isConnected := func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Es wird versucht den Boolwert einzulesen
		value, err := v8.NewValue(context.Isolate(), rrpcrequest.Context.IsConnected())
		if err != nil {
			// Der RPC Vorgang wird aufgrund eines Engine Fehlers abgebrochen
			writeRequestReturnResponse(request, &types.FunctionCallState{Error: "javascript engine error", State: "aborted"})

			// Es wird ein JS Throw ausgelöst
			utils.V8ContextThrow(info.Context(), "internal engine error")

			// Rückgabe ohne wert
			return nil
		}

		// Der Wert wird zurückgegeben
		return value
	}

	// Es wird ein neues Objekt erzeugt, dieses Objekt wird verwendet um den Aktuellen Request Darzustellen
	var rpcConnectionType string
	switch reqobj := rrpcrequest.Request.(type) {
	case types.HttpContext:
		// Der Type der Verbindung wird definiert
		rpcConnectionType = "http"

		// Die Cookies werden Extrahiert
		cookies := v8.NewObjectTemplate(context.Isolate())
		for _, item := range reqobj.Cookies {
			// Es wird ein neues Objekt erzeugt
			cookieObject := v8.NewObjectTemplate(context.Isolate())
			cookieObject.Set("Value", item.Value)
			cookieObject.Set("Domain", item.Domain)
			cookieObject.Set("Path", item.Path)
			cookieObject.Set("Expires", item.RawExpires)

			// Der Eintrag wird hinzugefügt
			if err := cookies.Set(item.Name, cookieObject); err != nil {
				return nil, utils.MakeV8Error("v8makeSharedFunctionObject", err)
			}
		}

		// Der Header wird vorbereitet
		headersTemplate := v8.NewObjectTemplate(context.Isolate())
		headers, err := headersTemplate.NewInstance(context)
		if err != nil {
			return nil, utils.MakeV8Error("v8makeSharedFunctionObject", err)
		}

		// Die Header werden extrahiert
		for k, v := range reqobj.Header {
			// Es wird ein neues Slices erzeugt
			sliceV8, err := context.RunScript("(function() { return []; })();", "slice.js")
			if err != nil {
				return nil, utils.MakeV8Error("v8makeSharedFunctionObject", err)
			}

			// Das Objekt wird ausgelesen
			sliceObject, err := sliceV8.AsObject()
			if err != nil {
				return nil, utils.MakeV8Error("v8makeSharedFunctionObject", err)
			}

			// Die Einzelnen Werte werden umgewandelt
			for _, value := range v {
				// Der Wert wird umgewandelt
				v8Value, err := v8.NewValue(context.Isolate(), value)
				if err != nil {
					return nil, utils.MakeV8Error("v8makeSharedFunctionObject", err)
				}

				// Der Wert wird hinzugefügt
				sliceObject.Object().MethodCall("push", v8Value)
			}

			// Der Eintrag wird hinzugefügt
			if err := headers.Set(k, sliceObject); err != nil {
				return nil, utils.MakeV8Error("v8makeSharedFunctionObject", err)
			}
		}

		// Das Http Objekt wird erzeugt
		httpObj := v8.NewObjectTemplate(context.Isolate())

		// Die Werte werden hinzugefügt
		if err := httpObj.Set("IsConnected", v8.NewFunctionTemplate(context.Isolate(), isConnected)); err != nil {
			return nil, utils.MakeV8Error("v8makeSharedFunctionObject", err)
		}
		if err := httpObj.Set("ContentLength", float64(reqobj.ContentLength)); err != nil {
			return nil, utils.MakeV8Error("v8makeSharedFunctionObject", err)
		}
		if err := httpObj.Set("Host", reqobj.Host); err != nil {
			return nil, utils.MakeV8Error("v8makeSharedFunctionObject", err)
		}
		if err := httpObj.Set("Proto", reqobj.Proto); err != nil {
			return nil, utils.MakeV8Error("v8makeSharedFunctionObject", err)
		}
		if err := httpObj.Set("RemoteAddr", reqobj.RemoteAddr); err != nil {
			return nil, utils.MakeV8Error("v8makeSharedFunctionObject", err)
		}
		if err := httpObj.Set("RequestURI", reqobj.RequestURI); err != nil {
			return nil, utils.MakeV8Error("v8makeSharedFunctionObject", err)
		}
		if err := httpObj.Set("Cookies", cookies); err != nil {
			return nil, utils.MakeV8Error("v8makeSharedFunctionObject", err)
		}

		// Das Finale Objekt wird erzeugt
		http, err := httpObj.NewInstance(context)
		if err != nil {
			return nil, utils.MakeV8Error("v8makeSharedFunctionObject", err)
		}

		// Die Header werden hinzugefügt
		if err := http.Set("Headers", headers); err != nil {
			return nil, utils.MakeV8Error("v8makeSharedFunctionObject", err)
		}

		// Das Objekt wird abgespeichert
		if err := obj.Set("http", http); err != nil {
			return nil, utils.MakeV8Error("v8makeSharedFunctionObject", err)
		}
	default:
		return nil, utils.MakeUnkownMethodeError("v8makeSharedFunctionObject")
	}

	// Der Wert wird eingelesen
	val, err := v8.NewValue(context.Isolate(), rpcConnectionType)
	if err != nil {
		return nil, utils.MakeV8Error("v8makeSharedFunctionObject", err)
	}

	// Der Eintrag wird im Objekt hinzugefügt
	if err := obj.Set("CallMethode", val); err != nil {
		return nil, utils.MakeV8Error("v8makeSharedFunctionObject", err)
	}

	// Rückgabe ohne Fehler
	return obj, nil
}

// Das This Objekt wird erstellt
func v8makeProxyForRPCCall(context *v8.Context, request *SharedFunctionRequestContext) (*v8.Object, error) {
	// Das Requestobjekt wird ersellt
	obj := v8.NewObjectTemplate(context.Isolate())

	// Die Funktionen werden hinzugefügt
	if err := obj.Set("proxyShieldConsoleLog", v8.NewFunctionTemplate(context.Isolate(), request.proxyShield_ConsoleLog)); err != nil {
		return nil, utils.V8ObjectWritingError()
	}
	if err := obj.Set("proxyShieldErrorLog", v8.NewFunctionTemplate(context.Isolate(), request.proxyShield_ErrorLog)); err != nil {
		return nil, utils.V8ObjectWritingError()
	}
	if err := obj.Set("clearInterval", v8.NewFunctionTemplate(context.Isolate(), request.proxyShield_ClearInterval)); err != nil {
		return nil, utils.V8ObjectWritingError()
	}
	if err := obj.Set("clearTimeout", v8.NewFunctionTemplate(context.Isolate(), request.proxyShield_ClearTimeout)); err != nil {
		return nil, utils.V8ObjectWritingError()
	}
	if err := obj.Set("setInterval", v8.NewFunctionTemplate(context.Isolate(), request.proxyShield_SetInterval)); err != nil {
		return nil, utils.V8ObjectWritingError()
	}
	if err := obj.Set("setTimeout", v8.NewFunctionTemplate(context.Isolate(), request.proxyShield_SetTimeout)); err != nil {
		return nil, utils.V8ObjectWritingError()
	}
	if err := obj.Set("resolve", v8.NewFunctionTemplate(context.Isolate(), request.resolveFunctionCallbackV8)); err != nil {
		return nil, utils.V8ObjectWritingError()
	}
	if err := obj.Set("reject", v8.NewFunctionTemplate(context.Isolate(), request.rejectFunctionCallbackV8)); err != nil {
		return nil, utils.V8ObjectWritingError()
	}
	if err := obj.Set("newPromise", v8.NewFunctionTemplate(context.Isolate(), request.proxyShield_NewPromise)); err != nil {
		return nil, utils.V8ObjectWritingError()
	}

	// Die Testfunktionen werden hinzugefügt
	if err := obj.Set("wait", v8.NewFunctionTemplate(context.Isolate(), request.testWait)); err != nil {
		return nil, utils.V8ObjectWritingError()
	}

	// Das Finale Objekt wird erstellt
	fobj, err := obj.NewInstance(context)
	if err != nil {
		return nil, utils.V8ObjectInstanceCreatingError()
	}

	// Rückgabe ohne Fehler
	return fobj, nil
}
