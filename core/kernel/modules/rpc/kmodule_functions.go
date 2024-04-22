package kmodulerpc

import (
	"fmt"
	"vnh1/static"
	"vnh1/types"
	"vnh1/utils"

	v8 "rogchap.com/v8go"
)

type RPCModule int

// Versucht eine Lokale RPC Funktion abzurufen
func __determineRPCFunction(kernel types.KernelInterface, funcsig *types.FunctionSignature) (types.SharedFunctionInterface, bool, error) {
	// Sollte keine VM ID / VM Name angegeben sein, oder die ID bzw der Name stimmen mit der Aktuellen VM überein, wird die Funktion in der Aktuellen VM gesucht
	if (funcsig.VMID == "" && funcsig.VMName == "") || funcsig.VMID == string(kernel.AsCoreVM().GetFingerprint()) || funcsig.VMName == kernel.AsCoreVM().GetVMName() {
		// Es wird versucht die Lokale Kernel Tabelle abzurufen
		table, isok := kernel.GloablRegisterRead("rpc").(map[string]types.SharedFunctionInterface)
		if !isok {
			return nil, false, fmt.Errorf("rpc table reading error")
		}

		// Es wird geprüft ob die Funktion innerhalb des Aktuellen Kernels Registriert wurde
		result, found := table[utils.FunctionOnlySignatureString(funcsig)]
		if !found {
			return nil, false, nil
		}

		// Rückgabe
		return result, true, nil
	}

	// Sollte eine VM ID angegeben wurden sein, wird versucht die VM anhand ihrer ID zu ermitteln
	var vmiface types.CoreVMInterface
	var vmFound bool
	var terr error
	switch {
	case funcsig.VMName == "" && funcsig.VMID != "":
		vmiface, vmFound, terr = kernel.GetCore().GetScriptContainerVMByID(string(funcsig.VMID))
	case funcsig.VMName != "" && funcsig.VMID == "":
		vmiface, terr = kernel.GetCore().GetScriptContainerByVMName(funcsig.VMName)
	default:
		terr = fmt.Errorf("crazy internal error")
	}

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if terr != nil {
		return nil, false, terr
	}

	// Es wird geprüft ob die VM ermittelt werden konnte
	// sollte die passende VM Lokal gefunden werden, wird versucht die Passende Funktion zu ermittlen
	if vmFound {
		// Es versucht die Funktion zu ermitteln (früher mittels vmiface.GetAllSharedFunctions())
		res, foundResult, err := vmiface.GetSharedFunctionBySignature(static.LOCAL, funcsig)
		if err != nil {
			return nil, false, fmt.Errorf("__determineRPCFunction: " + err.Error())
		}

		// Es wird geprüft ob ein Erebniss gefunden
		if !foundResult {
			return nil, false, nil
		}

		// Das Ergebniss wird zwischengespeichert
		return res, true, nil
	}

	// Es wird im Nodestack nach Verbindungen zu anderen Hosts geschaut
	// es wird geschaut ob es einen Host mit passender VM Id gibt,
	// sollte ein Passender Host ermittelt werden, wird der Funktionsauruf an diesen Übergeben
	// !! NOCHT NICHT IMPLEMENTIERT !!
	return nil, false, fmt.Errorf("unkown vm")
}

