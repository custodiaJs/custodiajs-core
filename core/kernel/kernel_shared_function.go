package kernel

import (
	"fmt"
	"log"
	"vnh1/utils"

	v8 "rogchap.com/v8go"
)

func _util_rpc_shareFunctionParmArrayReader(result *v8.Value) []string {
	// Das resultierende v8go.Value sollte ein Array sein
	if !result.IsArray() {
		log.Fatal("Result is not an array")
	}

	// Die Länge des Arrays ermitteln
	obj := result.Object()
	lengthJsValue, err := obj.Get("length")
	if err != nil {
		panic(err)
	}
	length := lengthJsValue.Integer()

	// Das Array wird abgearbeitet
	extrStr := make([]string, 0)
	for i := 0; i < int(length); i++ {
		value, err := result.Object().GetIdx(uint32(i))
		if err != nil {
			continue
		}
		extrStr = append(extrStr, value.String())
	}

	// Rückgabe
	return extrStr
}

func (o *Kernel) _kernel_rpc_functionIsSharing() *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(o.Context.Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Der Mutex wird angewendet und autoamtisch wieder Freigegeben
		o.mutex.Lock()
		defer o.mutex.Unlock()

		// Es wird geprüft ob die Funktion bereits Registriert wurde
		_, found := o.sharedLocalFunctions[info.Args()[0].String()]

		// Das Ergebniss wird erstellt
		returnV, err := v8.NewValue(o.Context.Isolate(), found)
		if err != nil {
			panic(err)
		}

		// Das Ergebnis wird zurückgegeben
		return returnV
	})
}

