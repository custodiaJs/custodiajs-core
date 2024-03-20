package vmdb

import (
	"fmt"
	"path/filepath"
	"vnh1/static"
)

type PythonModule struct {
}

type NodeJsModule struct {
}

type ScriptsFolder struct {
	PythonModules []*PythonModule
	NodeJsModules []*NodeJsModule
}

func tryToLoadPythonModule(path string) (*PythonModule, error) {
	// Es wird eine übersicht über den Dateiinhalt erstellt
	fmt.Println("TRYLOAD: " + path)
	return &PythonModule{}, nil
}

func tryToLoadNodeJsModule(path string) (*NodeJsModule, error) {
	return &NodeJsModule{}, nil
}

func loadScripts(path string) (*ScriptsFolder, error) {
	// Es werden alle Python Module eingelesen, sofern diese Vorhanden sind
	pythonModules := []*PythonModule{}
	pyFolder := filepath.Join(path, "python3")
	if static.FolderExists(pyFolder) {
		// Es werden alle Ordner aufgelistet welche im Python Ordner vorhanden sind
		pythonModulesFolders, err := static.ListAllFolders(pyFolder)
		if err != nil {
			return nil, fmt.Errorf("loadScripts: " + err.Error())
		}

		// Es werden alle Python Module der reihenach geladen
		for _, item := range pythonModulesFolders {
			// Es wird versucht das Python Module zu laden
			newPyScript, err := tryToLoadPythonModule(item)
			if err != nil {
				return nil, fmt.Errorf("loadScripts: " + err.Error())
			}

			// Das Script wird zwischengespeichert
			pythonModules = append(pythonModules, newPyScript)
		}
	}

	// Es werden alle NodeJs Module eingelesen, sofern diese Vorhanden sind
	nodeJsModules := []*NodeJsModule{}
	nodeJsFolder := filepath.Join(path, "nodejs")
	if static.FolderExists(nodeJsFolder) {
		// Es werden alle Ordner aufgelistet welche im Python Ordner vorhanden sind
		nodejsModuleFolder, err := static.ListAllFolders(pyFolder)
		if err != nil {
			return nil, fmt.Errorf("loadScripts: " + err.Error())
		}

		// Es werden alle NodeJs Module der reihenach geladen
		for _, item := range nodejsModuleFolder {
			// Es wird versucht das Python Module zu laden
			newNodeJsModule, err := tryToLoadNodeJsModule(item)
			if err != nil {
				return nil, fmt.Errorf("loadScripts: " + err.Error())
			}

			// Das Script wird zwischengespeichert
			nodeJsModules = append(nodeJsModules, newNodeJsModule)
		}
	}

	// Es wird geprüft ob überhaupt Script Module gefunden wurden
	if len(pythonModules) == 0 && len(nodeJsModules) == 0 {
		fmt.Println(pythonModules, nodeJsModules)
		return nil, fmt.Errorf("loadScripts: no scripts available")
	}

	// Das Rückgabe Objekt wird zurückgegeben
	newObj := &ScriptsFolder{
		PythonModules: pythonModules,
		NodeJsModules: nodeJsModules,
	}

	// Das Objekt wird zurückgegeben
	return newObj, nil
}
