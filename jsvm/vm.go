package jsvm

import (
	"fmt"

	"github.com/dop251/goja"
)

type JsVMConfig struct {
	EnableWebsockets bool
}

type JsVM struct {
	config       *JsVMConfig
	gojaVM       *goja.Runtime
	exports      *goja.Object
	scriptLoaded bool
}

func (o *JsVM) prepareVM() error {
	vnh1Obj := o.gojaVM.NewObject()
	vnh1Obj.Set("com", o.vnhCOMAPIEP)
	o.gojaVM.Set("vnh1", vnh1Obj)
	o.gojaVM.Set("exports", o.exports)
	return nil
}

func (o *JsVM) RunScript(script string) error {
	if o.scriptLoaded {
		return fmt.Errorf("LoadScript: always script loaded")
	}
	o.scriptLoaded = true
	_, err := o.gojaVM.RunString(script)
	if err != nil {
		panic(err)
	}
	return nil
}

func NewVM(config *JsVMConfig) (*JsVM, error) {
	// Die GoJA VM wird erstellt
	gojaVM := goja.New()

	// Das Basisobjekt wird erzeugt
	var vmObject *JsVM
	if config == nil {
		vmObject = &JsVM{config: &defaultConfig, gojaVM: gojaVM, scriptLoaded: false, exports: gojaVM.NewObject()}
	} else {
		vmObject = &JsVM{config: config, gojaVM: gojaVM, scriptLoaded: false, exports: gojaVM.NewObject()}
	}

	// Die Funktionen werden hinzugefügt
	if err := vmObject.prepareVM(); err != nil {
		return nil, fmt.Errorf("NewVM: " + err.Error())
	}

	// Das VM Objekt wird zurückgegeben
	return vmObject, nil
}
