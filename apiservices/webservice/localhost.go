package webservice

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"vnh1/utils"
)

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
	if utils.CHECK_SSL_LOCALHOST_ENABLE {
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
	webs := &Webservice{
		core: nil,
	}

	// Die Daten werden zurückgegeben
	return webs, nil
}
