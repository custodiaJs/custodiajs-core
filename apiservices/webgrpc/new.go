package webgrpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"vnh1/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewLocalService(family string, localport uint32, localCert *tls.Certificate) (*WebGrpcService, error) {
	// Parse das Zertifikat aus dem Schlüsselpaar
	x509Cert, err := x509.ParseCertificate(localCert.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("NewLocalWebservice: " + err.Error())
	}

	// Lade den Systemzertifikatsspeicher
	roots, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("NewLocalWebservice: " + err.Error())
	}

	// Verifiziere das Zertifikat gegen den Systemzertifikatsspeicher
	opts := x509.VerifyOptions{
		Roots: roots,
	}

	// Es wird geprüft ob es sich um ein Systembeakanntes Certificate handelt
	if _, err := x509Cert.Verify(opts); err != nil {
		if utils.CHECK_SSL_LOCALHOST_ENABLE {
			return nil, fmt.Errorf("NewLocalWebservice: " + err.Error())
		}
	}

	// Speichert alle Certs ab
	localHostNames := []string{"localhost", "127.0.0.1", "::1"}

	// Der Eigentliche Server wird estellt
	switch family {
	case "ipv4":
		r, err := New("127.0.0.1", localport, localHostNames, localCert)
		if err != nil {
			return nil, fmt.Errorf("NewLocalWebservice: " + err.Error())
		}
		r.isLocalhost = true
		return r, nil
	case "ipv6":
		r, err := New("[::1]", localport, localHostNames, localCert)
		if err != nil {
			return nil, fmt.Errorf("NewLocalWebservice: " + err.Error())
		}
		r.isLocalhost = true
		return r, nil
	default:
		return nil, fmt.Errorf("NewLocalWebservice: unkown ip family")
	}
}

func New(localIp string, localPort uint32, hostnames []string, localCert *tls.Certificate) (*WebGrpcService, error) {
	// Parse das Zertifikat aus dem Schlüsselpaar
	x509Cert, err := x509.ParseCertificate(localCert.Certificate[0])
	if err != nil {
		panic(err)
	}

	// Es wird geprüft ob es sich um ein gültiges HostCert handelt,
	for _, domain := range hostnames {
		if err := x509Cert.VerifyHostname(domain); err != nil {
			return nil, fmt.Errorf("NewSpeficAddressWebservice: invalid host cert")
		}
	}

	// Die Adresse wird erezugt
	addr := fmt.Sprintf("%s:%d", localIp, localPort)

	// Erstelle eine TLS-Konfiguration
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{*localCert}}

	// Erstelle gRPC-Serveroptionen mit der TLS-Konfiguration.
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(tlsConfig)),
	}

	// Starte den gRPC-Server mit den TLS-Optionen.
	server := grpc.NewServer(opts...)

	// Das httpapiObjekt wird zurückgegeben
	wgservice := &WebGrpcService{
		core:         nil,
		cert:         x509Cert,
		serverObj:    server,
		isLocalhost:  false,
		tlsConfig:    tlsConfig,
		localAddress: &LocalAddress{LocalIP: localIp, LocalPort: localPort},
	}

	// Log
	fmt.Printf("New webgrpc service created on: %s\n", addr)

	// Die Daten werden zurückgegeben
	return wgservice, nil
}
