// Author: fluffelpuff
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package kmodulerpc

import (
	"fmt"

	"github.com/CustodiaJS/custodiajs-core/global/static"
	"github.com/CustodiaJS/custodiajs-core/global/types"
	"github.com/CustodiaJS/custodiajs-core/global/utils"

	v8 "rogchap.com/v8go"
)

type RPCModule int

// Versucht eine Lokale RPC Funktion abzurufen
func __determineRPCFunction(kernel types.KernelInterface, funcsig *types.FunctionSignature) (types.SharedFunctionInterface, bool, error) {
	// Sollte keine VM ID / VM Name angegeben sein, oder die ID bzw der Name stimmen mit der Aktuellen VM überein, wird die Funktion in der Aktuellen VM gesucht
	if (funcsig.VMID == "" && funcsig.VMName == "") || funcsig.VMID == string(kernel.GetFingerprint()) || funcsig.VMName == kernel.AsVmInstance().GetManifest().Name {
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
	var vmiface types.VmInterface
	var vmFound bool
	var terr error
	switch {
	case funcsig.VMName == "" && funcsig.VMID != "":
		vmiface, vmFound, terr = kernel.GetCore().GetVmByID(string(funcsig.VMID), nil)
	case funcsig.VMName != "" && funcsig.VMID == "":
		vmiface, _, terr = kernel.GetCore().GetVmByName(funcsig.VMName, nil)
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
func (o *RPCModule) __rpcNewShareFunction(addPublic bool, kernel types.KernelInterface, context *v8.Context) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(context.Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Es wird geprüft ob 2 Argumente vorhanden sind
		if len(info.Args()) != 2 {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), "You have to explicitly specify 2 parameters")

			// Der Vorgang wird beendet
			return v8.Undefined(info.Context().Isolate())
		}

		// Es wird geprüft ob das Erste Argument ein String ist
		if !info.Args()[0].IsString() {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), "The first argument must be a string")

			// Der Vorgang wird beendet
			return v8.Undefined(info.Context().Isolate())
		}

		// Es wird geprüft ob das Zweite Argument eine Funktion ist
		if !info.Args()[1].IsAsyncFunction() {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), "The second argument must be an asynchronous function")

			// Der Vorgang wird beendet
			return v8.Undefined(info.Context().Isolate())
		}

		// Die Funktionssignatur wird eingegelesen
		funcSig, err := utils.ParseFunctionSignature(info.Args()[0].String())
		if err != nil {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), "Error creating the function signature")

			// Der Vorgang wird beendet
			return v8.Undefined(info.Context().Isolate())
		}

		// Die Funktion wird ausgelesen
		rpcFunc, err := info.Args()[1].AsFunction()
		if err != nil {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), "Error reading shared function, it is a JS Engine error")

			// Der Vorgang wird beendet
			return v8.Undefined(info.Context().Isolate())
		}

		// Es wird geprüft ob es sich um einen Zulässigen Funktionsnamen handelt
		if !utils.ValidateFunctionName(funcSig.FunctionName) {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), fmt.Sprintf("'%s' is not a legal function name, the syntax is not allowed", funcSig.FunctionName))

			// Der Vorgang wird beendet
			return v8.Undefined(info.Context().Isolate())
		}

		// Der Rückgabetype wird geprüft
		if !utils.ValidateDatatypeString(funcSig.ReturnType) && funcSig.ReturnType != "" {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), fmt.Sprintf("The return data type '%s' is not a valid data type", funcSig.ReturnType))

			// Der Vorgang wird beendet
			return v8.Undefined(info.Context().Isolate())
		}

		// Die Einzelnen Parametertypen werden geprüft
		for index, item := range funcSig.Params {
			// Der Einzelne Parameter wird geprüft
			if !utils.ValidateDatatypeString(item) {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				utils.V8ContextThrow(info.Context(), fmt.Sprintf("At the %d. Argument is not a valid data type", index))

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
			utils.V8ContextThrow(info.Context(), "An internal error occurred; the RPC sharing table could not be read")

			// Der Vorgang wird beendet
			return v8.Undefined(info.Context().Isolate())
		}

		// Es wird geprüft ob diese Funktion bereits registriert wurde
		if _, found := resolveData[utils.FunctionOnlySignatureString(funcSig)]; found {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), fmt.Sprintf("The function '%s' has already been shared", funcSig.FunctionName))

			// Der Vorgang wird beendet
			return v8.Undefined(info.Context().Isolate())
		}

		// Das Objekt wird erzeugt
		newSharedFunction := &SharedFunction{
			kernel:     kernel,
			name:       info.Args()[0].String(),
			v8Function: rpcFunc,
			signature:  funcSig,
		}

		// Die Geteilte Funktion wird erzeugt und abgespeichert
		resolveData[utils.FunctionOnlySignatureString(funcSig)] = newSharedFunction

		// Der Eintrag in der Datenbank wird geupdated
		if err := kernel.GloablRegisterWrite("rpc", resolveData); err != nil {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), "An error occurred while attempting to write to the RPC table")

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
				utils.V8ContextThrow(info.Context(), "An internal error occurred; the public RPC sharing table could not be read")

				// Der Vorgang wird beendet
				return v8.Undefined(info.Context().Isolate())
			}

			// Es wird geprüft ob diese Funktion bereits registriert wurde
			if _, found := resolvedPublicData[utils.FunctionOnlySignatureString(funcSig)]; found {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				utils.V8ContextThrow(info.Context(), fmt.Sprintf("The '%s' function has already been shared publicly", funcSig.FunctionName))

				// Der Vorgang wird beendet
				return v8.Undefined(info.Context().Isolate())
			}

			// Die Geteilte Funktion wird erzeugt und abgespeichert
			resolvedPublicData[utils.FunctionOnlySignatureString(funcSig)] = newSharedFunction

			// Der Eintrag in der Datenbank wird geupdated
			if err := kernel.GloablRegisterWrite("rpc_public", resolvedPublicData); err != nil {
				// Die Fehlermeldung wird erstellt und an JS zurückgegeben
				utils.V8ContextThrow(info.Context(), "An error occurred while attempting to write to the public RPC table")

				// Der Vorgang wird beendet
				return v8.Undefined(info.Context().Isolate())
			}

			// Es wird ein Kernel Signal verwendet um zu Signalisieren das eine neue geteilte RPC Funktion vorhanden ist
			kernel.Signal("rpc/global/add_share", funcSig)
		} else {
			// Es wird ein Kernel Signal verwendet um zu Signalisieren das eine neue geteilte RPC Funktion vorhanden ist
			kernel.Signal("rpc/local/add_share", funcSig)
		}

		// LOG
		kernel.LogPrint("RPC", "Add sharing function '%s'\n", funcSig.FunctionName)

		// Der Vorgang wurde ohne Fehler durchgeführt
		return v8.Undefined(info.Context().Isolate())
	})
}

