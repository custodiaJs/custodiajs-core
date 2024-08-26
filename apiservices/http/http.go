package http

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	"github.com/CustodiaJS/custodiajs-core/apiservices/http/middleware"
	"github.com/CustodiaJS/custodiajs-core/apiservices/http/middlewares"
	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/types"
)

// Erstellt einen API Service
func New(localIp string, localPort uint32, localCert *tls.Certificate, middlewareHandlers []middleware.MiddlewareFunction) (*HttpApiService, error) {
	// Parse das Zertifikat aus dem Schlüsselpaar
	x509Cert, err := x509.ParseCertificate(localCert.Certificate[0])
	if err != nil {
		panic(err)
	}

	// Die Adresse wird erezugt
	addr := fmt.Sprintf("%s:%d", localIp, localPort)

	// Erstelle eine TLS-Konfiguration
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{*localCert}}

	// Der Servermux Objekt wird erzeugt
	serverMux := http.NewServeMux()

	// Das Serverobjekt wird erzeugt
	httpServer := &http.Server{Addr: addr, TLSConfig: tlsConfig}

	// Das httpapiObjekt wird zurückgegeben
	webs := &HttpApiService{
		core:               nil,
		cert:               x509Cert,
		serverMux:          serverMux,
		serverObj:          httpServer,
		isLocalhost:        false,
		tlsConfig:          tlsConfig,
		middlewareHandlers: middlewareHandlers,
		localAddress:       &LocalAddress{LocalIP: localIp, LocalPort: localPort},
	}

	// Die Globale Middlware wird erzeugt
	globalMiddleware := middleware.GlobalMiddleware(webs.serverMux, webs.middlewareHandlers, webs.core)

	// Es wird ein neuer Handler Registriert
	webs.serverObj.Handler = webs.newSessionHandler(globalMiddleware)

	// Log
	fmt.Printf("New http created on: %s\n", addr)

	// Die Daten werden zurückgegeben
	return webs, nil
}

// Erstellt einen neuen Lokalen API Service
func NewLocalService(family string, localport uint32, localCert *tls.Certificate) (*HttpApiService, error) {
	// Parse das Zertifikat aus dem Schlüsselpaar
	x509Cert, err := x509.ParseCertificate(localCert.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("NewLocalWebservice: CERT_LOADING:: " + err.Error())
	}

	// Lade den Systemzertifikatsspeicher
	roots, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("NewLocalWebservice: SYSTEM_CERT_POOL:: " + err.Error())
	}

	// Speichert alle Certs ab
	localHostNames := []string{"localhost", "127.0.0.1", "::1"}

	// Es wird geprüft ob es sich um ein gültiges HostCert handelt,
	for _, domain := range localHostNames {
		if err := x509Cert.VerifyHostname(domain); err != nil {
			return nil, fmt.Errorf("NewSpeficAddressWebservice: invalid host cert '%s'", domain)
		}
	}

	// Verifiziere das Zertifikat gegen den Systemzertifikatsspeicher
	opts := x509.VerifyOptions{
		Roots: roots,
	}

	// Es wird geprüft ob es sich um ein Systembeakanntes Certificate handelt
	if _, err := x509Cert.Verify(opts); err != nil {
		if static.CHECK_SSL_LOCALHOST_ENABLE {
			return nil, fmt.Errorf("NewLocalWebservice: CERT_VERIFY:: " + err.Error())
		}
	}

	// Der Eigentliche Server wird estellt
	switch family {
	case "ipv4":
		r, err := New("127.0.0.1", localport, localCert, middleware.MiddlewareFunctionList{middlewares.IsLocalhostPOSTRequest, middleware.ValidateRequestDomain(localHostNames)})
		if err != nil {
			return nil, fmt.Errorf("NewLocalWebservice: " + err.Error())
		}
		r.isLocalhost = true
		return r, nil
	case "ipv6":
		r, err := New("[::1]", localport, localCert, middleware.MiddlewareFunctionList{middlewares.IsLocalhostPOSTRequest, middleware.ValidateRequestDomain(localHostNames)})
		if err != nil {
			return nil, fmt.Errorf("NewLocalWebservice: " + err.Error())
		}
		r.isLocalhost = true
		return r, nil
	default:
		return nil, fmt.Errorf("NewLocalWebservice: unkown ip family")
	}
}

// Verknüft den Core mit dem API Service
func (o *HttpApiService) SetupCore(coreObj types.CoreInterface) error {
	// Es wird geprüft ob der Core festgelegt wurde
	if o.core != nil {
		return fmt.Errorf("SetupCore: always linked with core")
	}

	// Das Objekt wird zwischengespeichert
	o.core = coreObj

	// Der Vorgang ist ohne fehler durchgeführt wurden
	return nil
}

// Führt den API Service aus
func (o *HttpApiService) Serve(closeSignal chan struct{}) error {
	// Die Basis Urls werden hinzugefügt
	o.serverMux.Handle("/", middleware.RequestMiddleware(o.httpIndex, IndexMiddlewares, o.core))

	// Gibt die einzelnenen VM Informationen aus
	o.serverMux.Handle("/vm", middleware.RequestMiddleware(o.httpVmInfo, VmMiddlewares, o.core))

	// Der VM-RPC Handler wird erstellt
	o.serverMux.Handle("/rpc", middleware.RequestMiddleware(o.httpRPC, RpcMiddlewares, o.core))

	// Der Webserver wird gestartet
	if err := o.serverObj.ListenAndServeTLS("", ""); err != nil {
		return fmt.Errorf("HttpApiService->Serve: " + err.Error())
	}

	// Der Vorgagn wurde ohne Fehler durchgeführt
	return nil
}
