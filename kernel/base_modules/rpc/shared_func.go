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
	"github.com/CustodiaJS/custodiajs-core/kernel/eventloop"
	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/utils"

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
func (o *SharedFunction) EnterFunctionCall(req *types.RpcRequest) *types.SpecificError {
	// Es wird geprüft ob die Aktuelle SharedFunction "o" NULL ist
	if o == nil {
		// Es wird ein 'Null RPC Call' Fehler zurückgegeben
		return utils.RPCFunctionCallNullSharedFunctionObject()
	}

	// Es wird geprüft ob der RPC Request "req" NULL ist
	if req == nil {
		// Es wird ein 'RPC Function is Null' Fehler zurückgegeben
		return utils.RPCFunctionCallNullRequest()
	}

	// Es wird geprüft ob das Req Objekt eine Verbindung besitzt
	if !req.Context.IsConnected() {
		// Es wird ein 'Connection closed' Fehler zurückgegeben
		return utils.MakeConnectionIsClosedError("SharedFunction->EnterFunctionCall")
	}

	// Es wird geprüft ob die Angeforderte Anzahl an Parametern vorhanden ist
	if len(req.Parms) != len(o.signature.Params) {
		// Es wird ein 'Parameters size invalid' Fehler zurückgegeben
		return utils.MakeRPCFunctionCallParametersNumberUnequal(uint(len(o.signature.Params)), uint(len(req.Parms)))
	}

	// Es wird ein neues Request Objekt
	request, err := newSharedFunctionRequestContext(o.kernel, o.signature.ReturnType, req)
	if err != nil {
		// Der Name der Aktuellen Funktion wird hinzugefügt
		err.AddCallerFunctionToHistory("SharedFunction->EnterFunctionCall")

		// Der Fehler wird zurückgegeben
		return err
	}

	// Es wird geprüft ob das Req Objekt eine Verbindung besitzt
	if !req.Context.IsConnected() {
		// Es wird ein 'Connection is closed' fehler zurückgegeben
		return utils.MakeConnectionIsClosedError("SharedFunction->EnterFunctionCall")
	}

	// Diese Funktion wird als Event ausgeführt
	event := func(_ *v8.Context, lopr types.KernelEventLoopContextInterface) {
		// Es wird geprüft ob die Verbindung getrennt wurde
		if !req.Context.IsConnected() {
			//o.kernel.LogPrint("RPC", "Process aborted, connection closed")
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
	}

	// Die Loop Aufgabe wird erzeugt
	kernelLoopOperation := eventloop.NewKernelEventLoopFunctionOperation(event)

	// Die Funktion wird an den Eventloop des Kernels übergeben
	if err := o.kernel.AddToEventLoop(kernelLoopOperation); err != nil {
		// Der Name der Aktuellen Funktion wird hinzugefügt
		err.AddCallerFunctionToHistory("SharedFunction->EnterFunctionCall")

		// Der Fehler wird zurückgegeben
		return err
	}

	// Der Vorgang wird in einer neuen Goroutine durchgeführt
	go func() {
		// Es wird darauf gewartet das die Loop Operation erfolgreich abgeschlossen wurde
		_, err := kernelLoopOperation.WaitOfResponse()
		if err != nil {
			panic(err)
		}
	}()

	// Das Ergebniss wird zurückgegeben
	return nil
}

func (o *SharedFunction) GetScriptVM() types.VmInterface {
	return o.kernel.AsCoreVM()
}
