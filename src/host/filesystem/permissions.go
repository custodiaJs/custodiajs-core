package filesystem

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

// setUserGroupPermission Setzt die Gruppenberechtigung für einen Socket
func SetUserGroupFilePermission(socketPath string, username string) error {
	user, err := user.Lookup(username)
	if err != nil {
		return fmt.Errorf("error looking up user: %v", err)
	}

	uid, _ := strconv.Atoi(user.Uid)

	// Setzen der User-ID als Besitzer des Sockets, Gruppen-ID bleibt unverändert
	if err := syscall.Chown(socketPath, uid, -1); err != nil {
		return fmt.Errorf("error setting user owner: %v", err)
	}

	// Lese-/Schreibzugriff nur für den Benutzer
	return os.Chmod(socketPath, 0600)
}

// setUserPermission Setzt die Benutzerberechtigung für einen Socket
func SetUserFilePermission(socketPath string, username string) error {
	user, err := user.Lookup(username)
	if err != nil {
		return fmt.Errorf("error looking up user: %v", err)
	}

	uid, _ := strconv.Atoi(user.Uid)

	// Setzen der User-ID als Besitzer des Sockets
	if err := syscall.Chown(socketPath, uid, -1); err != nil {
		return fmt.Errorf("error setting owner: %v", err)
	}

	// Nur Benutzer-Lese-/Schreibzugriff
	return os.Chmod(socketPath, 0600)
}