// Registriert eine neue Javascript Funktion als geteilt
func (o *RPCModule) __rpcNewShareFunction(addPublic bool, kernel types.KernelInterface, iso *v8.Isolate, _ *v8.Context) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Es wird geprüft ob 2 Argumente vorhanden sind
		if len(info.Args()) != 2 {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), "length of prameters invalid")

			// Der Vorgang wird beendet
			return v8.Undefined(info.Context().Isolate())
		}

		// Es wird geprüft ob das Erste Argument ein String ist
		if !info.Args()[0].IsString() {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), "the first argument isn't a string")

			// Der Vorgang wird beendet
			return v8.Undefined(info.Context().Isolate())
		}

		// Es wird geprüft ob das Zweite Argument eine Funktion ist
		if !info.Args()[1].IsAsyncFunction() {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), "the thir parameter isn't async function")

			// Der Vorgang wird beendet
			return v8.Undefined(info.Context().Isolate())
		}

		// Die Funktionssignatur wird eingegelesen
		funcSig, err := utils.ParseFunctionSignature(info.Args()[0].String())
		if err != nil {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), "invalid function signature")

			// Der Vorgang wird beendet
			return v8.Undefined(info.Context().Isolate())
		}

		// Die Funktion wird ausgelesen
		rpcFunc, err := info.Args()[1].AsFunction()
		if err != nil {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), "internal engine error")

			// Der Vorgang wird beendet
			return v8.Undefined(info.Context().Isolate())
		}

		// Der Quellcode der Funktion wird extrahiert
		functionSourceCode := rpcFunc.String()

		// Es wird geprüft ob es sich um einen Zulässigen Funktionsnamen handelt
		if !utils.ValidateFunctionName(funcSig.FunctionName) {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), "invalid function name")

			// Der Vorgang wird beendet
			return v8.Undefined(info.Context().Isolate())
		}

		// Der Rückgabetype wird geprüft
		if !utils.ValidateDatatypeString(funcSig.ReturnType) && funcSig.ReturnType != "" {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), fmt.Sprintf("invalid return type '%s'", funcSig.ReturnType))

			// Der Vorgang wird beendet
			return v8.Undefined(info.Context().Isolate())
		}

		// Die Einzelnen Parametertypen werden geprüft
		for index, item := range funcSig.Params {
			// Der Einzelne Parameter wird geprüft
			if !utils.ValidateDatatypeString(item) {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				fmt.Println(item)
				utils.V8ContextThrow(info.Context(), fmt.Sprintf("Failed to get element at index %d: %v", index, "invalid return datatype"))

				// Der Vorgang wird beendet
				return v8.Undefined(info.Context().Isolate())
			}
		}

		// Der Eintrag wird aus der Tabelle abgerufen
		rpcTable := kernel.GloablRegisterRead("rpc")

		// Es wird versucht die Tabelle abzurufen
		resolveData, isOk := rpcTable.(map[string]types.SharedFunctionInterface)
		if !isOk {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), "some type error 1")

			// Der Vorgang wird beendet
			return v8.Undefined(info.Context().Isolate())
		}

		// Es wird geprüft ob diese Funktion bereits registriert wurde
		if _, found := resolveData[utils.FunctionOnlySignatureString(funcSig)]; found {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), "some type error")

			// Der Vorgang wird beendet
			return v8.Undefined(info.Context().Isolate())
		}

		// Das Objekt wird erzeugt
		newSharedFunction := &SharedFunction{
			kernel:             kernel,
			name:               info.Args()[0].String(),
			parmTypes:          funcSig.Params,
			returnType:         funcSig.ReturnType,
			functionSourceCode: functionSourceCode,
		}

		// Die Geteilte Funktion wird erzeugt und abgespeichert
		resolveData[utils.FunctionOnlySignatureString(funcSig)] = newSharedFunction

		// Der Eintrag in der Datenbank wird geupdated
		if err := kernel.GloablRegisterWrite("rpc", resolveData); err != nil {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), "register writing error")

			// Der Vorgang wird beendet
			return v8.Undefined(info.Context().Isolate())
		}

		// Es wird geprüft ob der Eintrag auch in der Öffentlichen Datenbank hinzugefügt werden soll
		if addPublic {
			// Der Eintrag wird aus der Tabelle abgerufen
			rpcTable := kernel.GloablRegisterRead("rpc_public")

			// Es wird versucht die Tabelle abzurufen
			resolvedPublicData, isOk := rpcTable.(map[string]types.SharedFunctionInterface)
			if !isOk {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				utils.V8ContextThrow(info.Context(), "some type error 2")

				// Der Vorgang wird beendet
				return v8.Undefined(info.Context().Isolate())
			}

			// Es wird geprüft ob diese Funktion bereits registriert wurde
			if _, found := resolvedPublicData[utils.FunctionOnlySignatureString(funcSig)]; found {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				utils.V8ContextThrow(info.Context(), "some type error 3")

				// Der Vorgang wird beendet
				return v8.Undefined(info.Context().Isolate())
			}

			// Die Geteilte Funktion wird erzeugt und abgespeichert
			resolvedPublicData[utils.FunctionOnlySignatureString(funcSig)] = newSharedFunction

			// Der Eintrag in der Datenbank wird geupdated
			if err := kernel.GloablRegisterWrite("rpc_public", resolvedPublicData); err != nil {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				utils.V8ContextThrow(info.Context(), "register writing error")

				// Der Vorgang wird beendet
				return v8.Undefined(info.Context().Isolate())
			}
		}

		// LOG
		kernel.LogPrint("rpc", "Add sharing function '%s'\n", funcSig.FunctionName)

		// Der Vorgang wurde ohne Fehler durchgeführt
		return v8.Undefined(info.Context().Isolate())
	})
}

