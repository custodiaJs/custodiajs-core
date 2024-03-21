package vmdb

import (
	"fmt"
	"path/filepath"
	"vnh1/static"
)

type VmDBEntry struct {
	Path          string
	mainJSFile    *MainJsFile
	manifestFile  *ManifestFile
	signatureFile *SignatureFile
	nodeJsModules []*NodeJsModule
}

func (o *VmDBEntry) ValidateVM() bool {
	return true
}

func tryToLoadVM(path string) (*VmDBEntry, error) {
	// Die Kernpfade für die VM werden erstellt
	manifestVMJsonFilePath := filepath.Join(path, "manifest.json")
	nodeJsModulesPath := filepath.Join(path, "nodejs")
	signatureFilePath := filepath.Join(path, "signature")
	mainJSFilePath := filepath.Join(path, "main.js")

	// Es wird versucht die Manifestdatei einzulesen
	manifestFile, err := loadManifestFile(manifestVMJsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("tryToLoadVM: " + err.Error())
	}

	// Die Signatur wird geladen
	sigFile, err := loadSignatureFile(signatureFilePath)
	if err != nil {
		return nil, fmt.Errorf("tryToLoadVM: " + err.Error())
	}

	// Die MainJS Datei wird geladen
	mainJsFile, err := loadMainJsFile(mainJSFilePath)
	if err != nil {
		return nil, fmt.Errorf("tryToLoadVM: " + err.Error())
	}

	// Es wird geprüft ob die Manifestdatei Scripte angibt
	extractedNodejSModules := make([]*NodeJsModule, 0)
	if manifestFile.NodeJsEnable() {
		// Sollte der Scriptsordner nicht vorhanden sein, wird der Vorgang abgebrochen
		if !static.FolderExists(nodeJsModulesPath) {
			return nil, fmt.Errorf("tryToLoadVM: no scripts found")
		}

		// Es werden die Verfügbaren Scripte eingelesen
		scripts, err := tryToLoadNodeJsModules(nodeJsModulesPath)
		if err != nil {
			return nil, fmt.Errorf("tryToLoadVM: " + err.Error())
		}

		// Es wird geprüft ob NodeJs Scripte vorhanden sein müssen
		if manifestFile.manifest.NodeJS.Enable {
			// Es wird geprüft ob alle NodeJs Module welche das Manifest angibt, vorhanden sind
			if len(manifestFile.manifest.NodeJS.Modules) != len(scripts) {
				return nil, fmt.Errorf("tryToLoadVM: invalid vm container, not all scriptes avail")
			}

			// Speichert alle NodeJs Module ab, welche die Manifestdatei angibt
			validateNodeJsModules := make(map[string]bool)
			for _, scriptItem := range manifestFile.manifest.NodeJS.Modules {
				validateNodeJsModules[scriptItem.Name] = false
			}
		}
	} else {
		// Sollte der Scriptsordner vorhanden sein, wird der Vorgang abgebrochen
		if static.FolderExists(nodeJsModulesPath) {
			return nil, fmt.Errorf("tryToLoadVM: scripts not allowed")
		}
	}

	// Das Objekt wird erstellt
	newObject := &VmDBEntry{
		Path:          path,
		mainJSFile:    mainJsFile,
		manifestFile:  manifestFile,
		signatureFile: sigFile,
		nodeJsModules: extractedNodejSModules,
	}

	// Das Objekt wird zurückgegeben
	return newObject, nil
}