// Registriert eine neue Javascript Funktion als Loakl geteilte Funktion
func (o *RPCModule) rpcNewShareLocal(kernel types.KernelInterface, context *v8.Context) *v8.FunctionTemplate {
	return o.__rpcNewShareFunction(false, kernel, context)
}

// Registriert eine neue Javascript Funktion als öffentlich geteilte Funktion
func (o *RPCModule) rpcNewSharePublic(kernel types.KernelInterface, context *v8.Context) *v8.FunctionTemplate {
	return o.__rpcNewShareFunction(true, kernel, context)
}

// Ermöglicht es eine Lokale Javascript Funktion aufzurufen
func (o *RPCModule) rpcCall(kernel types.KernelInterface, iso *v8.Isolate) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Es wird geprüft ob mindestens 3 Parameter vorhanden sind
		if len(info.Args()) < 3 {
			// Es wird Javascript Fehler ausgelöst
			utils.V8ContextThrow(info.Context(), "You must provide at least 3 arguments. Argument 1 must be the function signature and possibly the target, argument 2 must be either a Config object or null, argument 3 must be an array, this array must contain the data to be transferred")

			// Rückgabe
			return v8.Undefined(info.Context().Isolate())
		}

		// Die Parameterreihenfolge wird geprüft
		if !info.Args()[0].IsString() {
			// Es wird Javascript Fehler ausgelöst
			utils.V8ContextThrow(info.Context(), "The first argument is not a string")

			// Rückgabe
			return v8.Undefined(info.Context().Isolate())
		}
		if !info.Args()[1].IsObject() {
			// Es wird Javascript Fehler ausgelöst
			utils.V8ContextThrow(info.Context(), "The second argument is not an object")

			// Rückgabe
			return v8.Undefined(info.Context().Isolate())
		}
		if !info.Args()[2].IsArray() {
			// Es wird Javascript Fehler ausgelöst
			utils.V8ContextThrow(info.Context(), "The third argument is not an array")

			// Rückgabe
			return v8.Undefined(info.Context().Isolate())
		}

		// Speichert die Container ID ab
		rpcSignatureStr := info.Args()[0].String()

		// Erstelle einen Promise Resolver
		resolver, err := v8.NewPromiseResolver(info.Context())
		if err != nil {
			// Es wird Javascript Fehler ausgelöst
			utils.V8ContextThrow(info.Context(), "Error attempting to create a promise 'rpcCall'")

			// Rückgabe
			return v8.Undefined(info.Context().Isolate())
		}

		// Gebe die Promise des Resolvers zurück
		//promise := resolver.GetPromise()

		// Es wird versucht die Signatur einzulesen
		funcSig, err := utils.ParseFunctionSignatureOptionalFunction(rpcSignatureStr)
		if err != nil {
			// Es wird Javascript Fehler ausgelöst
			utils.V8ContextThrow(info.Context(), fmt.Sprintf("'%s' is not a valid function signature", rpcSignatureStr))

			// Rückgabe
			return v8.Undefined(info.Context().Isolate())
		}

		// Es wird versucht das Konfigurationsobjekt auszulesen
		configObj, err := utils.V8ObjectToGoObject(info.Context(), info.Args()[1])
		if err != nil {
			// Es wird Javascript Fehler ausgelöst
			utils.V8ContextThrow(info.Context(), "An engine error occurred and the object could not be converted")

			// Rückgabe
			return v8.Undefined(info.Context().Isolate())
		}
		_ = configObj

		// Alle weiteren Werte werden ausgegegebn
		exportedParameters := make([]*types.FunctionParameterCapsle, 0)
		for place, item := range info.Args()[3:] {
			// Es wird versucht den Wert in einen Zulässigen Golang Wert umzuwandeln
			convertedFunctionParameterGoValue, err := utils.V8ValueToGoValue(info.Context(), item)
			if err != nil {
				// Es wird Javascript Fehler ausgelöst
				utils.V8ContextThrow(info.Context(), fmt.Sprintf("The %d entry in the array is not a permissible value that can be transferred using RPC", place))

				// Rückgabe
				return v8.Undefined(info.Context().Isolate())
			}

			// Der Parameter wird zwischengespeichert
			exportedParameters = append(exportedParameters, &types.FunctionParameterCapsle{Value: convertedFunctionParameterGoValue.Value, CType: convertedFunctionParameterGoValue.Type})
		}

		// Die Funktion wird ermittelt
		_, foundFunction, err := __determineRPCFunction(kernel, funcSig)
		if err != nil {
			// Es wird Javascript Fehler ausgelöst
			utils.V8ContextThrow(info.Context(), err.Error())

			// Rückgabe
			return v8.Undefined(info.Context().Isolate())
		}

		// Es wird geprüft ob die Funktion gefunden wurde
		if !foundFunction {
			// Es wird Javascript Fehler ausgelöst
			utils.V8ContextThrow(info.Context(), fmt.Sprintf("The function '%s' could not be determined", funcSig.FunctionName))

			// Rückgabe
			return v8.Undefined(info.Context().Isolate())
		}

		panic("not implemantated")

		// Rückgabe
		return resolver.Value
	})
}

// Gibt Details über eine Lokal geteilte Funktion in Javascript zurück
func (o *RPCModule) rpcGetDetails(_ types.KernelInterface, iso *v8.Isolate) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		return nil
	})
}

// Gibt an ob es sich um eine geteilte Lokale Funktion handelt
func (o *RPCModule) rpcIsShar(kernel types.KernelInterface, iso *v8.Isolate) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Es wird geprüft ob mindestens 1 Argument angegeben wurde
		if len(info.Args()) < 1 {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), "The 'rpcIsShare' function requires at least 1 argument")

			// Der Vorgang wird beendet
			return nil
		}

		// Es wird geprüft das nicht mehr als 999.999.999 Einträge auf dem Stack liegen
		if len(info.Args()) >= 999999999 {
			// Die Fehlermeldung wird erstellt und an JS zurückgegeben
			utils.V8ContextThrow(info.Context(), "You cannot specify more than 999999999 arguments")

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
