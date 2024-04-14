package kmoduledatabase

import (
	"fmt"
	"vnh1/types"

	v8 "rogchap.com/v8go"
)

type ModuleDatabase int

func (o *ModuleDatabase) Init(kernel types.KernelInterface) error {
	// Das Consolen Objekt wird definiert
	console := v8.NewObjectTemplate(kernel.Isolate())

	// Das Console Objekt wird final erzeugt
	consoleObj, err := console.NewInstance(kernel.ContextV8())
	if err != nil {
		return fmt.Errorf("Kernel->_new_kernel_load_console_module: " + err.Error())
	}

	// Die Consolen Funktionen werden hinzugefügt
	kernel.Global().Set("database", consoleObj)

	// Kein Fehler
	return nil
}

func (o *ModuleDatabase) GetName() string {
	return "console"
}

func NewDatabaseModule() *ModuleDatabase {
	return new(ModuleDatabase)
}