// Registriert eine neue Javascript Funktion als Loakl geteilte Funktion
func (o *RPCModule) rpcNewShareLocal(kernel types.KernelInterface, iso *v8.Isolate, context *v8.Context) *v8.FunctionTemplate {
	return o.__rpcNewShareFunction(false, kernel, iso, context)
}

// Registriert eine neue Javascript Funktion als öffentlich geteilte Funktion
func (o *RPCModule) rpcNewSharePublic(kernel types.KernelInterface, iso *v8.Isolate, context *v8.Context) *v8.FunctionTemplate {
	return o.__rpcNewShareFunction(true, kernel, iso, context)
}

// Ermöglicht es eine Lokale Javascript Funktion aufzurufen
func (o *RPCModule) rpcCall(kernel types.KernelInterface, iso *v8.Isolate, _ *v8.Context) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Erstelle einen Promise Resolver
		resolver, err := v8.NewPromiseResolver(info.Context())
		if err != nil {
			// Es wird Javascript Fehler ausgelöst
			utils.V8ContextThrow(info.Context(), "internal error, promise creating error")

			// Rückgabe
			return v8.Undefined(info.Context().Isolate())
		}

		// Gebe die Promise des Resolvers zurück
		promise := resolver.GetPromise()

		// Es wird geprüft ob mindestens 3 Parameter vorhanden sind
		if len(info.Args()) < 3 {
			// Es wird Javascript Fehler ausgelöst
			utils.V8ContextThrow(info.Context(), "invalid parameters, first must be the url of host and container, second must be a configuration and thrid must be a function name")

			// Rückgabe
			return v8.Undefined(info.Context().Isolate())
		}

		// Die Parameterreihenfolge wird geprüft
		if !info.Args()[0].IsString() {
			// Es wird Javascript Fehler ausgelöst
			utils.V8ContextThrow(info.Context(), "invalid parameters, first must be the url of host and container, second must be a configuration and thrid must be a function name")

			// Rückgabe
			return v8.Undefined(info.Context().Isolate())
		}
		if !info.Args()[1].IsObject() {
			// Es wird Javascript Fehler ausgelöst
			utils.V8ContextThrow(info.Context(), "invalid parameters, first must be the url of host and container, second must be a configuration and thrid must be a function name")

			// Rückgabe
			return v8.Undefined(info.Context().Isolate())
		}
		if !info.Args()[2].IsArray() {
			// Es wird Javascript Fehler ausgelöst
			utils.V8ContextThrow(info.Context(), "invalid parameters, first must be the url of host and container, second must be a configuration and thrid must be a function name")

			// Rückgabe
			return v8.Undefined(info.Context().Isolate())
		}

		// Speichert die Container ID ab
		rpcSignatureStr := info.Args()[0].String()

		// Es wird versucht die Signatur einzulesen
		funcSig, err := utils.ParseFunctionSignatureOptionalFunction(rpcSignatureStr)
		if err != nil {
			// Es wird Javascript Fehler ausgelöst
			utils.V8ContextThrow(info.Context(), "invalid function signature")

			// Rückgabe
			return v8.Undefined(info.Context().Isolate())
		}

		// Es wird versucht das Konfigurationsobjekt auszulesen
		configObj, err := utils.V8ObjectToGoObject(info.Context(), info.Args()[1])
		if err != nil {
			// Es wird Javascript Fehler ausgelöst
			utils.V8ContextThrow(info.Context(), "invalid function share owner container id")

			// Rückgabe
			return v8.Undefined(info.Context().Isolate())
		}
		_ = configObj

		// Alle weiteren Werte werden ausgegegebn
		exportedParameters := make([]*types.FunctionParameterCapsle, 0)
		for _, item := range info.Args()[3:] {
			// Es wird versucht den Wert in einen Zulässigen Golang Wert umzuwandeln
			convertedFunctionParameterGoValue, err := utils.V8ValueToGoValue(info.Context(), item)
			if err != nil {
				// Es wird Javascript Fehler ausgelöst
				utils.V8ContextThrow(info.Context(), "invalid parameters, first must be the url of host and container, second must be a configuration and thrid must be a function name")

				// Rückgabe
				return v8.Undefined(info.Context().Isolate())
			}

			// Der Parameter wird zwischengespeichert
			exportedParameters = append(exportedParameters, &types.FunctionParameterCapsle{Value: convertedFunctionParameterGoValue.Value, CType: convertedFunctionParameterGoValue.Type})
		}

		// Die Funktion wird ermittelt
		function, foundFunction, err := __determineRPCFunction(kernel, funcSig)
		if err != nil {
			// Es wird Javascript Fehler ausgelöst
			utils.V8ContextThrow(info.Context(), err.Error())

			// Rückgabe
			return v8.Undefined(info.Context().Isolate())
		}

		// Es wird geprüft ob die Funktion gefunden wurde
		if !foundFunction {
			// Es wird Javascript Fehler ausgelöst
			utils.V8ContextThrow(info.Context(), fmt.Sprintf("unkwon function %s", funcSig.FunctionName))

			// Rückgabe
			return v8.Undefined(info.Context().Isolate())
		}

		// Die Funktion wird ausgeführt
		go func() {
			// Die Funktion wird aufgerufen
			resultState, err := function.EnterFunctionCall(&types.RpcRequest{Parms: exportedParameters})
			if err != nil {
				val, _ := v8.NewValue(info.Context().Isolate(), fmt.Sprintf("function call throw:= %s", err.Error()))
				resolver.Reject(val)
				return
			}

			// Es wird geprüft ob ein Fehler aufgetreten ist
			if resultState.Error != "" {
				val, _ := v8.NewValue(info.Context().Isolate(), resultState.Error)
				resolver.Reject(val)
				return
			}

			// Es wird geprüft ob der Datensatz welcher zurückgegeben wurde, passend ist
			resolver.Resolve(v8.Undefined(info.Context().Isolate()))
		}()

		// Rückgabe
		return promise.Value
	})
}

