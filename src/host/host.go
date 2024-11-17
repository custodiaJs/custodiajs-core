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

package host

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

func VerifyLinuxSystem() error {
	return nil
}

func VerifyWindowsSystem() error {
	return nil
}

func VerifyAppleMacOSSystem() error {
	return nil
}

func VerifyBSDSystem() error {
	return nil
}

// checkAdmin überprüft, ob das Programm mit Administrator-Rechten ausgeführt wird
func CheckAdmin() bool {
	// Für Unix-basierte Systeme (Linux, macOS)
	return os.Geteuid() == 0
}
