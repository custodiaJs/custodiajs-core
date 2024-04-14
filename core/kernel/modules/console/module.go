package kmoduleconsole

import (
	"fmt"
	"log"
	"strings"
	"vnh1/types"

	v8 "rogchap.com/v8go"
)

type ModuleConsole int

func (o *ModuleConsole) _kernel_console_log(kernel types.KernelInterface) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(kernel.Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Es werden alle Stringwerte Extrahiert
		extracted := []string{}
		for _, item := range info.Args() {
			extracted = append(extracted, item.String())
		}

		// Es wird ein String aus der Ausgabe erzeugt
		outputStr := strings.Join(extracted, " ")

		// Die Ausgabe wird an den Console Cache übergeben
		kernel.Console().InfoLog(outputStr)
		log.Println(outputStr)

		// Rückgabe ohne Fehler
		return nil
	})
}

func (o *ModuleConsole) _kernel_console_error(kernel types.KernelInterface) *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(kernel.Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Es werden alle Stringwerte Extrahiert
		extracted := []string{}
		for _, item := range info.Args() {
			extracted = append(extracted, item.String())
		}

		// Es wird ein String aus der Ausgabe erzeugt
		outputStr := strings.Join(extracted, " ")

		// Die Ausgabe wird an den Console Cache übergeben
		kernel.Console().ErrorLog(outputStr)
		log.Printf("ERROR: %s\n", outputStr)

		// Rückgabe ohne Fehler
		return nil
	})
}

func (o *ModuleConsole) Init(kernel types.KernelInterface) error {
	// Das Consolen Objekt wird definiert
	console := v8.NewObjectTemplate(kernel.Isolate())
	console.Set("log", o._kernel_console_log(kernel), v8.ReadOnly)
	console.Set("error", o._kernel_console_error(kernel), v8.ReadOnly)

	// Das Console Objekt wird final erzeugt
	consoleObj, err := console.NewInstance(kernel.ContextV8())
	if err != nil {
		return fmt.Errorf("Kernel->_new_kernel_load_console_module: " + err.Error())
	}

	// Die Consolen Funktionen werden hinzugefügt
	kernel.Global().Set("console", consoleObj)

	// Kein Fehler
	return nil
}

func (o *ModuleConsole) GetName() string {
	return "console"
}

func NewConsoleModule() *ModuleConsole {
	return new(ModuleConsole)
}
