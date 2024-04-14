package kernel

import (
	"fmt"

	v8 "rogchap.com/v8go"
)

func (o *Kernel) _require(value string) (*v8.Value, error) {
	// Der Mutex wird verwendet
	o.mutex.Lock()
	defer o.mutex.Unlock()

	// Es wird geprüft ob es sich um ein bekanntes Module handelt
	importModule, foundIt := o.vmImports[value]
	if !foundIt {
		return nil, fmt.Errorf("unkown import")
	}

	/* Es wird ein neuer Kontext verwendet
	ctx := v8.NewContext(o.Isolate())

	// Die Exports Funktionen werden bereitsgestellt
	_, err := ctx.RunScript("const exports = {};", "require.js")
	if err != nil {
		return nil, fmt.Errorf("Kernel->_require: " + err.Error())
	}

	test := `
	exports.sayHello = {
		test : function() {
			console.log("Hello");
		},
	}
	`
	_, err = ctx.RunScript(test, "require.js")
	if err != nil {
		return nil, fmt.Errorf("Kernel->_require: " + err.Error())
	}

	// Die Schlüssel werden ausgelesen
	keysVal, err := ctx.RunScript("exports", "get_keys.js")
	if err != nil {
		log.Fatalf("Failed to get object keys: %v", err)
	}
	*/

	// Die Werte werden zurückgegeben
	return importModule, nil
}

func (o *Kernel) _setup_require() error {
	// Bereitstellen der 'require' Funktion im globalen Kontext
	requireFunc := v8.NewFunctionTemplate(o.Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Der Import wird abgerufen
		arg := info.Args()[0]
		modulePath := arg.String()

		// Es wird versucht das Module zu laden
		moduleValue, err := o._require(modulePath)
		if err != nil {
			o.KernelThrow(info.Context(), err.Error())
		}

		// Das Module wird zurückgegeben
		return moduleValue
	})

	// Die Funktion wird im Kontext registriert
	reqFunc := requireFunc.GetFunction(o.ContextV8())

	// Die Require Funktion wird registriert
	if err := o.Global().Set("require", reqFunc); err != nil {
		return fmt.Errorf("Kernel->_setup_require: " + err.Error())
	}

	// Es ist kein Fehler aufgetreten
	return nil
}
