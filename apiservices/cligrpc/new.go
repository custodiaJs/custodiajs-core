package cligrpc

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"vnh1/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

func NewTestTCP(certPath, keypath string, userRightState types.CLIUserRight) (*HostCliService, error) {
	// Das Host Cert wird geladen
	cert, err := os.ReadFile(certPath)
	if err != nil {
		panic(err)
	}

	// Der Private Schlüssel wird geladen
	key, err := os.ReadFile(keypath)
	if err != nil {
		panic(err)
	}

	// Erstelle ein TLS-Zertifikat aus den geladenen Dateien
	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		panic(err)
	}

	// Erstelle eine TLS-Konfiguration
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{tlsCert}}

	// Erstelle gRPC-Serveroptionen mit der TLS-Konfiguration.
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(tlsConfig)),
	}

	// Es wird ein neuer gRPC Server erstellt
	grpcServer := grpc.NewServer(opts...)

	// Starte den gRPC-Server auf dem angegebenen Port
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Das HostCLI Objekt wird erstellt
	hcs := &HostCliService{netListner: lis, grpcServer: grpcServer}

	// Das Objekt wird zurückgegeben
	switch runtime.GOOS {
	case "windows":
		fmt.Printf("New cli-grpc service created on '%s' (tcp)\n", ":50051")
	case "darwin", "linux":
		fmt.Printf("New cli-grpc service created on '%s' (tcp)\n", ":50051")
	}
	return hcs, nil
}
