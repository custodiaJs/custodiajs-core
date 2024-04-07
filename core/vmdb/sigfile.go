package vmdb

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"vnh1/utils"

	"github.com/gofrs/flock"
)

func loadSignatureFile(path string) (*SignatureFile, error) {
	// Es wird geprüft ob die Datei vorhanden ist
	if !utils.FileExists(path) {
		return nil, fmt.Errorf(fmt.Sprintf("loadSignatureFile: file '%s' not found", path))
	}

	// Create a new file lock
	fileLock := flock.New(path)

	// Try to lock the file
	locked, err := fileLock.TryLock()
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("loadSignatureFile: Error trying to lock the file: %v\n", err))
	}
	if !locked {
		return nil, fmt.Errorf("loadSignatureFile: unable to lock the file, it may be locked by another process")
	}

	// Die Datei wird geöffnet
	openFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("loadSignatureFile: %v", err))
	}

	// Der Owner sowie die Signatur wird ausgelsen
	var ownerCert string
	var ownerSignature string
	scanner := bufio.NewScanner(openFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "SIGNATURE:") {
			ownerSignature = strings.TrimSpace(line[len("SIGNATURE:"):])
		} else if strings.HasPrefix(line, "OWNERCERT:") {
			ownerCert = strings.TrimSpace(line[len("OWNERCERT:"):])
		}
	}

	// Der Inhalt der Datei wird eingelesen
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("loadSignatureFile: %v", err))
	}

	// Das Rückgabe Objekt wird erstellt
	newObject := &SignatureFile{
		osFile:         openFile,
		fileLock:       fileLock,
		ownerCert:      ownerCert,
		ownerSignature: ownerSignature,
	}

	// Das Objekt wird zurückgegeben
	return newObject, nil
}
