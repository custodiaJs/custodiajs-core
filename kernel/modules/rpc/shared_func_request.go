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
	"encoding/json"
	"fmt"
	"strings"
	"vnh1/static"
	"vnh1/types"
	"vnh1/utils"

	rpcrequest "vnh1/utils/rpc_request"

	v8 "rogchap.com/v8go"
)

// Sendet eine Erfolgreiche Antwort zurück
func (o *SharedFunctionRequest) resolveFunctionCall(info *v8.FunctionCallbackInfo) *v8.Value {
	// Es wird geprüft ob das Objekt zerstört wurde
	if o.IsClosedAndDestroyed() {
		panic("destroyed object")
	}

	// Es wird geprüft ob die v8.FunctionCallbackInfo "info" NULL ist
	// In diesem fall muss ein Panic ausgelöst werden, da es keine Logik gibt, wann diese Eintreten kann
	if info == nil {
		panic("vm information callback infor null error")
	}

	// Es wird geprüft ob SharedFunctionRequest "o" NULL ist
	if !validateSharedFunctionRequest(o) {
		// Es wird ein Exception zurückgegeben
		utils.V8ContextThrow(info.Context(), "invalid function share")

		// Undefined wird zurückgegeben
		return v8.Undefined(info.Context().Isolate())
	}

	// Speichert alle FunktionsStates ab
	resolves := &types.FunctionCallState{
		Return: make([]*types.FunctionCallReturnData, 0),
		State:  "ok",
	}

	// Die Einzelnen Parameter werden abgearbeitet
	for _, item := range info.Args() {
		// Das Datentyp wird ermittelt
		var responseData *types.FunctionCallReturnData
		if item == nil {
			responseData = &types.FunctionCallReturnData{CType: "null", Value: nil}
		} else if item.IsUndefined() || item.IsNull() {
			responseData = &types.FunctionCallReturnData{CType: "undefined", Value: nil}
		} else {
			switch {
			case item.IsString():
				responseData = &types.FunctionCallReturnData{CType: "string", Value: item.String()}
			case item.IsNumber():
				switch {
				case item.IsBigInt():
					responseData = &types.FunctionCallReturnData{CType: "number", Value: item.BigInt().String()}
				case item.IsInt32():
					responseData = &types.FunctionCallReturnData{CType: "number", Value: item.Int32()}
				case item.IsUint32():
					responseData = &types.FunctionCallReturnData{CType: "number", Value: item.Uint32()}
				case item.IsNumber():
					responseData = &types.FunctionCallReturnData{CType: "number", Value: item.Number()}
				default:
					responseData = &types.FunctionCallReturnData{CType: "number", Value: item.Integer()}
				}
			case item.IsBoolean():
				responseData = &types.FunctionCallReturnData{CType: "boolean", Value: item.Boolean()}
			case item.IsObject():
				fmt.Println("Wert ist ein Array:")
			case item.IsArray():
				fmt.Println("Wert ist ein Array:")
			case item.IsFunction():
				fmt.Println("Wert ist ein Array:")
			default:
				// Es wird ein Javascript Fehler zurückgegeben
				utils.V8ContextThrow(info.Context(), "unsupported datatype for shared function response")
				return nil
			}
		}

		// Der Eintrag wird abgespeichert
		resolves.Return = append(resolves.Return, responseData)
	}

	// Die Antwort wird zurückgesendet
	o.resolveChan <- resolves

	// Es wird Signalisiert dass eine Antwort gesendet wurde
	o._wasResponded = true

	// Es ist kein Fehler aufgetreten
	return nil
}

// Sendet eine Rejectantwort zurück
func (o *SharedFunctionRequest) rejectFunctionCall(info *v8.FunctionCallbackInfo) *v8.Value {
	// Es wird geprüft ob die v8.FunctionCallbackInfo "info" NULL ist
	// In diesem fall muss ein Panic ausgelöst werden, da es keine Logik gibt, wann diese Eintreten kann
	if info == nil {
		panic("vm information callback infor null error")
	}

	// Es wird geprüft ob SharedFunctionRequest "o" NULL ist
	if !validateSharedFunctionRequest(o) {
		// Es wird ein Exception zurückgegeben
		utils.V8ContextThrow(info.Context(), "invalid function share")

		// Undefined wird zurückgegeben
		return v8.Undefined(info.Context().Isolate())
	}

	// Die Einzelnen Parameter werden abgearbeitet
	extractedStrings := make([]string, 0)
	for _, item := range info.Args() {
		switch {
		case item.IsString():
			extractedStrings = append(extractedStrings, item.String())
		default:
			utils.V8ContextThrow(info.Context(), "unsupported datatype for shared function error response, only strings allowed")
			return nil
		}
	}

	// Der Finale Fehler wird gebaut
	finalErrorStr := ""
	if len(extractedStrings) > 0 {
		finalErrorStr = strings.Join(extractedStrings, " ")
	}

	// Die Antwort wird zurückgesendet
	o.resolveChan <- &types.FunctionCallState{Error: finalErrorStr, State: "failed"}

	// Es wird Signalisiert dass eine Antwort gesendet wurde
	o._wasResponded = true

	// Es ist kein Fehler aufgetreten
	return nil
}

