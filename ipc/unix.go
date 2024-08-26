package ipc

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateNewUnixSocket_NONE_ROOT_ADMIN(path string) (net.Listener, error) {
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

func CreateNewUnixSocket_ROOT_ADMIN(path string) (net.Listener, error) {
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

func CreateNewUnixClientSocketGRPC_ROOT_ADMIN(path string) (*grpc.ClientConn, error) {
	// Erstellen eines Dialer, um Unix Sockets zu verwenden
	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return net.Dial("unix", addr)
	}

	// Verbindung zum Server über Unix Socket herstellen
	conn, grpcDialError := grpc.Dial(path, grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if grpcDialError != nil {
		return nil, fmt.Errorf("CreateNewUnixClientSocketGRPC_NONE_ROOT_ADMIN: " + grpcDialError.Error())
	}

	// Die Zugriffsrechte für den Unix Path werden festgelegt
	if err := setUnixFilePermissionsForAll(path); err != nil {
		return nil, fmt.Errorf("CreateNewUnixClientSocketGRPC_ROOT_ADMIN: " + err.Error())
	}

	// Der Listener wird zurückgegeben
	return conn, nil
}

func CreateNewUnixClientSocketGRPC_NONE_ROOT_ADMIN(path string) (*grpc.ClientConn, error) {
	// Erstellen eines Dialer, um Unix Sockets zu verwenden
	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return net.Dial("unix", addr)
	}

	// Verbindung zum Server über Unix Socket herstellen
	conn, grpcDialError := grpc.Dial(path, grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if grpcDialError != nil {
		return nil, fmt.Errorf("CreateNewUnixClientSocketGRPC_NONE_ROOT_ADMIN: " + grpcDialError.Error())
	}

	// Die Zugriffsrechte für den Unix Path werden festgelegt
	if err := setUnixFilePermissionsForAll(path); err != nil {
		return nil, fmt.Errorf("CreateNewUnixClientSocketGRPC_NONE_ROOT_ADMINs: " + err.Error())
	}

	// Der Listener wird zurückgegeben
	return conn, nil
}

func CreateNewUnixSocket(path string, userRight types.IPCRight) (net.Listener, error) {
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
		unixSocket, err = CreateNewUnixSocket_NONE_ROOT_ADMIN(path)
	case static.ROOT_ADMIN:
		unixSocket, err = CreateNewUnixSocket_ROOT_ADMIN(path)
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

func CreateNewUnixSocketGRPC(path string, asRoot bool) (*grpc.ClientConn, error) {
	// Sollten Rootrechte benötigt werden
	if asRoot {
		if !(os.Geteuid() == 0) {
			return nil, fmt.Errorf("createNewUnixSocket: you don't have the rights you need")
		}
	}

	// Es wird ermittelt mit welchen Benutzerbrechtigungen der Vorgang durchgeführt werden soll
	var unixSocket *grpc.ClientConn
	var err error
	if !asRoot {
		unixSocket, err = CreateNewUnixClientSocketGRPC_NONE_ROOT_ADMIN(path)
	} else {
		unixSocket, err = CreateNewUnixClientSocketGRPC_ROOT_ADMIN(path)
	}

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if err != nil {
		return nil, fmt.Errorf("createNewUnixSocket: " + err.Error())
	}

	// Der Unix Socket wird zurückgegeben
	return unixSocket, nil
}
