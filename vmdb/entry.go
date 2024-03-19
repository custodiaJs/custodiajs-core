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

	fmt.Println("TRYTOLOAD", manifestVMJsonFile, configurationFile, signatureFile, mainJSFile)
	fmt.Println(scriptFolderPath)
	return nil, nil
}
