//go:build !windows
// +build !windows

package usergroups

import (
	"fmt"
	"os/user"
)

func ListAllUserGroups() {
	// Aktuellen Benutzer abrufen
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Fehler beim Abrufen des aktuellen Benutzers:", err)
		return
	}

	// Benutzername und Gruppen-ID anzeigen
	fmt.Println("Benutzername:", currentUser.Username)
	fmt.Println("Gruppen-ID:", currentUser.Gid)

	// Alle Gruppen des Benutzers abrufen
	groups, err := user.LookupGroupId(currentUser.Gid)
	if err != nil {
		fmt.Println("Fehler beim Abrufen der Gruppeninformationen:", err)
		return
	}
	fmt.Println("Gruppe:", groups.Name)
}
