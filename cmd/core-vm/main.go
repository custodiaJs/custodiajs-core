package main

import (
	"fmt"
	"os"
	"runtime"

	cmd "github.com/CustodiaJS/custodiajs-core/cmd"
	"github.com/CustodiaJS/custodiajs-core/crypto"
	"github.com/CustodiaJS/custodiajs-core/filesystem"
	"github.com/CustodiaJS/custodiajs-core/vm"
	"github.com/CustodiaJS/custodiajs-core/vmimage"
	"github.com/CustodiaJS/custodiajs-core/vmprocess"
)

func main() {
	// Maximale Anzahl von CPU-Kernen für die Go-Runtime festlegen
	runtime.GOMAXPROCS(1)

	// Es wird ermitelt in welchem Modus das Programm ausgeführt,
	// wenn das Programm mit Root rechten ausgeführt wird,
	// sprechen wir von einer System nahen ausführung.

	// Der Willkomensbildschrim wird angezeigt
	cmd.ShowBanner()

	// Es wird geprüft ob es sich um Unterstützes OS handelt
	cmd.OSSupportCheck()

	// Die Einstellungen der Prozess Inszanz werden abgerufen
	config, coreServiceSocketPath, err := cmd.LoadInstanceConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Es wird geprüft ob die Image Datei exestiert
	if !filesystem.FileExists(config.VmProcessParameters.VMImageFilePath) {
		fmt.Printf("Image file does not exist: %s\n", config.VmProcessParameters.VMImageFilePath)
		os.Exit(1)
	}

	// Es wird geprüft ob das Working Dir exestiert
	if !filesystem.FolderExists(config.VmProcessParameters.VMWorkingDir) {
		fmt.Printf("Working directory does not exist: %s\n", config.VmProcessParameters.VMWorkingDir)
		os.Exit(1)
	}

	// Es wird geprüft ob die Einzelnen HostKeys vorhanden sind
	for _, item := range config.VmProcessParameters.HostKeyCerts {
		if !filesystem.FileExists(item.FilePath) {
			fmt.Printf("Host key file does not exist: %s\n", item.FilePath)
			os.Exit(1)
		}
	}

	// Es wird versucht das Image zu laden
	vmImageInstance, vmImageLoadingErr := vmimage.TryToLoadVmImage(config.VmProcessParameters.VMImageFilePath)
	if vmImageLoadingErr != nil {
		fmt.Println(vmImageLoadingErr)
		os.Exit(1)
	}

	// Es wird ein neuer Crypto Store erzeugt
	cryptoStore := crypto.NewCoreVmCryptoStore()

	// Es wird versucht eine Verbindung mit dem Host Controller aufzubauen
	coreVmProcessInstance, instanceError := vmprocess.NewCoreVmClientProcess(false, coreServiceSocketPath, cryptoStore, vmImageInstance.GetManifest())
	if instanceError != nil {
		fmt.Println(instanceError)
		os.Exit(1)
	}

	// Es wird eine neue VM erstellt
	vmInstance, vmErr := vm.NewCoreVM(coreVmProcessInstance, config.VmProcessParameters.VMWorkingDir, vmImageInstance, "")
	if vmErr != nil {
		fmt.Println(vmErr)
		os.Exit(1)
	}

	// Die VM wird gestartet und am leben erhalten
	vmInstance.Serve(nil)
}