// Gibt an ob das Objekt zerstört wurde
func (o *SharedFunctionRequest) IsClosedAndDestroyed() bool {
	// Rückgabe des Destroyed Wertes
	return o._destroyed
}

// Räumt auf und Zerstört das Objekt
func (o *SharedFunctionRequest) clearAndDestroy() {
	// Es wird geprüft ob das Objekt zerstört wurde
	if o.IsClosedAndDestroyed() {
		panic("destroyed object")
	}

	// Kernel Log
	o.kernel.LogPrint(fmt.Sprintf("RPC(%s)", o._rprequest.ProcessLog.GetID()), "request closed")
}

// Wird ausgeführt wenn die Funktion zuende aufgerufen wurde
func (o *SharedFunctionRequest) functionCallFinal() error {
	// Es wird geprüft ob das Objekt zerstört wurde
	if o.IsClosedAndDestroyed() {
		panic("destroyed object")
	}

	// Es wird geprüft ob eine Antwort gesendet wurde
	if !o.hasResponse() {
		// Der Timer zum abbrechen des Vorganges wird gestartet
		o.startTimeoutTimer()
	}

	// Kernel Log
	o.kernel.LogPrint(fmt.Sprintf("RPC(%s)", o._rprequest.ProcessLog.GetID()), "function call finalized")

	// Es wird nichts zurückgegeben
	return nil
}

// Wird ausgeführt wenn ein Throw durch die Funktion ausgelöst wird
func (o *SharedFunctionRequest) functionCallException(msg string) error {
	// Es wird geprüft ob das Objekt zerstört wurde
	if o.IsClosedAndDestroyed() {
		panic("destroyed object")
	}

	// Die Antwort wird zurückgesendet
	o.resolveChan <- &types.FunctionCallState{Error: msg, State: "exception"}

	// Es wird Signalisiert dass eine Antwort gesendet wurde
	o._wasResponded = true

	// Rückgabe
	return nil
}

// Proxy Shielded, Set Timeout funktion
func (o *SharedFunctionRequest) proxyShield_SetTimeout(info *v8.FunctionCallbackInfo) *v8.Value {
	fmt.Println("D: TIMER")
	return v8.Undefined(info.Context().Isolate())
}

// Proxy Shielded, Set Interval funktion
func (o *SharedFunctionRequest) proxyShield_SetInterval(info *v8.FunctionCallbackInfo) *v8.Value {
	fmt.Println("D: TIMER")
	return v8.Undefined(info.Context().Isolate())
}

// Proxy Shielded, Clear Timeout funktion
func (o *SharedFunctionRequest) proxyShield_ClearTimeout(info *v8.FunctionCallbackInfo) *v8.Value {
	fmt.Println("D: TIMER")
	return v8.Undefined(info.Context().Isolate())
}

// Proxy Shielded, Clear Interval funktion
func (o *SharedFunctionRequest) proxyShield_ClearInterval(info *v8.FunctionCallbackInfo) *v8.Value {
	fmt.Println("D: TIMER")
	return v8.Undefined(info.Context().Isolate())
}

// Signalisiert dass ein neuer Promises erzeugt wurde und gibt die Entsprechenden Funktionen zurück
func (o *SharedFunctionRequest) proxyShield_NewPromise(info *v8.FunctionCallbackInfo) *v8.Value {
	v8Object := v8.NewObjectTemplate(info.Context().Isolate())
	v8Object.Set("resolveProxy", v8.NewFunctionTemplate(info.Context().Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
		o.kernel.LogPrint(fmt.Sprintf("RPC(%s)", o._rprequest.ProcessLog.GetID()), "Promise was resolved")
		return nil
	}))
	v8Object.Set("rejectProxy", v8.NewFunctionTemplate(info.Context().Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
		o.kernel.LogPrint(fmt.Sprintf("RPC(%s)", o._rprequest.ProcessLog.GetID()), "Promise was rejected")
		return nil
	}))

	// Log
	o.kernel.LogPrint(fmt.Sprintf("RPC(%s)", o._rprequest.ProcessLog.GetID()), "New Promise registrated")

	// Das Objekt wird in ein Wert umgewandelt
	obj, _ := v8Object.NewInstance(info.Context())

	// Das Objekt wird zurückgegeben
	return obj.Value
}

