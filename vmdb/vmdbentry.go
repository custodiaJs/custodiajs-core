package vmdb

import (
	"fmt"
	"path/filepath"
	"strings"
	"vnh1/utils"
)

func (o *VmDBEntry) GetManifest() *Manifest {
	return o.manifestFile.GetManifestObject()
}

func (o *VmDBEntry) GetOwner() string {
	return strings.ToLower(o.manifestFile.manifest.Owner)
}

func (o *VmDBEntry) GetRepoURL() string {
	return strings.ToLower(o.manifestFile.manifest.RepoURL)
}

func (o *VmDBEntry) GetMode() string {
	return strings.ToLower(o.manifestFile.manifest.Mode)
}

func (o *VmDBEntry) ValidateVM() bool {
	return true
}

func (o *VmDBEntry) GetVMName() string {
	return strings.ToLower(o.manifestFile.manifest.Name)
}

func (o *VmDBEntry) GetVMContainerMerkleHash() string {
	return strings.ToLower(o.vmContainerMerkleHash)
}

func (o *VmDBEntry) GetBaseSize() uint64 {
	return o.containerBaseSize
}

func (p *VmDBEntry) GetWhitelist() []Whitelist {
	return p.manifestFile.GetManifestObject().Whitelist
}

func (o *VmDBEntry) GetRootMemberIDS() []*CAMemberData {
	ret := make([]*CAMemberData, 0)
	for _, item := range o.manifestFile.manifest.HostCAMember {
		ret = append(ret, &CAMemberData{
			Fingerprint: item.Fingerprint,
			Type:        item.Type,
			ID:          item.ID,
		})
	}
	return ret
}

func (o *VmDBEntry) GetMainCodeFile() *MainJsFile {
	return o.mainJSFile
}

func (o *VmDBEntry) GetAllowedHttpSources() map[string]bool {
	a := make(map[string]bool)
	a["*.com"] = true
	return a
}

func (o *VmDBEntry) GetAllDatabaseServices() []*VMDatabaseData {
	vmdlist := make([]*VMDatabaseData, 0)
	for _, item := range o.manifestFile.GetAllDatabaseServices() {
		vmdlist = append(vmdlist, &VMDatabaseData{
			Type:     item.Type,
			Host:     item.Host,
			Port:     item.Port,
			Username: item.Username,
			Password: item.Password,
			Database: item.Database,
			Alias:    item.Alias,
		})
	}
	return vmdlist
}

func (o *VmDBEntry) GetAllExternalServices() []*VmExternalService {
	servicesList := make([]*VmExternalService, 0)
	for _, item := range o.manifestFile.GetManifestObject().Services.External {
		servicesList = append(servicesList, &VmExternalService{MinVersion: uint(item.MinVersion), Name: item.Name, Required: item.Required})
	}
	return servicesList
}

func (o *VmDBEntry) GetAllExperimentalWebservices() []*VmExperimentalWebservice {
	return []*VmExperimentalWebservice{}
}

func tryToLoadVM(path string) (*VmDBEntry, error) {
	// Die Kernpfade für die VM werden erstellt
	manifestVMJsonFilePath := filepath.Join(path, "manifest.json")
	signatureFilePath := filepath.Join(path, "signature")
	mainJSFilePath := filepath.Join(path, "main.js")

	// Es wird versucht die Manifestdatei einzulesen
	manifestFile, err := loadManifestFile(manifestVMJsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("tryToLoadVM: " + err.Error())
	}

	// Es wird geprüft ob es sich um eine gültige ManifestFile handelt
	if err := manifestFile.ValidateWithState(); err != nil {
		return nil, fmt.Errorf("VmDBEntry->tryToLoadVM: " + err.Error())
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

	// Es wird eine Hashliste aus allen Hashes erstellt
	mergedHashList := []string{manifestFile.GetFileHash(), mainJsFile.fileHash}

	// Die Hashliste wird Sortiert
	sortedHashList, err := utils.SortHexStrings(mergedHashList)
	if err != nil {
		return nil, fmt.Errorf("tryToLoadVM: " + err.Error())
	}

	// Es wird ein Merkelhash aus der Sortierten Liste erstellt
	merkleRoot, err := utils.BuildMerkleRoot(sortedHashList)
	if err != nil {
		return nil, fmt.Errorf("tryToLoadVM: " + err.Error())
	}

	// Die Gesamtgröße der VM wird ermittelt
	containerBaseSize := mainJsFile.GetFileSize() + manifestFile.GetFileSize()

	// Das Objekt wird erstellt
	newObject := &VmDBEntry{
		Path:                  path,
		mainJSFile:            mainJsFile,
		manifestFile:          manifestFile,
		signatureFile:         sigFile,
		vmContainerMerkleHash: merkleRoot,
		containerBaseSize:     containerBaseSize,
	}

	// Das Objekt wird zurückgegeben
	return newObject, nil
}
