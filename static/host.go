package static

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// IsRunningInContainer prüft, ob das Programm in einem Container unter Linux läuft.
func IsRunningInContainer() bool {
	// Prüfe zunächst, ob das Betriebssystem Linux ist
	if runtime.GOOS != "linux" {
		return false
	}

	// Lese den Inhalt von /proc/1/cgroup, da dies in Containern einzigartige Pfade enthält.
	data, err := os.ReadFile("/proc/1/cgroup")
	if err != nil {
		return false
	}

	// Suche nach eindeutigen Zeichenfolgen in den Cgroup-Pfaden, die auf Containerisierung hinweisen könnten.
	content := string(data)
	if strings.Contains(content, "docker") || strings.Contains(content, "kubepods") {
		return true
	}

	// Das Programm wird nicht in einem Container ausgeführt
	return false
}

// readFileContent versucht, eine spezifische Zeile aus einer Datei zu extrahieren, die mit dem Präfix beginnt.
func readFileContent(filePath, prefix string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, prefix) {
			// Extrahiere den Wert ohne Anführungszeichen
			return strings.Trim(line[len(prefix):], "\""), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("nicht gefunden")
}

// DetectLinuxDist versucht, die Linux-Distribution zu erkennen, indem mehrere Dateien geprüft werden.
func DetectLinuxDist() (string, error) {
	// Versuche, /etc/os-release zu lesen
	if dist, err := readFileContent("/etc/os-release", "PRETTY_NAME="); err == nil {
		return dist, nil
	}

	// Versuche, /etc/lsb-release zu lesen
	if dist, err := readFileContent("/etc/lsb-release", "DISTRIB_DESCRIPTION="); err == nil {
		return dist, nil
	}

	// Versuche, /etc/debian_version für Debian spezifisch zu lesen
	if data, err := os.ReadFile("/etc/debian_version"); err == nil {
		return "Debian " + string(data), nil
	}

	return "Unbekannte Distribution", nil
}
