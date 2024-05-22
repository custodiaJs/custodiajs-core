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
	rpcrequest "vnh1/utils/rpc_request"

	v8 "rogchap.com/v8go"
)

// Gibt den Namen der Funktion zurück
func (o *SharedFunction) GetName() string {
	// Es wird geprüft ob die Aktuelle SharedFunction "o" NULL ist
	if o == nil {
		return ""
	}

	// Rückgabe
	return o.name
}

// Gibt die Parameterdatentypen welche die Funktion erwartet zurück
func (o *SharedFunction) GetParmTypes() []string {
	// Es wird geprüft ob die Aktuelle SharedFunction "o" NULL ist
	if o == nil {
		return make([]string, 0)
	}

	// Rückgabe
	return o.signature.Params
}

// Gibt den Rückgabedatentypen zurück
func (o *SharedFunction) GetReturnDatatype() string {
	// Es wird geprüft ob die Aktuelle SharedFunction "o" NULL ist
	if o == nil {
		return ""
	}

	// Rückgabe
	return o.signature.ReturnType
}

// Fügt ein Event hinzu, dieses Event wird ausgeführt sobald eine neuer Funktionsaufruf angenommen wurde
func (o *SharedFunction) AddOnRequestFunction(funcv8 *v8.Function) error {
	// Die Funktion wird abgespeichert
	o.eventOnRequest = append(o.eventOnRequest, funcv8)

	// Log
	o.kernel.LogPrint(fmt.Sprintf("RPC(%s)", o.signature.FunctionName), "New 'OnRequest' event added for function %s", o.name)

	// Rückgabe ohne Fehler
	return nil
}

// Fügt ein Event hinzu, dieses Event wird ausgeführt sobald ein neuer Funktionsaufruf Fehlschlägt, bevor er genau zugeordnet werden kann
func (o *SharedFunction) AddOnRequestFailFunction(funcv8 *v8.Function) error {
	// Die Funktion wird abgespeichert
	o.eventOnRequestFail = append(o.eventOnRequest, funcv8)

	// Log
	o.kernel.LogPrint(fmt.Sprintf("RPC(%s)", o.signature.FunctionName), "New 'OnRequestFail' event added for function %s", o.name)

	// Rückgabe ohne Fehler
	return nil
}

// Wird verwendet um sicherzustellen dass alle Mikroaufgaben eines RPC Aufrufes durchgeführt wurden
func (o *SharedFunction) _callInKernelEventLoopCheck(ctx *v8.Context, prom *v8.Promise, request *SharedFunctionRequest) error {
	switch prom.State() {
	case v8.Pending:
		// Planen Sie die nächste Überprüfung, ohne aktives Warten zu verwenden
		go func() {
			time.Sleep(1 * time.Millisecond)
			o.kernel.AddFunctionCallToEventLoop(func(ctx *v8.Context) error {
				return o._callInKernelEventLoopCheck(ctx, prom, request)
			})
		}()
	case v8.Fulfilled:
		ctx.PerformMicrotaskCheckpoint()
	case v8.Rejected:
		ctx.PerformMicrotaskCheckpoint()
	}
	return nil
}

// Wird verwendet um die Funktion innerhalb des Kernels aufzurufen
func (o *SharedFunction) _callInKernelEventLoop(ctx *v8.Context, request *SharedFunctionRequest, req *types.RpcRequest) error {
	// Die Parameter werden umgewandelt
	convertedValues, err := convertRequestParametersToV8Parameters(ctx.Isolate(), o.signature.Params, req.Parms)
	if err != nil {
		return fmt.Errorf("EnterFunctionCall: " + err.Error())
	}

	// Das Request Objekt wird erstellt
	requestObj, err := makeSharedFunctionObject(ctx, request, req)
	if err != nil {
		return fmt.Errorf("EnterFunctionCall: " + err.Error())
	}

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
		return fmt.Errorf("EnterFunctionCall: " + err.Error())
	}

	// Es wird geprüft ob es sich um eine Funktion handelt,
	// wenn ja wird die Funktion extrahiert
	proxFunction, err := procxyFunction.AsFunction()
	if err != nil {
		return fmt.Errorf("EnterFunctionCall: " + err.Error())
	}

	// Das Proxy Objekt wird erzeugt
	proxyObject, err := makeProxyForRPCCall(ctx, request)
	if err != nil {
		return fmt.Errorf("EnterFunctionCall: " + err.Error())
	}

	// Die Argumente für den Proxy werden erstellt
	proxArguments := []v8.Valuer{o.v8Function, proxyObject}
	proxArguments = append(proxArguments, finalArguments...)

	// Die Funktion wird ausgdeneführt
	result, err := proxFunction.Call(v8.Undefined(ctx.Isolate()), proxArguments...)
	if err != nil {
		panic(err)
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

	// Es wird ein neuer Eintrag zu der Event Schleife hinzugefügt
	o.kernel.AddFunctionCallToEventLoop(func(ctx *v8.Context) error { return o._callInKernelEventLoopCheck(ctx, prom, request) })

	// Es ist kein Fehler aufgetreten
	return nil
}

