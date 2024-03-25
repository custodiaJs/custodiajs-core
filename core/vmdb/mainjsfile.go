package vmdb

import (
	"encoding/hex"
	"fmt"
	"vnh1/static"

	"github.com/gofrs/flock"
)

type MainJsFile struct {
	fileLock *flock.Flock
	fileSize uint64
	fileHash string
	filePath string
}

func (o *MainJsFile) GetFileSize() uint64 {
	return o.fileSize
}

func loadMainJsFile(path string) (*MainJsFile, error) {
	// Es wird geprüft ob die Datei vorhanden ist
	if !static.FileExists(path) {
		return nil, fmt.Errorf(fmt.Sprintf("loadMainJsFile: file '%s' not found", path))
	}

	// Create a new file lock
	fileLock := flock.New(path)

	// Try to lock the file
	locked, err := fileLock.TryLock()
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Error trying to lock the file: %v\n", err))
	}
	if !locked {
		return nil, fmt.Errorf("loadMainJsFile: unable to lock the file, it may be locked by another process")
	}

	// Es wird ein Hash aus der Datei erzeugt
	fileHash, err := static.HashFile(path)
	if err != nil {
		return nil, fmt.Errorf("loadMainJsFile: " + err.Error())
	}

	// Die Größe der Datei wird ermittelt
	fsize, err := static.GetFileSize(path)
	if err != nil {
		return nil, fmt.Errorf("loadMainJsFile: " + err.Error())
	}

	// Die Daten werden zusammengefasst
	newObj := &MainJsFile{
		fileHash: hex.EncodeToString(fileHash),
		fileLock: fileLock,
		fileSize: uint64(fsize),
		filePath: path,
	}

	// Das Objekt wird zurückgegeben
	return newObj, nil
}
