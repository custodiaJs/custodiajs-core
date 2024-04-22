package kmodulehttp

import (
	"fmt"
	"vnh1/types"

	v8 "rogchap.com/v8go"
)

type ModuleHttp struct {
	context *v8.Context
}

func (o *ModuleHttp) Init(kernel types.KernelInterface, iso *v8.Isolate, context *v8.Context) error {
	// Das Consolen Objekt wird definiert
	console := v8.NewObjectTemplate(iso)

	// Der Kontext wird abgespeichert
	o.context = v8.NewContext(iso)

	// Das Console Objekt wird final erzeugt
	consoleObj, err := console.NewInstance(o.context)
	if err != nil {
		return fmt.Errorf("Kernel->_new_kernel_load_console_module: " + err.Error())
	}

	// Das Objekt wird als Import Registriert
	if err := kernel.AddImportModule("http", consoleObj.Value); err != nil {
		return fmt.Errorf("ModuleHttp->Init: " + err.Error())
	}

	// Kein Fehler
	return nil
}

func (o *ModuleHttp) GetName() string {
	return "http"
}

func (o *ModuleHttp) OnlyForMain() bool {
	return false
}

func NewHttpModule() *ModuleHttp {
	return new(ModuleHttp)
}
