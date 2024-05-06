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
	"vnh1/types"

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
	return o.parmTypes
}

// Gibt den Rückgabedatentypen zurück
func (o *SharedFunction) GetReturnDatatype() string {
	// Es wird geprüft ob die Aktuelle SharedFunction "o" NULL ist
	if o == nil {
		return ""
	}

	// Rückgabe
	return o.returnType
}

// Gibt den Datentypen zurück
func (o *SharedFunction) GetReturnDType() string {
	// Es wird geprüft ob die Aktuelle SharedFunction "o" NULL ist
	if o == nil {
		return ""
	}

	// Rückgabe
	return o.returnType
}

// Fügt ein Event hinzu, dieses Event wird ausgeführt sobald eine neuer Funktionsaufruf angenommen wurde
func (o *SharedFunction) AddOnRequestFunction(funcv8 *v8.Function) error {
	// Der Mutex wird verwendet
	o.mutex.Lock()
	defer o.mutex.Unlock()

	// Die Funktion wird abgespeichert
	o.eventOnRequest = append(o.eventOnRequest, funcv8)

	// Log
	o.kernel.LogPrint("rpc", "New 'OnRequest' event added for function %s", o.name)

	// Rückgabe ohne Fehler
	return nil
}

// Fügt ein Event hinzu, dieses Event wird ausgeführt sobald ein neuer Funktionsaufruf Fehlschlägt, bevor er genau zugeordnet werden kann
func (o *SharedFunction) AddOnRequestFailFunction(funcv8 *v8.Function) error {
	// Der Mutex wird verwendet
	o.mutex.Lock()
	defer o.mutex.Unlock()

	// Die Funktion wird abgespeichert
	o.eventOnRequestFail = append(o.eventOnRequest, funcv8)

	// Log
	o.kernel.LogPrint("rpc", "New 'OnRequestFail' event added for function %s", o.name)

	// Rückgabe ohne Fehler
	return nil
}

// Wird verwendet um die Funktion innerhalb des Kernels aufzurufen
func (o *SharedFunction) _callInKernelEventLoop(ctx *v8.Context, request *SharedFunctionRequest, req *types.RpcRequest) error {
	// Die Parameter werden umgewandelt
	convertedValues, err := convertRequestParametersToV8Parameters(ctx.Isolate(), o.parmTypes, req.Parms)
	if err != nil {
		return fmt.Errorf("EnterFunctionCall: " + err.Error())
	}

	// Das Request Objekt wird erstellt
	requestObj, err := makeSharedFunctionObject(ctx, request)
	if err != nil {
		return fmt.Errorf("EnterFunctionCall: " + err.Error())
	}

	// Die Finalen Argumente werden erstellt
	finalArguments := make([]v8.Valuer, 0)
	finalArguments = append(finalArguments, requestObj)
	finalArguments = append(finalArguments, convertedValues...)

	// Die Funktion wird aufgerufen
	promisse, err := o.v8Function.Call(v8.Undefined(ctx.Isolate()), finalArguments...)
	if err != nil {
		panic(err)
	}

	// Es wird geprüft ob es sich um eine Promise handelt
	if !promisse.IsPromise() {
		return fmt.Errorf("is invalid code")
	}

	// Die Rückgabe wird in ein Promise umgewandelt
	goPromisse, err := promisse.AsPromise()
	if err != nil {
		panic(err)
	}

	// Wird ausgeführt wenn die Funktion zuende ausgeführt wurde
	thenCatch := func(_ *v8.FunctionCallbackInfo) *v8.Value {
		// Es wird Siganlisiert dass die Funktion erfolgreich zuende ausgeführt wurde
		request.functionIsDoneSignal()

		// Es wird nichts zurückgegeben
		return nil
	}

	// Wird ausgeführt wenn die Funktion einen Fehler auslöst
	throwCatch := func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Sollte kein Argument vorhanden sein, ein NULL übergeben
		if len(info.Args()) < 1 {
			request.functionHasThrowSigal("")
		}

		// Der Fehler wird zurückgegeben
		request.functionHasThrowSigal(info.Args()[0].String())

		// Rückgabe
		return nil
	}

	// Wird ausgeführt wenn der Vorgang erfolgreich war
	goPromisse.Then(thenCatch, throwCatch)

	// Rückgabe
	return nil
}

