package vmdb

import (
	"fmt"
	"path/filepath"
	"strings"
	"vnh1/static"
)

type VmDBEntry struct {
	Path                  string
	mainJSFile            *MainJsFile
	manifestFile          *ManifestFile
	signatureFile         *SignatureFile
	nodeJsModules         []*NodeJsModule
	vmContainerMerkleHash string
	containerBaseSize     uint64
}

func (o *VmDBEntry) ValidateVM() bool {
	return true
}

func (o *VmDBEntry) GetVMName() string {
	return o.manifestFile.manifest.Name
}

func (o *VmDBEntry) GetVMContainerMerkleHash() string {
	return o.vmContainerMerkleHash
}

func (o *VmDBEntry) GetBaseSize() uint64 {
	return o.containerBaseSize
}

func (o *VmDBEntry) GetTotalNodeJsModules() uint64 {
	return uint64(len(o.nodeJsModules))
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

		// Es wird geprüft ob alle NodeJs Module welche das Manifest angibt, vorhanden sind
		if len(manifestFile.manifest.NodeJS.Modules) != len(scripts) {
			return nil, fmt.Errorf("tryToLoadVM: invalid vm container, not all scriptes avail")
		}

		// Speichert alle NodeJs Module ab, welche die Manifestdatei angibt
		validateNodeJsModules := make(map[string]bool)
		for _, scriptItem := range manifestFile.manifest.NodeJS.Modules {
			validateNodeJsModules[scriptItem.Name] = false
		}

		// Jedes NodeJS Modul wird geprüft
		unkownModules := make([]string, 0)
		for _, item := range scripts {
			// Es wird geprüft ob es sich um ein bekanntes Module handelt
			val, found := validateNodeJsModules[item.name]
			if !found {
				unkownModules = append(unkownModules, item.name)
				continue
			}

			// Der Wert des eintrags muss false sein, ansonsten liegt ein Massiver fehler vor
			if val {
				return nil, fmt.Errorf("tryToLoadVM: broken vm container, unkown fatal error")
			}

			// Der Wert des Modules wird auf geprüft gesetzt
			validateNodeJsModules[item.name] = true

			// Das Modul wird abgespeichert
			extractedNodejSModules = append(extractedNodejSModules, item)
		}

		// Es dürfen keine Unbekannten Module vorhanden sein
		if len(unkownModules) != 0 {
			return nil, fmt.Errorf("tryToLoadVM: broken vm container, nodejs module '%s' not in manifest file", strings.Join(unkownModules, ", "))
		}
	} else {
		// Sollte der Scriptsordner vorhanden sein, wird der Vorgang abgebrochen
		if static.FolderExists(nodeJsModulesPath) {
			return nil, fmt.Errorf("tryToLoadVM: scripts not allowed")
		}
	}

	// Es wird eine Hashliste aus allen Hashes erstellt
	mergedHashList := []string{manifestFile.fileHash, mainJsFile.fileHash}
	for _, item := range extractedNodejSModules {
		mergedHashList = append(mergedHashList, item.merkleRoot)
	}

	// Die Hashliste wird Sortiert
	sortedHashList, err := static.SortHexStrings(mergedHashList)
	if err != nil {
		return nil, fmt.Errorf("tryToLoadVM: " + err.Error())
	}

	// Es wird ein Merkelhash aus der Sortierten Liste erstellt
	merkleRoot, err := static.BuildMerkleRoot(sortedHashList)
	if err != nil {
		return nil, fmt.Errorf("tryToLoadVM: " + err.Error())
	}

	// Die Gesamtgröße der VM wird ermittelt
	containerBaseSize := mainJsFile.GetFileSize() + manifestFile.GetFileSize()
	for _, item := range extractedNodejSModules {
		containerBaseSize += item.GetBaseSize()
	}

	// Das Objekt wird erstellt
	newObject := &VmDBEntry{
		Path:                  path,
		mainJSFile:            mainJsFile,
		manifestFile:          manifestFile,
		signatureFile:         sigFile,
		nodeJsModules:         extractedNodejSModules,
		vmContainerMerkleHash: merkleRoot,
		containerBaseSize:     containerBaseSize,
	}

	// Das Objekt wird zurückgegeben
	return newObject, nil
}