func (o *Kernel) _kernel_rpc_shareLocalFunction() *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(o.Context.Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Die Parameter werden geprüft
		var sharedFunction *v8.Function
		parameterTypes := make([]string, 0)
		var functionName string
		returnType := "none"
		var err error
		if len(info.Args()) == 2 {
			// Es wird geprüft ob als erstes ein String angegeben wurde
			if !info.Args()[0].IsString() {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				o._kernel_throw(info.Context(), "invalid parameter chain")

				// Der Vorgang wird beendet
				return nil
			}

			// Es wird geprüft ob als nächstes eine Funktion vorhanden ist
			if !info.Args()[1].IsAsyncFunction() {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				o._kernel_throw(info.Context(), "invalid parameter chain")

				// Der Vorgang wird beendet
				return nil
			}

			// Der Funktionsname wird extrahiert
			functionName = info.Args()[0].String()

			// Die Funktion wird Extrahiert
			sharedFunction, err = info.Args()[1].AsFunction()
			if err != nil {
				panic(err)
			}
		} else if len(info.Args()) >= 3 {
			// Es wird geprüft ob als erstes ein String angegeben wurde
			if !info.Args()[0].IsString() {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				o._kernel_throw(info.Context(), "invalid parameter chain")

				// Der Vorgang wird beendet
				return nil
			}

			// Es wird geprüft ob als zweites ein Array nur mit Strings angegeben wurde
			if !info.Args()[1].IsArray() {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				o._kernel_throw(info.Context(), "invalid parameter chain")

				// Der Vorgang wird beendet
				return nil
			}

			// Es wird geprüft ob ein Rückgabetyp angegeben wurde
			if len(info.Args()) >= 4 {
				// Es wird geprüft ob als Drittes ein String vorhanden ist, welcher Angibt was für ein Datentyp zurückgegeben wird
				if !info.Args()[2].IsString() {
					// Die Fehlermeldung wird erstellt und an JS zurückgegeben
					o._kernel_throw(info.Context(), "invalid parameter chain")

					// Der Vorgang wird beendet
					return nil
				}

				// Es wird geprüft ob als nächstes eine Funktion vorhanden ist
				if !info.Args()[3].IsAsyncFunction() {
					// Die Fehlermeldung wird erstellt und an JS zurückgegeben
					o._kernel_throw(info.Context(), "invalid parameter chain")

					// Der Vorgang wird beendet
					return nil
				}

				// Der Funktionsname wird extrahiert
				functionName = info.Args()[0].String()

				// Die Parametertypen werden ausgelsen
				parameterTypes = _util_rpc_shareFunctionParmArrayReader(info.Args()[1])

				// Der Rückgabe Type wird ausgelsen
				returnType = info.Args()[2].String()

				// Die Funktion wird Extrahiert
				sharedFunction, err = info.Args()[3].AsFunction()
				if err != nil {
					panic(err)
				}
			} else {
				// Es wird geprüft ob als nächstes eine Funktion vorhanden ist
				if !info.Args()[2].IsAsyncFunction() {
					// Die Fehlermeldung wird erstellt und an JS zurückgegeben
					o._kernel_throw(info.Context(), "invalid parameter chain")

					// Der Vorgang wird beendet
					return nil
				}

				// Der Funktionsname wird extrahiert
				functionName = info.Args()[0].String()

				// Die Parametertypen werden ausgelsen
				parameterTypes = _util_rpc_shareFunctionParmArrayReader(info.Args()[1])

				// Die Funktion wird Extrahiert
				sharedFunction, err = info.Args()[2].AsFunction()
				if err != nil {
					panic(err)
				}
			}
		} else {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			o._kernel_throw(info.Context(), "invalid parameter chain")

			// Der Vorgang wird beendet
			return nil
		}

		// Es wird geprüft ob es sich um einen Zulässigen Funktionsnamen handelt
		if !utils.ValidateFunctionName(functionName) {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			o._kernel_throw(info.Context(), "invalid function name")

			// Der Vorgang wird beendet
			return nil
		}

		// Der Rückgabetype wird geprüft
		if !utils.ValidateDatatypeString(returnType) && returnType != "none" {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			o._kernel_throw(info.Context(), fmt.Sprintf("invalid return type '%s'", returnType))

			// Der Vorgang wird beendet
			return nil
		}

		// Die Einzelnen Parametertypen werden geprüft
		for index, item := range parameterTypes {
			// Der Einzelne Parameter wird geprüft
			if !utils.ValidateDatatypeString(item) {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				o._kernel_throw(info.Context(), fmt.Sprintf("Failed to get element at index %d: %v", index, "invalid return datatype"))

				// Der Vorgang wird beendet
				return nil
			}
		}

		// Der Mutex wird angewendet
		o.mutex.Lock()
		defer o.mutex.Unlock()

		// Es wird geprüft ob diese Funktion bereits registriert wurde
		if _, found := o.sharedLocalFunctions[functionName]; found {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			o._kernel_throw(info.Context(), "some type error")

			// Der Vorgang wird beendet
			return nil
		}

		// Die Geteilte Funktion wird erzeugt
		o.sharedLocalFunctions[functionName] = &SharedLocalFunction{
			callFunction: sharedFunction,
			name:         info.Args()[0].String(),
			parmTypes:    parameterTypes,
			returnType:   returnType,
			v8VM:         o.Context,
		}

		// Die Funktion wird im Core registriert
		fmt.Println("VM:SHARE_LOCAL_FUNCTION:", functionName, "")

		// Der Vorgang wurde ohne Fehler durchgeführt
		return nil
	})
}

func (o *Kernel) _kernel_rpc_sharePublicFunction() *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(o.Context.Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Es wird geprüft ob diese Funktion bereits registriert wurde
		if _, found := o.sharedPublicFunctions[info.Args()[0].String()]; found {
			typeErr, err := v8.NewValue(o.Context.Isolate(), "some type error")
			if err != nil {
				panic(err)
			}
			info.Context().Isolate().ThrowException(typeErr)
			return nil
		}

		// Die Funktion wird zwischengespeichert
		o.sharedPublicFunctions[info.Args()[0].String()] = &SharedPublicFunction{
			//callFunction: function,
			name: info.Args()[0].String(),
			//parmTypes: info.Args()[0].Object(),
			v8VM: o.Context,
		}

		// Die Funktion wird im Core registriert
		fmt.Println("VM:SHARE_LOCAL_FUNCTION:", info.Args()[0].String(), "")
		fmt.Println(info.Args())

		// Der Vorgang wurde ohne Fehler durchgeführt
		return nil
	})
}

func (o *Kernel) _kernel_rpc_call_local() *v8.FunctionTemplate {
	return nil
}

func (o *Kernel) _kernel_rpc_call_remote() *v8.FunctionTemplate {
	return nil
}
