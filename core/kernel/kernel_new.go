package kernel

import (
	"fmt"
	"sync"
	"vnh1/core/consolecache"

	v8 "rogchap.com/v8go"
)

func NewKernel(consoleCache *consolecache.ConsoleOutputCache) (*Kernel, error) {
	// DIe VM wird erezugt
	iso := v8.NewIsolate()

	// Das Kernelobjekt wird erzeugt
	kernelObj := &Kernel{
		sharedLocalFunctions:  make(map[string]*SharedLocalFunction),
		sharedPublicFunctions: make(map[string]*SharedPublicFunction),
		mutex:                 &sync.Mutex{},
		Console:               consoleCache,
		Context:               v8.NewContext(iso),
	}

	// Die Konsolen Funktionen werden geladen
	if err := kernelObj._new_kernel_load_console_module(); err != nil {
		return nil, fmt.Errorf("Kernel->NewKernel: " + err.Error())
	}

	// Die RPC Funktionen werden geladen
	if err := kernelObj._new_kernel_load_rpc_module(); err != nil {
		return nil, fmt.Errorf("Kernel->NewKernel: " + err.Error())
	}

	// Die SQL Funktionen werden geladen
	if err := kernelObj._new_kernel_load_sql_module(); err != nil {
		return nil, fmt.Errorf("Kernel->NewKernel: " + err.Error())
	}

	// Die Netzwerkfunktionen werden geladen
	if err := kernelObj._new_kernel_load_network_module(); err != nil {
		return nil, fmt.Errorf("Kernel->NewKernel: " + err.Error())
	}

	// Die Kryptofunktionen werden geladen
	if err := kernelObj._new_kernel_load_crypto_module(); err != nil {
		return nil, fmt.Errorf("Kernel->NewKernel: " + err.Error())
	}

	// Die Webserver Funktionen werden bereitgestellt
	if err := kernelObj._new_kernel_experimental_load_webserver(); err != nil {
		return nil, fmt.Errorf("Kernel->NewKernel: " + err.Error())
	}

	// Die Ident Funktionen werden bereitgestellt
	if err := kernelObj._new_kernel_load_ident_module(); err != nil {
		return nil, fmt.Errorf("Kernel->NewKernel: " + err.Error())
	}

	// Die Allgemeinen VM Funktionen werden bereitgestellt
	if err := kernelObj._new_kernel_load_vm_base(); err != nil {
		return nil, fmt.Errorf("Kernel->NewKernel: " + err.Error())
	}

	// Die Externen Module werden geladen
	if err := kernelObj._new_kernel_load_extmods(); err != nil {
		return nil, fmt.Errorf("Kernel->NewKernel: " + err.Error())
	}

	// Das Kernelobejkt wird zurückgegeben
	return kernelObj, nil
}
