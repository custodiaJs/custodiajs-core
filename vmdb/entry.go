package vmdb

import (
	"fmt"
	"path/filepath"
)

type VmDBEntry struct {
	Path   string
	vmHash string
}

func (o *VmDBEntry) ValidateVM() bool {
	return true
}

func tryToLoadVM(path string) (*VmDBEntry, error) {
	// Die Kernpfade für die VM werden erstellt
	manifestVMJsonFile := filepath.Join(path, "vm.manifest")
	scriptFolderPath := filepath.Join(path, "scripts")
	configurationFile := filepath.Join(path, "config")
	signatureFile := filepath.Join(path, "signature")
	mainJSFile := filepath.Join(path, "main.js")

	// Es wird versucht die Manifestdatei einzulesen

	// Es wird geprüft ob der Scriptsordner vorhanden ist, wenn ja werden diese eingelesen
	if existsDir(scriptFolderPath) {

	}

	fmt.Println(scriptFolderPath)
	return nil, nil
}
