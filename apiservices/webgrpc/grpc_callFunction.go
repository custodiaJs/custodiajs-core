package webgrpc

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
	"vnh1/grpc/publicgrpc"
	"vnh1/types"
	"vnh1/utils"

	"github.com/dop251/goja"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func (s *GrpcServer) CallFunction(ctx context.Context, in *publicgrpc.RPCFunctionCall) (*publicgrpc.RPCResponse, error) {
	// Die Informationen über den Client werden ermittelt
	p, ok := peer.FromContext(ctx)
	if !ok {
		log.Println("Keine Peer-Informationen gefunden")
		return nil, status.Error(codes.Internal, "Fehler beim Abrufen der Peer-Informationen")
	}

	// Es wird eine neue Process Log Session erzeugt
	proc := utils.NewProcLogSession()
	proc.LogPrint("gRPC: %s", utils.FormatConsoleText(types.VALIDATE_INCOMMING_REMOTE_FUNCTION_CALL_REQUEST_FROM, p.Addr.String()))

	// Die TLS Informationen werden abgerufen
	tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		log.Println("Keine TLS-Authentifizierungsinformationen gefunden")
		return nil, status.Error(codes.Unauthenticated, "Keine TLS-Authentifizierungsinformationen")
	}
	_ = tlsInfo

	// Es wird geprüft ob es sich um eine bekannte VM handelt
	proc.LogPrint("gRPC: determine the script container '%s'\n", strings.ToLower(in.ContainerId))
	foundedVM, err := s.core.GetScriptContainerVMByID(in.ContainerId)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("unkown script container")
	}

	// Es wird versucht die Passende Funktion zu ermitteln
	proc.LogPrint("gRPC: &[%s]: determine the function '%s'\n", foundedVM.GetVMName(), in.FunctionName)
	var foundFunction types.SharedFunctionInterface
	for _, item := range foundedVM.GetLocalSharedFunctions() {
		if item.GetName() == in.FunctionName {
			foundFunction = item
			break
		}
	}

	// Es wird geprüft ob eine Funktion gefunden wurde
	if foundFunction == nil {
		for _, item := range foundedVM.GetPublicSharedFunctions() {
			if item.GetName() == in.FunctionName {
				foundFunction = item
				break
			}
		}
	}

	// Sollte keine Passende Funktion gefunden werden, wird der Vorgang abgebrochen
	if foundFunction == nil {
		return nil, fmt.Errorf("unkown function")
	}

	// Es wird ermitelt ob die Datentypen korrekt sind
	if len(foundFunction.GetParmTypes()) != len(in.Parms) {
		return nil, fmt.Errorf("invalid parms")
	}

	// Die Einzelnen Parameter werden geprüft und abgearbeitet
	proc.LogPrint("gRPC: &[%s]: convert function '%s' parameters\n", foundedVM.GetVMName(), foundFunction.GetName())
	extractedValues := make([]types.FunctionParameterBundleInterface, 0)
	for x := range foundFunction.GetParmTypes() {
		// Es wird versucht den Datentypen umzuwandeln
		switch v := in.Parms[x].Value.(type) {
		case *publicgrpc.RPCFunctionParameter_StringValue:
			if foundFunction.GetParmTypes()[x] != "string" {
				return nil, fmt.Errorf("invalid function parm datatype 1")
			}

			// Der Eintrag wird erzeugt
			newEntry := &FunctionParameterCapsle{Value: v.StringValue, CType: "string"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case *publicgrpc.RPCFunctionParameter_Int64Value:
			if foundFunction.GetParmTypes()[x] != "number" {
				return nil, fmt.Errorf("invalid function parm datatype 2")
			}

			// Der Eintrag wird erzeugt
			newEntry := &FunctionParameterCapsle{Value: v.Int64Value, CType: "number"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case *publicgrpc.RPCFunctionParameter_Uint64Value:
			if foundFunction.GetParmTypes()[x] != "number" {
				return nil, fmt.Errorf("invalid function parm datatype 3")
			}

			// Der Eintrag wird erzeugt
			newEntry := &FunctionParameterCapsle{Value: v.Uint64Value, CType: "number"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case *publicgrpc.RPCFunctionParameter_FloatValue:
			if foundFunction.GetParmTypes()[x] != "number" {
				return nil, fmt.Errorf("invalid function parm datatype 4")
			}

			// Der Eintrag wird erzeugt
			newEntry := &FunctionParameterCapsle{Value: v.FloatValue, CType: "number"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case *publicgrpc.RPCFunctionParameter_BooleanValue:
			if foundFunction.GetParmTypes()[x] != "boolean" {
				return nil, fmt.Errorf("invalid function parm datatype 5")
			}

			// Der Eintrag wird erzeugt
			newEntry := &FunctionParameterCapsle{Value: v.BooleanValue, CType: "boolean"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case *publicgrpc.RPCFunctionParameter_BytesValue:
			if foundFunction.GetParmTypes()[x] != "bytes" {
				return nil, fmt.Errorf("invalid function parm datatype 6")
			}

			// Der Eintrag wird erzeugt
			newEntry := &FunctionParameterCapsle{Value: v.BytesValue, CType: "byted"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case *publicgrpc.RPCFunctionParameter_TimestampValue:
			if foundFunction.GetParmTypes()[x] != "timestamp" {
				return nil, fmt.Errorf("invalid function parm datatype 7")
			}

			// Der Eintrag wird erzeugt
			newEntry := &FunctionParameterCapsle{Value: v.TimestampValue, CType: "timestamp"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case *publicgrpc.RPCFunctionParameter_ObjectValue:
			if foundFunction.GetParmTypes()[x] != "object" {
				return nil, fmt.Errorf("invalid function parm datatype 8")
			}

			// Der Eintrag wird erzeugt
			newEntry := &FunctionParameterCapsle{Value: v.ObjectValue, CType: "object"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		case *publicgrpc.RPCFunctionParameter_ArrayValue:
			if foundFunction.GetParmTypes()[x] != "array" {
				return nil, fmt.Errorf("invalid function parm datatype 9")
			}

			// Der Eintrag wird erzeugt
			newEntry := &FunctionParameterCapsle{Value: v.ArrayValue, CType: "array"}

			// Die Daten werden hinzugefügt
			extractedValues = append(extractedValues, newEntry)
		default:
			return nil, fmt.Errorf("unsuported datatype")
		}
	}

	// Die Funktion wird aufgerufen
	proc.LogPrint("gRPC: &[%s]: call function '%s'\n", foundedVM.GetVMName(), foundFunction.GetName())
	result, err := foundFunction.EnterFunctionCall(nil, &RpcRequest{parms: extractedValues})
	if err != nil {
		proc.LogPrint("RPC: &[%s]: call function '%s' error\n\t%s\n", foundedVM.GetVMName(), foundFunction.GetName(), err)
		return nil, fmt.Errorf("calling error")
	}
	proc.LogPrintSuccs("gRPC: &[%s]: function '%s' call, done\n", foundedVM.GetVMName(), foundFunction.GetName())

	// Der Rückgabewert wird ermittelt und geprüft
	var data *publicgrpc.RPCResponseData
	if result == nil {
		data = &publicgrpc.RPCResponseData{Value: nil}
	} else if result.ExportType() == goja.Undefined().ExportType() && result.Export() == nil {
		data = &publicgrpc.RPCResponseData{Value: nil}
	} else {
		switch result.ExportType().Kind() {
		case reflect.Bool:
			data = &publicgrpc.RPCResponseData{Value: &publicgrpc.RPCResponseData_BooleanValue{BooleanValue: result.ToBoolean()}}
		case reflect.Int64:
			data = &publicgrpc.RPCResponseData{Value: &publicgrpc.RPCResponseData_Int64Value{Int64Value: result.ToInteger()}}
		case reflect.Float64:
			data = &publicgrpc.RPCResponseData{Value: &publicgrpc.RPCResponseData_FloatValue{FloatValue: result.ToFloat()}}
		case reflect.String:
			data = &publicgrpc.RPCResponseData{Value: &publicgrpc.RPCResponseData_StringValue{StringValue: result.String()}}
		case reflect.Slice:
			// Das Interface Objekt wird erstellt
			slicedObject, isConverted := result.Export().([]interface{})
			if !isConverted {
				return nil, fmt.Errorf("invalid object datatype, slice")
			}

			// Die Daten werden in gRPC Array Wert umgewandelt
			arrayValue := &publicgrpc.ArrayValue{}
			for _, item := range slicedObject {
				rpcParam := convertToRPCFunctionParameter(item)
				arrayValue.Values = append(arrayValue.Values, rpcParam)
			}

			// Die Antwort wird erstellt
			data = &publicgrpc.RPCResponseData{Value: &publicgrpc.RPCResponseData_ArrayValue{ArrayValue: arrayValue}}
		case reflect.Map:
			// Das Objekt wird in eine Map umgewandelt
			mapObjected, isConverted := result.Export().(map[string]interface{})
			if !isConverted {
				return nil, fmt.Errorf("invalid object datatype, object")
			}

			// Das Objekt wird Convertiert
			convertedMap, err := mapToStruct(mapObjected)
			if err != nil {
				return nil, fmt.Errorf("map converting error, object")
			}

			// Die Antwort wird gebaut
			data = &publicgrpc.RPCResponseData{Value: &publicgrpc.RPCResponseData_ObjectValue{ObjectValue: convertedMap}}
		case reflect.Func:
			return nil, fmt.Errorf("function return not allowed in web grpc")
		default:
			fmt.Println(result.ExportType())
			return nil, fmt.Errorf("EnterFunctionCall: unsupported datatype")
		}
	}

	// Das Ergebniss wird zurückgesendet
	proc.LogPrint("gRPC: &[%s]: done, send response\n", foundedVM.GetVMName())
	return &publicgrpc.RPCResponse{Result: "success", Data: data, Error: ""}, nil
}