// Ruft die Geteilte Funktion auf
func (o *SharedFunction) EnterFunctionCall(req *types.RpcRequest) (*types.FunctionCallReturn, error) {
	// Es wird geprüft ob die Aktuelle SharedFunction "o" NULL ist
	if o == nil {
		return nil, &types.MultiError{ErrorValue: fmt.Errorf("EnterFunctionCall: o is null"), VmErrorValue: fmt.Errorf("internal error")}
	}

	// Es wird geprüft ob der RPC Request "req" NULL ist
	if req == nil {
		return nil, &types.MultiError{ErrorValue: fmt.Errorf("EnterFunctionCall: req is null"), VmErrorValue: fmt.Errorf("request error")}
	}

	// Es wird geprüft ob das Req Objekt eine Verbindung besitzt
	if !rpcrequest.ConnectionIsOpen(req) {
		o.kernel.LogPrint("RPC", "Process aborted, connection closed")
		return nil, nil
	}

	// Es wird geprüft ob die Angeforderte Anzahl an Parametern vorhanden ist
	if len(req.Parms) != len(o.signature.Params) {
		return nil, fmt.Errorf("EnterFunctionCall: invalid parameters")
	}

	// Es wird ein neues Request Objekt
	request := newSharedFunctionRequest(o.kernel, o.signature.ReturnType, req)

	// Es wird geprüft ob das Req Objekt eine Verbindung besitzt
	if !rpcrequest.ConnectionIsOpen(req) {
		o.kernel.LogPrint("RPC", "Process aborted, connection closed")
		return nil, nil
	}

	// Die Funktion wird an den Eventloop des Kernels übergeben
	err := o.kernel.AddFunctionCallToEventLoop(func(ctx *v8.Context) error {
		return o._callInKernelEventLoop(ctx, request, req)
	})
	if err != nil {
		return nil, err
	}

	// Die Daten werden aus dem Request ausgelesen
	returnData, err := request.waitOfResponse()
	if err != nil {
		// Es wird geprüft ob die Verbindung getrennt wurde, wenn ja, wird der Vorgang abgebrochen
		if !rpcrequest.ConnectionIsOpen(req) {
			o.kernel.LogPrint("RPC", "Process aborted, connection closed")
			return nil, nil
		}

		// Es handelt sich um einen Mysteriösen Fehler
		return nil, &types.MultiError{ErrorValue: fmt.Errorf("EnterFunctionCall: " + err.Error()), VmErrorValue: fmt.Errorf("internal error")}
	}

	// Es wird geprüft ob daten zurückgelifert wurden
	if returnData == nil {
		// Es wird geprüft ob die Verbindung getrennt wurde, wenn ja, wird der Vorgang abgebrochen
		if !rpcrequest.ConnectionIsOpen(req) {
			o.kernel.LogPrint("RPC", "Process aborted, connection closed")
			return nil, nil
		}

		// Es handelt sich um einen Mysteriösen Fehler
		return nil, &types.MultiError{ErrorValue: fmt.Errorf("EnterFunctionCall: invalid returned data"), VmErrorValue: fmt.Errorf("internal error")}
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

	// Log
	o.kernel.LogPrint(fmt.Sprintf("RPC(%s)", req.ProcessLog.GetID()), "'%s' has return", o.name)

	// Das Ergebniss wird zurückgegeben
	return returnObject, nil
}

// Wandelt die RPC Argumente in V8 Argumente für den Aktuellen Context um
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
