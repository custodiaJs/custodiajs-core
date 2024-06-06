// Author: fluffelpuff
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package kmodulerpc

import (
	"fmt"
	"time"
	"vnh1/types"
	"vnh1/utils/eventloop"
	rpcrequest "vnh1/utils/rpc_request"

	v8 "rogchap.com/v8go"
)

// GetName gibt den Namen der aktuellen SharedFunction zurück.
// Wenn die SharedFunction null ist, wird ein leerer String zurückgegeben.
func (o *SharedFunction) GetName() string {
	// Es wird geprüft ob die Aktuelle SharedFunction "o" NULL ist
	if o == nil {
		return ""
	}

	// Rückgabe
	return o.name
}

// GetParmTypes gibt die Parameter-Datentypen der aktuellen SharedFunction zurück.
// Wenn die SharedFunction null ist, wird ein leerer Slice von Strings zurückgegeben.
func (o *SharedFunction) GetParmTypes() []string {
	// Es wird geprüft ob die Aktuelle SharedFunction "o" NULL ist
	if o == nil {
		return make([]string, 0)
	}

	// Rückgabe
	return o.signature.Params
}

// GetReturnDatatype gibt den Rückgabedatentyp der aktuellen SharedFunction zurück.
// Wenn die SharedFunction null ist, wird ein leerer String zurückgegeben.
func (o *SharedFunction) GetReturnDatatype() string {
	// Es wird geprüft ob die Aktuelle SharedFunction "o" NULL ist
	if o == nil {
		return ""
	}

	// Rückgabe
	return o.signature.ReturnType
}

// AddOnRequestFunction fügt eine neue Funktion hinzu, die aufgerufen wird, wenn eine Anforderung eingeht.
// Die hinzugefügte Funktion wird in der eventOnRequest-Liste gespeichert und eine Log-Meldung wird ausgegeben.
func (o *SharedFunction) AddOnRequestFunction(funcv8 *v8.Function) error {
	// Die Funktion wird abgespeichert
	o.eventOnRequest = append(o.eventOnRequest, funcv8)

	// Log
	o.kernel.LogPrint(fmt.Sprintf("RPC(%s)", o.signature.FunctionName), "New 'OnRequest' event added for function %s", o.name)

	// Rückgabe ohne Fehler
	return nil
}

// AddOnRequestFailFunction fügt eine neue Funktion hinzu, die aufgerufen wird, wenn eine Anforderung fehlschlägt.
// Die hinzugefügte Funktion wird in der eventOnRequestFail-Liste gespeichert und eine Log-Meldung wird ausgegeben.
func (o *SharedFunction) AddOnRequestFailFunction(funcv8 *v8.Function) error {
	// Die Funktion wird abgespeichert
	o.eventOnRequestFail = append(o.eventOnRequest, funcv8)

	// Log
	o.kernel.LogPrint(fmt.Sprintf("RPC(%s)", o.signature.FunctionName), "New 'OnRequestFail' event added for function %s", o.name)

	// Rückgabe ohne Fehler
	return nil
}

