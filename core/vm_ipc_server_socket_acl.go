package core

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"strconv"
	"syscall"

	"github.com/CustodiaJS/custodiajs-core/log"
)

// createAclListeners erstellt für jede ACL in der Liste einen _AclListener und gibt sie zurück.
func createAclListeners(aclList []*ACL, basePath string) ([]*_AclListener, error) {
	var listeners []*_AclListener

	for i, acl := range aclList {
		socketPath := fmt.Sprintf("%s/socket_%d.sock", basePath, i)

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
		if acl.Groupname != nil && acl.Username != nil {
			log.DebugLogPrint("VM-IPC Socket created: %s - %s", *acl.Username, *acl.Groupname)
		} else if acl.Username != nil {
			log.DebugLogPrint("VM-IPC Socket created: %s", *acl.Username)
		} else {
			log.DebugLogPrint("VM-IPC Socket created: %s", *acl.Groupname)
		}
	}

	return listeners, nil
}

// createListenerWithACL erstellt einen Unix-Socket und wendet die ACL-Einstellungen an.
func createListenerWithACL(socketPath string, acl ACL) (net.Listener, error) {
	// Existierende Datei entfernen, falls vorhanden
	_ = os.Remove(socketPath)

	// Unix Listener erstellen
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("error creating unix listener: %v", err)
	}

	// Berechtigungen setzen basierend auf ACL
	if acl.Username != nil && acl.Groupname == nil {
		err = setUserPermission(socketPath, *acl.Username)
	} else if acl.Username != nil && acl.Groupname != nil {
		err = setUserGroupPermission(socketPath, *acl.Username)
	}
	if err != nil {
		listener.Close()
		return nil, fmt.Errorf("error setting ACL permissions: %v", err)
	}

	return listener, nil
}

func setUserPermission(socketPath string, username string) error {
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

func setUserGroupPermission(socketPath string, username string) error {
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
		fmt.Printf("Warning: could not resolve primary group for user %s\n", username)
	}

	return acl, nil
}
