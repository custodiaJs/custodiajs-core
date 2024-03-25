package jsvm

import (
	"fmt"

	"github.com/dop251/goja"
)

type SharedFunctionInterface interface {
}

type SharedLocalFunction struct {
	callFunction goja.Callable
}

type SharedPublicFunction struct {
	callFunction goja.Callable
}

func (o *JsVM) functionIsSharing(functionName string) bool {
	// Es wird geprüft ob die Funktion bereits Registriert wurde
	_, found := o.sharedFunctions[functionName]

	// Das Ergebniss wird zurückgegeben
	return found
}

func (o *JsVM) sharLocalFunction(funcName string, totalParms []string, function goja.Callable) error {
	// Es wird geprüft ob diese Funktion bereits registriert wurde
	if _, found := o.sharedFunctions[funcName]; found {
		return fmt.Errorf("function always registrated")
	}

	// Die Funktion wird zwischengespeichert
	o.sharedFunctions[funcName] = &SharedLocalFunction{
		callFunction: function,
	}

	// Die Funktion wird im Core registriert
	fmt.Println("VM:SHARE_LOCAL_FUNCTION:", funcName, totalParms)
	o.coreService.RegisterSharedLocalFunction(o, funcName, totalParms, function)

	// Der Vorgang wurde ohne Fehler durchgeführt
	return nil
}

func (o *JsVM) sharePublicFunction(funcName string, totalParms []string, function goja.Callable) error {
	// Es wird geprüft ob diese Funktion bereits registriert wurde
	if _, found := o.sharedFunctions[funcName]; found {
		return fmt.Errorf("function always registrated")
	}

	// Die Funktion wird zwischengespeichert
	o.sharedFunctions[funcName] = &SharedPublicFunction{
		callFunction: function,
	}

	// Der Vorgang wurde ohne Fehler durchgeführt
	return nil
}
