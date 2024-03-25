package webservice

import (
	"crypto/tls"
	"crypto/x509"
	"embed"
	"fmt"
	"net/http"
	"path"
	"strings"
	"vnh1/static"
)

//go:embed www/*
var wwwEmbedDir embed.FS

type Webservice struct {
}

type NodeStateResponse struct {
	ID   string `json:"id"`
	Wert string `json:"wert"`
}

func (o *Webservice) _handler_index(w http.ResponseWriter, r *http.Request) {
	// Ermittle den Pfad der angeforderten Datei
	requestedFile := strings.TrimPrefix(r.URL.Path, "/console/name")

	// Behandle den Sonderfall, wenn kein spezifischer Dateiname angegeben wurde oder "/console/name/"
	if requestedFile == "" || requestedFile == "/" {
		requestedFile = "/index.html" // Standarddatei
	}

	// Pfad korrigieren, um auf das Verzeichnis innerhalb des eingebetteten Dateisystems zu verweisen
	filePath := "www" + requestedFile

	// Versuche, die angeforderte Datei aus den eingebetteten Ressourcen zu lesen
	fileContent, err := wwwEmbedDir.ReadFile(filePath)
	if err != nil {
		// Datei nicht gefunden oder ein anderer Fehler
		http.Error(w, "Datei nicht gefunden", http.StatusNotFound)
		return
	}

	// Content-Type setzen basierend auf der Dateierweiterung
	switch path.Ext(filePath) {
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".html":
		w.Header().Set("Content-Type", "text/html")
	}

	// Schreibe den Inhalt der Datei in die Antwort
	w.Write(fileContent)
}

func (o *Webservice) Serve(closeSignal chan struct{}) error {
	// Die Basis Urls werden hinzugefügt
	http.HandleFunc("/", o._handler_index)

	// Der HTTP Server wird gestartet
	if err := http.ListenAndServe(":8080", nil); err != nil {
		return fmt.Errorf("Serve: " + err.Error())
	}

	// Der Vorgagn wurde ohne Fehler durchgeführt
	return nil
}

func NewLocalWebservice(ipv4 bool, ipv6 bool, localCert *tls.Certificate) (*Webservice, error) {
	// Parse das Zertifikat aus dem Schlüsselpaar
	x509Cert, err := x509.ParseCertificate(localCert.Certificate[0])
	if err != nil {
		panic(err)
	}

	// Es wird geprüft ob es sich um ein gültiges HostCert handelt,
	// localhost muss in dem Zertifikat vorhanden sein
	// Definiere die zu überprüfenden Domains
	domainsToCheck := []string{"localhost", "127.0.0.1", "::1"}
	for _, domain := range domainsToCheck {
		if err := x509Cert.VerifyHostname(domain); err != nil {
			return nil, fmt.Errorf("NewCore: invalid host cert")
		}
	}

	// Sollte die Funktion nicht deaktiviert wurden sein, so wird jetzt geprüft ob der Host das Verwendete Cert kennt und Validieren kann
	if static.CHECK_SSL_LOCALHOST_ENABLE {
		// Lade den Systemzertifikatsspeicher
		roots, err := x509.SystemCertPool()
		if err != nil {
			panic(err)
		}

		// Verifiziere das Zertifikat gegen den Systemzertifikatsspeicher
		opts := x509.VerifyOptions{
			Roots: roots,
		}

		// Es wird geprüft ob es sich um ein Systembeakanntes Certificate handelt
		if _, err := x509Cert.Verify(opts); err != nil {
			return nil, fmt.Errorf("NewLocalWebservice: " + err.Error())
		}
	} else {
		fmt.Println("Warning: SSL verification for localhost has been completely disabled during compilation.\nThis may lead to unexpected issues, as programs or websites might not be able to communicate with the VNH1 service anymore.\nIf you have downloaded and installed VNH1 and are seeing this message, please be aware that you are not using an official build.")
	}

	// Das Webservice Objekt wird zurückgegeben
	return &Webservice{}, nil
}
