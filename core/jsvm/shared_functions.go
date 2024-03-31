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

type FunctionCallReturn struct {
	CType string
	Value interface{}
}

func (o *FunctionCallReturn) GetType() string {
	return o.CType
}

func (o *FunctionCallReturn) GetValue() any {
	return o.Value
}

func (o *SharedLocalFunction) GetName() string {
	return o.name
}

func (o *SharedLocalFunction) GetParmTypes() []string {
	return o.parmTypes
}

func (o *SharedLocalFunction) EnterFunctionCall(parms types.RpcRequestInterface) (types.FunctionCallReturnInterface, error) {
	// Es wird geprüft ob die Angeforderte Anzahl an Parametern vorhanden ist
	if len(parms.GetParms()) != len(o.parmTypes) {
		return nil, fmt.Errorf("EnterFunctionCall: invalid parameters")
	}

	// Es wird versucht die Paraemter in den Richtigen GoJa Datentypen umzuwandeln
	convertedValues := make([]goja.Value, 0)
	for hight, item := range parms.GetParms() {
		// Es wird geprüft ob der Datentyp gewünscht ist
		if item.GetType() != o.parmTypes[hight] {
			return nil, fmt.Errorf("EnterFunctionCall: not same parameter")
		}

		// Es wird geprüft ob es sich um einen Zulässigen Datentypen handelt
		switch item.GetType() {
		case "boolean":
			gojaValue := o.gojaVM.ToValue(item.GetValue())
			if _, ok := gojaValue.Export().(bool); !ok {
				return nil, fmt.Errorf("EnterFunctionCall: invalid boolean data")
			}
			convertedValues = append(convertedValues, gojaValue)
		case "number":
			gojaValue := o.gojaVM.ToValue(item.GetValue())
			if reflect.TypeOf(gojaValue.Export()).Kind() != reflect.Int64 && reflect.TypeOf(gojaValue.Export()).Kind() != reflect.Float64 {
				return nil, fmt.Errorf("EnterFunctionCall: invalid number data")
			}
			convertedValues = append(convertedValues, gojaValue)
		case "string":
			gojaValue := o.gojaVM.ToValue(item.GetValue())
			if _, ok := gojaValue.Export().(string); !ok {
				return nil, fmt.Errorf("EnterFunctionCall: invalid string")
			}
			convertedValues = append(convertedValues, gojaValue)
		case "array":
			gojaValue := o.gojaVM.ToValue(item.GetValue())
			if _, ok := gojaValue.Export().([]interface{}); !ok {
				return nil, fmt.Errorf("EnterFunctionCall: invalid array")
			}
			convertedValues = append(convertedValues, gojaValue)
		case "object":
			gojaValue := o.gojaVM.ToValue(item.GetValue())
			_, ok := gojaValue.Export().(map[string]interface{})
			if !ok && reflect.TypeOf(gojaValue.Export()).Kind() != reflect.Struct {
				return nil, fmt.Errorf("EnterFunctionCall: invalid object")
			}
			convertedValues = append(convertedValues, gojaValue)
		case "bytes":
			gojaValue := o.gojaVM.ToValue(item.GetValue())
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

	// Der Rückgabewert wird ermittelt und geprüft
	var resultReturn types.FunctionCallReturnInterface
	if result == nil {
		resultReturn = &FunctionCallReturn{CType: "null", Value: nil}
	} else if result.ExportType() == goja.Undefined().ExportType() && result.Export() == nil {
		resultReturn = &FunctionCallReturn{CType: "undefined", Value: nil}
	} else {
		switch result.ExportType().Kind() {
		case reflect.Bool:
			resultReturn = &FunctionCallReturn{CType: "boolean", Value: result.ToBoolean()}
		case reflect.Int64:
			resultReturn = &FunctionCallReturn{CType: "number", Value: result.ToInteger()}
		case reflect.Float64:
			resultReturn = &FunctionCallReturn{CType: "number", Value: result.ToFloat()}
		case reflect.String:
			resultReturn = &FunctionCallReturn{CType: "string", Value: result.ToString()}
		case reflect.Array:
			resultReturn = &FunctionCallReturn{CType: "array", Value: result.Export()}
		case reflect.Map:
			resultReturn = &FunctionCallReturn{CType: "object", Value: result.ToBoolean()}
		default:
			return nil, fmt.Errorf("EnterFunctionCall: unsupported datatype")
		}
	}

	// Das Ergebniss wird zurückgegeben
	return resultReturn, nil
}

func (o *SharedPublicFunction) GetName() string {
	return o.name
}

func (o *SharedPublicFunction) GetParmTypes() []string {
	return o.parmTypes
}

func (o *SharedPublicFunction) EnterFunctionCall(parms types.RpcRequestInterface) (types.FunctionCallReturnInterface, error) {
	// Es wird geprüft ob die Angeforderte Anzahl an Parametern vorhanden ist
	if len(parms.GetParms()) != len(o.parmTypes) {
		return nil, fmt.Errorf("EnterFunctionCall: invalid parameters")
	}

	// Es wird versucht die Paraemter in den Richtigen GoJa Datentypen umzuwandeln
	convertedValues := make([]goja.Value, 0)
	for _, item := range parms.GetParms() {
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

	// Der Rückgabewert wird ermittelt und geprüft
	var resultReturn *FunctionCallReturn
	if result == nil {
		resultReturn = &FunctionCallReturn{CType: "null", Value: nil}
	} else if result.ExportType() == goja.Undefined().ExportType() && result.Export() == nil {
		resultReturn = &FunctionCallReturn{CType: "undefined", Value: nil}
	} else {
		switch result.ExportType().Kind() {
		case reflect.Bool:
			resultReturn = &FunctionCallReturn{CType: "boolean", Value: result.ToBoolean()}
		case reflect.Int64:
			resultReturn = &FunctionCallReturn{CType: "number", Value: result.ToInteger()}
		case reflect.Float64:
			resultReturn = &FunctionCallReturn{CType: "number", Value: result.ToFloat()}
		case reflect.String:
			resultReturn = &FunctionCallReturn{CType: "string", Value: result.ToString()}
		case reflect.Array:
			resultReturn = &FunctionCallReturn{CType: "array", Value: result.Export()}
		case reflect.Map:
			resultReturn = &FunctionCallReturn{CType: "object", Value: result.ToBoolean()}
		default:
			return nil, fmt.Errorf("EnterFunctionCall: unsupported datatype")
		}
	}

	// Das Ergebniss wird zurückgegeben
	return resultReturn, nil
}
