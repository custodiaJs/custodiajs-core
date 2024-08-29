package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
	"runtime"
	"strings"

	http "github.com/CustodiaJS/custodiajs-core/apiservices/http"
	"github.com/CustodiaJS/custodiajs-core/apiservices/localgrpc"
	"github.com/CustodiaJS/custodiajs-core/core"
	"github.com/CustodiaJS/custodiajs-core/filesystem"
	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/utils"
)

// Wird verwendet alle Verzeichnisse zu ermitteln
func GetPathsAndDirs() (types.HOST_CRYPTOSTORE_WATCH_DIR_PATH, types.VM_DB_DIR_PATH, types.LOG_DIR, types.HOST_CONFIG_FILE_PATH, types.HOST_CONFIG_PATH) {
	switch gos := runtime.GOOS; gos {
	case "darwin":
		// Erzeugt den Cryptostore Path
		cryptoStorePath := types.HOST_CRYPTOSTORE_WATCH_DIR_PATH(path.Join(string(static.DARWIN_HOST_CONFIG_DIR_PATH), "crypstore"))

		// Erzeugt den Host Config File Path
		hostConfigFilePath := types.HOST_CONFIG_FILE_PATH(path.Join(string(static.DARWIN_HOST_CONFIG_DIR_PATH), "config.json"))

		// Gibt die Pfade zurück
		return cryptoStorePath, static.DARWIN_DEFAULT_HOST_VM_DB_DIR_PATH, static.DARWIN_DEFAULT_LOGGING_DIR_PATH, hostConfigFilePath, static.DARWIN_HOST_CONFIG_DIR_PATH
	case "linux":
		// Erzeugt den Cryptostore Path
		cryptoStorePath := types.HOST_CRYPTOSTORE_WATCH_DIR_PATH(path.Join(string(static.LINUX_HOST_CONFIG_DIR_PATH), "crypstore"))

		// Erzeugt den Host Config File Path
		hostConfigFilePath := types.HOST_CONFIG_FILE_PATH(path.Join(string(static.LINUX_HOST_CONFIG_DIR_PATH), "config.json"))

		// Gibt die Pfade zurück
		return cryptoStorePath, static.LINUX_DEFAULT_HOST_VM_DB_DIR_PATH, static.LINUX_DEFAULT_LOGGING_DIR_PATH, hostConfigFilePath, static.LINUX_HOST_CONFIG_DIR_PATH
	default:
		panic("LoadHostKeyPair: unsupported os")
	}
}

// Gibt den Pfad für die UnixSockets bzw Named Pipes zurück
func GetSocketOrPipeNameOrAddress(root bool) types.SOCKET_PATH {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" || runtime.GOOS == "freebsd" || runtime.GOOS == "openbsd" || runtime.GOOS == "netbsd" {
		if !root {
			return types.SOCKET_PATH(static.NONE_ROOT_UNIX_SOCKET)
		} else {
			return types.SOCKET_PATH(static.ROOT_UNIX_SOCKET)
		}
	} else {
		return ""
	}
}

// Wird verwendet um die HostCliServices vorzubereiten
func NewCLIHostSockets(withRoot bool) ([]*localgrpc.HostAPIService, error) {
	// Speichert alle Verfügabren API Instanzen ab
	apiInstances := make([]*localgrpc.HostAPIService, 0)

	// Es wird versucht die Testunit zu erzeugen
	cliAPIInstance, cliAPIInstanceError := localgrpc.New(GetSocketOrPipeNameOrAddress(withRoot), static.NONE_ROOT_ADMIN)

	// Es wird geprüft ob ein Fehler aufgetretn ist
	if cliAPIInstanceError != nil {
		return nil, fmt.Errorf("NewCLIHostSockets: " + cliAPIInstanceError.Error())
	}

	// Die CLI Instanz wird zwischengsepeichert
	apiInstances = append(apiInstances, cliAPIInstance)

	// Es wird geprüft ob Root gewünscht wird
	if withRoot {
		// Sollte Root gewünscht sein, wird Zusätzlich der Root CLI Socket hinzugefügt
		cliRootAPIInstance, cliRootAPIInstanceError := localgrpc.New(GetSocketOrPipeNameOrAddress(false), static.NONE_ROOT_ADMIN)

		// Es wird geprüft ob ein Fehler aufgetretn ist
		if cliRootAPIInstanceError != nil {
			return nil, fmt.Errorf("NewCLIHostSockets: " + cliRootAPIInstanceError.Error())
		}

		// Die CLI API Instanz wird zwischengspeichert
		apiInstances = append(apiInstances, cliRootAPIInstance)
	}

	// Gibt die API Sockets zurück
	return apiInstances, nil
}

