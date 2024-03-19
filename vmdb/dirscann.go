package vmdb

import (
	"fmt"
	"os"
	"path/filepath"
)

func scanVmDir(rpath string) ([]string, error) {
	// Liste, um die Pfade der gefundenen Verzeichnisse zu speichern
	var dirs []string

	// filepath.Walk durchläuft das Dateisystem beginnend bei startDir
	startDir := filepath.Join(rpath, "vms")
	err := filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Überprüfen, ob es sich um ein Verzeichnis handelt
		if info.IsDir() {
			// Das Startverzeichnis selbst überspringen
			if path != startDir {
				dirs = append(dirs, path)
			}
		}
		return nil
	})

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if err != nil {
		return nil, fmt.Errorf("LoadAllVirtualMachines: " + err.Error())
	}

	return dirs, nil
}
