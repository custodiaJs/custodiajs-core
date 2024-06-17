package kernel

import (
	"fmt"

	"github.com/CustodiaJS/custodiajs-core/kernelmodules/extmodules"
	cgowrapper "github.com/CustodiaJS/custodiajs-core/kernelmodules/extmodules/cgo_wrapper"
	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/utils"

	v8 "rogchap.com/v8go"
)

type ExtModuleLink struct {
	exteModuleLink *extmodules.ExternalModule
	name           string
}

func (o *ExtModuleLink) addGlobalFunc(extModFunc *cgowrapper.CGOWrappedLibModuleFunction, iso *v8.Isolate, context *v8.Context) error {
	// Die Funktion wird erzeugt
	funcTemplate := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Die Einzelnen Parameter werden ausgewertet
		for _, item := range info.Args() {
			// Das Datentyp wird ermittelt
			var responseData *types.FunctionCallReturnData
			if item == nil {
				responseData = &types.FunctionCallReturnData{Type: "null", Value: nil}
			} else if item.IsUndefined() || item.IsNull() {
				responseData = &types.FunctionCallReturnData{Type: "undefined", Value: nil}
			} else {
				switch {
				case item.IsString():
					responseData = &types.FunctionCallReturnData{Type: "string", Value: item.String()}
				case item.IsNumber():
					switch {
					case item.IsBigInt():
						responseData = &types.FunctionCallReturnData{Type: "number", Value: item.BigInt().String()}
					case item.IsInt32():
						responseData = &types.FunctionCallReturnData{Type: "number", Value: item.Int32()}
					case item.IsUint32():
						responseData = &types.FunctionCallReturnData{Type: "number", Value: item.Uint32()}
					case item.IsNumber():
						responseData = &types.FunctionCallReturnData{Type: "number", Value: item.Number()}
					default:
						responseData = &types.FunctionCallReturnData{Type: "number", Value: item.Integer()}
					}
				case item.IsBoolean():
					responseData = &types.FunctionCallReturnData{Type: "boolean", Value: item.Boolean()}
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
			_ = responseData
		}

		// Die CGO Wrapped Funktion wird aufgerufen
		res, err := extModFunc.Call()
		if err != nil {
			// Es wird geprüft ob es sich um CGO Panic Error handelt
			switch err.(type) {
			case *types.ExtModCGOPanic:
				// Der VM wird Signalisiert das ein Fehler aufgetreten ist
				utils.V8ContextThrow(info.Context(), "cgo panic by call function from external module")

				// Der Vorgang wird abgebrochen
				return nil
			default:
				// Der VM wird Signalisiert das ein Fehler aufgetreten ist
				utils.V8ContextThrow(info.Context(), err.Error())

				// Der Vorgang wird abgebrochen
				return nil
			}
		}
		_ = res // NOT IMPLEM

		//fmt.Println("A: ", res)

		// Es gibt keine Werte welche zurückgegeben werden können
		return nil
	})

	// Die Funktion wird hinzugefügt
	funcObj := funcTemplate.GetFunction(context)
	if err := context.Global().Set(extModFunc.GetName(), funcObj); err != nil {
		panic(err)
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

func (o *ExtModuleLink) addGlobalImport(extModImport *extmodules.ExternModuleImport, kernel types.KernelInterface) error {
	_ = extModImport
	_ = kernel
	// Es ist kein Fehler aufgetreten
	return nil
}

func (o *ExtModuleLink) addGlobalObject(extModObject *extmodules.ExternModuleObject, kernel types.KernelInterface) error {
	// Es ist kein Fehler aufgetreten
	_ = extModObject
	_ = kernel
	return nil
}

func (o *ExtModuleLink) addEvent(extModEvent *extmodules.ExternModuleEvent, kernel types.KernelInterface) error {
	// Es ist kein Fehler aufgetreten
	_ = extModEvent
	_ = kernel
	return nil
}

func (o *ExtModuleLink) Init(kernel types.KernelInterface, iso *v8.Isolate, context *v8.Context) error {
	// Die Globalen Funktionen werden Exportiert
	for _, item := range o.exteModuleLink.GetGlobalFunctions() {
		if err := o.addGlobalFunc(item, iso, context); err != nil {
			return fmt.Errorf("ExtModuleLink->Init: " + err.Error())
		}
	}

	// Die Imports werden verfügbar gemacht
	for _, item := range o.exteModuleLink.GetImports() {
		if err := o.addGlobalImport(item, kernel); err != nil {
			return fmt.Errorf("ExtModuleLink->Init: " + err.Error())
		}
	}

	// Die Globalen Objekte werden verfügbar gemacht
	for _, item := range o.exteModuleLink.GetGlobalObjects() {
		if err := o.addGlobalObject(item, kernel); err != nil {
			return fmt.Errorf("ExtModuleLink->Init: " + err.Error())
		}
	}

	// Die Events, auf welche die Lib reagieren soll, werden verfügbar gemacht
	for _, item := range o.exteModuleLink.GetEventTriggers() {
		if err := o.addEvent(item, kernel); err != nil {
			return fmt.Errorf("ExtModuleLink->Init: " + err.Error())
		}
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

func (o *ExtModuleLink) OnlyForMain() bool {
	return false
}

func (o *ExtModuleLink) GetName() string {
	return o.name
}

func LinkWithExternalModule(extmodule *extmodules.ExternalModule) (*ExtModuleLink, error) {
	// Das Modul wird erzeugt
	vat := &ExtModuleLink{name: extmodule.GetName(), exteModuleLink: extmodule}

	// Das Modul wird zurückgegeben
	return vat, nil
}
