package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/dop251/goja"
)

func main() {
	vm := goja.New() // Eine neue VM Instanz erstellen

	// Ein 'console' Objekt erstellen und dieses dem globalen Namespace hinzufügen
	console := vm.NewObject()
	logFunc := func(call goja.FunctionCall) goja.Value {
		fmt.Println(call.Argument(0).String())
		return nil
	}

	vnh1com := func(call goja.FunctionCall) goja.Value {
		fmt.Println(call.Argument(0).String())
		return nil
	}

	console.Set("log", logFunc)
	console.Set("log", logFunc)
	vm.Set("console", console)

	vnh1Obj := vm.NewObject()
	vnh1Obj.Set("com", vnh1com)
	vm.Set("vnh1", vnh1Obj)
	vm.Set("exports", vm.NewObject())

	// JavaScript-Code aus einer Datei laden
	jsCode, err := ioutil.ReadFile("index.js")
	if err != nil {
		log.Fatalf("Fehler beim Laden der JavaScript-Datei: %v", err)
	}

	// Den geladenen JavaScript-Code ausführen
	_, err = vm.RunString(string(jsCode))
	if err != nil {
		log.Fatalf("Fehler beim Ausführen des JavaScript-Codes: %v", err)
	}
}
