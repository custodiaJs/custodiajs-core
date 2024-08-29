package kmoduledatabase

import (
	"fmt"

	"github.com/CustodiaJS/custodiajs-core/global/types"

	v8 "rogchap.com/v8go"
)

type ModuleDatabase int

func (o *ModuleDatabase) Init(kernel types.KernelInterface, iso *v8.Isolate, context *v8.Context) error {
	// Das Consolen Objekt wird definiert
	console := v8.NewObjectTemplate(iso)

	// Das Console Objekt wird final erzeugt
	consoleObj, err := console.NewInstance(context)
	if err != nil {
		return fmt.Errorf("Kernel->_new_kernel_load_console_module: " + err.Error())
	}

	// Die Consolen Funktionen werden hinzugefügt
	context.Global().Set("database", consoleObj)

	// Kein Fehler
	return nil
}

func (o *ModuleDatabase) OnlyForMain() bool {
	return false
}

func (o *ModuleDatabase) GetName() string {
	return "console"
}

func NewDatabaseModule() *ModuleDatabase {
	return new(ModuleDatabase)
}
