package vmdb

import (
	"encoding/json"
	"fmt"
	"os"
	"vnh1/utils"
)

type NodeJsModule struct {
	merkleRoot string
	files      []utils.FileInfo
	baseSize   uint64
	name       string
}

func (o *NodeJsModule) GetBaseSize() uint64 {
	return o.baseSize
}

func (o *NodeJsModule) GetName() string {
	return o.name
}

func tryToLoadNodeJsModules(path string) ([]*NodeJsModule, error) {
	// Es wird geprüft ob es sich um einen gültigen Path handelt
	if !utils.FolderExists(path) {
		return nil, fmt.Errorf("tryToLoadNodeJsModule: no nodejs modules folder found")
	}

	// Es wird eine übersicht über alle Ordner estellt
	folders, err := utils.ListAllFolders(path)
	if err != nil {
		return nil, err
	}

	// Die einzelnen Ordner werden abgeabeitet
	loadedNodeJsModules := make([]*NodeJsModule, 0)
	for _, folderPath := range folders {
		// Es wird eine Übersicht über den Ordnerinhalt erstellt
		folderOverview, err := utils.WalkDir(folderPath, true)
		if err != nil {
			return nil, fmt.Errorf("tryToLoadNodeJsModules: " + err.Error())
		}

		// Es wird geprüft ob die benötigten Dateien vorhanden sind
		packageJsonFile := false
		unsortedHashList := make([]string, 0)
		for _, item := range folderOverview {
			// Der Hash der Datei wird der Liste hinzugefügt
			unsortedHashList = append(unsortedHashList, item.FileHash)

			// Es wird geprüft ob es sich um die Package.JSON handelt
			if item.Name == "package.json" {
				// Lese die Datei
				data, err := os.ReadFile(item.Path)
				if err != nil {
					return nil, fmt.Errorf("konnte die Datei nicht lesen: %w", err)
				}

				// Parse das JSON in die Config-Struktur
				var config Config
				if err := json.Unmarshal(data, &config); err != nil {
					return nil, fmt.Errorf("konnte das JSON nicht parsen: %w", err)
				}

				// Überprüfe, ob der spezifische scripts-Eintrag existiert
				command, exists := config.Scripts["start:vm:child"]
				if !exists || command != "node build/index.js" {
					return nil, fmt.Errorf("der 'start:vm:child'-enty dosen't found")
				}

				// Es wird angegeben dass die Package.json Datei vorhanden und korrekt ist
				packageJsonFile = true
			}
		}

		// Es wird geprüft ob die Package.json Datei gefunden wurde
		if !packageJsonFile {
			return nil, fmt.Errorf("tryToLoadNodeJsModules: isnt a nodejs module")
		}

		// Die Liste wird Sortiert
		sortedHashList, err := utils.SortHexStrings(unsortedHashList)
		if err != nil {
			return nil, fmt.Errorf("tryToLoadNodeJsModules: " + err.Error())
		}

		// Es wird ein Merkle Root estellt
		merkleRoot, err := utils.BuildMerkleRoot(sortedHashList)
		if err != nil {
			return nil, fmt.Errorf("tryToLoadNodeJsModules: " + err.Error())
		}

		// Das Objekt wird zwischegespeichert
		loadedNodeJsModules = append(loadedNodeJsModules, &NodeJsModule{merkleRoot: merkleRoot, files: folderOverview, name: utils.ExtractFileName(folderPath)})
	}

	// Die NodeJS Module werden zurückgegeben
	return loadedNodeJsModules, nil
}
