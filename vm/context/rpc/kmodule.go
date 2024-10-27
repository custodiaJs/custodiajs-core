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

package rpcjsprocessor

import (
	"fmt"

	"github.com/CustodiaJS/custodiajs-core/global/types"

	v8 "rogchap.com/v8go"
)

// Initalisiert einen das Kernel Modul
func (o *RPCModule) Init(kernel types.KernelInterface, iso *v8.Isolate, context *v8.Context) error {
	// Es wird versucht ein Global Register Eintrag zu erzeugen
	slfmap := make(map[string]types.SharedFunctionInterface)
	if err := kernel.GloablRegisterWrite("rpc", slfmap); err != nil {
		return fmt.Errorf("")
	}
	spfmap := make(map[string]types.SharedFunctionInterface)
	if err := kernel.GloablRegisterWrite("rpc_public", spfmap); err != nil {
		return fmt.Errorf("")
	}

	// Die RPC (Remote Function Call) funktionen werden bereitgestellt
	rpc := v8.NewObjectTemplate(iso)
	rpc.Set("IsShare", o.rpcIsShar(kernel, iso), v8.ReadOnly)
	rpc.Set("Call", o.rpcCall(kernel, iso), v8.ReadOnly)
	rpc.Set("SharePublic", o.rpcNewSharePublic(kernel, context), v8.ReadOnly)
	rpc.Set("ShareLocal", o.rpcNewShareLocal(kernel, context), v8.ReadOnly)
	rpc.Set("GetDetails", o.rpcGetDetails(kernel, iso), v8.ReadOnly)

	// Das RPC Objekt wird final erzeugt
	rpcObj, err := rpc.NewInstance(context)
	if err != nil {
		return fmt.Errorf("Kernel->_new_kernel_load_rpc_module: " + err.Error())
	}

	// Das Objekt wird als Import Registriert
	if err := kernel.AddImportModule("rpc", rpcObj.Value); err != nil {
		return fmt.Errorf("ModuleHttp->Init: " + err.Error())
	}

	// Kein Fehler
	return nil
}

// Gibt den Namen des Kernel Modules zurück
func (o *RPCModule) GetName() string {
	return "rpc"
}

// Gibt an, ob das Modul nur für den Main Context ist
func (o *RPCModule) OnlyForMain() bool {
	return true
}

// Erstellt ein neues RPC Module
func NewRPCModule() *RPCModule {
	return new(RPCModule)
}
