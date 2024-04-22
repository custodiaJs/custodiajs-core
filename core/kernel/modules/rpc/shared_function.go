package kmodulerpc

import (
	"fmt"
	"vnh1/types"

	v8 "rogchap.com/v8go"
)

// Gibt den Namen der Funktion zurück
func (o *SharedFunction) GetName() string {
	return o.name
}

// Gibt die Parameterdatentypen welche die Funktion erwartet zurück
func (o *SharedFunction) GetParmTypes() []string {
	return o.parmTypes
}

// Gibt den Rückgabedatentypen zurück
func (o *SharedFunction) GetReturndatentype() string {
	return o.returnType
}

// Gibt den Datentypen zurück
func (o *SharedFunction) GetReturnDType() string {
	return o.returnType
}

// Ruft die Geteilte Funktion auf
func (o *SharedFunction) EnterFunctionCall(req *types.RpcRequest) (*types.FunctionCallState, error) {
	// Es wird geprüft ob die Angeforderte Anzahl an Parametern vorhanden ist
	if len(req.Parms) != len(o.parmTypes) {
		return nil, fmt.Errorf("EnterFunctionCall: invalid parameters")
	}

	// Es wird ein neuer Context und ein neues Isolate beim Kernel angefordert
	iso, context, err := o.kernel.GetNewIsolateContext()
	if err != nil {
		return nil, fmt.Errorf("EnterFunctionCall: " + err.Error())
	}

	// Das Requestobjekt wird ersellt
	obj := v8.NewObjectTemplate(iso)

	// Es wird versucht die Paraemter in den Richtigen v8 Datentypen umzuwandeln
	convertedValues := make([]v8.Valuer, 0)
	for hight, item := range req.Parms {
		// Es wird geprüft ob der Datentyp gewünscht ist
		if item.CType != o.parmTypes[hight] {
			return nil, fmt.Errorf("EnterFunctionCall: not same parameter")
		}

		// Es wird geprüft ob es sich um einen Zulässigen Datentypen handelt
		switch item.CType {
		case "boolean":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		case "number":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		case "string":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		case "array":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		case "object":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		case "bytes":
			val, err := v8.NewValue(iso, item.Value)
			if err != nil {
				return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
			}
			convertedValues = append(convertedValues, val)
		default:
			return nil, fmt.Errorf("EnterFunctionCall: unsuported datatype")
		}
	}

	// Es wird ein neuer Request erzeugt
	request := NewSharedFunctionRequest(req)

	// Die Resolve Funktion wird festgelegt
	obj.Set("SendResponse", v8.NewFunctionTemplate(iso, request.SendResponse))

	// Die Senderror Funktion wird festgelegt
	obj.Set("SendError", v8.NewFunctionTemplate(iso, request.SendError))

	// Die Reject Funktion wird festgelegt
	obj.Set("Reject", v8.NewFunctionTemplate(iso, request.Reject))

	// Das Finale Objekt wird erstellt
	fobj, err := obj.NewInstance(context)
	if err != nil {
		return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
	}

	// Die Message Funktionen werden hinzugefügt
	writeMessageFunction := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		return nil
	})

	// Die Finalen Argumente werden erstellt
	finalArguments := make([]v8.Valuer, 0)
	finalArguments = append(finalArguments, fobj)
	finalArguments = append(finalArguments, convertedValues...)

	// Die Funktion wird eingelesen
	result, err := context.RunScript(o.functionSourceCode, "rpc_request.js")
	if err != nil {
		return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
	}

	// Die Funktion wird ausgelesen
	resultFunction, err := result.AsFunction()
	if err != nil {
		return nil, fmt.Errorf("SharedLocalFunction->EnterFunctionCall: " + err.Error())
	}

	// Die Funktion wird aufgerufen
	go func() {
		// Die Funktion wird aufgerufen
		_, err := resultFunction.Call(v8.Null(iso), finalArguments...)
		if err != nil {
			request.resolveChan <- &types.FunctionCallState{State: "failed", Error: err.Error()}
			return
		}
	}()

	// Es wird auf eine Antwort gewartet
	resolve := <-request.resolveChan

	// Sollte kein Reoslve Empfangen wurden sein, wird ein Fehler zurückgegeben, sowohl an den Absender sowohl auch an die Funktion
	if resolve == nil {
		return &types.FunctionCallState{State: "ok", Return: []*types.FunctionCallReturnData{}}, nil
	}

	// Der Context und die ISO werden zerstört
	context.Close()
	iso.Dispose()

	// Das Ergebniss wird zurückgegeben
	return resolve, nil
}
