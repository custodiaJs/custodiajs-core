package kernel

import (
	"log"
	"strings"

	v8 "rogchap.com/v8go"
)

func (o *Kernel) _kernel_console_log() *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(o.Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Es werden alle Stringwerte Extrahiert
		extracted := []string{}
		for _, item := range info.Args() {
			extracted = append(extracted, item.String())
		}

		// Es wird ein String aus der Ausgabe erzeugt
		outputStr := strings.Join(extracted, " ")

		// Die Ausgabe wird an den Console Cache übergeben
		o.Console.InfoLog(outputStr)
		log.Println(outputStr)

		// Rückgabe ohne Fehler
		return nil
	})
}

func (o *Kernel) _kernel_console_error() *v8.FunctionTemplate {
	return v8.NewFunctionTemplate(o.Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
		// Es werden alle Stringwerte Extrahiert
		extracted := []string{}
		for _, item := range info.Args() {
			extracted = append(extracted, item.String())
		}

		// Es wird ein String aus der Ausgabe erzeugt
		outputStr := strings.Join(extracted, " ")

		// Die Ausgabe wird an den Console Cache übergeben
		o.Console.ErrorLog(outputStr)
		log.Printf("ERROR: %s\n", outputStr)

		// Rückgabe ohne Fehler
		return nil
	})
}