// Proxy Shielded, Console Log Funktion
func (o *SharedFunctionRequest) proxyShield_ConsoleLog(info *v8.FunctionCallbackInfo) *v8.Value {
	// Es werden alle Stringwerte Extrahiert
	extracted := convertArguments(info)

	// Es wird ein String aus der Ausgabe erzeugt
	outputStr := strings.Join(extracted, " ")

	// Die Ausgabe wird an den Console Cache übergeben
	o.kernel.Console().Log(fmt.Sprintf("RPC(%s):-$ %s", strings.ToUpper(o._rprequest.ProcessLog.GetID()), outputStr))

	// Rückgabe ohne Fehler
	return v8.Undefined(info.Context().Isolate())
}

// Proxy Shielded, Console Log Funktion
func (o *SharedFunctionRequest) proxyShield_ErrorLog(info *v8.FunctionCallbackInfo) *v8.Value {
	// Es werden alle Stringwerte Extrahiert
	extracted := convertArguments(info)

	// Es wird ein String aus der Ausgabe erzeugt
	outputStr := strings.Join(extracted, " ")

	// Die Ausgabe wird an den Console Cache übergeben
	o.kernel.Console().ErrorLog(fmt.Sprintf("RPC(%s):-$ %s", strings.ToUpper(o._rprequest.ProcessLog.GetID()), outputStr))

	// Rückgabe ohne Fehler
	return v8.Undefined(info.Context().Isolate())
}

// Wartet auf eine Antwort
func (o *SharedFunctionRequest) waitOfResponse() (*types.FunctionCallState, error) {
	// Es wird ein neuer Response Waiter erzeugt
	responseWaiter, err := newRequestResponseWaiter(o)
	if err != nil {
		return nil, err
	}

	// Es wird auf den Status gewartet
	finalResolve, err := responseWaiter.WaitOfState()
	if err != nil {
		return nil, err
	}

	// Das Ergebniss wird zurückgegeben
	return finalResolve, nil
}

// Gibt an ob eine Antwort verfügbar ist
func (o *SharedFunctionRequest) hasResponse() bool {
	// Der Mutex wird verwendet
	return o._wasResponded
}

// Startet den Timer, welcher den Vorgang nach erreichen des Timeouts, abbricht
func (o *SharedFunctionRequest) startTimeoutTimer() {
}

// Erstellt einen neuen SharedFunctionRequest
func newSharedFunctionRequest(kernel types.KernelInterface, returnDatatype string, rpcRequest *types.RpcRequest) *SharedFunctionRequest {
	// Das Rückgabeobjekt wird erstellt
	returnObject := &SharedFunctionRequest{
		resolveChan:     make(chan *types.FunctionCallState),
		_returnDataType: returnDatatype,
		_wasResponded:   false,
		_rprequest:      rpcRequest,
		kernel:          kernel,
	}

	// Das Objekt wird zurückgegeben
	return returnObject
}

// Überprüft ob ein SharedFunctionRequest korrekt aufgebaut ist
func validateSharedFunctionRequest(o *SharedFunctionRequest) bool {
	// Sollte die SharedFunctionRequest "o" NULL sein, wird ein False zurückgegeben
	if o == nil {
		return false
	}

	// Es wird geprüft ob die Resolve Chain NULL ist
	if o.resolveChan == nil {
		return false
	}

	// Es handelt sich um ein zulässiges Objekt
	return true
}

// Wandelt die Funktionsargumente in Strings um
func convertArguments(info *v8.FunctionCallbackInfo) []string {
	// Es werden alle Stringwerte Extrahiert
	extracted := []string{}
	for _, item := range info.Args() {
		if item.IsObject() && !item.IsArray() {
			obj, err := utils.V8ObjectToGoObject(info.Context(), item)
			if err != nil {
				utils.V8ContextThrow(info.Context(), "internal error by converting, value")
				return nil
			}
			encoded, err := json.Marshal(obj)
			if err != nil {
				utils.V8ContextThrow(info.Context(), "internal error by converting, value")
				return nil
			}
			extracted = append(extracted, string(encoded))
		} else if item.IsArray() {
			obj, err := utils.V8ArrayToGoArray(info.Context(), item)
			if err != nil {
				utils.V8ContextThrow(info.Context(), "internal error by converting, value")
				return nil
			}

			var extra []interface{}
			for _, item := range obj {
				extra = append(extra, item.Value)
			}

			encoded, err := json.Marshal(extra)
			if err != nil {
				utils.V8ContextThrow(info.Context(), "internal error by converting, value")
				return nil
			}

			extracted = append(extracted, string(encoded))
		} else if item.IsFunction() {
			if item.IsAsyncFunction() {
				extracted = append(extracted, fmt.Sprintf("ASYNC:=%p", item))
			} else {
				extracted = append(extracted, fmt.Sprintf("SYNC:=%p", item))
			}
		} else {
			extracted = append(extracted, item.String())
		}
	}
	return extracted
}

