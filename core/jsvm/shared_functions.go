package jsvm

import (
	"fmt"

	"github.com/dop251/goja"
)

type SharedLocalFunction struct {
	callFunction goja.Callable
	name         string
	parmTypes    []string
}

type SharedPublicFunction struct {
	callFunction goja.Callable
	name         string
	parmTypes    []string
}

func (o *SharedLocalFunction) GetName() string {
	return o.name
}

func (o *SharedLocalFunction) GetParmTypes() []string {
	return o.parmTypes
}

func (o *SharedPublicFunction) GetName() string {
	return o.name
}

func (o *SharedPublicFunction) GetParmTypes() []string {
	return o.parmTypes
}

func (o *JsVM) functionIsSharing(functionName string) bool {
	// Es wird geprüft ob die Funktion bereits Registriert wurde
	_, found := o.sharedLocalFunctions[functionName]

	// Das Ergebniss wird zurückgegeben
	return found
}

func (o *JsVM) shareLocalFunction(funcName string, parmTypes []string, function goja.Callable) error {
	// Es wird geprüft ob diese Funktion bereits registriert wurde
	if _, found := o.sharedLocalFunctions[funcName]; found {
		return fmt.Errorf("function always registrated")
	}

	// Die Funktion wird zwischengespeichert
	o.sharedLocalFunctions[funcName] = &SharedLocalFunction{
		callFunction: function,
		name:         funcName,
		parmTypes:    parmTypes,
	}

	// Die Funktion wird im Core registriert
	fmt.Println("VM:SHARE_LOCAL_FUNCTION:", funcName, parmTypes)

	// Der Vorgang wurde ohne Fehler durchgeführt
	return nil
}

func (o *JsVM) sharePublicFunction(funcName string, parmTypes []string, function goja.Callable) error {
	// Es wird geprüft ob diese Funktion bereits registriert wurde
	if _, found := o.sharedPublicFunctions[funcName]; found {
		return fmt.Errorf("function always registrated")
	}

	// Die Funktion wird zwischengespeichert
	o.sharedPublicFunctions[funcName] = &SharedPublicFunction{
		callFunction: function,
		name:         funcName,
		parmTypes:    parmTypes,
	}

	// Die Funktion wird im Core registriert
	fmt.Println("VM:SHARE_PUBLIC_FUNCTION:", funcName, parmTypes)

	// Der Vorgang wurde ohne Fehler durchgeführt
	return nil
}
