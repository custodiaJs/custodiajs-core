package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"vnh1/jsvm"
)

func main() {
	// Laden des Dateiinhalts
	content, err := ioutil.ReadFile("indes.js")
	if err != nil {
		log.Fatalf("Fehler beim Lesen der Datei: %v", err)
	}

	test, err := jsvm.NewVM(nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	test.RunScript(string(content))
}
