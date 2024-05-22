package kmodulecrypto

import (
	"fmt"
	"vnh1/types"

	v8 "rogchap.com/v8go"
)

type ModuleCrypto int

func (o *ModuleCrypto) Init(kernel types.KernelInterface, iso *v8.Isolate, context *v8.Context) error {
	// Das Consolen Objekt wird definiert
	console := v8.NewObjectTemplate(iso)

	// Das Console Objekt wird final erzeugt
	consoleObj, err := console.NewInstance(context)
	if err != nil {
		return fmt.Errorf("Kernel->_new_kernel_load_console_module: " + err.Error())
	}

	// Die Consolen Funktionen werden hinzugefügt
	context.Global().Set("crypto", consoleObj)

	// Kein Fehler
	return nil
}

func (o *ModuleCrypto) OnlyForMain() bool {
	return false
}

func (o *ModuleCrypto) GetName() string {
	return "console"
}

func NewCryptoModule() *ModuleCrypto {
	return new(ModuleCrypto)
}
