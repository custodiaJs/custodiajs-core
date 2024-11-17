package utils

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/sha3"
)

func HashFile(filePath string) ([]byte, error) {
	// Öffne die Datei zum Lesen
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Erstelle eine neue Instanz des SHA3-256 Hashers
	hasher := sha3.New256()

	// Lese die Datei in Teilen und aktualisiere den Hasher nach jedem gelesenen Teil
	buf := make([]byte, 4096) // Puffer für das Lesen in Teilen; Größe kann angepasst werden
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		// Aktualisiere den Hasher mit dem Inhalt des aktuellen Teils
		hasher.Write(buf[:n])
	}

	// Finalisiere den Hash-Prozess und erhalte den resultierenden Hash
	return hasher.Sum(nil), nil
}

func HashOSFile(file *os.File) ([]byte, error) {
	// Lesezeiger zurücksetzen
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("HashOSFile: " + err.Error())
	}

	// Erstelle eine neue Instanz des SHA3-256 Hashers
	hasher := sha3.New256()

	// Lese die Datei in Teilen und aktualisiere den Hasher nach jedem gelesenen Teil
	buf := make([]byte, 4096) // Puffer für das Lesen in Teilen; Größe kann angepasst werden
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		// Aktualisiere den Hasher mit dem Inhalt des aktuellen Teils
		hasher.Write(buf[:n])
	}

	// Lesezeiger zurücksetzen
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("HashOSFile: " + err.Error())
	}

	// Finalisiere den Hash-Prozess und erhalte den resultierenden Hash
	return hasher.Sum(nil), nil
}

func HashOfString(data string) string {
	hasher := sha3.New256()
	hasher.Write([]byte(data))
	varh := hasher.Sum(nil)
	return hex.EncodeToString(varh)
}
