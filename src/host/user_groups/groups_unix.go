//go:build !windows
// +build !windows

// Author: fluffelpuff
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
