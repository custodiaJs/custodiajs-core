package static

import (
	"os"
	"path/filepath"
)

// GetFileSize gibt die Größe einer Datei in Bytes zurück.
func GetFileSize(filePath string) (int64, error) {
	// Öffne die Datei
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		// Gibt einen Fehler zurück, falls das Öffnen fehlschlägt
		return 0, err
	}

	// Die Größe der Datei in Bytes
	return fileInfo.Size(), nil
}

// ExtractFileName nimmt einen Dateipfad als Eingabe und gibt den Dateinamen zurück.
func ExtractFileName(filePath string) string {
	return filepath.Base(filePath)
}
