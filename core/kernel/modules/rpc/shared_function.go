package kmodulerpc

import (
	"fmt"
	"vnh1/types"

	v8 "rogchap.com/v8go"
)

func (o *SharedFunction) GetName() string {
	return o.name
}

func (o *SharedFunction) GetParmTypes() []string {
	return o.parmTypes
}

func (o *SharedFunction) EnterFunctionCall(req types.RpcRequestInterface) (*types.FunctionCallState, error) {
	// Es wird geprüft ob die Angeforderte Anzahl an Parametern vorhanden ist
	if len(req.GetParms()) != len(o.parmTypes) {
		return nil, fmt.Errorf("EnterFunctionCall: invalid parameters")
	}

	// Es wird versucht die Paraemter in den Richtigen v8 Datentypen umzuwandeln
	convertedValues := make([]v8.Valuer, 0)
	for hight, item := range req.GetParms() {
		// Es wird geprüft ob der Datentyp gewünscht ist
		if item.GetType() != o.parmTypes[hight] {
			return nil, fmt.Errorf("EnterFunctionCall: not same parameter")
		}

		// Es wird geprüft ob es sich um einen Zulässigen Datentypen handelt
		switch item.GetType() {
		case "boolean":
			val, err := v8.NewValue(o.v8VM.Isolate(), item.GetValue())
			if err != nil {
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		case "number":
			val, err := v8.NewValue(o.v8VM.Isolate(), item.GetValue())
			if err != nil {
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		case "string":
			val, err := v8.NewValue(o.v8VM.Isolate(), item.GetValue())
			if err != nil {
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		case "array":
			val, err := v8.NewValue(o.v8VM.Isolate(), item.GetValue())
			if err != nil {
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		case "object":
			val, err := v8.NewValue(o.v8VM.Isolate(), item.GetValue())
			if err != nil {
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		case "bytes":
			val, err := v8.NewValue(o.v8VM.Isolate(), item.GetValue())
			if err != nil {
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		default:
			return nil, fmt.Errorf("EnterFunctionCall: unsuported datatype")
		}
	}

	// Das Requestobjekt wird ersellt
	obj := v8.NewObjectTemplate(o.v8VM.Isolate())

	// Es wird ein neuer Request erzeugt
	request := NewSharedFunctionRequest(req)

	// Die Resolve Funktion wird festgelegt
	obj.Set("SendResponse", v8.NewFunctionTemplate(o.v8VM.Isolate(), request.SendResponse))

	// Die Senderror Funktion wird festgelegt
	obj.Set("SendError", v8.NewFunctionTemplate(o.v8VM.Isolate(), request.SendError))

	// Die Reject Funktion wird festgelegt
	obj.Set("Reject", v8.NewFunctionTemplate(o.v8VM.Isolate(), request.Reject))

	// Das Finale Objekt wird erstellt
	fobj, err := obj.NewInstance(o.v8VM)
	if err != nil {
		return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
	}

	// Die Finalen Argumente werden erstellt
	finalArguments := make([]v8.Valuer, 0)
	finalArguments = append(finalArguments, fobj)
	finalArguments = append(finalArguments, convertedValues...)

	// Die Funktion wird ausgeführt
	go func() {
		_, err = o.callFunction.Call(v8.Null(o.v8VM.Isolate()), finalArguments...)
		if err != nil {
			request.resolveChan <- &types.FunctionCallState{State: "failed", Error: err.Error()}
		}
	}()

	// Es wird auf eine Antwort gewartet
	resolve := <-request.resolveChan

	// Sollte kein Reoslve Empfangen wurden sein, wird ein Fehler zurückgegeben, sowohl an den Absender sowohl auch an die Funktion
	if resolve == nil {
		return &types.FunctionCallState{State: "ok", Return: []*types.FunctionCallReturnData{}}, nil
	}

	// Das Ergebniss wird zurückgegeben
	return resolve, nil
}
