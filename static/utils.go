package static

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/sha3"
)

// Funktion zum Berechnen des SHA3-Hashes einer Datei, ohne sie komplett in den Speicher zu laden
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

// Funktion, die prüft, ob eine Datei existiert
func FileExists(filePath string) bool {
	f, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	if f.IsDir() {
		return false
	}

	return true // Kein Fehler oder ein anderer Fehler, der nicht bedeutet, dass die Datei nicht existiert
}

func FolderExists(filePath string) bool {
	f, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	if !f.IsDir() {
		return false
	}

	return true // Kein Fehler oder ein anderer Fehler, der nicht bedeutet, dass die Datei nicht existiert
}

func ListAllFolders(rpath string) ([]string, error) {
	// Liste, um die Pfade der gefundenen Verzeichnisse zu speichern
	var dirs []string

	// Lese den Inhalt des Startverzeichnisses
	files, err := os.ReadDir(rpath)
	if err != nil {
		return nil, fmt.Errorf("ListAllFolders: %w", err)
	}

	// Durchlaufe den Inhalt des Verzeichnisses
	for _, file := range files {
		// Überprüfen, ob es sich um ein Verzeichnis handelt
		if file.IsDir() {
			// Füge den Pfad des Verzeichnisses zur Liste hinzu
			dirs = append(dirs, filepath.Join(rpath, file.Name()))
		}
	}

	return dirs, nil
}

func ScanVmDir(rpath string) ([]string, error) {
	// Liste, um die Pfade der gefundenen Verzeichnisse zu speichern
	var dirs []string

	// Pfad zum Startverzeichnis
	startDir := filepath.Join(rpath, "vms")

	// Lese den Inhalt des Startverzeichnisses
	files, err := os.ReadDir(startDir)
	if err != nil {
		return nil, fmt.Errorf("ScanVmDir: %w", err)
	}

	// Durchlaufe den Inhalt des Verzeichnisses
	for _, file := range files {
		// Überprüfen, ob es sich um ein Verzeichnis handelt
		if file.IsDir() {
			// Füge den Pfad des Verzeichnisses zur Liste hinzu
			dirs = append(dirs, filepath.Join(startDir, file.Name()))
		}
	}

	return dirs, nil
}

// Ließt Bytes aus einer Datei aus
func ReadFileBytes(file *os.File) ([]byte, error) {
	// Lesezeiger zurücksetzen
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("HashOSFile: " + err.Error())
	}

	// Lese den gesamten Inhalt der Datei
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Lesezeiger zurücksetzen
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("HashOSFile: " + err.Error())
	}

	// Die Daten werden zurückgegeben
	return data, nil
}
