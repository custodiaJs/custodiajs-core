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
}

func (o *VmDBEntry) ValidateVM() bool {
	return true
}

func tryToLoadVM(path string) (*VmDBEntry, error) {
	// Die Kernpfade für die VM werden erstellt
	manifestVMJsonFilePath := filepath.Join(path, "manifest.json")
	scriptFolderPath := filepath.Join(path, "scripts")
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
	if manifestFile.ScriptsEnable() {
		// Sollte der Scriptsordner nicht vorhanden sein, wird der Vorgang abgebrochen
		if !static.FolderExists(scriptFolderPath) {
			return nil, fmt.Errorf("tryToLoadVM: no scripts found")
		}

		// Es werden die Verfügbaren Scripte eingelesen
		scripts, err := loadScripts(scriptFolderPath)
		if err != nil {
			return nil, fmt.Errorf("tryToLoadVM: " + err.Error())
		}

		// Es wird geprüft ob Python Scripts vorhanden sein müssen
		if manifestFile.manifest.Scripts.Python.Enable {
			// Es wird geprüft ob alle Module welche das Manifest angibt, vorhanden sind
			if len(manifestFile.manifest.Scripts.Python.Modules) != len(scripts.NodeJsModules) {
				//return nil, fmt.Errorf("tryToLoadVM: invalid vm container, not all scriptes avail")
			}

			// Es wird geprüft ob es für jedes Python Script welche im Manifest angegeben wurde, vorhanden ist
			for _, item := range manifestFile.manifest.Scripts.Python.Modules {
				fmt.Println(item.Alias)
			}
		}
	} else {
		// Sollte der Scriptsordner vorhanden sein, wird der Vorgang abgebrochen
		if static.FolderExists(scriptFolderPath) {
			return nil, fmt.Errorf("tryToLoadVM: scripts not allowed")
		}
	}

	// Das Objekt wird erstellt
	newObject := &VmDBEntry{
		Path:          path,
		mainJSFile:    mainJsFile,
		manifestFile:  manifestFile,
		signatureFile: sigFile,
	}

	// Das Objekt wird zurückgegeben
	return newObject, nil
}
