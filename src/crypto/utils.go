// Author: fluffelpuff
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
