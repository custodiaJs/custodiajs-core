package vmdb

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/CustodiaJS/custodiajs-core/filesystem"
	"github.com/CustodiaJS/custodiajs-core/utils"

	"github.com/gofrs/flock"
)

func (o *ManifestFile) GetManifestObject() *Manifest {
	return o.manifest
}

func (o *ManifestFile) GetFileHash() string {
	return strings.ToLower(o.fileHash)
}

func (o *ManifestFile) GetAllDatabaseServices() []*VMEntryBaseData {
	vmdlist := make([]*VMEntryBaseData, 0)
	for _, item := range o.manifest.Databases {
		vmdlist = append(vmdlist, &VMEntryBaseData{
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

func (o *ManifestFile) GetFileSize() uint64 {
	return o.fSize
}

func (o *ManifestFile) ValidateWithState() error {
	return nil
}

func loadManifestFile(path string) (*ManifestFile, error) {
	// Es wird geprüft ob die Datei vorhanden ist
	if !filesystem.FileExists(path) {
		return nil, fmt.Errorf(fmt.Sprintf("loadManifestFile: file '%s' not found", path))
	}

	// Create a new file lock
	fileLock := flock.New(path)

	// Try to lock the file
	locked, err := fileLock.TryLock()
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("loadSignatureFile: Error trying to lock the file: %v\n", err))
	}
	if !locked {
		return nil, fmt.Errorf("loadManifestFile: unable to lock the file, it may be locked by another process")
	}

	// Die Datei wird geöffnet
	openFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// Es wird ein Hash aus der Datei erzeugt
	fileHash, err := utils.HashOSFile(openFile)
	if err != nil {
		return nil, fmt.Errorf("loadManifestFile: " + err.Error())
	}

	// Die Datei wird eingelesen
	readedFileBytes, err := filesystem.ReadFileBytes(openFile)
	if err != nil {
		return nil, fmt.Errorf("loadManifestFile: " + err.Error())
	}

	// Das Manifestobjekt wird eingelesen
	var config Manifest
	if err := json.Unmarshal(readedFileBytes, &config); err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("loadManifestFile: converting to object %s", err.Error()))
	}

	// Die Größe der Datei wird ermittelt
	fsize, err := filesystem.GetFileSize(path)
	if err != nil {
		return nil, fmt.Errorf("loadMainJsFile: " + err.Error())
	}

	// Das Objekt wird gebaut
	resolvObject := &ManifestFile{
		osFile:   openFile,
		fileLock: fileLock,
		fileHash: hex.EncodeToString(fileHash),
		manifest: &config,
		fSize:    uint64(fsize),
	}

	// Das Objekt wird zurückgegeben
	return resolvObject, nil
}
