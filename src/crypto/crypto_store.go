package crypto

import (
	"crypto/tls"
	"fmt"
	"path"

	cenvxcore "github.com/custodia-cenv/cenvx-core/src"
	"github.com/custodia-cenv/cenvx-core/src/host/filesystem"
)

func TryToLoad(hostCryptoStoreDirPath cenvxcore.HOST_CRYPTOSTORE_WATCH_DIR_PATH) (*CryptoStore, error) {
	// Speichert alle Verfügbaren Zertifikate und Private Keys ab
	hostCertsAndPrivKeys := make([]*HostCertAndOrPrivateKey, 0)

	// LOG
	fmt.Printf("Try to load localhost api certificate and private key from '%s'\n", hostCryptoStoreDirPath)

	// Es wird versucht das Localhost Zertifikat sowie den Privaten Schlüssel zu laden
	localhostCertPath := path.Join(string(hostCryptoStoreDirPath), "localhost.pem")
	if !filesystem.FileExists(localhostCertPath) {
		return nil, fmt.Errorf("TryToLoad: localhost certificate and private key not found")
	}

	// Es wird versucht das Zertifikat sowie den Privaten Schlüssel zu laden
	localhostTLSCert, err := loadCertAndPrivateKeyFromOneFile(localhostCertPath)
	if err != nil {
		return nil, fmt.Errorf("TryToLoad: " + err.Error())
	}

	// LOG
	fmt.Printf("Try to load host keypairs and certificates from '%s'\n", hostCryptoStoreDirPath)

	// Es wird geprüpft ob der SSL Ordner vorhanden ist
	// sollte der SSL Ordner vorhanden sein, werden alle SSL Zertfikate Paare geladen
	sslPath := path.Join(string(hostCryptoStoreDirPath), "ssl")
	if filesystem.FolderExists(sslPath) {
		// Es werden alle Dateien im Ordner aufgelistet
		sslFiles, err := filesystem.WalkDir(sslPath, true)
		if err != nil {
			return nil, fmt.Errorf("LoadHostKeyPairs: " + err.Error())
		}

		// Es werden alle PEM Dateien eingelesen
		for _, certItem := range sslFiles {
			// Es wird geprüft ob es sich um eine .pem Datei handelt, wenn nicht wird sie übersprungen
			if certItem.Extension != ".pem" {
				continue
			}

			// Es wird versucht die Datei einzulesen
			tlsCert, err := loadCertAndPrivateKeyFromOneFile(certItem.Path)
			if err != nil {
				return nil, fmt.Errorf("TryToLoad: " + err.Error())
			}

			// Das Cert Pair wird zwischenegspeicehrt
			hostCertsAndPrivKeys = append(hostCertsAndPrivKeys, &HostCertAndOrPrivateKey{HostTLSKey: tlsCert})

			// Log
			fmt.Printf(" -> Host Certificate Keypair %s from %s added\n", certItem.FileHash, certItem.Path)
		}
	}

	// Das Crypto Store objekt wird zurückgegebn
	return &CryptoStore{localhostIdentPairs: hostCertsAndPrivKeys, localhostTLSCert: localhostTLSCert}, nil
}

func NewVmInstanceCryptoStore() *VmCryptoStore {
	return &VmCryptoStore{CryptoStore: &CryptoStore{localhostIdentPairs: make([]*HostCertAndOrPrivateKey, 0), localhostTLSCert: nil}}
}

func (o *CryptoStore) GetLocalhostAPICertificate() *tls.Certificate {
	return o.localhostTLSCert
}
