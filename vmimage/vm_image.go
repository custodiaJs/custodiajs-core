package vmimage

import (
	"archive/zip"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/CustodiaJS/custodiajs-core/utils"
	"golang.org/x/crypto/sha3"
)

func (o *VmImage) GetManifest() *Manifest {
	return o.manifest
}

func (o *VmImage) GetMain() *MainJsFile {
	return o.mainFile
}

func TryToLoadVmImage(vmImageFilePath string) (*VmImage, error) {
	// Öffne die ZIP-Datei
	r, err := zip.OpenReader(vmImageFilePath)
	if err != nil {
		fmt.Println("Fehler beim Öffnen der ZIP-Datei:", err)
		return nil, err
	}
	defer r.Close()

	// Flags zum Überprüfen der Existenz der Verzeichnisse und Dateien
	var manifestFile *Manifest

	// Speichert die MainJsFile Datei ab
	var mainJsFile *MainJsFile

	// Speichert die Image Signatur ab
	var imageSignature *ImageSignature

	// Iteriere über die Dateien im ZIP-Archiv
	for _, file := range r.File {
		switch {
		// Prüfe auf die manifest.json Datei
		case file.Name == "manifest.json":
			// Öffne die manifest.json Datei
			rc, err := file.Open()
			if err != nil {
				rc.Close()
				return nil, fmt.Errorf("failed to open manifest.json: %v", err)
			}

			// Lese die JSON-Daten ein, begrenzt auf 10 MB (10 * 1024 * 1024 Bytes)
			const maxSize = 10 * 1024 * 1024
			limitedReader := io.LimitReader(rc, maxSize)

			// Das Manifestobjekt wird eingelesen
			decoder := json.NewDecoder(limitedReader)
			if err := decoder.Decode(&manifestFile); err != nil {
				return nil, fmt.Errorf("failed to decode manifest.json: %v", err)
			}

			// SHA3-256-Hash berechnen
			hash := sha3.New256()
			if _, err := io.Copy(hash, rc); err != nil {
				return nil, fmt.Errorf("failed to hash file: %v", err)
			}

			// Der Hash wird erzeugt und in der Datei zwischengspeichert
			manifestFile.filehash = hex.EncodeToString(hash.Sum(nil))

			// Die Datei wird geschlossen
			rc.Close()
		// Prüfe auf die main.js Datei
		case file.Name == "main.js":
			// Öffne die main.js Datei
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open main.js: %v", err)
			}

			// Lese die Datei ein, begrenzt auf 5 MB (5 * 1024 * 1024 Bytes)
			const maxMainJSSize = 5 * 1024 * 1024
			limitedReader := io.LimitReader(rc, maxMainJSSize)

			// Lese den Inhalt in einen String ein
			content, err := io.ReadAll(limitedReader)
			if err != nil {
				rc.Close()
				return nil, fmt.Errorf("failed to read main.js: %v", err)
			}

			// Es wird ein Hash aus dem Script erzeugt
			scriptHash := utils.HashOfString(string(content))

			// Das Finale MainJsFile Objekt wird erzeugt
			mainJsFile = &MainJsFile{
				fileHash: scriptHash,
				content:  string(content),
				fileSize: file.UncompressedSize64,
			}

			// Die Datei wird geschlossen
			rc.Close()
		// Prüfe auf die Signaturdatei mit dem Namen "signature"
		case file.Name == "signature":
			imageSignature = &ImageSignature{}
		// Prüfe auf den modules-Ordner oder Dateien darin
		case strings.HasPrefix(file.Name, "modules/"):
			if file.FileInfo().IsDir() {
				fmt.Println("MODULE")
			}
		// Prüfe auf den trusted_crypto-Ordner oder Dateien darin
		case strings.HasPrefix(file.Name, "trusted_crypto/"):
			if file.FileInfo().IsDir() {
				fmt.Println("TRUSTED_CRYPTO")
			}
		// Es handelt sich um eine nicht zulässige Datei
		default:
			return nil, fmt.Errorf("invalid file %s", file.Name)
		}
	}

	// Zusätzliche Logik, wenn bestimmte Elemente fehlen
	if manifestFile == nil || mainJsFile == nil || imageSignature == nil {
		return nil, fmt.Errorf("CustodiaJsImage file must contain manifest.json, main.js, and the signature file 'signature'")
	}

	// Das VMImage wird gebaut und zurückgegeben
	vmImage := &VmImage{
		mainFile:  mainJsFile,
		manifest:  manifestFile,
		signature: imageSignature,
	}

	// Das Image wird zurückgegeben
	return vmImage, nil
}