// EnterFunctionCall führt einen Funktionsaufruf innerhalb der SharedFunction-Instanz durch.
// Es überprüft die Gültigkeit der Parameter und Verbindungen, erstellt ein Request-Objekt,
// und übergibt die Operation an die Kernel-Eventschleife. Der Aufruf wird in einer neuen Goroutine
// ausgeführt, die auf die Verarbeitung der Eventschleife wartet und das Ergebnis oder einen Fehler zurückgibt.
func (o *SharedFunction) EnterFunctionCall(req *types.RpcRequest) error {
	// Es wird geprüft ob die Aktuelle SharedFunction "o" NULL ist
	if o == nil {
		return &types.MultiError{ErrorValue: fmt.Errorf("EnterFunctionCall: o is null"), VmErrorValue: fmt.Errorf("internal error")}
	}

	// Es wird geprüft ob der RPC Request "req" NULL ist
	if req == nil {
		return &types.MultiError{ErrorValue: fmt.Errorf("EnterFunctionCall: req is null"), VmErrorValue: fmt.Errorf("request error")}
	}

	// Es wird geprüft ob das Req Objekt eine Verbindung besitzt
	if !rpcrequest.ConnectionIsOpen(req) {
		o.kernel.LogPrint("RPC", "Process aborted, connection closed")
		return nil
	}

	// Es wird geprüft ob die Angeforderte Anzahl an Parametern vorhanden ist
	if len(req.Parms) != len(o.signature.Params) {
		return fmt.Errorf("EnterFunctionCall: invalid parameters")
	}

	// Es wird ein neues Request Objekt
	request := newSharedFunctionRequest(o.kernel, o.signature.ReturnType, req)

	// Es wird geprüft ob das Req Objekt eine Verbindung besitzt
	if !rpcrequest.ConnectionIsOpen(req) {
		o.kernel.LogPrint("RPC", "Process aborted, connection closed")
		return nil
	}

	// Die Loop Aufgabe wird erzeugt
	kernelLoopOperation := eventloop.NewKernelEventLoopFunctionOperation(func(_ *v8.Context, lopr *types.KernelLoopOperation) {
		err := functionCallInEventloopInit(o, request, req)
		if err != nil {
			lopr.SetError(err)
			return
		}
		lopr.SetResult(nil)
	})

	// Die Funktion wird an den Eventloop des Kernels übergeben
	if err := o.kernel.AddToEventLoop(kernelLoopOperation); err != nil {
		return fmt.Errorf("EnterFunctionCall: " + err.Error())
	}

	// Der Vorgang wird in einer neuen Goroutine durchgeführt
	go func() {
		// Es wird darauf gewartet das die Loop Operation erfolgreich abgeschlossen wurde
		_, err := kernelLoopOperation.WaitOfResponse()
		if err != nil {
			panic(err)
		}

		// Die Daten werden aus dem Request ausgelesen
		returnData, err := request.waitOfResponse()
		if err != nil {
			// Es wird geprüft ob die Verbindung getrennt wurde, wenn ja, wird der Vorgang abgebrochen
			if !rpcrequest.ConnectionIsOpen(req) {
				o.kernel.LogPrint("RPC", "Process aborted, connection closed")
				return
			}

			// Es handelt sich um einen Mysteriösen Fehler
			return //&types.MultiError{ErrorValue: fmt.Errorf("EnterFunctionCall: " + err.Error()), VmErrorValue: fmt.Errorf("internal error")}
		}

		// Es wird geprüft ob daten zurückgelifert wurden
		if returnData == nil {
			// Es wird geprüft ob die Verbindung getrennt wurde, wenn ja, wird der Vorgang abgebrochen
			if !rpcrequest.ConnectionIsOpen(req) {
				o.kernel.LogPrint("RPC", "Process aborted, connection closed")
				return
			}

			// Es handelt sich um einen Mysteriösen Fehler
			return //&types.MultiError{ErrorValue: fmt.Errorf("EnterFunctionCall: invalid returned data"), VmErrorValue: fmt.Errorf("internal error")}
		}

		// Diese Funktion wird aufgerufen, sobald die Antwort Übermittelt wurde
		resolveTransmittedData := func() {
			request.clearAndDestroy()
		}

		// Diese Funktion wird aufgerufen, wenn das übermitteln der Daten fehlgeschlagen ist
		rejectTransmittedData := func() {
		}

		// Das Rückgabe Objekt wird erstellt
		returnObject := &types.FunctionCallReturn{
			FunctionCallState: returnData,
			Resolve:           resolveTransmittedData,
			Reject:            rejectTransmittedData,
		}

		// Das Rückgabe Objekt wird zurückgegeben
		req.Resolve <- returnObject

		// Log
		o.kernel.LogPrint(fmt.Sprintf("RPC(%s)", req.ProcessLog.GetID()), "'%s' has return", o.name)
	}()

	// Das Ergebniss wird zurückgegeben
	return nil
}

