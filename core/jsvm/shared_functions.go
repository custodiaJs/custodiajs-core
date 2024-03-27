package jsvm

import (
	"fmt"

	"github.com/dop251/goja"
)

type SharedLocalFunction struct {
	gojaVM       *goja.Runtime
	callFunction goja.Callable
	name         string
	parmTypes    []string
}

type SharedPublicFunction struct {
	gojaVM       *goja.Runtime
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

func (o *SharedLocalFunction) EnterFunctionCall(parms ...interface{}) (interface{}, error) {
	// Es wird geprüft ob die Angeforderte Anzahl an Parametern vorhanden ist
	if len(parms) != len(o.parmTypes) {
		return nil, fmt.Errorf("EnterFunctionCall: invalid parameters")
	}

	// Es wird versucht die Paraemter in den Richtigen GoJa Datentypen umzuwandeln
	convertedValues := make([]goja.Value, 0)
	for _, item := range parms {
		// Der Wert wird umgewandelt
		gojaValue := o.gojaVM.ToValue(item)

		// Es wird geprüft ob es sich um einen Zulässigen Datentypen handelt
		switch v := gojaValue.Export().(type) {
		case string:
		case uint64:
		case bool:
		case goja.ArrayBuffer:
		default:
			return nil, fmt.Errorf("EnterFunctionCall: unsupported datatype %T", v)
		}

		// Der Wert wird zwischengespeichert
		convertedValues = append(convertedValues, gojaValue)
	}

	// Die Funktion wird aufgerufen
	result, err := o.callFunction(nil, convertedValues...)
	if err != nil {
		return nil, fmt.Errorf("EnterFunctionCall: " + err.Error())
	}

	fmt.Println("CALL FUNCTION", result)
	return nil, nil
}

func (o *SharedPublicFunction) EnterFunctionCall(parms ...interface{}) (interface{}, error) {
	return nil, nil
}
