package utils

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/sha3"
)

func HashFile(filePath string) ([]byte, error) {
	// Öffne die Datei zum Lesen
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Erstelle eine neue Instanz des SHA3-256 Hashers
	hasher := sha3.New256()

	// Lese die Datei in Teilen und aktualisiere den Hasher nach jedem gelesenen Teil
	buf := make([]byte, 4096) // Puffer für das Lesen in Teilen; Größe kann angepasst werden
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		// Aktualisiere den Hasher mit dem Inhalt des aktuellen Teils
		hasher.Write(buf[:n])
	}

	// Finalisiere den Hash-Prozess und erhalte den resultierenden Hash
	return hasher.Sum(nil), nil
}

func HashOSFile(file *os.File) ([]byte, error) {
	// Lesezeiger zurücksetzen
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("HashOSFile: " + err.Error())
	}

	// Erstelle eine neue Instanz des SHA3-256 Hashers
	hasher := sha3.New256()

	// Lese die Datei in Teilen und aktualisiere den Hasher nach jedem gelesenen Teil
	buf := make([]byte, 4096) // Puffer für das Lesen in Teilen; Größe kann angepasst werden
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		// Aktualisiere den Hasher mit dem Inhalt des aktuellen Teils
		hasher.Write(buf[:n])
	}

	// Lesezeiger zurücksetzen
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("HashOSFile: " + err.Error())
	}

	// Finalisiere den Hash-Prozess und erhalte den resultierenden Hash
	return hasher.Sum(nil), nil
}

func HashOfString(data string) string {
	hasher := sha3.New256()
	hasher.Write([]byte(data))
	varh := hasher.Sum(nil)
	return hex.EncodeToString(varh)
}

func hashPair(hash1, hash2 []byte) []byte {
	pair := append(hash1, hash2...)
	hashed := sha3.Sum256(pair)
	return hashed[:]
}

func BuildMerkleRoot(hexHashes []string) (string, error) {
	if len(hexHashes) == 0 {
		return "", fmt.Errorf("keine Hashes zur Verarbeitung vorhanden")
	}

	// Konvertiere die Hex-Strings in Bytes.
	var nodes [][]byte
	for _, hash := range hexHashes {
		hashBytes, err := hex.DecodeString(hash)
		if err != nil {
			return "", fmt.Errorf("fehler beim Dekodieren des Hex-Strings: %v", err)
		}
		nodes = append(nodes, hashBytes)
	}

	// Baue den Merkle Tree.
	for len(nodes) > 1 {
		var newLevel [][]byte
		for i := 0; i < len(nodes); i += 2 {
			if i+1 < len(nodes) {
				newLevel = append(newLevel, hashPair(nodes[i], nodes[i+1]))
			} else {
				// Für ungerade Anzahl, das letzte Element replizieren.
				newLevel = append(newLevel, hashPair(nodes[i], nodes[i]))
			}
		}
		nodes = newLevel
	}

	return hex.EncodeToString(nodes[0]), nil
}

func BuildStringHashChain(values ...string) (string, error) {
	hashes := make([]string, 0)
	for _, item := range values {
		hashes = append(hashes, HashOfString(item))
	}
	return BuildMerkleRoot(hashes)
}

func ComputeTlsCertFingerprint(tlsCert *tls.Certificate) []byte {
	if len(tlsCert.Certificate) == 0 {
		return nil
	}

	x509Cert, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return nil
	}

	// Berechne den Fingerprint des Zertifikats (hier weiterhin SHA-256)
	hash := sha3.New256()
	_, err = hash.Write(x509Cert.Raw)
	if err != nil {
		return nil
	}
	fingerprintBytes := hash.Sum(nil)

	// Der Fingerabdruck werden zurückgegeben
	return fingerprintBytes
}
