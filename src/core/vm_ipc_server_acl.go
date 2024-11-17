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

package core

import (
	"fmt"
	"net"
	"os"
	"os/user"

	cenvxcore "github.com/custodia-cenv/cenvx-core/src"
	"github.com/custodia-cenv/cenvx-core/src/host/filesystem"
	"github.com/custodia-cenv/cenvx-core/src/log"
)

// createAclListeners erstellt für jede ACL in der Liste einen _AclListener und gibt sie zurück.
func createAclListeners(aclList []*ACL) ([]*_AclListener, error) {
	var listeners []*_AclListener

	for i, acl := range aclList {
		var socketPath cenvxcore.CoreVmIpcSocketPath
		if acl.Username == nil && acl.Groupname == nil {
			socketPath = cenvxcore.CoreVmIpcRootSocketPath
		} else if acl.Username != nil && acl.Groupname == nil {
			socketPath = cenvxcore.GetCoreSpeficSocketUserPath(*acl.Username)
		} else if acl.Username != nil && acl.Groupname != nil {
			socketPath = cenvxcore.GetCoreSpeficSocketUserAndGroupPath(*acl.Username, *acl.Groupname)
		} else if acl.Username == nil && acl.Groupname != nil {
			socketPath = cenvxcore.GetCoreSpeficSocketUserGroupPath(*acl.Groupname)
		} else {
			panic("unkown acl config")
		}

		// Listener für Unix-Socket erstellen
		listener, err := createListenerWithACL(socketPath, *acl)
		if err != nil {
			return nil, fmt.Errorf("error creating listener for ACL %d: %v", i, err)
		}

		// _AclListener erstellen und hinzufügen
		aclListener := &_AclListener{
			Listener: listener,
			AclRule:  acl,
		}
		listeners = append(listeners, aclListener)

		// Berechtigungen setzen basierend auf ACL
		if acl.Username == nil && acl.Groupname == nil {
			// Zugriff nur für root erlauben
			err = os.Chmod(string(socketPath), 0600)
			if err != nil {
				listener.Close()
				return nil, fmt.Errorf("error setting root-only permissions: %v", err)
			}
			log.DebugLogPrint("VM-IPC Socket created: root -> %s", socketPath)
		} else if acl.Groupname != nil && acl.Username != nil {
			log.DebugLogPrint("VM-IPC Socket created: %s - %s -> %s", *acl.Username, *acl.Groupname, socketPath)
		} else if acl.Username != nil {
			log.DebugLogPrint("VM-IPC Socket created: %s -> %s", *acl.Username, socketPath)
		} else {
			log.DebugLogPrint("VM-IPC Socket created: %s -> %s", *acl.Groupname, socketPath)
		}
	}

	return listeners, nil
}

// createListenerWithACL erstellt einen Unix-Socket und wendet die ACL-Einstellungen an.
func createListenerWithACL(socketPath cenvxcore.CoreVmIpcSocketPath, acl ACL) (net.Listener, error) {
	// Existierende Datei entfernen, falls vorhanden
	_ = os.Remove(string(socketPath))

	// Unix Listener erstellen
	listener, err := net.Listen("unix", string(socketPath))
	if err != nil {
		return nil, fmt.Errorf("error creating unix listener: %v", err)
	}

	// Berechtigungen setzen basierend auf ACL
	if acl.Username != nil && acl.Groupname == nil {
		err = filesystem.SetUserFilePermission(string(socketPath), *acl.Username)
	} else if acl.Username != nil && acl.Groupname != nil {
		err = filesystem.SetUserGroupFilePermission(string(socketPath), *acl.Username)
	}
	if err != nil {
		listener.Close()
		return nil, fmt.Errorf("error setting ACL permissions: %v", err)
	}

	return listener, nil
}

// createACLForCurrentUser erstellt ein ACL-Objekt für den aktuellen Benutzer.
func createACLForCurrentUser() (*ACL, error) {
	// Hole den aktuellen Benutzer
	currentUser, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("error getting current user: %v", err)
	}

	// Setze den Benutzernamen
	username := currentUser.Username
	acl := &ACL{
		Username: &username,
	}

	// Versuche, die Gruppeninformationen zu setzen
	groupID := currentUser.Gid
	group, err := user.LookupGroupId(groupID)
	if err == nil {
		// Wenn die Gruppe erfolgreich abgerufen wurde, setze sie in ACL
		groupname := group.Name
		acl.Groupname = &groupname
	} else {
		// Falls die Gruppe nicht gefunden wird, gebe eine Warnung aus und lasse Groupname nil
		log.LogError("Warning: could not resolve primary group for user %s\n", username)
	}

	return acl, nil
}
