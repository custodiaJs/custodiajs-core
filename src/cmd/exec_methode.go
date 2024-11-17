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

package cmd

import (
	"os"
	"os/user"
	"path"
	"runtime"

	cenvxcore "github.com/custodia-cenv/cenvx-core/src"
	"github.com/custodia-cenv/cenvx-core/src/host"
	"github.com/custodia-cenv/cenvx-core/src/host/filesystem"
	"github.com/custodia-cenv/cenvx-core/src/log"
)

// Wird verwendet alle Verzeichnisse zu ermitteln
func GetPathsAndDirs() (cenvxcore.HOST_CRYPTOSTORE_WATCH_DIR_PATH, cenvxcore.LOG_DIR, cenvxcore.HOST_CONFIG_FILE_PATH, cenvxcore.HOST_CONFIG_PATH) {
	// Erzeugt den Cryptostore Path
	cryptoStorePath := cenvxcore.HOST_CRYPTOSTORE_WATCH_DIR_PATH(path.Join(string(cenvxcore.HOST_CONFIG_DIR_PATH), "crypstore"))

	// Erzeugt den Host Config File Path
	hostConfigFilePath := cenvxcore.HOST_CONFIG_FILE_PATH(path.Join(string(cenvxcore.HOST_CONFIG_DIR_PATH), "config.json"))

	// Gibt die Pfade zurück
	return cryptoStorePath, cenvxcore.DEFAULT_LOGGING_DIR_PATH, hostConfigFilePath, cenvxcore.HOST_CONFIG_DIR_PATH
}

// Gibt den Pfad für die UnixSockets bzw Named Pipes zurück
func GetSocketOrPipeNameOrAddress(root bool) cenvxcore.SOCKET_PATH {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" || runtime.GOOS == "freebsd" || runtime.GOOS == "openbsd" || runtime.GOOS == "netbsd" {
		if !root {
			return cenvxcore.SOCKET_PATH(cenvxcore.NONE_ROOT_UNIX_SOCKET)
		} else {
			return cenvxcore.SOCKET_PATH(cenvxcore.ROOT_UNIX_SOCKET)
		}
	} else {
		return ""
	}
}

// Wird verwendet um zu ermitteln ob es sich um ein unterstützes OS handelt
func OSSupportCheck() {
	// Es wird geprüft ob es sich um Unterstützes OS handelt
	switch runtime.GOOS {
	case "linux":
		if err := host.VerifyLinuxSystem(); err != nil {
			log.InfoLogPrint(err.Error())
			os.Exit(1)
		}
	case "windows":
		if err := host.VerifyWindowsSystem(); err != nil {
			log.InfoLogPrint(err.Error())
			os.Exit(1)
		}
	case "darwin":
		if err := host.VerifyAppleMacOSSystem(); err != nil {
			log.InfoLogPrint(err.Error())
			os.Exit(1)
		}
	case "freebsd", "openbsd", "netbsd":
		if err := host.VerifyBSDSystem(); err != nil {
			log.InfoLogPrint(err.Error())
			os.Exit(1)
		}
	default:
		log.InfoLogPrint("It is an unsupported operating system.")
		os.Exit(1)
	}
}

