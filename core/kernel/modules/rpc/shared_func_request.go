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
	"sync"
	"vnh1/types"
	"vnh1/utils"

	v8 "rogchap.com/v8go"
)

// Sendet eine Erfolgreiche Antwort zurück
func (o *SharedFunctionRequest) SendResponse(info *v8.FunctionCallbackInfo) *v8.Value {
	// Es wird geprüft ob das Objekt zerstört wurde
	if o.IsClosedAndDestroyed() {
		panic("destroyed object")
	}

	// Der Mutex wird verwendet
	o.mutex.Lock()
	defer o.mutex.Unlock()

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
	resolves := &types.FunctionCallState{State: "ok", Return: make([]*types.FunctionCallReturnData, 0)}

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

	// Es wird geprüft ob ein Rückgabewert vorhanden ist, wenn nicht wird ein Undefined zurückgegeben
	if len(resolves.Return) == 0 {
		resolves.Return = append(resolves.Return, &types.FunctionCallReturnData{CType: "undefined", Value: nil})
	}

	// Die Antwort wird zurückgesendet
	o.resolveChan <- resolves

	// Es ist kein Fehler aufgetreten
	return nil
}

// Sendet eine Fehlerantwort zurück
func (o *SharedFunctionRequest) SendError(info *v8.FunctionCallbackInfo) *v8.Value {
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

	// Es ist kein Fehler aufgetreten
	return nil
}

// Sendet eine Rejectantwort zurück
func (o *SharedFunctionRequest) Reject(info *v8.FunctionCallbackInfo) *v8.Value {
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

	// Es ist kein Fehler aufgetreten
	return nil
}

// Gibt an ob das Objekt zerstört wurde
func (o *SharedFunctionRequest) IsClosedAndDestroyed() bool {
	// Der Mutex wird verwendet
	o.mutex.Lock()
	defer o.mutex.Unlock()

	// Rückgabe des Destroyed Wertes
	return o._destroyed
}

// Räumt auf und Zerstört das Objekt
func (o *SharedFunctionRequest) ClearAndDestroy() {
	// Es wird geprüft ob das Objekt zerstört wurde
	if o.IsClosedAndDestroyed() {
		panic("destroyed object")
	}

	// Der Mutex wird verwendet
	o.mutex.Lock()
	defer o.mutex.Unlock()
}

// Signalisiert dass die Ausführende Funktion fertigestellt wurde
func (o *SharedFunctionRequest) functionIsDoneSignal() {
	// Es wird geprüft ob das Objekt zerstört wurde
	if o.IsClosedAndDestroyed() {
		panic("destroyed object")
	}

	// Der Mutex wird verwendet
	o.mutex.Lock()
	defer o.mutex.Unlock()
}

// Signalisiert das beim Ausführen der RPC Funktion Throw aufgetreten ist
func (o *SharedFunctionRequest) functionHasThrowSigal(errorvalue string) {
	// Es wird geprüft ob das Objekt zerstört wurde
	if o.IsClosedAndDestroyed() {
		panic("destroyed object")
	}

	// Der Mutex wird verwendet
	o.mutex.Lock()
	defer o.mutex.Unlock()

	// Die Antwort wird zurückgesendet
	o.resolveChan <- &types.FunctionCallState{Error: errorvalue, State: "exception"}
}

// Erstellt einen neuen SharedFunctionRequest
func newSharedFunctionRequest(kernel types.KernelInterface) *SharedFunctionRequest {
	return &SharedFunctionRequest{resolveChan: make(chan *types.FunctionCallState), mutex: &sync.Mutex{}, kernel: kernel}
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
