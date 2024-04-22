package localgrpc

import (
	"fmt"
	"net"
	"os"
	"vnh1/static"
	"vnh1/types"
)

func createNewUnixSocket_NONE_ROOT_ADMIN(path string) (net.Listener, error) {
	// Es wird geprüft ob die Datei vorhanden ist, wenn ja wird versucht die Datei zu löschen
	if err := deleteFileIfExists(path); err != nil {
		return nil, fmt.Errorf("createNewUnixSocket_NONE_ROOT_ADMIN: " + err.Error())
	}

	// Der UnixSocket wird estellt
	lis, err := net.Listen("unix", path)
	if err != nil {
		return nil, fmt.Errorf("createNewUnixSocket_NONE_ROOT_ADMIN: " + err.Error())
	}

	// Die Zugriffsrechte für den Unix Path werden festgelegt
	if err := setUnixFilePermissionsForAll(path); err != nil {
		return nil, fmt.Errorf("createNewUnixSocket_NONE_ROOT_ADMIN: " + err.Error())
	}

	// Der Listener wird zurückgegeben
	return lis, nil
}

func createNewUnixSocket_ROOT_ADMIN(path string) (net.Listener, error) {
	// Es wird geprüft ob die Datei vorhanden ist, wenn ja wird versucht die Datei zu löschen
	if err := deleteFileIfExists(path); err != nil {
		return nil, fmt.Errorf("createNewUnixSocket_ROOT_ADMIN: " + err.Error())
	}

	// Der UnixSocket wird estellt
	lis, err := net.Listen("unix", path)
	if err != nil {
		return nil, fmt.Errorf("createNewUnixSocket_ROOT_ADMIN: " + err.Error())
	}

	// Die Zugriffsrechte für den Unix Path werden festgelegt
	if err := setUnixFileOwnerToRoot(path); err != nil {
		return nil, fmt.Errorf("createNewUnixSocket_ROOT_ADMIN: " + err.Error())
	}

	// Der Listener wird zurückgegeben
	return lis, nil
}

func createNewUnixSocket(path string, userRight types.CLIUserRight) (net.Listener, error) {
	// Sollten Rootrechte benötigt werden
	if userRight == static.ROOT_ADMIN {
		if !(os.Geteuid() == 0) {
			return nil, fmt.Errorf("createNewUnixSocket: you don't have the rights you need")
		}
	}

	// Es wird ermittelt mit welchen Benutzerbrechtigungen der Vorgang durchgeführt werden soll
	var unixSocket net.Listener
	var err error
	switch userRight {
	case static.NONE_ROOT_ADMIN:
		unixSocket, err = createNewUnixSocket_NONE_ROOT_ADMIN(path)
	case static.ROOT_ADMIN:
		unixSocket, err = createNewUnixSocket_ROOT_ADMIN(path)
	default:
		return nil, fmt.Errorf("createNewUnixSocket: unkown user right")
	}

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if err != nil {
		return nil, fmt.Errorf("createNewUnixSocket: " + err.Error())
	}

	// Der Unix Socket wird zurückgegeben
	return unixSocket, nil
}