// Diese Funktion überprüft ob alle benötigten Ordner vorhanden sind
func CheckFolderAndFileStructureOnHost() {
	// Log
	log.InfoLogPrint("Folder and file structure checking...")

	// Es werden alle Pfade abgerufen welche notwendig sind
	hostCryptoStoreDirPath, logDirectoryPath, hostConfigFile, hostConfigBaseDirectoryPath := GetPathsAndDirs()

	// Gibt an wieviele Dateien / Ordner nicht gefunden wurden, muss bei 0 stehen
	totalFoldersNotFound := uint(0)

	// Es wird geprüft ob der Host Config Ordner vorhanden ist
	hasConfigDir := false
	if filesystem.FolderExists(string(hostConfigBaseDirectoryPath)) {
		hasConfigDir = true
	} else {
		log.LogError(" -> Host config directory %s not found\n", hostConfigBaseDirectoryPath)
		totalFoldersNotFound = totalFoldersNotFound + 1
	}

	// Sollte der Config Ordner vorhanden sein, wird seine Substruktur geprüft
	if hasConfigDir {
		// Es wird geprüft ob die Host Config vorhanden ist
		if !filesystem.FileExists(string(hostConfigFile)) {
			log.LogError(" -> Config file %s not found\n", hostConfigFile)
			totalFoldersNotFound = totalFoldersNotFound + 1
		}

		// Es wird geprüft ob die Host Config vorhanden ist
		if !filesystem.FileExists(string(hostConfigFile)) {
			log.LogError(" -> Config file %s not found\n", hostConfigFile)
			totalFoldersNotFound = totalFoldersNotFound + 1
		}

		// Es wird geprüft ob der CryptoStore Ordner vorhanden ist
		hasCryptostoreDirectory := false
		if !filesystem.FolderExists(string(hostCryptoStoreDirPath)) {
			log.LogError(" -> Cryptostore directory %s not found\n", hostCryptoStoreDirPath)
			totalFoldersNotFound = totalFoldersNotFound + 1
		} else {
			hasCryptostoreDirectory = true
		}

		// Es wird geprüft ob die Unterordner des Cryptostore komplett sind
		if hasCryptostoreDirectory {
			// Die Pfade werden erzeugt
			localhostStore, trustedStore, localhostCertKeyFiles := path.Join(string(hostCryptoStoreDirPath), "localhost"), path.Join(string(hostCryptoStoreDirPath), "trusted"), path.Join(string(hostCryptoStoreDirPath), "localhost.pem")

			// Es wirdn geprüft ob der Localhost Ordner vorhanden ist
			if !filesystem.FolderExists(string(localhostStore)) {
				log.LogError(" -> Cryptostore 'localhost' directory %s not found\n", localhostStore)
				totalFoldersNotFound = totalFoldersNotFound + 1
			}

			// Es wird geprüft ob der Trusted ordner vorhanden ist
			if !filesystem.FolderExists(string(trustedStore)) {
				log.LogError(" -> Cryptostore 'trusted' directory %s not found\n", trustedStore)
				totalFoldersNotFound = totalFoldersNotFound + 1
			}

			// Es wird geprüft ob das Localhost Zertifikat sowie der Private Schlüssel vorhanden sind
			if !filesystem.FileExists(localhostCertKeyFiles) {
				log.LogError(" -> Cryptostore has not localhost API-Certificate Keypair %s found\n", localhostCertKeyFiles)
				totalFoldersNotFound = totalFoldersNotFound + 1
			}
		}
	}

	// Es wird geprüft ob das Log Direcotry vorhanden ist
	if !filesystem.FolderExists(string(logDirectoryPath)) {
		// Es wird versucht den Ordner zu erstellen
		if err := filesystem.CreateDirectory(string(logDirectoryPath)); err != nil {
			log.LogError(" -> Logging directory %s not found\n", logDirectoryPath)
			log.LogError(" -> Creating error %s\n", err.Error())
			totalFoldersNotFound = totalFoldersNotFound + 1
		}
	}

	// Es wird geprüft ob 'totalNotFound' gleich 0 ist,
	// wenn nicht wird der Vorgang abgebrochen
	if totalFoldersNotFound != 0 {
		log.LogError("The folder structure is not complete or it is confirmed, the startup process was aborted")
		os.Exit(1)
	}
}

// Zeigt die Aktuellen Host Informationen an
func PrintHostInformations() {
	// Die Linux Informationen werden ausgelesen
	hostInfo, err := host.DetectLinuxDist()
	if err != nil {
		panic(err)
	}

	// Die Host Informationen werden angezigt
	log.InfoLogPrint("Host OS:", hostInfo)

	// Es wird ermittelt ob das Programm in einem Container ausgeführt wird
	isRunningInLinuxContainer := host.IsRunningInContainer()

	// Die Info wird angezeigt
	if isRunningInLinuxContainer {
		log.InfoLogPrint("Running in container: yes")
	} else {
		log.InfoLogPrint("Running in container: no")
	}
}

// Wird verwendet um zu ermittelt ob der Host derzeit im Root Modus ausgeführt wird
func IsRunningAsRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		return false
	}
	uid := currentUser.Uid
	return uid == "0"
}
