package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"vnh1/types"

	v8 "rogchap.com/v8go"
)

// Wandelt ein V8 Array in ein Golang Array um
func V8ArrayToGoArray(ctx *v8.Context, v8value *v8.Value) ([]*types.ExportedV8Value, error) {
	// Der V8 Wert wird in ein Go Obekt umgewnadelt
	obj, err := v8value.AsObject()
	if err != nil {
		return nil, fmt.Errorf("V8ArrayToGoArray: " + err.Error())
	}

	// Die Anzahl der Einträge wird ermittelt
	lengthJsValue, err := obj.Get("length")
	if err != nil {
		panic(err)
	}

	// Das Array wird abgearbeitet
	extrValues := make([]*types.ExportedV8Value, 0)
	for i := 0; i < int(lengthJsValue.Integer()); i++ {
		// Der Wert wird abgerufen
		value, err := obj.GetIdx(uint32(i))
		if err != nil {
			continue
		}

		// Der V8 Wert wird in ein Go Wert umgewandelt
		extrValue, err := V8ValueToGoValue(ctx, value)
		if err != nil {
			return nil, fmt.Errorf("V8ArrayToGoArray: " + err.Error())
		}

		// Der wert wird Zwischengespeichert
		extrValues = append(extrValues, extrValue)
	}

	// Rückgabe
	return extrValues, nil
}

// Wandelt ein V8 Objekt in eine Golang Objekt um
func V8ObjectToGoObject(ctx *v8.Context, val v8.Valuer) (map[string]interface{}, error) {
	// Es wird versucht das Objekt in einen String umzuwandeln
	strvalue, err := v8.JSONStringify(ctx, val)
	if err != nil {
		return nil, fmt.Errorf("V8ObjectToGoObject: " + err.Error())
	}

	// Der Sting wird mittels JSON in ein Go Objekt umgewandelt
	var goObject map[string]interface{}
	err = json.Unmarshal([]byte(strvalue), &goObject)
	if err != nil {
		fmt.Println(strvalue, val)
		log.Fatalf("Error parsing JSON: %s", err)
	}

	// Das GoObjekt wird zurückgegben
	return goObject, nil
}

// Wandelt ein V8 Wert in ein Golang Wert um
func V8ValueToGoValue(ctx *v8.Context, val *v8.Value) (*types.ExportedV8Value, error) {
	switch {
	case val.IsUndefined():
		return &types.ExportedV8Value{Type: "undefined", Value: nil}, nil
	case val.IsNull():
		return &types.ExportedV8Value{Type: "null", Value: nil}, nil
	case val.IsBoolean():
		return &types.ExportedV8Value{Type: "bool", Value: val.Boolean()}, nil
	case val.IsNumber():
		return &types.ExportedV8Value{Type: "number", Value: val.Number()}, nil
	case val.IsString():
		return &types.ExportedV8Value{Type: "string", Value: val.String()}, nil
	case val.IsArray():
		// Das V8 Objekt wird in ein Go Array umgewadelt
		resolve, err := V8ArrayToGoArray(ctx, val)
		if err != nil {
			return nil, fmt.Errorf("V8ValueToGoValue: %v", err)
		}

		// Rückgabe
		return &types.ExportedV8Value{Type: "array", Value: resolve}, nil
	case val.IsObject():
		// Es wird ein Objekt erzeugt
		v8Object, err := val.AsObject()
		if err != nil {
			return nil, fmt.Errorf("V8ValueToGoValue: %v", err)
		}

		// Das V8 Objekt wird in ein Go Objekt umgewandelt
		resolve, err := V8ObjectToGoObject(ctx, v8Object)
		if err != nil {
			return nil, fmt.Errorf("V8ValueToGoValue: %v", err)
		}

		// Rückgabe
		return &types.ExportedV8Value{Type: "object", Value: resolve}, nil
	default:
		return nil, fmt.Errorf("V8ValueToGoValue: unsupported type")
	}
}
