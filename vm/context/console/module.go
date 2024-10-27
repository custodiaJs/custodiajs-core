package contextconsole

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/CustodiaJS/custodiajs-core/global/types"
	"github.com/CustodiaJS/custodiajs-core/global/utils"

	v8 "rogchap.com/v8go"
)

type ModuleConsole int

func (o *ModuleConsole) _kernel_console_log(kernel types.KernelInterface, iso *v8.Isolate, _ *v8.Context) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Es werden alle Stringwerte Extrahiert
		extracted := []string{}
		for _, item := range info.Args() {
			if item.IsObject() && !item.IsArray() {
				obj, err := utils.V8ObjectToGoObject(info.Context(), item)
				if err != nil {
					utils.V8ContextThrow(info.Context(), "internal error by converting, value")
					return nil
				}
				encoded, err := json.Marshal(obj)
				if err != nil {
					utils.V8ContextThrow(info.Context(), "internal error by converting, value")
					return nil
				}
				extracted = append(extracted, string(encoded))
			} else if item.IsArray() {
				obj, err := utils.V8ArrayToGoArray(info.Context(), item)
				if err != nil {
					utils.V8ContextThrow(info.Context(), "internal error by converting, value")
					return nil
				}

				var extra []interface{}
				for _, item := range obj {
					extra = append(extra, item.Value)
				}

				encoded, err := json.Marshal(extra)
				if err != nil {
					utils.V8ContextThrow(info.Context(), "internal error by converting, value")
					return nil
				}

				extracted = append(extracted, string(encoded))
			} else if item.IsFunction() {
				if item.IsAsyncFunction() {
					extracted = append(extracted, fmt.Sprintf("ASYNC:=%p", item))
				} else {
					extracted = append(extracted, fmt.Sprintf("SYNC:=%p", item))
				}
			} else {
				extracted = append(extracted, item.String())
			}
		}

		// Es wird ein String aus der Ausgabe erzeugt
		outputStr := strings.Join(extracted, " ")

		// Die Ausgabe wird an den Console Cache übergeben
		kernel.Console().InfoLog(outputStr)
		//kernel.LogPrint("Console", "%s\n", outputStr)

		// Rückgabe ohne Fehler
		return nil
	})
}

func (o *ModuleConsole) _kernel_console_error(kernel types.KernelInterface, iso *v8.Isolate, _ *v8.Context) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Es werden alle Stringwerte Extrahiert
		extracted := []string{}
		for _, item := range info.Args() {
			if item.IsObject() && !item.IsArray() {
				obj, err := utils.V8ObjectToGoObject(info.Context(), item)
				if err != nil {
					utils.V8ContextThrow(info.Context(), "internal error by converting, value")
					return nil
				}
				encoded, err := json.Marshal(obj)
				if err != nil {
					utils.V8ContextThrow(info.Context(), "internal error by converting, value")
					return nil
				}
				extracted = append(extracted, string(encoded))
			} else if item.IsArray() {
				obj, err := utils.V8ArrayToGoArray(info.Context(), item)
				if err != nil {
					utils.V8ContextThrow(info.Context(), "internal error by converting, value")
					return nil
				}

				var extra []interface{}
				for _, item := range obj {
					extra = append(extra, item.Value)
				}

				encoded, err := json.Marshal(extra)
				if err != nil {
					utils.V8ContextThrow(info.Context(), "internal error by converting, value")
					return nil
				}

				extracted = append(extracted, string(encoded))
			} else if item.IsFunction() {
				if item.IsAsyncFunction() {
					extracted = append(extracted, fmt.Sprintf("ASYNC:=%p", item))
				} else {
					extracted = append(extracted, fmt.Sprintf("SYNC:=%p", item))
				}
			} else {
				extracted = append(extracted, item.String())
			}
		}

		// Es wird ein String aus der Ausgabe erzeugt
		outputStr := strings.Join(extracted, " ")

		// Die Ausgabe wird an den Console Cache übergeben
		kernel.Console().ErrorLog(outputStr)
		//kernel.LogPrint("Console", "ERROR: %s\n", outputStr)

		// Rückgabe ohne Fehler
		return nil
	})
}

func (o *ModuleConsole) Init(kernel types.KernelInterface, iso *v8.Isolate, context *v8.Context) error {
	// Das Consolen Objekt wird definiert
	console := v8.NewObjectTemplate(iso)
	console.Set("log", o._kernel_console_log(kernel, iso, context), v8.ReadOnly)
	console.Set("error", o._kernel_console_error(kernel, iso, context), v8.ReadOnly)

	// Das Console Objekt wird final erzeugt
	consoleObj, err := console.NewInstance(context)
	if err != nil {
		return fmt.Errorf("Kernel->_new_kernel_load_console_module: " + err.Error())
	}

	// Die Consolen Funktionen werden hinzugefügt
	context.Global().Set("console", consoleObj)

	// Kein Fehler
	return nil
}

func (o *ModuleConsole) OnlyForMain() bool {
	return false
}

func (o *ModuleConsole) GetName() string {
	return "console"
}

func NewConsoleModule() *ModuleConsole {
	return new(ModuleConsole)
}
