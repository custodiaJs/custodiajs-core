package webservice

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"vnh1/utils"

	"github.com/soheilhy/cmux"
)

func NewLocalWebservice(family string, localport uint32, localCert *tls.Certificate) (*Webservice, error) {
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
		fmt.Println("WARNING: The certificate used cannot be checked locally: ignored")
	}

	// Speichert alle Certs ab
	localHostNames := []string{"localhost", "127.0.0.1", "::1"}

	// Der Eigentliche Server wird estellt
	switch family {
	case "ipv4":
		r, err := NewSpeficAddressWebservice("127.0.0.1", localport, localHostNames, localCert)
		if err != nil {
			return nil, fmt.Errorf("NewLocalWebservice: " + err.Error())
		}
		r.isLocalhost = true
		return r, nil
	case "ipv6":
		r, err := NewSpeficAddressWebservice("[::1]", localport, localHostNames, localCert)
		if err != nil {
			return nil, fmt.Errorf("NewLocalWebservice: " + err.Error())
		}
		r.isLocalhost = true
		return r, nil
	default:
		return nil, fmt.Errorf("NewLocalWebservice: unkown ip family")
	}
}

func NewSpeficAddressWebservice(localIp string, localPort uint32, hostnames []string, localCert *tls.Certificate) (*Webservice, error) {
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

	// Create the main listener.
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("NewSpeficAddressWebservice: " + err.Error())
	}

	// Erstelle eine TLS-Konfiguration
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{*localCert}}

	// Wende die TLS-Konfiguration auf den Listener an
	tlsListener := tls.NewListener(l, tlsConfig)

	tcpMux := cmux.New(tlsListener)
	grpcSocket := tcpMux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpSocket := tcpMux.Match(cmux.HTTP1())

	// Der Servermux Objekt wird erzeugt
	serverMux := http.NewServeMux()

	// Das Serverobjekt wird erzeugt
	httpServer := &http.Server{Handler: serverMux}

	// Das Webservice Objekt wird zurückgegeben
	webs := &Webservice{
		core:         nil,
		cert:         x509Cert,
		serverMux:    serverMux,
		serverObj:    httpServer,
		isLocalhost:  false,
		httpSocket:   httpSocket,
		grpcSocket:   grpcSocket,
		tcpMux:       tcpMux,
		localAddress: &LocalAddress{LocalIP: localIp, LocalPort: localPort},
	}

	// Log
	fmt.Printf("New Webservice API-Socket created on: %s\n", addr)

	// Die Daten werden zurückgegeben
	return webs, nil
}
