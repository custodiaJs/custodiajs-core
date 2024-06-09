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

// Wandelt die Funktionsargumente in Strings um
func ConvertV8ValuesToString(ctx *v8.Context, args []*v8.Value) ([]string, error) {
	// Es werden alle Stringwerte Extrahiert
	extracted := []string{}
	for _, item := range args {
		// Es wird geprüft ob es sich um ein Objekt oder um ein Array handelt
		switch {
		case item.IsObject() && !item.IsArray():
			// Es wird versucht das V8GO Objekt in ein Go Struct umzuwandeln
			obj, err := V8ObjectToGoObject(ctx, item)
			if err != nil {
				return nil, fmt.Errorf("internal error by converting, value")
			}

			// Das Objekt wird in JSON Umgewandelt
			encoded, err := json.Marshal(obj)
			if err != nil {
				return nil, fmt.Errorf("internal error by converting, value")
			}

			// Der JSON Wert wird zwischengespeichert
			extracted = append(extracted, string(encoded))
		case item.IsArray():
			// Das V8Go Array wird in ein Go Array umgewandelt
			obj, err := V8ArrayToGoArray(ctx, item)
			if err != nil {
				return nil, fmt.Errorf("internal error by converting, value")
			}

			// Die Einzelnene Einträge werden abgeabreitet
			var extra []interface{}
			for _, item := range obj {
				extra = append(extra, item.Value)
			}

			// Die Einträge werden Dekodiert
			encoded, err := json.Marshal(extra)
			if err != nil {
				return nil, fmt.Errorf("internal error by converting, value")
			}

			// Das Go Array wird zwischengspeichert
			extracted = append(extracted, string(encoded))
		case item.IsFunction():
			// Es wird geprüft ob es sich um eine Asynchrone oder um eine Synchrone Funktion handelt
			if item.IsAsyncFunction() {
				extracted = append(extracted, fmt.Sprintf("ASYNC:=%p", item))
			} else {
				extracted = append(extracted, fmt.Sprintf("SYNC:=%p", item))
			}
		default:
			extracted = append(extracted, item.String())
		}
	}

	// Die Exportierten Werte werden zurückgegeben
	return extracted, nil
}

// Wandelt V8GO Daten in Go Daten um
func ConvertV8DataToGoData(args []*v8.Value) ([]*types.ExportedV8Value, error) {
	// Speichert alle FunktionsStates ab
	returnValues := make([]*types.ExportedV8Value, 0)

	// Die Einzelnen Parameter werden abgearbeitet
	for _, item := range args {
		// Das Datentyp wird ermittelt
		var responseData *types.ExportedV8Value
		if item == nil {
			responseData = &types.ExportedV8Value{Type: "null", Value: nil}
		} else if item.IsUndefined() || item.IsNull() {
			responseData = &types.ExportedV8Value{Type: "undefined", Value: nil}
		} else {
			switch {
			case item.IsString():
				responseData = &types.ExportedV8Value{Type: "string", Value: item.String()}
			case item.IsNumber():
				switch {
				case item.IsBigInt():
					responseData = &types.ExportedV8Value{Type: "number", Value: item.BigInt().String()}
				case item.IsInt32():
					responseData = &types.ExportedV8Value{Type: "number", Value: item.Int32()}
				case item.IsUint32():
					responseData = &types.ExportedV8Value{Type: "number", Value: item.Uint32()}
				case item.IsNumber():
					responseData = &types.ExportedV8Value{Type: "number", Value: item.Number()}
				default:
					responseData = &types.ExportedV8Value{Type: "number", Value: item.Integer()}
				}
			case item.IsBoolean():
				responseData = &types.ExportedV8Value{Type: "boolean", Value: item.Boolean()}
			case item.IsObject():
				fmt.Println("Wert ist ein Array:")
			case item.IsArray():
				fmt.Println("Wert ist ein Array:")
			case item.IsFunction():
				fmt.Println("Wert ist ein Array:")
			default:
				// Es wird ein Javascript Fehler zurückgegeben
				return nil, fmt.Errorf("unsupported datatype for shared function response")
			}
		}

		// Der Eintrag wird abgespeichert
		returnValues = append(returnValues, responseData)
	}

	// Die Werte werden zurückgegeben
	return returnValues, nil
}