// Gibt Details über eine Lokal geteilte Funktion in Javascript zurück
func (o *RPCModule) rpcGetDetails(_ types.KernelInterface, iso *v8.Isolate, _ *v8.Context) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		return nil
	})
}

// Gibt an ob es sich um eine geteilte Lokale Funktion handelt
func (o *RPCModule) rpcIsShar(kernel types.KernelInterface, iso *v8.Isolate, _ *v8.Context) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Es wird geprüft ob mindestens 1 Argument angegeben wurde
		if len(info.Args()) < 1 {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), "some type erro 1r")

			// Der Vorgang wird beendet
			return nil
		}

		// Es wird geprüft das nicht mehr als 999.999.999 Einträge auf dem Stack liegen
		if len(info.Args()) >= 999999999 {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), "to many arguments")

			// Der Vorgang wird beendet
			return nil
		}

		// Es wird geprüft ob es sich bei den Argumenten um Strings handelt
		results := make([]bool, 0)
		for _, item := range info.Args() {
			// Es wird geprüft ob es sich um einen String handelt
			if !item.IsString() {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				utils.V8ContextThrow(info.Context(), "only strings allowed, minimum = 1")

				// Der Vorgang wird beendet
				return nil
			}

			// Es wird versucht den Funktionspointer auszulesen
			funcDecodedPtr, err := utils.ParseFunctionSignature(item.String())
			if err != nil {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				utils.V8ContextThrow(info.Context(), err.Error())

				// Der Vorgang wird beendet
				return nil
			}

			// Es wird versucht die Funktion zu ermitteln
			_, found, err := __determineRPCFunction(kernel, funcDecodedPtr)
			if err != nil {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				utils.V8ContextThrow(info.Context(), err.Error())

				// Der Vorgang wird beendet
				return nil
			}

			// Das Ergebniss wird zwischengespeichert
			results = append(results, found)
		}

		// Es wird ein Exception ausgelöst sollte nicht mindestens 1 Wert auf dem Results Slice liegen
		if len(results) < 1 {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			fmt.Println(results)
			utils.V8ContextThrow(info.Context(), "unkown internal error")

			// Der Vorgang wird beendet
			return nil
		}

		// Es wird wir ermittel ob mehr als 1 Wert auf dem Slice liegt oder 1
		if len(results) == 1 {
			// Der Wert wird eingelesen
			returnV, err := v8.NewValue(iso, results[0])
			if err != nil {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				utils.V8ContextThrow(info.Context(), "internal v8 converting error")

				// Der Vorgang wird beendet
				return nil
			}

			// Der Wert wird zurückgegeben
			return returnV
		} else {
			// Der Wert wird eingelesen
			returnV, err := v8.NewValue(iso, results)
			if err != nil {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				utils.V8ContextThrow(info.Context(), "internal v8 converting error")

				// Der Vorgang wird beendet
				return nil
			}

			// Der Wert wird zurückgegeben
			return returnV
		}
	})
}
