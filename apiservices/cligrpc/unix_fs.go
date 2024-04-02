package cligrpc

import (
	"os"
)

func deleteFileIfExists(filepath string) error {
	// Überprüfe, ob die Datei existiert
	_, err := os.Stat(filepath)
	if err == nil {
		err := os.Remove(filepath)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func setUnixFilePermissionsForAll(filepath string) error {
	err := os.Chmod(filepath, 0777)
	if err != nil {
		return err
	}
	return nil
}

func setUnixFileOwnerToRoot(filepath string) error {
	err := os.Chown(filepath, 0, 0)
	if err != nil {
		return err
	}
	return nil
}
