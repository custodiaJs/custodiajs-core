package vmdb

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"vnh1/static"

	"github.com/gofrs/flock"
)

type ManifestFile struct {
	osFile   *os.File
	fileLock *flock.Flock
	fileHash string
	manifest *Manifest
}

func (o *ManifestFile) GetManifestObject() *Manifest {
	return o.manifest
}

func (o *ManifestFile) GetFileHash() string {
	return o.fileHash
}

func (o *ManifestFile) NodeJsEnable() bool {
	return o.manifest.NodeJS.Enable
}

func (o *ManifestFile) GetNodeJsScriptAlias() []string {
	resolv := []string{}
	for _, item := range o.manifest.NodeJS.Modules {
		resolv = append(resolv, item.Alias)
	}
	return resolv
}

func loadManifestFile(path string) (*ManifestFile, error) {
	// Es wird geprüft ob die Datei vorhanden ist
	if !static.FileExists(path) {
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
	fileHash, err := static.HashOSFile(openFile)
	if err != nil {
		return nil, fmt.Errorf("loadManifestFile: " + err.Error())
	}

	// Die Datei wird eingelesen
	readedFileBytes, err := static.ReadFileBytes(openFile)
	if err != nil {
		return nil, fmt.Errorf("loadManifestFile: " + err.Error())
	}

	// Das Manifestobjekt wird eingelesen
	var config Manifest
	if err := json.Unmarshal(readedFileBytes, &config); err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("loadManifestFile: converting to object %s", err.Error()))
	}

	// Das Objekt wird gebaut
	resolvObject := &ManifestFile{
		osFile:   openFile,
		fileLock: fileLock,
		fileHash: hex.EncodeToString(fileHash),
		manifest: &config,
	}

	// Das Objekt wird zurückgegeben
	return resolvObject, nil
}
