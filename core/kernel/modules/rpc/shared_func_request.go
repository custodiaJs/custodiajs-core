package kmodulerpc

import (
	"fmt"
	"vnh1/types"

	v8 "rogchap.com/v8go"
)

func (o *SharedFunctionRequest) SendResponse(info *v8.FunctionCallbackInfo) *v8.Value {
	// Speichert alle FunktionsStates ab
	resolves := &types.FunctionCallState{
		Return: make([]*types.FunctionCallReturnData, 0),
	}

	// Die Einzelnen Parameter werden abgearbeitet
	for _, item := range info.Args() {
		// Das Datentyp wird ermittelt
		var responseData *types.FunctionCallReturnData
		if item == nil {
			responseData = &types.FunctionCallReturnData{CType: "null", Value: nil}
		} else if item.IsUndefined() || item.IsNull() {
			responseData = &types.FunctionCallReturnData{CType: "undefined", Value: nil}
		} else {
			switch {
			case item.IsString():
				responseData = &types.FunctionCallReturnData{CType: "string", Value: item.String()}
			case item.IsNumber():
				switch {
				case item.IsBigInt():
					responseData = &types.FunctionCallReturnData{CType: "number", Value: item.BigInt().String()}
				case item.IsInt32():
					responseData = &types.FunctionCallReturnData{CType: "number", Value: item.Int32()}
				case item.IsUint32():
					responseData = &types.FunctionCallReturnData{CType: "number", Value: item.Uint32()}
				case item.IsNumber():
					responseData = &types.FunctionCallReturnData{CType: "number", Value: item.Number()}
				default:
					responseData = &types.FunctionCallReturnData{CType: "number", Value: item.Integer()}
				}
			case item.IsBoolean():
				responseData = &types.FunctionCallReturnData{CType: "boolean", Value: item.Boolean()}
			case item.IsObject():
				fmt.Println("Wert ist ein Array:")
			case item.IsArray():
				fmt.Println("Wert ist ein Array:")
			case item.IsFunction():
				fmt.Println("Wert ist ein Array:")
			default:
				return nil
			}
		}

		// Der Eintrag wird abgespeichert
		resolves.Return = append(resolves.Return, responseData)
	}

	// Es wird geprüft ob ein Rückgabewert vorhanden ist, wenn nicht wird ein Undefined zurückgegeben
	if len(resolves.Return) == 0 {
		resolves.Return = append(resolves.Return, &types.FunctionCallReturnData{CType: "undefined", Value: nil})
	}

	// Der Stauts wird aktualisiert
	resolves.State = "ok"

	// Die Antwort wird zurückgesendet
	o.resolveChan <- resolves

	// Es ist kein Fehler aufgetreten
	return nil
}

func (o *SharedFunctionRequest) SendError(info *v8.FunctionCallbackInfo) *v8.Value {
	o.resolveChan <- &types.FunctionCallState{}
	fmt.Println("RESPONSE")
	return nil
}

func (o *SharedFunctionRequest) Reject(info *v8.FunctionCallbackInfo) *v8.Value {
	o.resolveChan <- &types.FunctionCallState{}
	fmt.Println("RESPONSE")
	return nil
}

func NewSharedFunctionRequest(req types.RpcRequestInterface) *SharedFunctionRequest {
	return &SharedFunctionRequest{resolveChan: make(chan *types.FunctionCallState), parms: req}
}
