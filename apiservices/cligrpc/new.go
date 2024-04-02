package cligrpc

import (
	"fmt"
	"net"
	"runtime"
	"vnh1/types"

	"google.golang.org/grpc"
)

func New(unixOrWinNamedPipeAddr string, userRightState types.CLIUserRight) (*HostCliService, error) {
	// Es wird passend zum Hostos der Richtige Listener erzeugt
	var cliSocket net.Listener
	var err error
	switch runtime.GOOS {
	case "windows":
		err = fmt.Errorf("not supported os")
	case "darwin":
		cliSocket, err = createNewUnixSocket(unixOrWinNamedPipeAddr, userRightState)
	case "linux":
		cliSocket, err = createNewUnixSocket(unixOrWinNamedPipeAddr, userRightState)
	default:
		err = fmt.Errorf("unkown os")
	}

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if err != nil {
		return nil, fmt.Errorf("New: " + err.Error())
	}

	// Es wird ein neuer gRPC Server erstellt
	grpcServer := grpc.NewServer()

	// Das HostCLI Objekt wird erstellt
	hcs := &HostCliService{netListner: cliSocket, grpcServer: grpcServer}

	// Das Objekt wird zurückgegeben
	switch runtime.GOOS {
	case "windows":
		fmt.Printf("New cli-grpc service created on: '%s' (windows named pipe) \n", unixOrWinNamedPipeAddr)
	case "darwin", "linux":
		fmt.Printf("New cli-grpc service created on: '%s' (unix socket)\n", unixOrWinNamedPipeAddr)
	}
	return hcs, nil
}
