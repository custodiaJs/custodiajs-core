package kmodulerpc

import (
	"fmt"
	"log"
	"vnh1/types"
	"vnh1/utils"

	v8 "rogchap.com/v8go"
)

type RPCModule int

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

func (o *RPCModule) _kernel_rpc_NewShareLocalFunction(kernel types.KernelInterface) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(kernel.ContextV8().Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
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
				kernel.KernelThrow(info.Context(), "invalid parameter chain")

				// Der Vorgang wird beendet
				return nil
			}

			// Es wird geprüft ob als nächstes eine Funktion vorhanden ist
			if !info.Args()[1].IsAsyncFunction() {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				kernel.KernelThrow(info.Context(), "invalid parameter chain")

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
				kernel.KernelThrow(info.Context(), "invalid parameter chain")

				// Der Vorgang wird beendet
				return nil
			}

			// Es wird geprüft ob als zweites ein Array nur mit Strings angegeben wurde
			if !info.Args()[1].IsArray() {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				kernel.KernelThrow(info.Context(), "invalid parameter chain")

				// Der Vorgang wird beendet
				return nil
			}

			// Es wird geprüft ob ein Rückgabetyp angegeben wurde
			if len(info.Args()) >= 4 {
				// Es wird geprüft ob als Drittes ein String vorhanden ist, welcher Angibt was für ein Datentyp zurückgegeben wird
				if !info.Args()[2].IsString() {
					// Die Fehlermeldung wird erstellt und an JS zurückgegeben
					kernel.KernelThrow(info.Context(), "invalid parameter chain")

					// Der Vorgang wird beendet
					return nil
				}

				// Es wird geprüft ob als nächstes eine Funktion vorhanden ist
				if !info.Args()[3].IsAsyncFunction() {
					// Die Fehlermeldung wird erstellt und an JS zurückgegeben
					kernel.KernelThrow(info.Context(), "invalid parameter chain")

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
					kernel.KernelThrow(info.Context(), "invalid parameter chain")

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
			kernel.KernelThrow(info.Context(), "invalid parameter chain")

			// Der Vorgang wird beendet
			return nil
		}

		// Es wird geprüft ob es sich um einen Zulässigen Funktionsnamen handelt
		if !utils.ValidateFunctionName(functionName) {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			kernel.KernelThrow(info.Context(), "invalid function name")

			// Der Vorgang wird beendet
			return nil
		}

		// Der Rückgabetype wird geprüft
		if !utils.ValidateDatatypeString(returnType) && returnType != "none" {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			kernel.KernelThrow(info.Context(), fmt.Sprintf("invalid return type '%s'", returnType))

			// Der Vorgang wird beendet
			return nil
		}

		// Die Einzelnen Parametertypen werden geprüft
		for index, item := range parameterTypes {
			// Der Einzelne Parameter wird geprüft
			if !utils.ValidateDatatypeString(item) {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				kernel.KernelThrow(info.Context(), fmt.Sprintf("Failed to get element at index %d: %v", index, "invalid return datatype"))

				// Der Vorgang wird beendet
				return nil
			}
		}

		// Es wird versucht die Tabelle abzurufen
		table, isok := kernel.GloablRegisterRead("rpc_local").(map[string]*SharedLocalFunction)
		if !isok {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			kernel.KernelThrow(info.Context(), "some type error 1")

			// Der Vorgang wird beendet
			return nil
		}

		// Es wird geprüft ob diese Funktion bereits registriert wurde
		if _, found := table[functionName]; found {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			kernel.KernelThrow(info.Context(), "some type error")

			// Der Vorgang wird beendet
			return nil
		}

		// Die Geteilte Funktion wird erzeugt
		table[functionName] = &SharedLocalFunction{
			callFunction: sharedFunction,
			name:         info.Args()[0].String(),
			parmTypes:    parameterTypes,
			returnType:   returnType,
			v8VM:         kernel.ContextV8(),
		}

		// Der Eintrag in der Datenbank wird geupdated
		if err := kernel.GloablRegisterWrite("rpc_local", table); err != nil {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			kernel.KernelThrow(info.Context(), "some type error")

			// Der Vorgang wird beendet
			return nil
		}

		// Die Funktion wird im Core registriert
		fmt.Println("VM:SHARE_LOCAL_FUNCTION:", functionName, "")

		// Der Vorgang wurde ohne Fehler durchgeführt
		return nil
	})
}

