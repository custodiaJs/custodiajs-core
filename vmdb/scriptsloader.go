package vmdb

import (
	"fmt"
	"vnh1/static"
)

type NodeJsModule struct {
}

func tryToLoadNodeJsModules(path string) ([]*NodeJsModule, error) {
	// Es wird geprüft ob es sich um einen gültigen Path handelt
	if !static.FolderExists(path) {
		return nil, fmt.Errorf("tryToLoadNodeJsModule: no nodejs modules folder found")
	}

	// Es wird geprüft ob die 'package.json' Datei vorhanden ist
	if !static.FileExists("package.json") {
		return nil, fmt.Errorf("tryToLoadNodeJsModules: isnt a nodejs module")
	}

	// Es wird eine Übersicht über den Ordnerinhalt erstellt

	// Die NodeJS Module werden zurückgegeben
	return nodeJsModules, nil
}
