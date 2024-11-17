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
