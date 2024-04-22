package kmodulerpc

import (
	"fmt"
	"vnh1/types"

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
	rpc.Set("IsShare", o.rpcIsShar(kernel, iso, context), v8.ReadOnly)
	rpc.Set("Call", o.rpcCall(kernel, iso, context), v8.ReadOnly)
	rpc.Set("SharePublic", o.rpcNewSharePublic(kernel, iso, context), v8.ReadOnly)
	rpc.Set("ShareLocal", o.rpcNewShareLocal(kernel, iso, context), v8.ReadOnly)
	rpc.Set("GetDetails", o.rpcGetDetails(kernel, iso, context), v8.ReadOnly)

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
