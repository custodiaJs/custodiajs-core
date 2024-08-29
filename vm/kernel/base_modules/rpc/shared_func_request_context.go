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

	"github.com/CustodiaJS/custodiajs-core/global/types"
	"github.com/CustodiaJS/custodiajs-core/global/utils"
	"github.com/CustodiaJS/custodiajs-core/global/utils/grsbool"

	v8 "rogchap.com/v8go"
)

// Sendet eine Erfolgreiche Antwort zurück
func (o *SharedFunctionRequestContext) resolveFunctionCallbackV8(info *v8.FunctionCallbackInfo) *v8.Value {
	// Es wird geprüft ob die v8.FunctionCallbackInfo "info" NULL ist
	// In diesem fall muss ein Panic ausgelöst werden, da es keine Logik gibt, wann diese Eintreten kann
	if info == nil {
		// Der Fehler wird erzeugt
		err := utils.MakeV8FunctionCallbackInfoIsNullError("SharedFunctionRequestContext->resolveFunctionCallbackV8")

		// Es wird ein Panic ausgelöst
		panic(err.GoProcessError.Error())
	}

	// Es wird geprüft ob SharedFunctionRequestContext "o" NULL ist
	if !validateSharedFunctionRequestContextObject(o) {
		// Der Fehler wird erzeugt
		err := utils.MakeSharedFunctionRequestContextObjectError("SharedFunctionRequestContext->resolveFunctionCallbackV8")

		// Es wird ein V8 Throw erzeugt und ausgeführt
		utils.V8ContextThrow(info.Context(), err.LocalJSVMError.Error())

		// Es wird ein Nill zurückgegeben
		return nil
	}

	// Es wird geprüft ob der Context geschlossen wurde
	if requestContextIsClosedAndDestroyed(o) {
		// Der Fehler wird erzeugt
		err := utils.MakeRPCRequestContextIsClosedAndDestroyed("SharedFunctionRequestContext->resolveFunctionCallbackV8")

		// Es wird ein V8 Throw erzeugt und ausgeführt
		utils.V8ContextThrow(info.Context(), err.LocalJSVMError.Error())

		// Es wird ein Nill zurückgegeben
		return nil
	}

	// Es wird versucht zu Signalisieren dass ein Antwortvorgang durchgeführt wird,
	// durch diese Funktion wird verhindert dass ein anderer Vorgang eine Antwort in diesem Context schreibt
	if err := trySignalWriteOperation(o); err != nil {
		// Es wird ein V8 Throw erzeugt und ausgeführt
		utils.V8ContextThrow(info.Context(), err.LocalJSVMError.Error())

		// Es wird ein Nill zurückgegeben
		return nil
	}

	// Es wird geprüft ob der Vorgang bereits beantwortet wurde,
	// wenn ja wird ein Fehler zurückgegeben dass der Vorgang bereits beantwortet wurde
	if o.wasResponsed() {
		// Der Fehler wird erzeugt
		err := utils.MakeRPCRequestAlwaysResponsedError("writeRequestReturnResponse")

		// Es wird ein V8 Throw erzeugt und ausgeführt
		utils.V8ContextThrow(info.Context(), err.LocalJSVMError.Error())

		// Es wird ein Nill zurückgegeben
		return nil
	}

	// Es wird geprüft ob es sich um eine Remote Verbindung handelt

	// Die Argumente werden umgewandelt
	convertedArguments, err := utils.ConvertV8DataToGoData(info.Args())
	if err != nil {
		// Der Fehler wird erzeugt
		err := utils.MakeV8ToGoConvertingError("writeRequestReturnResponse")

		// Es wird ein V8 Throw erzeugt und ausgeführt
		utils.V8ContextThrow(info.Context(), err.LocalJSVMError.Error())

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

	// Die Antwort wird geschrieben
	if err := writeRequestReturnResponse(o, resolves); err != nil {
		utils.V8ContextThrow(info.Context(), err.LocalJSVMError.Error())
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

	// Es wird geprüft ob SharedFunctionRequestContext "o" NULL ist
	if !validateSharedFunctionRequestContextObject(o) {
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
	// Es wird ein neues V8 Objekt erzeugt
	v8Object := v8.NewObjectTemplate(info.Context().Isolate())

	// Die Proxy Resolve und Reject funktion wird ausgeführt
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
	obj, err := v8Object.NewInstance(info.Context())
	if err != nil {
		// Der Fehler wird erzeugt
		err := utils.MakeNewRPCSharedFunctionNewV8ObjectInstanceError("SharedFunctionRequestContext->proxyShield_NewPromise", err)

		// Der Fehler wird als V8 Throw ausgeführt
		utils.V8ContextThrow(info.Context(), err.LocalJSVMError.Error())

		// Es wird ein Null zurückgegeben
		return nil
	}

	// Das Objekt wird zurückgegeben
	return obj.Value
}

// Proxy Shielded, Console Log Funktion
func (o *SharedFunctionRequestContext) proxyShield_ConsoleLog(info *v8.FunctionCallbackInfo) *v8.Value {
	// Es werden alle Stringwerte Extrahiert
	extracted, err := utils.ConvertV8ValuesToString(info.Context(), info.Args())
	if err != nil {
		// Der Fehler wird erzeugt
		err := utils.MakeV8ConvertValueToStringError("SharedFunctionRequestContext->proxyShield_ConsoleLog", err)

		// Der Fehler wird als V8 Throw ausgeführt
		utils.V8ContextThrow(info.Context(), err.LocalJSVMError.Error())

		// Es wird ein Null zurückgegeben
		return nil
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
		// Der Fehler wird erzeugt
		err := utils.MakeV8ConvertValueToStringError("SharedFunctionRequestContext->proxyShield_ErrorLog", err)

		// Der Fehler wird als V8 Throw ausgeführt
		utils.V8ContextThrow(info.Context(), err.LocalJSVMError.Error())

		// Es wird ein Null zurückgegeben
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
	// Der Boolwert wird zurückgegeben
	return o._wasResponded.Bool()
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
	// Es wird geprüft ob das Info Objekt nicht null ist
	if info == nil {
		panic("SharedFunctionRequestContext->testWait: info is null")
	}

	// Es wird geprüft ob SharedFunctionRequestContext "o" NULL ist
	if !validateSharedFunctionRequestContextObject(o) {
		// Der Fehler wird erzeugt
		err := utils.MakeSharedFunctionRequestContextObjectError("SharedFunctionRequestContext->testWait")

		// Der Fehler wird als V8 Throw ausgeführt
		utils.V8ContextThrow(info.Context(), err.LocalJSVMError.Error())

		// Es wird ein Null zurückgegeben
		return nil
	}

	// Es wird ermittelt ob ein Argument angegeben wurde
	if len(info.Args()) < 1 {
		// Der Fehler wird erzeugt
		err := utils.MakeV8MissingParameters("SharedFunctionRequestContext->testWait", 1, len(info.Args()))

		// Der Fehler wird als V8 Throw ausgeführt
		utils.V8ContextThrow(info.Context(), err.LocalJSVMError.Error())

		// Es wird ein Null zurückgegeben
		return nil
	}
	if len(info.Args()) > 1 {
		// Der Fehler wird erzeugt
		err := utils.MakeV8MissingParameters("SharedFunctionRequestContext->testWait", 1, len(info.Args()))

		// Der Fehler wird als V8 Throw ausgeführt
		utils.V8ContextThrow(info.Context(), err.LocalJSVMError.Error())

		// Es wird ein Null zurückgegeben
		return nil
	}

	// Es muss sich um ein uint32 handeln
	if !info.Args()[0].IsUint32() {
		// Der Fehler wird erzeugt
		err := utils.MakeV8InvalidParameterDatatype("SharedFunctionRequestContext->testWait", 0, "uint32")

		// Der Fehler wird als V8 Throw ausgeführt
		utils.V8ContextThrow(info.Context(), err.LocalJSVMError.Error())

		// Es wird ein Null zurückgegeben
		return nil
	}

	// Erstelle einen Promise Resolver
	resolver, err := v8.NewPromiseResolver(info.Context())
	if err != nil {
		// Der Fehler wird erzeugt
		err := utils.MakeV8PromiseCreatingError("SharedFunctionRequestContext->testWait", err)

		// Der Fehler wird als V8 Throw ausgeführt
		utils.V8ContextThrow(info.Context(), err.LocalJSVMError.Error())

		// Es wird ein Null zurückgegeben
		return nil
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
func newSharedFunctionRequestContext(kernel types.KernelInterface, returnDatatype string, rpcRequest *types.RpcRequest) (*SharedFunctionRequestContext, *types.SpecificError) {
	// Es wird geprüft ob der Kernel nill ist
	if kernel == nil {
		return nil, utils.MakeNewRPCSharedFunctionContextKernelIsNullError("newSharedFunctionRequestContext")
	}

	// Es wird geprüft ob es sich bei dem ReturnDatatype um einen zulässigen Datentypen handelt
	if !utils.ValidateDatatypeString(returnDatatype) {
		return nil, utils.MakeNewRPCSharedFunctionContextReturnDatatypeStringIsInvalidError("newSharedFunctionRequestContext", returnDatatype)
	}

	// Es wird geprüft ob der rpcRequest nill ist
	if rpcRequest == nil {
		return nil, utils.MakeNewRPCSharedFunctionContextRPCRequestIsNullError("newSharedFunctionRequestContext")
	}

	// Sollte es sich um eine Remote Verbindung handeln, wird geprüft ob diese Geschlossen wurde, wenn ja wird der Vorgang abgebrochen
	if !rpcRequest.Context.IsConnected() {
		return nil, utils.MakeConnectionIsClosedError("newSharedFunctionRequestContext")
	}

	// Das Rückgabeobjekt wird erstellt
	returnObject := &SharedFunctionRequestContext{
		_wasResponded:   grsbool.NewGrsbool(false),
		_destroyed:      grsbool.NewGrsbool(false),
		_returnDataType: returnDatatype,
		_rprequest:      rpcRequest,
		kernel:          kernel,
	}

	// Es wird geprüft ob es sich um ein gültiges ShareFunctionRecquestContext Objekt handelt
	if !validateSharedFunctionRequestContextObject(returnObject) {
		return nil, utils.MakeNewRPCSharedFunctionInvalidContextObjectError("newSharedFunctionRequestContext")
	}

	// Das Objekt wird zurückgegeben
	return returnObject, nil
}

// Wird verwendet um zu Signalisieren dass ein Schreibvorgang durchgeführt wird
func trySignalWriteOperation(o *SharedFunctionRequestContext) *types.SpecificError {
	_ = o
	return nil
}

// Gibt an ob das Objekt zerstört wurde
func requestContextIsClosedAndDestroyed(o *SharedFunctionRequestContext) bool {
	// Es wird geprüft ob SharedFunctionRequestContext "o" NULL ist
	if !validateSharedFunctionRequestContextObject(o) {
		panic("SharedFunctionRequestContext(o) is null")
	}

	// Der Wert wird zurückgegeben
	return o._destroyed.Bool()
}

// Sendet die Antwort zurück und setzt den Vorgang auf erfolgreich
func writeRequestReturnResponse(o *SharedFunctionRequestContext, returnv *types.FunctionCallState) *types.SpecificError {
	// Es wird geprüft ob SharedFunctionRequestContext "o" NULL ist
	if !validateSharedFunctionRequestContextObject(o) {
		return utils.MakeSharedFunctionRequestContextObjectError("writeRequestReturnResponse")
	}

	// Es wird geprüft ob der Context geschlossen wurden
	if requestContextIsClosedAndDestroyed(o) {
		return utils.MakeRPCRequestContextIsClosedAndDestroyed("writeRequestReturnResponse")
	}

	// Es wird geprüft ob der FunctionCallState NULL ist
	if returnv == nil {
		return utils.MakeSharedFunctionCallStateError("writeRequestReturnResponse")
	}

	// Sollte es sich um eine Remote Verbindung handeln, wird geprüft ob diese Geschlossen wurde, wenn ja wird der Vorgang abgebrochen
	if !o._rprequest.Context.IsConnected() {
		return utils.MakeConnectionIsClosedError("writeRequestReturnResponse")
	}

	// Es wird geprüft ob der Response bereits beantwortet wurde
	if o.wasResponsed() {
		return utils.MakeRPCRequestAlwaysResponsedError("writeRequestReturnResponse")
	}

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
	if err := o._rprequest.WriteResponse(returnObject); err != nil {
		// Die Funktion wir hinzugefügt
		err := utils.MakeRPCResolvingDataError("writeRequestReturnResponse")

		// Der Fehler wird zurückgegeben
		return err
	}

	// Es ist kein Fehler aufgetreten
	return nil
}
