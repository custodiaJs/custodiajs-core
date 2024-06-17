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
	"strings"
	"time"
	"vnh1/saftychan"
	"vnh1/types"
	"vnh1/utils"
	rpcrequest "vnh1/utils/rpc_request"

	v8 "rogchap.com/v8go"
)

// Sendet eine Erfolgreiche Antwort zurück
func (o *SharedFunctionRequestContext) resolveFunctionCallbackV8(info *v8.FunctionCallbackInfo) *v8.Value {
	// Es wird geprüft ob die v8.FunctionCallbackInfo "info" NULL ist
	// In diesem fall muss ein Panic ausgelöst werden, da es keine Logik gibt, wann diese Eintreten kann
	if info == nil {
		panic("vm information callback null error")
	}

	// Es wird geprüft ob SharedFunctionRequestContext "o" NULL ist
	if !validateSharedFunctionRequestContext(o) {
		panic("SharedFunctionRequestContext 'o' is empty")
	}

	// Es wird geprüft ob das Objekt zerstört wurde
	if o.responseChan.IsClosed() || !rpcrequest.ConnectionIsOpen(o._rprequest) || requestContextIsClosedAndDestroyed(o) {
		// Der Grund für das Abbrechen des Vorganges wird ermittelt
		switch {
		case requestContextIsClosedAndDestroyed(o):
			utils.V8ContextThrow(info.Context(), "It is not possible to respond to an already fully closed RPC request.")
		case o.responseChan.IsClosed() && rpcrequest.ConnectionIsOpen(o._rprequest):
			utils.V8ContextThrow(info.Context(), "An answer has already been sent, the process was aborted.")
		case !rpcrequest.ConnectionIsOpen(o._rprequest):
			utils.V8ContextThrow(info.Context(), "Resolving could not be performed, the connection was terminated.")
		default:
			utils.V8ContextThrow(info.Context(), "Resolving could not be performed, reason: unknown.")
		}

		// Rückgabe
		return nil
	}

	// Es wird geprüft ob der Vorgang bereits beantwortet wurde,
	// wenn ja wird ein Fehler zurückgegeben dass der Vorgang bereits beantwortet wurde
	if o.wasResponsed() {
		// Es wird ein V8 Throw erzeugt und ausgeführt
		utils.V8ContextThrow(info.Context(), "")

		// Es wird ein Nill zurückgegeben
		return nil
	}

	// Die Argumente werden umgewandelt
	convertedArguments, err := utils.ConvertV8DataToGoData(info.Args())
	if err != nil {
		// Es wird ein V8 Throw erzeugt und ausgeführt
		utils.V8ContextThrow(info.Context(), err.Error())

		// Es wird ein Nill zurückgegeben
		return nil
	}

	// Speichert alle FunktionsStates ab
	resolves := &types.FunctionCallState{
		Return: make([]*types.FunctionCallReturnData, 0),
		State:  "ok",
	}

	// Die Einzelnen Parameter werden abgearbeitet
	for _, item := range convertedArguments {
		resolves.Return = append(resolves.Return, (*types.FunctionCallReturnData)(item))
	}

	// Es wird geprüft ob das ResponseChan geschlossen wurde
	if o.responseChan.IsClosed() || !rpcrequest.ConnectionIsOpen(o._rprequest) {
		// Der Grund für das Abbrechen des Vorganges wird ermittelt
		switch {
		case o.responseChan.IsClosed() && rpcrequest.ConnectionIsOpen(o._rprequest):
			utils.V8ContextThrow(info.Context(), "An answer has already been sent, the process was aborted.")
		case !rpcrequest.ConnectionIsOpen(o._rprequest):
			utils.V8ContextThrow(info.Context(), "Resolving could not be performed, the connection was terminated.")
		default:
			utils.V8ContextThrow(info.Context(), "Resolving could not be performed, reason: unknown.")
		}

		// Rückgabe
		return nil
	}

	// Die Antwort wird geschrieben
	if err := writeRequestReturnResponse(o, resolves); err != nil {
		switch err := err.(type) {
		case *types.SpecificError:
			utils.V8ContextThrow(info.Context(), err.LocalJSVMError.Error())
		default:
			utils.V8ContextThrow(info.Context(), err.Error())
		}
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

// Sendet eine Rejectantwort zurück
func (o *SharedFunctionRequestContext) rejectFunctionCallbackV8(info *v8.FunctionCallbackInfo) *v8.Value {
	// Es wird geprüft ob die v8.FunctionCallbackInfo "info" NULL ist
	// In diesem fall muss ein Panic ausgelöst werden, da es keine Logik gibt, wann diese Eintreten kann
	if info == nil {
		panic("vm information callback infor null error")
	}

	// Es wird geprüft ob das Objekt zerstört wurde
	if o.responseChan.IsClosed() || !rpcrequest.ConnectionIsOpen(o._rprequest) || requestContextIsClosedAndDestroyed(o) {
		// Der Grund für das Abbrechen des Vorganges wird ermittelt
		switch {
		case requestContextIsClosedAndDestroyed(o):
			utils.V8ContextThrow(info.Context(), "It is not possible to respond to an already fully closed RPC request.")
		case o.responseChan.IsClosed() && rpcrequest.ConnectionIsOpen(o._rprequest):
			utils.V8ContextThrow(info.Context(), "An answer has already been sent, the process was aborted.")
		case !rpcrequest.ConnectionIsOpen(o._rprequest):
			utils.V8ContextThrow(info.Context(), "Resolving could not be performed, the connection was terminated.")
		default:
			utils.V8ContextThrow(info.Context(), "Resolving could not be performed, reason: unknown.")
		}

		// Rückgabe
		return v8.Undefined(info.Context().Isolate())
	}

	// Es wird geprüft ob SharedFunctionRequestContext "o" NULL ist
	if !validateSharedFunctionRequestContext(o) {
		// Es wird ein Exception zurückgegeben
		utils.V8ContextThrow(info.Context(), "invalid function share")

		// Undefined wird zurückgegeben
		return v8.Undefined(info.Context().Isolate())
	}

	// Es wird geprüft ob der Vorgang bereits beantwortet wurde,
	// wenn ja wird ein Fehler zurückgegeben dass der Vorgang bereits beantwortet wurde
	if o.wasResponsed() {
		utils.V8ContextThrow(info.Context(), "")
		return nil
	}

	// Die Einzelnen Parameter werden abgearbeitet
	extractedStrings, err := utils.ConvertV8ValuesToString(info.Context(), info.Args())
	if err != nil {
		utils.V8ContextThrow(info.Context(), err.Error())
		return nil
	}

	// Der Finale Fehler wird gebaut
	finalErrorStr := ""
	if len(extractedStrings) > 0 {
		finalErrorStr = strings.Join(extractedStrings, " ")
	}

	// Es wird geprüft ob das ResponseChan geschlossen wurde
	if o.responseChan.IsClosed() || !rpcrequest.ConnectionIsOpen(o._rprequest) {
		// Der Fehler wird ermittelt
		switch {
		case o.responseChan.IsClosed() && rpcrequest.ConnectionIsOpen(o._rprequest):
			utils.V8ContextThrow(info.Context(), "An answer has already been sent, the process was aborted.")
		case !rpcrequest.ConnectionIsOpen(o._rprequest):
			utils.V8ContextThrow(info.Context(), "Resolving could not be performed, the connection was terminated.")
		default:
			utils.V8ContextThrow(info.Context(), "Resolving could not be performed, reason: unknown.")
		}

		// Rückgabe
		return v8.Undefined(info.Context().Isolate())
	}

	// Die Antwort wird zurückgesendet
	if err := writeRequestReturnResponse(o, &types.FunctionCallState{Error: finalErrorStr, State: "failed"}); err != nil {
		utils.V8ContextThrow(info.Context(), "Resolving could not be performed, the connection was terminated.")
	}

	// Es ist kein Fehler aufgetreten
	return v8.Undefined(info.Context().Isolate())
}

// Räumt auf und Zerstört das Objekt
func (o *SharedFunctionRequestContext) clearAndDestroy() {
	// Es wird geprüft ob das Objekt zerstört wurde
	if requestContextIsClosedAndDestroyed(o) {
		panic("destroyed object")
	}

	// Kernel Log
	o.kernel.LogPrint(fmt.Sprintf("RPC(%s)", o._rprequest.ProcessLog.GetID()), "request closed")
}

// Wird ausgeführt wenn die Funktion zuende aufgerufen wurde
func (o *SharedFunctionRequestContext) functionCallFinal() error {
	// Es wird geprüft ob das Objekt zerstört wurde
	if requestContextIsClosedAndDestroyed(o) {
		panic("destroyed object")
	}

	// Es wird geprüft ob eine Antwort gesendet wurde
	if !o.wasResponsed() {
		// Es wird geprüft ob ein Timeout angegeben wurde, wenn ja wird dieser gestartet
		if o.hasTimeout() {
			// Der Timeout Timer wird gestartet
			if err := o.startTimeoutTimer(); err != nil {
				return utils.TimeoutFunctionCallError()
			}
		}
	}

	// Kernel Log
	o.kernel.LogPrint(fmt.Sprintf("RPC(%s)", o._rprequest.ProcessLog.GetID()), "function call finalized")

	// Es ist kein Fehler aufgetreten
	return nil
}

// Wird ausgeführt wenn ein Throw durch die Funktion ausgelöst wird
func (o *SharedFunctionRequestContext) functionCallException(msg string) error {
	// Es wird geprüft ob das Objekt zerstört wurde
	if requestContextIsClosedAndDestroyed(o) {
		panic("destroyed object")
	}

	// Die Antwort wird zurückgesendet
	if err := writeRequestReturnResponse(o, &types.FunctionCallState{Error: msg, State: "exception"}); err != nil {
		return utils.MakeHttpRequestIsClosedBeforeException()
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

// Proxy Shielded, Set Timeout funktion
func (o *SharedFunctionRequestContext) proxyShield_SetTimeout(info *v8.FunctionCallbackInfo) *v8.Value {
	fmt.Println("D: TIMER")
	return v8.Undefined(info.Context().Isolate())
}

// Proxy Shielded, Set Interval funktion
func (o *SharedFunctionRequestContext) proxyShield_SetInterval(info *v8.FunctionCallbackInfo) *v8.Value {
	fmt.Println("D: TIMER")
	return v8.Undefined(info.Context().Isolate())
}

// Proxy Shielded, Clear Timeout funktion
func (o *SharedFunctionRequestContext) proxyShield_ClearTimeout(info *v8.FunctionCallbackInfo) *v8.Value {
	fmt.Println("D: TIMER")
	return v8.Undefined(info.Context().Isolate())
}

// Proxy Shielded, Clear Interval funktion
func (o *SharedFunctionRequestContext) proxyShield_ClearInterval(info *v8.FunctionCallbackInfo) *v8.Value {
	fmt.Println("D: TIMER")
	return v8.Undefined(info.Context().Isolate())
}

// Signalisiert dass ein neuer Promises erzeugt wurde und gibt die Entsprechenden Funktionen zurück
func (o *SharedFunctionRequestContext) proxyShield_NewPromise(info *v8.FunctionCallbackInfo) *v8.Value {
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
func (o *SharedFunctionRequestContext) proxyShield_ConsoleLog(info *v8.FunctionCallbackInfo) *v8.Value {
	// Es werden alle Stringwerte Extrahiert
	extracted, err := utils.ConvertV8ValuesToString(info.Context(), info.Args())
	if err != nil {
		utils.V8ContextThrow(info.Context(), err.Error())
	}

	// Es wird ein String aus der Ausgabe erzeugt
	outputStr := strings.Join(extracted, " ")

	// Die Ausgabe wird an den Console Cache übergeben
	o.kernel.Console().Log(fmt.Sprintf("RPC(%s):-$ %s", strings.ToUpper(o._rprequest.ProcessLog.GetID()), outputStr))

	// Rückgabe ohne Fehler
	return v8.Undefined(info.Context().Isolate())
}

// Proxy Shielded, Console Log Funktion
func (o *SharedFunctionRequestContext) proxyShield_ErrorLog(info *v8.FunctionCallbackInfo) *v8.Value {
	// Es werden alle Stringwerte Extrahiert
	extracted, err := utils.ConvertV8ValuesToString(info.Context(), info.Args())
	if err != nil {
		utils.V8ContextThrow(info.Context(), err.Error())
		return nil
	}

	// Es wird ein String aus der Ausgabe erzeugt
	outputStr := strings.Join(extracted, " ")

	// Die Ausgabe wird an den Console Cache übergeben
	o.kernel.Console().ErrorLog(fmt.Sprintf("RPC(%s):-$ %s", strings.ToUpper(o._rprequest.ProcessLog.GetID()), outputStr))

	// Rückgabe ohne Fehler
	return v8.Undefined(info.Context().Isolate())
}

// Gibt an ob eine Antwort verfügbar ist
func (o *SharedFunctionRequestContext) wasResponsed() bool {
	// Der Mutex wird verwendet
	return o._wasResponded
}

// Startet den Timer, welcher den Vorgang nach erreichen des Timeouts, abbricht
func (o *SharedFunctionRequestContext) startTimeoutTimer() error {
	return nil
}

// Gibt an ob der Request eine Timeout angabe hat
func (po *SharedFunctionRequestContext) hasTimeout() bool {
	return false
}

// Wird für Tests verwendet um den RPC aufruf zu stoppen bis die Verbindung geschlossen wurde
func (o *SharedFunctionRequestContext) testWait(info *v8.FunctionCallbackInfo) *v8.Value {
	// Es wird ermittelt ob ein Argument angegeben wurde
	if len(info.Args()) < 1 {
		utils.V8ContextThrow(info.Context(), "to few arguments")
		return nil
	}
	if len(info.Args()) > 1 {
		utils.V8ContextThrow(info.Context(), "to many arguments")
		return nil
	}

	// Es muss sich um ein uint32 handeln
	if !info.Args()[0].IsUint32() {
		utils.V8ContextThrow(info.Context(), "only integer allowed")
		return nil
	}

	// Erstelle einen Promise Resolver
	resolver, err := v8.NewPromiseResolver(info.Context())
	if err != nil {
		// Es wird Javascript Fehler ausgelöst
		utils.V8ContextThrow(info.Context(), "Error attempting to create a promise 'rpcCall'")

		// Rückgabe
		return v8.Undefined(info.Context().Isolate())
	}

	// Es wird eine Goroutine ausgeführt, diese Wartet X Millisekunden
	go func(res *v8.PromiseResolver, wtime uint32, iso *v8.Isolate) {
		// Es wird 'wtime' * Millisecond gewartet
		if wtime == 0 {
			time.Sleep(1 * time.Millisecond)
		} else {
			time.Sleep(time.Duration(wtime) * time.Millisecond)
		}

		// Es wird Signalisiert, das der Vorgang erfolgreich ausgeführt wurde
		res.Resolve(v8.Null(info.Context().Isolate()))
	}(resolver, info.Args()[0].Uint32(), info.Context().Isolate())

	// Das Promise wird zurückgegeben
	promise := resolver.GetPromise()

	// Rückgabe
	return promise.Value
}

// Erstellt einen neuen SharedFunctionRequestContext
func newSharedFunctionRequestContext(kernel types.KernelInterface, returnDatatype string, rpcRequest *types.RpcRequest) (*SharedFunctionRequestContext, error) {
	// Das Rückgabeobjekt wird erstellt
	returnObject := &SharedFunctionRequestContext{
		//resolveChan:     make(chan *types.FunctionCallState),
		responseChan:    saftychan.NewFunctionCallStateChan(),
		_returnDataType: returnDatatype,
		_wasResponded:   false,
		_rprequest:      rpcRequest,
		kernel:          kernel,
	}

	// Das Objekt wird zurückgegeben
	return returnObject, nil
}

// Gibt an ob das Objekt zerstört wurde
func requestContextIsClosedAndDestroyed(o *SharedFunctionRequestContext) bool {
	return o._destroyed
}

// Sendet die Antwort zurück und setzt den Vorgang auf erfolgreich
func writeRequestReturnResponse(o *SharedFunctionRequestContext, returnv *types.FunctionCallState) error {
	// Diese Funktion wird aufgerufen, sobald die Antwort Übermittelt wurde
	resolveTransmittedData := func() {
		o.clearAndDestroy()
	}

	// Diese Funktion wird aufgerufen, wenn das übermitteln der Daten fehlgeschlagen ist
	rejectTransmittedData := func() {
	}

	// Das Rückgabe Objekt wird erstellt
	returnObject := &types.FunctionCallReturn{
		FunctionCallState: returnv,
		Resolve:           resolveTransmittedData,
		Reject:            rejectTransmittedData,
	}

	// Die Antwort wird an den Eigentlichen Request zurückgegeben
	if err := o._rprequest.Resolve(returnObject); err != nil {
		switch err := err.(type) {
		case *types.SpecificError:
			return err
		default:
			return fmt.Errorf("writeRequestReturnResponse: " + err.Error())
		}
	}

	// Es ist kein Fehler aufgetreten
	return nil
}