// Die Funktion wird erstellt
func makeSharedFunctionObject(context *v8.Context, request *SharedFunctionRequest, rrpcrequest *types.RpcRequest) (*v8.Object, error) {
	// Das Requestobjekt wird ersellt
	objTemplate := v8.NewObjectTemplate(context.Isolate())

	// Die Resolve Funktion wird festgelegt
	if err := objTemplate.Set("Resolve", v8.NewFunctionTemplate(context.Isolate(), request.resolveFunctionCall)); err != nil {
		return nil, fmt.Errorf("makeSharedFunctionObject: " + err.Error())
	}

	// Die Reject Funktion wird festgelegt
	if err := objTemplate.Set("Reject", v8.NewFunctionTemplate(context.Isolate(), request.rejectFunctionCall)); err != nil {
		return nil, fmt.Errorf("makeSharedFunctionObject: " + err.Error())
	}

	// Das Objekt wird erzeugt
	obj, err := objTemplate.NewInstance(context)
	if err != nil {
		return nil, fmt.Errorf("makeSharedFunctionObject: " + err.Error())
	}

	// Es wird ein neues Objekt erzeugt, dieses Objekt wird verwendet um den Aktuellen Request Darzustellen
	var rpcType string
	switch rrpcrequest.RequestType {
	case static.HTTP_REQUEST:
		// Es wird geprüft ob der http Request vorhanden ist
		if !rpcrequest.IsHttpRequest(rrpcrequest) {
			return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: invalid http request, has no http request data")
		}

		// Der Type wird Signalisiert
		rpcType = "http"

		// Die Einzelnenen Werte werden umgewandelt
		isConnected := func(info *v8.FunctionCallbackInfo) *v8.Value {
			value, err := v8.NewValue(context.Isolate(), rrpcrequest.HttpRequest.IsConnected.Bool())
			if err != nil {
				panic(err)
			}
			return value
		}

		// Die Cookies werden Extrahiert
		cookies := v8.NewObjectTemplate(context.Isolate())
		for _, item := range rrpcrequest.HttpRequest.Cookies {
			// Es wird ein neues Objekt erzeugt
			cookieObject := v8.NewObjectTemplate(context.Isolate())
			cookieObject.Set("Value", item.Value)
			cookieObject.Set("Domain", item.Domain)
			cookieObject.Set("Path", item.Path)
			cookieObject.Set("Expires", item.RawExpires)

			// Der Eintrag wird hinzugefügt
			if err := cookies.Set(item.Name, cookieObject); err != nil {
				panic(err)
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}
		}

		// Der Header wird vorbereitet
		headersTemplate := v8.NewObjectTemplate(context.Isolate())
		headers, err := headersTemplate.NewInstance(context)
		if err != nil {
			return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
		}

		// Die Header werden extrahiert
		for k, v := range rrpcrequest.HttpRequest.Header {
			// Es wird ein neues Slices erzeugt
			sliceV8, err := context.RunScript("(function() { return []; })();", "slice.js")
			if err != nil {
				fmt.Println(err)
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}

			// Das Objekt wird ausgelesen
			sliceObject, err := sliceV8.AsObject()
			if err != nil {
				panic(err)
			}

			// Die Einzelnen Werte werden umgewandelt
			for _, value := range v {
				// Der Wert wird umgewandelt
				v8Value, err := v8.NewValue(context.Isolate(), value)
				if err != nil {
					panic(err)
				}

				// Der Wert wird hinzugefügt
				sliceObject.Object().MethodCall("push", v8Value)
			}

			// Der Eintrag wird hinzugefügt
			if err := headers.Set(k, sliceObject); err != nil {
				panic(err)
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}
		}

		// Das Http Objekt wird erzeugt
		httpObj := v8.NewObjectTemplate(context.Isolate())

		// Die Werte werden hinzugefügt
		if err := httpObj.Set("IsConnected", v8.NewFunctionTemplate(context.Isolate(), isConnected)); err != nil {
			panic(err)
			return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
		}
		if err := httpObj.Set("ContentLength", float64(rrpcrequest.HttpRequest.ContentLength)); err != nil {
			return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
		}
		if err := httpObj.Set("Host", rrpcrequest.HttpRequest.Host); err != nil {
			return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
		}
		if err := httpObj.Set("Proto", rrpcrequest.HttpRequest.Proto); err != nil {
			return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
		}
		if err := httpObj.Set("RemoteAddr", rrpcrequest.HttpRequest.RemoteAddr); err != nil {
			return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
		}
		if err := httpObj.Set("RequestURI", rrpcrequest.HttpRequest.RequestURI); err != nil {
			return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
		}
		if err := httpObj.Set("Cookies", cookies); err != nil {
			return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
		}

		// Das Finale Objekt wird erzeugt
		http, err := httpObj.NewInstance(context)
		if err != nil {
			return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
		}

		// Die Header werden hinzugefügt
		if err := http.Set("Headers", headers); err != nil {
			return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
		}

		// Das Objekt wird abgespeichert
		if err := obj.Set("http", http); err != nil {
			return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
		}
	case static.WEBSOCKET_REQUEST:
		// Der Type wird Signalisiert
		rpcType = "ws"
	case static.IPC_REQUEST:
		// Der Type wird Signalisiert
		rpcType = "ipc"
	default:
		return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: unkown methode")
	}

	// Der Wert wird eingelesen
	val, err := v8.NewValue(context.Isolate(), rpcType)
	if err != nil {
		return nil, fmt.Errorf("makeSharedFunctionObject: " + err.Error())
	}

	// Der Eintrag wird im Objekt hinzugefügt
	if err := obj.Set("CallMethode", val); err != nil {
		return nil, fmt.Errorf("makeSharedFunctionObject: " + err.Error())
	}

	// Rückgabe ohne Fehler
	return obj, nil
}

