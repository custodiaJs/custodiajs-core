package kernel

import (
	"fmt"
	"vnh1/extmodules"
	"vnh1/types"

	v8 "rogchap.com/v8go"
)

type ExtModuleLink struct {
	exteModuleLink *extmodules.ExternalModule
	name           string
}

func (o *ExtModuleLink) addGlobalFunc(extModFunc *extmodules.ExternModuleFunction, kernel types.KernelInterface) error {
	// Die Funktion wird erzeugt
	funcTemplate := v8.NewFunctionTemplate(kernel.ContextV8().Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Die Funktion aus der Lib wird aufgerufen
		state, result, err := extModFunc.Call()
		if err != nil {
			kernel.KernelThrow(kernel.ContextV8(), "internal linking error: "+err.Error())
			return nil
		}

		fmt.Println(state, result)

		// Es ist kein Fehler aufgetreten
		return nil
	})

	// Die Funktion wird hinzugefügt
	funcObj := funcTemplate.GetFunction(kernel.ContextV8())
	if err := kernel.Global().Set(extModFunc.GetName(), funcObj); err != nil {
		panic(err)
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

func (o *ExtModuleLink) Init(kernel types.KernelInterface) error {
	// Die Globalen Funktionen werden Exportiert
	for _, item := range o.exteModuleLink.GetGlobalFunctions() {
		// Es wird versucht die Globale Funktion hinzuzufügen
		if err := o.addGlobalFunc(item, kernel); err != nil {
			return fmt.Errorf("ExtModuleLink->Init: " + err.Error())
		}
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

func (o *ExtModuleLink) GetName() string {
	return o.name
}

func LinkWithExternalModule(extmodule *extmodules.ExternalModule) (*ExtModuleLink, error) {
	// Das Modul wird erzeugt
	vat := &ExtModuleLink{name: extmodule.GetName(), exteModuleLink: extmodule}

	// Das Modul wird zurückgegeben
	return vat, nil
}
