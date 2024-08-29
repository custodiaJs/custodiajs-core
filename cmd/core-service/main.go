package main

import (
	"fmt"
	"os"
	"runtime"

	cmd "github.com/CustodiaJS/custodiajs-core/cmd"
	"github.com/CustodiaJS/custodiajs-core/core"
	"github.com/CustodiaJS/custodiajs-core/core/ipnetwork"
	"github.com/CustodiaJS/custodiajs-core/crypto"
)

func main() {
	// Maximale Anzahl von CPU-Kernen für die Go-Runtime festlegen
	runtime.GOMAXPROCS(1)

	// Der Willkomensbildschrim wird angezeigt
	cmd.ShowBanner()

	// Es wird geprüft ob es sich um Unterstützes OS handelt
	cmd.OSSupportCheck()

	// Es wird geprüft ob die Benötigten Ordner vorhanden sind,
	// sollten nicht alle Ordner vorhanden sein, wird der Vorgang abegrbrochen
	cmd.CheckFolderAndFileStructureOnHost()

	// Die Default Pfade werden ermittelt
	hostCryptoStoreDirPath, _, logDirectoryPath, _, _ := cmd.GetPathsAndDirs()

	// Es wird versucht den CryptoStore zu laden,
	// sollte kein Crypto Store vorhanden sein,
	// wird versucht einer zu erstellen
	cryptoStore, cryptoStoreError := crypto.TryToLoad(hostCryptoStoreDirPath)
	if cryptoStoreError != nil {
		fmt.Println(cryptoStoreError.Error())
		os.Exit(1)
	}

	// Die CLI Sockets werden vorbereitet
	cliSockets, cliSocketsError := cmd.NewCLIHostSockets(false)
	if cliSocketsError != nil {
		panic(cliSocketsError)
	}

	// Der Host Netzwerk Controller wird erstellt
	ipnetcon := ipnetwork.NewHostNetworkManagmentUnit()

	// Der Core wird erzeugt
	coreInstance, coreInstanceError := core.NewCore(cryptoStore, logDirectoryPath, ipnetcon)
	if coreInstanceError != nil {
		panic(coreInstanceError)
	}

	// Die API's werden vorbereitet
	if apiSocketsError := cmd.SetupHostapi(coreInstance); apiSocketsError != nil {
		panic(apiSocketsError)
	}

	// Die API Socket Instanzen werden dem Core hinzugefügt
	for _, coreApiInstance := range cliSockets {
		if err := coreInstance.AddAPISocket(coreApiInstance, nil); err != nil {
			panic(err)
		}
	}

	// Das Hauptprogramm wird offen gehalten
	cmd.RunCoreConsoleOrBackgroundService(coreInstance)
}
