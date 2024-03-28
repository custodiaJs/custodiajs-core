package jsvm

import (
	"fmt"
	"reflect"
	"vnh1/types"

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

func (o *SharedLocalFunction) EnterFunctionCall(parms ...*types.FunctionParameterCapsle) (interface{}, error) {
	// Es wird geprüft ob die Angeforderte Anzahl an Parametern vorhanden ist
	if len(parms) != len(o.parmTypes) {
		return nil, fmt.Errorf("EnterFunctionCall: invalid parameters")
	}

	// Es wird versucht die Paraemter in den Richtigen GoJa Datentypen umzuwandeln
	convertedValues := make([]goja.Value, 0)
	for hight, item := range parms {
		// Es wird geprüft ob der Datentyp gewünscht ist
		if item.CapsleType != o.parmTypes[hight] {
			return nil, fmt.Errorf("EnterFunctionCall: not same parameter")
		}

		// Es wird geprüft ob es sich um einen Zulässigen Datentypen handelt
		switch item.CapsleType {
		case "boolean":
			gojaValue := o.gojaVM.ToValue(item.Value)
			if _, ok := gojaValue.Export().(bool); !ok {
				return nil, fmt.Errorf("EnterFunctionCall: invalid boolean data")
			}
			convertedValues = append(convertedValues, gojaValue)
		case "number":
			gojaValue := o.gojaVM.ToValue(item.Value)
			if reflect.TypeOf(gojaValue.Export()).Kind() != reflect.Int &&
				reflect.TypeOf(gojaValue.Export()).Kind() != reflect.Float64 {
				return nil, fmt.Errorf("EnterFunctionCall: invalid number data")
			}
			convertedValues = append(convertedValues, gojaValue)
		case "string":
			gojaValue := o.gojaVM.ToValue(item.Value)
			if _, ok := gojaValue.Export().(string); !ok {
				return nil, fmt.Errorf("EnterFunctionCall: invalid string")
			}
			convertedValues = append(convertedValues, gojaValue)
		case "array":
			gojaValue := o.gojaVM.ToValue(item.Value)
			if _, ok := gojaValue.Export().([]interface{}); !ok {
				return nil, fmt.Errorf("EnterFunctionCall: invalid array")
			}
			convertedValues = append(convertedValues, gojaValue)
		case "object":
			gojaValue := o.gojaVM.ToValue(item.Value)
			_, ok := gojaValue.Export().(map[string]interface{})
			if !ok && reflect.TypeOf(gojaValue.Export()).Kind() != reflect.Struct {
				return nil, fmt.Errorf("EnterFunctionCall: invalid object")
			}
			convertedValues = append(convertedValues, gojaValue)
		case "bytes":
			gojaValue := o.gojaVM.ToValue(item.Value)
			if _, ok := gojaValue.Export().([]byte); !ok {
				return nil, fmt.Errorf("EnterFunctionCall: invalid byte array")
			}
			convertedValues = append(convertedValues, gojaValue)
		default:
			return nil, fmt.Errorf("EnterFunctionCall: unsuported datatype")
		}
	}

	// Die Funktion wird aufgerufen
	result, err := o.callFunction(nil, convertedValues...)
	if err != nil {
		return nil, fmt.Errorf("EnterFunctionCall: " + err.Error())
	}

	// Das Ergebniss wird zurückgegeben
	return result.Export(), nil
}

func (o *SharedPublicFunction) GetName() string {
	return o.name
}

func (o *SharedPublicFunction) GetParmTypes() []string {
	return o.parmTypes
}

func (o *SharedPublicFunction) EnterFunctionCall(parms ...*types.FunctionParameterCapsle) (interface{}, error) {
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

	// Das Ergebniss wird zurückgegeben
	return result.Export(), nil
}
