package filesystem

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/CustodiaJS/custodiajs-core/utils"
)

type FileInfo struct {
	Path             string // Der vollständige Pfad der Datei
	Extension        string
	Name             string    // Der Name der Datei
	Size             int64     // Die Größe der Datei in Bytes
	ModificationTime time.Time // Das letzte Änderungsdatum der Datei
	FileHash         string
}

func WalkDir(dirPath string, withHash bool) ([]FileInfo, error) {
	var files []FileInfo
	absDirPath, err := filepath.Abs(dirPath) // Ermittelt den absoluten Pfad des Startordners
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(absDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if absPath, err := filepath.Abs(path); err == nil && absPath != absDirPath {
			// Sollte es sich um einen Ordner handeln, wird ein Hash aus dem Namen erzeugt, ansonsten wird der Dateiinhalt gehast
			var hash string
			if info.IsDir() {
				hash = utils.HashOfString(info.Name())
			} else {
				fHash, err := utils.HashFile(path)
				if err != nil {
					return err
				}
				hash = hex.EncodeToString(fHash)
			}

			files = append(files, FileInfo{
				Path:             path,
				Name:             info.Name(),
				Extension:        filepath.Ext(path),
				Size:             info.Size(),
				ModificationTime: info.ModTime(),
				FileHash:         hash,
			})
		}
		return nil
	})

	return files, err
}

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

func CreateDirectory(path string) error {
	// Versucht, den Ordner zu erstellen.
	// MkdirAll erstellt auch alle nötigen übergeordneten Verzeichnisse.
	err := os.MkdirAll(path, os.ModePerm) // os.ModePerm setzt die Berechtigungen auf 0777
	if err != nil {
		return fmt.Errorf("CreateDirectory: %v", err)
	}
	return nil
}