// Wird verwendet um die Host API Services bereizustellen
func SetupHostAPIServices(coreinst *core.Core) error {
	// Der Lokale Crypto Store wird abgerufen
	localhostAPICert := coreinst.GetLocalhostCryptoStore(nil)

	// Der Localhost http wird erzeugt
	localhostWebserviceV6, err := http.NewLocalService("ipv6", 8080, localhostAPICert.GetLocalhostAPICertificate())
	if err != nil {
		panic(err)
	}
	localhostWebserviceV4, err := http.NewLocalService("ipv4", 8080, localhostAPICert.GetLocalhostAPICertificate())
	if err != nil {
		panic(err)
	}

	// Der Localhost http wird hinzugefügt
	if err := coreinst.AddAPISocket(localhostWebserviceV6, nil); err != nil {
		panic(err)
	}
	if err := coreinst.AddAPISocket(localhostWebserviceV4, nil); err != nil {
		panic(err)
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

// Wird verwendet um zu ermitteln ob es sich um ein unterstützes OS handelt
func OSSupportCheck() {
	// Es wird geprüft ob es sich um Unterstützes OS handelt
	switch runtime.GOOS {
	case "linux":
		if err := utils.VerifyLinuxSystem(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case "windows":
		if err := utils.VerifyWindowsSystem(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case "darwin":
		if err := utils.VerifyAppleMacOSSystem(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case "freebsd", "openbsd", "netbsd":
		if err := utils.VerifyBSDSystem(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	default:
		fmt.Println("It is an unsupported operating system.")
		os.Exit(1)
	}
}

// Diese Funktion überprüft ob alle benötigten Ordner vorhanden sind
func CheckFolderAndFileStructureOnHost() {
	// Log
	fmt.Println("Folder and file structure checking...")

	// Es werden alle Pfade abgerufen welche notwendig sind
	hostCryptoStoreDirPath, vmDatabaseDirectoryPath, logDirectoryPath, hostConfigFile, hostConfigBaseDirectoryPath := GetPathsAndDirs()

	// Gibt an wieviele Dateien / Ordner nicht gefunden wurden, muss bei 0 stehen
	totalFoldersNotFound := uint(0)

	// Es wird geprüft ob der Host Config Ordner vorhanden ist
	hasConfigDir := false
	if filesystem.FolderExists(string(hostConfigBaseDirectoryPath)) {
		hasConfigDir = true
	} else {
		fmt.Printf(" -> Host config directory %s not found\n", hostConfigBaseDirectoryPath)
		totalFoldersNotFound = totalFoldersNotFound + 1
	}

	// Sollte der Config Ordner vorhanden sein, wird seine Substruktur geprüft
	if hasConfigDir {
		// Es wird geprüft ob die Host Config vorhanden ist
		if !filesystem.FileExists(string(hostConfigFile)) {
			fmt.Printf(" -> Config file %s not found\n", hostConfigFile)
			totalFoldersNotFound = totalFoldersNotFound + 1
		}

		// Es wird geprüft ob die Host Config vorhanden ist
		if !filesystem.FileExists(string(hostConfigFile)) {
			fmt.Printf(" -> Config file %s not found\n", hostConfigFile)
			totalFoldersNotFound = totalFoldersNotFound + 1
		}

		// Es wird geprüft ob der CryptoStore Ordner vorhanden ist
		hasCryptostoreDirectory := false
		if !filesystem.FolderExists(string(hostCryptoStoreDirPath)) {
			fmt.Printf(" -> Cryptostore directory %s not found\n", hostCryptoStoreDirPath)
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
				fmt.Printf(" -> Cryptostore 'localhost' directory %s not found\n", localhostStore)
				totalFoldersNotFound = totalFoldersNotFound + 1
			}

			// Es wird geprüft ob der Trusted ordner vorhanden ist
			if !filesystem.FolderExists(string(trustedStore)) {
				fmt.Printf(" -> Cryptostore 'trusted' directory %s not found\n", trustedStore)
				totalFoldersNotFound = totalFoldersNotFound + 1
			}

			// Es wird geprüft ob das Localhost Zertifikat sowie der Private Schlüssel vorhanden sind
			if !filesystem.FileExists(localhostCertKeyFiles) {
				fmt.Printf(" -> Cryptostore has not localhost API-Certificate Keypair %s found\n", localhostCertKeyFiles)
				totalFoldersNotFound = totalFoldersNotFound + 1
			}
		}
	}

	// Es wird geprüft ob die VM-DB vorhanden ist
	if !filesystem.FolderExists(string(vmDatabaseDirectoryPath)) {
		fmt.Printf(" -> VM-Database directory %s not found\n", vmDatabaseDirectoryPath)
		totalFoldersNotFound = totalFoldersNotFound + 1
	}

	// Es wird geprüft ob das Log Direcotry vorhanden ist
	if !filesystem.FolderExists(string(logDirectoryPath)) {
		// Es wird versucht den Ordner zu erstellen
		if err := filesystem.CreateDirectory(string(logDirectoryPath)); err != nil {
			fmt.Printf(" -> Logging directory %s not found\n", logDirectoryPath)
			fmt.Printf(" -> Creating error %s\n", err.Error())
			totalFoldersNotFound = totalFoldersNotFound + 1
		}
	}

	// Es wird geprüft ob 'totalNotFound' gleich 0 ist,
	// wenn nicht wird der Vorgang abgebrochen
	if totalFoldersNotFound != 0 {
		fmt.Println("The folder structure is not complete or it is confirmed, the startup process was aborted")
		os.Exit(1)
	}
}

// Zeigt die Aktuellen Host Informationen an
func PrintHostInformations() {
	// Die Linux Informationen werden ausgelesen
	hostInfo, err := utils.DetectLinuxDist()
	if err != nil {
		panic(err)
	}

	// Die Host Informationen werden angezigt
	fmt.Println("Host OS:", hostInfo)

	// Es wird ermittelt ob das Programm in einem Container ausgeführt wird
	isRunningInLinuxContainer := utils.IsRunningInContainer()

	// Die Info wird angezeigt
	if isRunningInLinuxContainer {
		fmt.Println("Running in container: yes")
	} else {
		fmt.Println("Running in container: no")
	}
}

// Wird verwendet um die Parameter auszlesen (Core-VM)
func ReadParametersVmInstance() (*types.VmInstanceProcessParameters, error) {
	// Nutzungshinweis
	usage := `Usage of ./build/core-vm:
	--image string
		  Path to VM-Image
	--workdir string
		  Path to working directory
	--hostkeycert string string [--alias string]
		  Host key certificate with algorithm and file path. Optionally specify an alias for the key with --alias, required for ecdsa and rsa algorithms.`

	// Die verfügbaren Argumente werden als Variablen deklariert
	var vmImageFilePath string // Gibt den vollständigen Pfad zur Datei an
	var vmWorkingDir string    // Gibt den vollständigen Pfad zum Arbeitsverzeichnis an
	var hostKeyCerts []types.HostKeyCert
	disableCoreCrypto := false

	// Manuelles Parsen der Argumente
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "--image=") {
			vmImageFilePath = strings.TrimPrefix(args[i], "--image=")
		} else if args[i] == "--image" {
			if i+1 >= len(args) {
				return nil, errors.New("fehlender Wert für --image")
			}
			vmImageFilePath = args[i+1]
			i++
		} else if strings.HasPrefix(args[i], "--workdir=") {
			vmWorkingDir = strings.TrimPrefix(args[i], "--workdir=")
		} else if args[i] == "--workdir" {
			if i+1 >= len(args) {
				return nil, errors.New("fehlender Wert für --workdir")
			}
			vmWorkingDir = args[i+1]
			i++
		} else if args[i] == "--hostkeycert" {
			if i+2 >= len(args) {
				return nil, errors.New("invalid use of --hostkeycert: algorithm and file path required")
			}
			algorithm := args[i+1]
			filePath := args[i+2]
			var alias string

			// Überprüfen, ob ein optionales alias angegeben wurde
			if i+3 < len(args) && strings.HasPrefix(args[i+3], "alias=") {
				alias = strings.TrimPrefix(args[i+3], "alias=")
				i += 3 // Überspringe den Algorithmus, den Dateipfad und den Alias
			} else {
				i += 2 // Überspringe nur den Algorithmus und den Dateipfad
			}

			// Prüfen, ob ein Alias erforderlich ist
			if (strings.ToLower(algorithm) == "ecdsa" || strings.ToLower(algorithm) == "rsa") && alias == "" {
				return nil, fmt.Errorf("alias is required for algorithm %s but not provided", algorithm)
			}

			hostKeyCerts = append(hostKeyCerts, types.HostKeyCert{
				Algorithm: algorithm,
				FilePath:  filePath,
				Alias:     alias,
			})
		} else if args[i] == "--disablecorecrypto" {
			disableCoreCrypto = true
		} else {
			return nil, fmt.Errorf("unkown argument: %s\n%s", args[i], usage)
		}
	}

	// Es wird geprüft, ob die Pflichtparameter vorhanden sind
	if vmImageFilePath == "" || vmWorkingDir == "" {
		return nil, fmt.Errorf("required parameters are missing: Please make sure that --image and --workdir are set")
	}

	// Rückgabe der eingelesenen Parameter als VmInstanceProcessParameters (angenommen, dass dieser Typ existiert)
	params := &types.VmInstanceProcessParameters{
		VMImageFilePath:   vmImageFilePath,
		VMWorkingDir:      vmWorkingDir,
		HostKeyCerts:      hostKeyCerts,
		DisableCoreCrypto: disableCoreCrypto,
	}

	// Die Parameter werden zurückgegeben
	return params, nil
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

// Wird verwender um die Instanzeinstellungen abzurufen
func LoadInstanceConfig() (*types.VmInstanceInstanceConfig, types.SOCKET_PATH, error) {
	// Die Programmparameter werden ausgelesen
	parms, err := ReadParametersVmInstance()
	if err != nil {
		return nil, "", fmt.Errorf("LoadInstanceConfig: " + err.Error())
	}

	// Der Socketpath wird ermittelt
	socketPath := GetSocketOrPipeNameOrAddress(IsRunningAsRoot())

	// Das Objket wird zurückgegeben
	return &types.VmInstanceInstanceConfig{VmProcessParameters: parms}, socketPath, nil
}
