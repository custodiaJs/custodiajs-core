package vmdb

import (
	"encoding/hex"
	"fmt"
	"vnh1/static"

	"github.com/gofrs/flock"
)

type MainJsFile struct {
	fileLock       *flock.Flock
	mainJsFileHash string
	filePath       string
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

	// Die Daten werden zusammengefasst
	newObj := &MainJsFile{
		mainJsFileHash: hex.EncodeToString(fileHash),
		fileLock:       fileLock,
		filePath:       path,
	}

	// Das Objekt wird zurückgegeben
	return newObj, nil
}