// Das This Objekt wird erstellt
func makeProxyForRPCCall(context *v8.Context, request *SharedFunctionRequest) (*v8.Object, error) {
	// Das Requestobjekt wird ersellt
	obj := v8.NewObjectTemplate(context.Isolate())

	// Die Funktionen werden hinzugefügt
	if err := obj.Set("proxyShieldConsoleLog", v8.NewFunctionTemplate(context.Isolate(), request.proxyShield_ConsoleLog)); err != nil {
		return nil, fmt.Errorf("makeProxyForRPCCall: " + err.Error())
	}
	if err := obj.Set("proxyShieldErrorLog", v8.NewFunctionTemplate(context.Isolate(), request.proxyShield_ErrorLog)); err != nil {
		return nil, fmt.Errorf("makeProxyForRPCCall: " + err.Error())
	}
	if err := obj.Set("clearInterval", v8.NewFunctionTemplate(context.Isolate(), request.proxyShield_ClearInterval)); err != nil {
		return nil, fmt.Errorf("makeProxyForRPCCall: " + err.Error())
	}
	if err := obj.Set("clearTimeout", v8.NewFunctionTemplate(context.Isolate(), request.proxyShield_ClearTimeout)); err != nil {
		return nil, fmt.Errorf("makeProxyForRPCCall: " + err.Error())
	}
	if err := obj.Set("setInterval", v8.NewFunctionTemplate(context.Isolate(), request.proxyShield_SetInterval)); err != nil {
		return nil, fmt.Errorf("makeProxyForRPCCall: " + err.Error())
	}
	if err := obj.Set("setTimeout", v8.NewFunctionTemplate(context.Isolate(), request.proxyShield_SetTimeout)); err != nil {
		return nil, fmt.Errorf("makeProxyForRPCCall: " + err.Error())
	}
	if err := obj.Set("resolve", v8.NewFunctionTemplate(context.Isolate(), request.resolveFunctionCall)); err != nil {
		return nil, fmt.Errorf("makeProxyForRPCCall: " + err.Error())
	}
	if err := obj.Set("reject", v8.NewFunctionTemplate(context.Isolate(), request.rejectFunctionCall)); err != nil {
		return nil, fmt.Errorf("makeProxyForRPCCall: " + err.Error())
	}
	if err := obj.Set("newPromise", v8.NewFunctionTemplate(context.Isolate(), request.proxyShield_NewPromise)); err != nil {
		return nil, fmt.Errorf("makeProxyForRPCCall: " + err.Error())
	}

	// Das Finale Objekt wird erstellt
	fobj, err := obj.NewInstance(context)
	if err != nil {
		return nil, fmt.Errorf("makeProxyForRPCCall: " + err.Error())
	}

	// Rückgabe ohne Fehler
	return fobj, nil
}