// Ruft die Geteilte Funktion auf
func (o *SharedFunction) EnterFunctionCall(req *types.RpcRequest) (*types.FunctionCallState, error) {
	// Es wird geprüft ob die Aktuelle SharedFunction "o" NULL ist
	if o == nil {
		return nil, &types.MultiError{ErrorValue: fmt.Errorf("EnterFunctionCall: o is null"), VmErrorValue: fmt.Errorf("internal error")}
	}

	// Es wird geprüft ob der RPC Request "req" NULL ist
	if req == nil {
		return nil, &types.MultiError{ErrorValue: fmt.Errorf("EnterFunctionCall: req is null"), VmErrorValue: fmt.Errorf("request error")}
	}

	// Es wird geprüft ob die Angeforderte Anzahl an Parametern vorhanden ist
	if len(req.Parms) != len(o.parmTypes) {
		return nil, fmt.Errorf("EnterFunctionCall: invalid parameters")
	}

	// Es wird ein neues Request Objekt
	request := newSharedFunctionRequest(o.kernel)

	// Die Funktion wird an den Eventloop des Kernels übergeben
	err := o.kernel.AddFunctionCallToEventLoop(func(ctx *v8.Context) error {
		return o._callInKernelEventLoop(ctx, request, req)
	})

	// Es wird geprüft ob ein Fehler beim ausführen der Funktion aufgetreten ist
	if err != nil {
		return nil, err
	}

	// Die Daten werden aus dem Request ausgelesen
	returnData := <-request.resolveChan

	// Log
	o.kernel.LogPrint("rpc", "Incomming remote function call request for '%s' from '%s' has return", o.name, "<source>")

	// Es es wird zursicherheit geprüft ob die Daten abgerufen werden konnten
	if returnData == nil {
		return &types.FunctionCallState{State: "ok", Return: []*types.FunctionCallReturnData{}}, nil
	}

	// Das Ergebniss wird zurückgegeben
	return returnData, nil
}

// Wandelt die RPC Argumente in V8 Argumente für den Aktuellen Context um
func convertRequestParametersToV8Parameters(iso *v8.Isolate, parmTypes []string, reqparms []*types.FunctionParameterCapsle) ([]v8.Valuer, error) {
	// Es wird versucht die Paraemter in den Richtigen v8 Datentypen umzuwandeln
	convertedValues := make([]v8.Valuer, 0)
	for hight, item := range reqparms {
		// Es wird geprüft ob der Datentyp gewünscht ist
		if item.CType != parmTypes[hight] {
			return nil, fmt.Errorf("EnterFunctionCall: not same parameter")
		}

		// Es wird geprüft ob es sich um einen Zulässigen Datentypen handelt
		switch item.CType {
		case "boolean":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		case "number":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		case "string":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		case "array":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		case "object":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		case "bytes":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		default:
			return nil, fmt.Errorf("EnterFunctionCall: unsuported datatype")
		}
	}

	// Rückgabe ohne Fehler
	return convertedValues, nil
}

// Die Funktion wird erstellt
func makeSharedFunctionObject(context *v8.Context, request *SharedFunctionRequest) (*v8.Object, error) {
	// Das Requestobjekt wird ersellt
	obj := v8.NewObjectTemplate(context.Isolate())

	// Die Resolve Funktion wird festgelegt
	obj.Set("SendResponse", v8.NewFunctionTemplate(context.Isolate(), request.SendResponse))

	// Die Senderror Funktion wird festgelegt
	obj.Set("SendError", v8.NewFunctionTemplate(context.Isolate(), request.SendError))

	// Die Reject Funktion wird festgelegt
	obj.Set("Reject", v8.NewFunctionTemplate(context.Isolate(), request.Reject))

	// Das Finale Objekt wird erstellt
	fobj, err := obj.NewInstance(context)
	if err != nil {
		return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
	}

	// Rückgabe ohne Fehler
	return fobj, nil
}