func (o *RPCModule) _kernel_rpc_NewSharePublicFunction(kernel types.KernelInterface) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(kernel.ContextV8().Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
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
				kernel.KernelThrow(info.Context(), "invalid parameter chain")

				// Der Vorgang wird beendet
				return nil
			}

			// Es wird geprüft ob als nächstes eine Funktion vorhanden ist
			if !info.Args()[1].IsAsyncFunction() {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				kernel.KernelThrow(info.Context(), "invalid parameter chain")

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
				kernel.KernelThrow(info.Context(), "invalid parameter chain")

				// Der Vorgang wird beendet
				return nil
			}

			// Es wird geprüft ob als zweites ein Array nur mit Strings angegeben wurde
			if !info.Args()[1].IsArray() {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				kernel.KernelThrow(info.Context(), "invalid parameter chain")

				// Der Vorgang wird beendet
				return nil
			}

			// Es wird geprüft ob ein Rückgabetyp angegeben wurde
			if len(info.Args()) >= 4 {
				// Es wird geprüft ob als Drittes ein String vorhanden ist, welcher Angibt was für ein Datentyp zurückgegeben wird
				if !info.Args()[2].IsString() {
					// Die Fehlermeldung wird erstellt und an JS zurückgegeben
					kernel.KernelThrow(info.Context(), "invalid parameter chain")

					// Der Vorgang wird beendet
					return nil
				}

				// Es wird geprüft ob als nächstes eine Funktion vorhanden ist
				if !info.Args()[3].IsAsyncFunction() {
					// Die Fehlermeldung wird erstellt und an JS zurückgegeben
					kernel.KernelThrow(info.Context(), "invalid parameter chain")

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
					kernel.KernelThrow(info.Context(), "invalid parameter chain")

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
			kernel.KernelThrow(info.Context(), "invalid parameter chain")

			// Der Vorgang wird beendet
			return nil
		}

		// Es wird geprüft ob es sich um einen Zulässigen Funktionsnamen handelt
		if !utils.ValidateFunctionName(functionName) {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			kernel.KernelThrow(info.Context(), "invalid function name")

			// Der Vorgang wird beendet
			return nil
		}

		// Der Rückgabetype wird geprüft
		if !utils.ValidateDatatypeString(returnType) && returnType != "none" {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			kernel.KernelThrow(info.Context(), fmt.Sprintf("invalid return type '%s'", returnType))

			// Der Vorgang wird beendet
			return nil
		}

		// Die Einzelnen Parametertypen werden geprüft
		for index, item := range parameterTypes {
			// Der Einzelne Parameter wird geprüft
			if !utils.ValidateDatatypeString(item) {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				kernel.KernelThrow(info.Context(), fmt.Sprintf("Failed to get element at index %d: %v", index, "invalid return datatype"))

				// Der Vorgang wird beendet
				return nil
			}
		}

		// Es wird versucht die Tabelle abzurufen
		table, isok := kernel.GloablRegisterRead("rpc_public").(map[string]*SharedPublicFunction)
		if !isok {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			kernel.KernelThrow(info.Context(), "some type error")

			// Der Vorgang wird beendet
			return nil
		}

		// Es wird geprüft ob diese Funktion bereits registriert wurde
		if _, found := table[functionName]; found {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			kernel.KernelThrow(info.Context(), "some type error")

			// Der Vorgang wird beendet
			return nil
		}

		// Die Geteilte Funktion wird erzeugt
		table[functionName] = &SharedPublicFunction{
			callFunction: sharedFunction,
			name:         info.Args()[0].String(),
			parmTypes:    parameterTypes,
			returnType:   returnType,
			v8VM:         kernel.ContextV8(),
		}

		// Der Eintrag in der Datenbank wird geupdated
		if err := kernel.GloablRegisterWrite("rpc_public", table); err != nil {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			kernel.KernelThrow(info.Context(), "some type error")

			// Der Vorgang wird beendet
			return nil
		}

		// Die Funktion wird im Core registriert
		fmt.Println("VM:SHARE_LOCAL_FUNCTION:", functionName, "")

		// Der Vorgang wurde ohne Fehler durchgeführt
		return nil
	})
}

func (o *RPCModule) _kernel_rpc_CallLocal(kernel types.KernelInterface) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(kernel.ContextV8().Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
		return nil
	})
}

