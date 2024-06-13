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
	"vnh1/eventloop"
	"vnh1/types"
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
	request := newSharedFunctionRequestContext(o.kernel, o.signature.ReturnType, req)

	// Es wird geprüft ob das Req Objekt eine Verbindung besitzt
	if !rpcrequest.ConnectionIsOpen(req) {
		o.kernel.LogPrint("RPC", "Process aborted, connection closed")
		return nil
	}

	// Die Loop Aufgabe wird erzeugt
	kernelLoopOperation := eventloop.NewKernelEventLoopFunctionOperation(func(_ *v8.Context, lopr types.KernelEventLoopContextInterface) {
		// Es wird geprüft ob die Verbindung getrennt wurde
		if !rpcrequest.ConnectionIsOpen(req) {
			o.kernel.LogPrint("RPC", "Process aborted, connection closed")
			return
		}

		// Die Initalisierung des Funktionsaufrufes wird vorbereitet
		err := functionCallInEventloopInit(o, request, req)
		if err != nil {
			lopr.SetError(err)
			return
		}

		// Es wird Signalisiert dass der Vorgang erfolgreich war
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
