package crypto

import (
	"crypto/tls"
	"encoding/pem"
	"fmt"
	"os"
)

func loadCertAndPrivateKeyFromOneFile(filepath string) (*tls.Certificate, error) {
	// Das Cert wird eingelesn
	pemData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("loadCertAndPrivateKeyFromOneFile: 1:// " + err.Error())
	}

	// Lädt die PEM Blöcke
	var certData, privKeyData []byte
	for block, rest := pem.Decode(pemData); block != nil; block, rest = pem.Decode(rest) {
		switch block.Type {
		case "CERTIFICATE":
			// Überprüfen, ob bereits ein Zertifikat gefunden wurde
			if certData != nil {
				return nil, fmt.Errorf("mehrere Zertifikate gefunden")
			}
			// Zertifikat-Daten speichern
			certData = pem.EncodeToMemory(block)
		case "PRIVATE KEY", "EC PRIVATE KEY":
			// Überprüfen, ob bereits ein Schlüssel gefunden wurde
			if privKeyData != nil {
				return nil, fmt.Errorf("mehrere private Schlüssel gefunden")
			}
			// Private Schlüssel-Daten speichern
			privKeyData = pem.EncodeToMemory(block)
		default:
			// Unbekannten Blocktyp behandeln
			return nil, fmt.Errorf("unbekannter Blocktyp: %s", block.Type)
		}
	}

	// Es wird geprüft ob ein Certifikat und ein Privater Schlüssel vorhanden sind
	if certData == nil || privKeyData == nil {
		return nil, fmt.Errorf("LoadHostKeyPairs: invlaid host certificate key pair %s", filepath)
	}

	// Es wird versucht das Cert und den Privaten Schlüssel einzulesen
	tlsCert, err := tls.X509KeyPair(certData, privKeyData)
	if err != nil {
		return nil, fmt.Errorf("LoadHostKeyPairs: 2:// " + err.Error())
	}

	// Das Cert wird zurückgegeben
	return &tlsCert, nil
}