// callInKernelEventLoopCheck überprüft den Status eines Promises in der Kernel-Eventschleife.
// Bei einem Pending-Promise plant es die nächste Überprüfung ohne aktives Warten.
// Bei einem Rejected-Promise führt es einen Microtask-Checkpoint durch.
func callInKernelEventLoopCheck(o *SharedFunction, ctx *v8.Context, prom *v8.Promise, request *SharedFunctionRequest, req *types.RpcRequest) error {
	// Der Stauts des Objektes wird ermittelt
	switch prom.State() {
	case v8.Pending:
		// Planen Sie die nächste Überprüfung, ohne aktives Warten zu verwenden
		go func() {
			// Es wird 1ne Milisekunde gewartet
			time.Sleep(1 * time.Millisecond)

			// Es wird ein neues Event zum Kernel hinzugefügt
			o.kernel.AddToEventLoop(eventloop.NewKernelEventLoopFunctionOperation(func(ctx *v8.Context, klopr *types.KernelLoopOperation) {
				callInKernelEventLoopCheck(o, ctx, prom, request, req)
			}))
		}()
	case v8.Rejected:
		ctx.PerformMicrotaskCheckpoint()
	}

	// Keine Rückgabe
	return nil
}

// functionCallInEventloopFinall führt den abschließenden Schritt eines Funktionsaufrufs durch.
// Es fügt einen neuen Eintrag zur Eventschleife hinzu, prüft den Promise-Status und behandelt etwaige Fehler.
// Bei Erfolg wird das Ergebnis der Operation signalisiert.
func functionCallInEventloopFinall(o *SharedFunction, request *SharedFunctionRequest, req *types.RpcRequest, prom *v8.Promise) error {
	// Es wird ermittelt ob die Verbindung getrennt wurde
	if !req.HttpRequest.IsConnected.Bool() {
		return fmt.Errorf("connection closed")
	}

	// Die Eventloop Funktion wird erzeugt
	eventloopFunction := eventloop.NewKernelEventLoopFunctionOperation(func(ctx *v8.Context, klopr *types.KernelLoopOperation) {
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
		return fmt.Errorf("functionCallInEventloopFinall: " + err.Error())
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

// functionCallInEventloopPromiseOperation verarbeitet das Ergebnis eines Funktionsaufrufs, der ein Promise zurückgibt.
// Es prüft, ob die Verbindung noch besteht, behandelt das Promise und führt die finalen Schritte des Funktionsaufrufs durch.
// Bei Erfolg wird das Ergebnis der Operation signalisiert.
func functionCallInEventloopPromiseOperation(o *SharedFunction, request *SharedFunctionRequest, req *types.RpcRequest, result *v8.Value) error {
	// Es wird ermittelt ob die Verbindung getrennt wurde
	if !req.HttpRequest.IsConnected.Bool() {
		return fmt.Errorf("connection closed")
	}

	// Die Eventloop Funktion wird erzeugt
	eventloopFunction := eventloop.NewKernelEventLoopFunctionOperation(func(ctx *v8.Context, klopr *types.KernelLoopOperation) {
		// Es wird ermittelt ob die Verbindung getrennt wurde
		if !req.HttpRequest.IsConnected.Bool() {
			klopr.SetError(fmt.Errorf("connection closed"))
			return
		}

		// Es wird geprüft ob es sich um ein Promises handelt
		if !result.IsPromise() {
			panic("isnr promise")
		}

		// Das Promises Objekt wird erzeugt
		prom, err := result.AsPromise()
		if err != nil {
			panic(err)
		}

		// Wird ausgeführt wenn die Funktion erfolgreich aufgerufen wurde
		prom.Then(func(info *v8.FunctionCallbackInfo) *v8.Value {
			request.functionCallFinal()
			return v8.Undefined(info.Context().Isolate())
		}, func(info *v8.FunctionCallbackInfo) *v8.Value {
			request.functionCallException(info.Args()[0].String())
			return v8.Undefined(info.Context().Isolate())
		})

		// Der 5. Schritt des Funktionsaufrufes wird durchgeführt
		if err := functionCallInEventloopFinall(o, request, req, prom); err != nil {
			return
		}

		// Signalisiert dass der Vorgang erfolgreich war
		klopr.SetResult(nil)
	})

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if err := o.kernel.AddToEventLoop(eventloopFunction); err != nil {
		return fmt.Errorf("functionCallInEventloopPromiseOperation: " + err.Error())
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

// functionCallInEventloop führt den vorbereiteten Funktionsaufruf innerhalb der Eventschleife aus.
// Es prüft, ob die Verbindung noch besteht, führt die Funktion aus und behandelt das Ergebnis.
// Bei Erfolg wird das Ergebnis der Operation signalisiert.
func functionCallInEventloop(o *SharedFunction, request *SharedFunctionRequest, req *types.RpcRequest, proxFunction *v8.Function, proxArguments []v8.Valuer) error {
	// Es wird ermittelt ob die Verbindung getrennt wurde
	if !req.HttpRequest.IsConnected.Bool() {
		return fmt.Errorf("connection closed")
	}

	// Die Eventloop Funktion wird erzeugt
	eventloopFunction := eventloop.NewKernelEventLoopFunctionOperation(func(ctx *v8.Context, klopr *types.KernelLoopOperation) {
		// Es wird ermittelt ob die Verbindung getrennt wurde
		if !req.HttpRequest.IsConnected.Bool() {
			klopr.SetError(fmt.Errorf("connection closed"))
			return
		}

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
		return fmt.Errorf("functionCallInEventloop: " + err.Error())
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

// functionCallInEventloopProxyObjectPrepare bereitet den Proxy-Objekt-Funktionsaufruf innerhalb der Eventschleife vor.
// Es erstellt die finalen Argumente, setzt den JavaScript-Code für den Proxy-Wrap,
// führt die Funktion in der Eventschleife aus und behandelt mögliche Fehler.
// Bei Erfolg wird das Ergebnis der Operation signalisiert.
func functionCallInEventloopProxyObjectPrepare(o *SharedFunction, request *SharedFunctionRequest, req *types.RpcRequest, requestObj *v8.Object, convertedValues []v8.Valuer) error {
	// Es wird ermittelt ob die Verbindung getrennt wurde
	if !req.HttpRequest.IsConnected.Bool() {
		return fmt.Errorf("connection closed")
	}

	// Die Eventloop Funktion wird erzeugt
	eventloopFunction := eventloop.NewKernelEventLoopFunctionOperation(func(ctx *v8.Context, klopr *types.KernelLoopOperation) {
		// Die Finalen Argumente werden erstellt
		finalArguments := make([]v8.Valuer, 0)
		finalArguments = append(finalArguments, requestObj)
		finalArguments = append(finalArguments, convertedValues...)

		// Legt den JS Code fest, dieser wird als Wrapper ausgeführt
		code := `
		(funct, proxyobject, ...parms) => {
			console = { log: proxyobject.proxyShieldConsoleLog, error: proxyobject.proxyShieldErrorLog };
			clearInterval = () => proxyobject.clearInterval();
			clearTimeout = () => proxyobject.clearTimeout();
			setInterval = () => proxyobject.setInterval();
			setTimeout = () => proxyobject.setTimeout();
			Resolve = (...parms) =>  proxyobject.resolve(...parms);
			Promise = class vnh1Promise extends Promise {
				constructor(executor) {
					const {resolveProxy, rejectProxy} = proxyobject.newPromise();
					const wrappedExecutor = (resolve, reject) => {
						executor(
							(value) => {
								resolveProxy();
								resolve(value);
							},
							(reason) => {
								rejectProxy();
								reject(reason);
							}
						);
					};
					super(wrappedExecutor);
				}
			}
			return funct(...parms);
		}`

		// Der Code für die Proxy Shield Funktion wird ersteltl
		procxyFunction, err := ctx.RunScript(code, "rpc_function_call_proxy_shield.js")
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
		proxyObject, err := makeProxyForRPCCall(ctx, request)
		if err != nil {
			return
		}

		// Die Argumente für den Proxy werden erstellt
		proxArguments := []v8.Valuer{o.v8Function, proxyObject}
		proxArguments = append(proxArguments, finalArguments...)

		// Es wird ermittelt ob die Verbindung getrennt wurde
		if !req.HttpRequest.IsConnected.Bool() {
			klopr.SetError(fmt.Errorf("connection closed"))
			return
		}

		// Der 3. Schritt des Funktionsaufrufes wird durchgeführt
		if err := functionCallInEventloop(o, request, req, proxFunction, proxArguments); err != nil {
			return
		}

		// Signalisiert dass der Vorgang erfolgreich war
		klopr.SetResult(nil)
	})

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if err := o.kernel.AddToEventLoop(eventloopFunction); err != nil {
		return fmt.Errorf("functionCallInEventloopProxyObjectPrepare: " + err.Error())
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

// functionCallInEventloopInit initialisiert einen Funktionsaufruf innerhalb der Eventschleife.
// Es prüft, ob die Verbindung besteht, wandelt die Parameter um, erstellt ein Request-Objekt,
// und führt die vorbereitenden Schritte des Funktionsaufrufs durch.
// Die Funktion wird zur Eventschleife hinzugefügt und das Ergebnis des Aufrufs wird verarbeitet.
func functionCallInEventloopInit(o *SharedFunction, request *SharedFunctionRequest, req *types.RpcRequest) error {
	// Es wird ermittelt ob die Verbindung getrennt wurde
	if !req.HttpRequest.IsConnected.Bool() {
		return fmt.Errorf("connection closed")
	}

	// Die Eventloop Funktion wird erzeugt
	eventloopFunction := eventloop.NewKernelEventLoopFunctionOperation(func(ctx *v8.Context, klopr *types.KernelLoopOperation) {
		// Die Parameter werden umgewandelt
		convertedValues, err := convertRequestParametersToV8Parameters(ctx.Isolate(), o.signature.Params, req.Parms)
		if err != nil {
			return
		}

		// Das Request Objekt wird erstellt
		requestObj, err := makeSharedFunctionObject(ctx, request, req)
		if err != nil {
			return
		}

		// Es wird ermittelt ob die Verbindung getrennt wurde
		if !req.HttpRequest.IsConnected.Bool() {
			klopr.SetError(fmt.Errorf("connection closed"))
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
		return err
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

// convertRequestParametersToV8Parameters wandelt die RPC-Argumente in V8-Argumente für den aktuellen Kontext um.
// Es überprüft die Datentypen und konvertiert sie in die entsprechenden V8-Typen.
// Bei einem Fehler wird eine entsprechende Fehlermeldung zurückgegeben.
func convertRequestParametersToV8Parameters(iso *v8.Isolate, parmTypes []string, reqparms []*types.FunctionParameterCapsle) ([]v8.Valuer, error) {
	// Es wird versucht die Paraemter in den Richtigen v8 Datentypen umzuwandeln
	convertedValues := make([]v8.Valuer, 0)
	for hight, item := range reqparms {
		// Es wird geprüft ob der Datentyp gewünscht ist
		if item.CType != parmTypes[hight] {
			return nil, fmt.Errorf("convertRequestParametersToV8Parameters: not same parameter")
		}

		// Es wird geprüft ob es sich um einen Zulässigen Datentypen handelt
		switch item.CType {
		case "boolean":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, fmt.Errorf("convertRequestParametersToV8Parameters: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		case "number":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, fmt.Errorf("convertRequestParametersToV8Parameters: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		case "string":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, fmt.Errorf("convertRequestParametersToV8Parameters: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		case "array":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, fmt.Errorf("convertRequestParametersToV8Parameters: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		case "object":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, fmt.Errorf("convertRequestParametersToV8Parameters: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		case "bytes":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, fmt.Errorf("convertRequestParametersToV8Parameters: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		default:
			return nil, fmt.Errorf("convertRequestParametersToV8Parameters: unsuported datatype")
		}
	}

	// Rückgabe ohne Fehler
	return convertedValues, nil
}