func (o *RPCModule) _kernel_rpc_CallRemote(kernel types.KernelInterface) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(kernel.ContextV8().Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
		return nil
	})
}

func (o *RPCModule) _kernel_rpc_GetFunctionDetailsLocal(kernel types.KernelInterface) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(kernel.ContextV8().Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
		return nil
	})
}

func (o *RPCModule) _kernel_rpc_GetFunctionDetailsRemote(kernel types.KernelInterface) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(kernel.ContextV8().Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
		return nil
	})
}

func (o *RPCModule) _kernel_rpc_IsShareLocal(kernel types.KernelInterface) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(kernel.ContextV8().Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Es wird versucht die Tabelle abzurufen
		table, isok := kernel.GloablRegisterRead("rpc_local").(map[string]*SharedLocalFunction)
		if !isok {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			kernel.KernelThrow(info.Context(), "some type erro 1r")

			// Der Vorgang wird beendet
			return nil
		}

		// Es wird geprüft ob die Funktion bereits Registriert wurde
		_, found := table[info.Args()[0].String()]

		// Das Ergebniss wird erstellt
		returnV, err := v8.NewValue(kernel.ContextV8().Isolate(), found)
		if err != nil {
			panic(err)
		}

		// Das Ergebnis wird zurückgegeben
		return returnV
	})
}

func (o *RPCModule) _kernel_rpc_IsShareRemote(kernel types.KernelInterface) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(kernel.ContextV8().Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Es wird versucht die Tabelle abzurufen
		table, isok := kernel.GloablRegisterRead("rpc_public").(map[string]*SharedPublicFunction)
		if !isok {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			kernel.KernelThrow(info.Context(), "some type error 2")

			// Der Vorgang wird beendet
			return nil
		}

		// Es wird geprüft ob die Funktion bereits Registriert wurde
		_, found := table[info.Args()[0].String()]

		// Das Ergebniss wird erstellt
		returnV, err := v8.NewValue(kernel.ContextV8().Isolate(), found)
		if err != nil {
			panic(err)
		}

		// Das Ergebnis wird zurückgegeben
		return returnV
	})
}

func (o *RPCModule) Init(kernel types.KernelInterface) error {
	// Es wird versucht ein Global Register Eintrag zu erzeugen
	slfmap := make(map[string]*SharedLocalFunction)
	if err := kernel.GloablRegisterWrite("rpc_local", slfmap); err != nil {
		return fmt.Errorf("")
	}
	spfmap := make(map[string]*SharedPublicFunction)
	if err := kernel.GloablRegisterWrite("rpc_public", spfmap); err != nil {
		return fmt.Errorf("")
	}

	// Die RPC (Remote Function Call) funktionen werden bereitgestellt
	rpc := v8.NewObjectTemplate(kernel.Isolate())
	rpc.Set("CallLocal", o._kernel_rpc_CallLocal(kernel), v8.ReadOnly)
	rpc.Set("CallRemote", o._kernel_rpc_CallRemote(kernel), v8.ReadOnly)
	rpc.Set("IsShareLocal", o._kernel_rpc_IsShareLocal(kernel), v8.ReadOnly)
	rpc.Set("IsShareRemote", o._kernel_rpc_IsShareRemote(kernel), v8.ReadOnly)
	rpc.Set("NewShareLocal", o._kernel_rpc_NewShareLocalFunction(kernel), v8.ReadOnly)
	rpc.Set("NewSharePublic", o._kernel_rpc_NewSharePublicFunction(kernel), v8.ReadOnly)
	rpc.Set("GetFunctionDetailsLocal", o._kernel_rpc_GetFunctionDetailsLocal(kernel), v8.ReadOnly)
	rpc.Set("GetFunctionDetailsRemote", o._kernel_rpc_GetFunctionDetailsRemote(kernel), v8.ReadOnly)

	// Das RPC Objekt wird final erzeugt
	rpcObj, err := rpc.NewInstance(kernel.ContextV8())
	if err != nil {
		return fmt.Errorf("Kernel->_new_kernel_load_rpc_module: " + err.Error())
	}

	// Das RFC Modul wird hinzugefügt
	kernel.Global().Set("rpc", rpcObj)

	// Kein Fehler
	return nil
}

func (o *RPCModule) GetName() string {
	return "rpc"
}

func NewRPCModule() *RPCModule {
	return new(RPCModule)
}
