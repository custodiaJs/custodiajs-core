package webgrpc

import (
	"fmt"
	"vnh1/grpc/publicgrpc"

	"google.golang.org/protobuf/types/known/structpb"
)

func convertToRPCFunctionParameter(val interface{}) *publicgrpc.RPCFunctionParameter {
	param := &publicgrpc.RPCFunctionParameter{}
	switch v := val.(type) {
	case string:
		param.Value = &publicgrpc.RPCFunctionParameter_StringValue{StringValue: v}
	case int64:
		param.Value = &publicgrpc.RPCFunctionParameter_Int64Value{Int64Value: v}
	case uint64:
		param.Value = &publicgrpc.RPCFunctionParameter_Uint64Value{Uint64Value: v}
	case float64:
		param.Value = &publicgrpc.RPCFunctionParameter_FloatValue{FloatValue: v}
	case bool:
		param.Value = &publicgrpc.RPCFunctionParameter_BooleanValue{BooleanValue: v}
	case []byte:
		param.Value = &publicgrpc.RPCFunctionParameter_BytesValue{BytesValue: v}
	// Füge weitere Typkonvertierungen nach Bedarf hinzu
	default:
		fmt.Println("Unbekannter Typ:", v)
	}
	return param
}

func mapToStruct(mapObj map[string]interface{}) (*structpb.Struct, error) {
	structObj := &structpb.Struct{
		Fields: make(map[string]*structpb.Value),
	}

	for key, val := range mapObj {
		structVal, err := structpb.NewValue(val)
		if err != nil {
			return nil, fmt.Errorf("mapToStruct: Fehler bei der Konvertierung des Werts für Schlüssel '%s': %v", key, err)
		}
		structObj.Fields[key] = structVal
	}

	return structObj, nil
}
