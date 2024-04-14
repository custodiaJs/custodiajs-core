package kernel

import (
	"fmt"

	v8 "rogchap.com/v8go"
)

func (o *Kernel) _new_kernel_load_console_module() error {
	// Das Consolen Objekt wird definiert
	console := v8.NewObjectTemplate(o.Isolate())
	console.Set("log", o._kernel_console_log(), v8.ReadOnly)
	console.Set("error", o._kernel_console_error(), v8.ReadOnly)

	// Das Console Objekt wird final erzeugt
	consoleObj, err := console.NewInstance(o.Context)
	if err != nil {
		return fmt.Errorf("Kernel->_new_kernel_load_console_module: " + err.Error())
	}

	// Die Consolen Funktionen werden hinzugefügt
	o.Global().Set("console", consoleObj)

	// Kein Fehler
	return nil
}

func (o *Kernel) _new_kernel_load_rpc_module() error {
	// Die RPC (Remote Function Call) funktionen werden bereitgestellt
	rpc := v8.NewObjectTemplate(o.Isolate())
	rpc.Set("CallLocal", o._kernel_rpc_call_local(), v8.ReadOnly)
	rpc.Set("CallRemote", o._kernel_rpc_call_remote(), v8.ReadOnly)
	rpc.Set("ShareLocal", o._kernel_rpc_shareLocalFunction(), v8.ReadOnly)
	rpc.Set("SharePublic", o._kernel_rpc_sharePublicFunction(), v8.ReadOnly)
	rpc.Set("IsShare", o._kernel_rpc_functionIsSharing(), v8.ReadOnly)

	// Das RPC Objekt wird final erzeugt
	rpcObj, err := rpc.NewInstance(o.Context)
	if err != nil {
		return fmt.Errorf("Kernel->_new_kernel_load_rpc_module: " + err.Error())
	}

	// Das RFC Modul wird hinzugefügt
	o.Global().Set("rpc", rpcObj)

	// Kein Fehler
	return nil
}

func (o *Kernel) _new_kernel_load_sql_module() error {
	// Die SQL funktionen werden bereitgestellt
	sql := v8.NewObjectTemplate(o.Isolate())

	// Das RPC Objekt wird final erzeugt
	sqlObj, err := sql.NewInstance(o.Context)
	if err != nil {
		return fmt.Errorf("Kernel->_new_kernel_load_rpc_module: " + err.Error())
	}

	// Das RFC Modul wird hinzugefügt
	o.Global().Set("sql", sqlObj)

	// Kein Fehler
	return nil
}

func (o *Kernel) _new_kernel_load_network_module() error {
	// Die Net funktionen werden bereitgestellt
	net := v8.NewObjectTemplate(o.Isolate())

	// Das Net Objekt wird final erzeugt
	netObj, err := net.NewInstance(o.Context)
	if err != nil {
		return fmt.Errorf("Kernel->_new_kernel_load_rpc_module: " + err.Error())
	}

	// Das Network Modul wird hinzugefügt
	o.Global().Set("net", netObj)

	// Kein Fehler
	return nil
}

func (o *Kernel) _new_kernel_load_crypto_module() error {
	// Die Crypto funktionen werden bereitgestellt
	crypto := v8.NewObjectTemplate(o.Isolate())

	// Das Crypto Objekt wird final erzeugt
	cryptoObj, err := crypto.NewInstance(o.Context)
	if err != nil {
		return fmt.Errorf("Kernel->_new_kernel_load_rpc_module: " + err.Error())
	}

	// Das Network Modul wird hinzugefügt
	o.Global().Set("crypto", cryptoObj)

	// Kein Fehler
	return nil
}

func (o *Kernel) _new_kernel_experimental_load_webserver() error {
	// Die Webserver funktionen werden bereitgestellt
	experimentalWebserver := v8.NewObjectTemplate(o.Isolate())

	// Das Crypto Objekt wird final erzeugt
	experimentalWebserverObj, err := experimentalWebserver.NewInstance(o.Context)
	if err != nil {
		return fmt.Errorf("Kernel->_new_kernel_load_rpc_module: " + err.Error())
	}

	// Das Webserver Modul wird hinzugefügt
	o.Global().Set("webserver", experimentalWebserverObj)

	// Kein Fehler
	return nil
}

func (o *Kernel) _new_kernel_load_ident_module() error {
	// Die Ident funktionen werden bereitgestellt
	identModule := v8.NewObjectTemplate(o.Isolate())

	// Das Ident Modul Objekt wird final erzeugt
	identModuleObj, err := identModule.NewInstance(o.Context)
	if err != nil {
		return fmt.Errorf("Kernel->_new_kernel_load_rpc_module: " + err.Error())
	}

	// Das Webserver Modul wird hinzugefügt
	o.Global().Set("ident", identModuleObj)

	// Kein Fehler
	return nil
}

func (o *Kernel) _new_kernel_load_vm_base() error {
	// Die Ident funktionen werden bereitgestellt
	identModule := v8.NewObjectTemplate(o.Isolate())

	// Das Ident Modul Objekt wird final erzeugt
	identModuleObj, err := identModule.NewInstance(o.Context)
	if err != nil {
		return fmt.Errorf("Kernel->_new_kernel_load_rpc_module: " + err.Error())
	}

	// Das Webserver Modul wird hinzugefügt
	o.Global().Set("ident", identModuleObj)

	// Kein Fehler
	return nil
}
