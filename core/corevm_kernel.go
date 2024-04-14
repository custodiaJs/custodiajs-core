package core

/*
func root_fshare(methode string, functionName string, parameterTypes goja.Value, function goja.Value, runtime *goja.Runtime, vm *JsVM) goja.Value {
	// Die JS Funktion wird geprüft und abgespeichert
	jsFunc, ok := goja.AssertFunction(function)
	if !ok {
		panic(runtime.NewTypeError("Zweites Argument ist keine Funktion"))
	}

	// Es wird geprüft ob eine Funktion mit dem Namen bereits geteilt wird
	if vm.functionIsSharing(functionName) {
		// Der Vorgang wird abgebrochen
		panic(runtime.NewTypeError("Zweites Argument ist keine Funktion"))
	}

	// Die Anzahl der Parameter werden ermittelt
	paramCount := function.ToObject(runtime).Get("length")

	// Die Parametertypen werden ausgelesen
	functionParms, isList := parameterTypes.Export().([]interface{})
	if !isList {
		panic("invalid data type")
	}

	// Die Einzelnen Parameter werden geprüft und extrahiert
	extractedData := []string{}
	for _, item := range functionParms {
		stringItem := item.(string)
		if !utils.ValidateDatatypeString(stringItem) {
			panic("not allowed datatype: " + stringItem)
		}
		extractedData = append(extractedData, stringItem)
	}

	// Es wird geprüft ob genausoviele Parametertypen angegeben wurden, wie Parameter vorhanden sind
	if len(extractedData) != int(paramCount.ToInteger()) {
		panic(fmt.Sprintf("not all parameters has datatypes: %d %d", paramCount.ToInteger(), len(extractedData)))
	}

	// Die Funktion wird geteilt
	switch strings.ToLower(methode) {
	case "local":
		if err := vm.shareLocalFunction(functionName, extractedData, jsFunc); err != nil {
			panic(runtime.NewTypeError(err.Error()))
		}
	case "public":
		if err := vm.sharePublicFunction(functionName, extractedData, jsFunc); err != nil {
			panic(runtime.NewTypeError(err.Error()))
		}
	default:
		panic(runtime.NewTypeError("Zweites: " + methode))
	}

	// Es wird ein Undefined zurückgegeben
	return goja.Undefined()
}


func (o *JsVM) functionIsSharing(functionName string) bool {
	// Es wird geprüft ob die Funktion bereits Registriert wurde
	_, found := o.sharedLocalFunctions[functionName]

	// Das Ergebniss wird zurückgegeben
	return found
}

func (o *JsVM) shareLocalFunction(funcName string, parmTypes []string, function goja.Callable) error {
	// Es wird geprüft ob diese Funktion bereits registriert wurde
	if _, found := o.sharedLocalFunctions[funcName]; found {
		return fmt.Errorf("function always registrated")
	}

	// Die Funktion wird zwischengespeichert
	o.sharedLocalFunctions[funcName] = &SharedLocalFunction{
		callFunction: function,
		name:         funcName,
		parmTypes:    parmTypes,
		v8VM:         o.v8VM,
	}

	// Die Funktion wird im Core registriert
	fmt.Println("VM:SHARE_LOCAL_FUNCTION:", funcName, parmTypes)

	// Der Vorgang wurde ohne Fehler durchgeführt
	return nil
}

func (o *JsVM) sharePublicFunction(funcName string, parmTypes []string, function goja.Callable) error {
	// Es wird geprüft ob diese Funktion bereits registriert wurde
	if _, found := o.sharedPublicFunctions[funcName]; found {
		return fmt.Errorf("function always registrated")
	}

	// Die Funktion wird zwischengespeichert
	o.sharedPublicFunctions[funcName] = &SharedPublicFunction{
		callFunction: function,
		name:         funcName,
		parmTypes:    parmTypes,
		v8VM:         o.v8VM,
	}

	// Die Funktion wird im Core registriert
	fmt.Println("VM:SHARE_PUBLIC_FUNCTION:", funcName, parmTypes)

	// Der Vorgang wurde ohne Fehler durchgeführt
	return nil
}
*/
