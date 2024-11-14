package main

import (
	"fmt"
	"os"
	"runtime"

	cmd "github.com/CustodiaJS/custodiajs-core/cmd"
	"github.com/CustodiaJS/custodiajs-core/core"
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
	hostCryptoStoreDirPath, _, _, _ := cmd.GetPathsAndDirs()

	// Es wird versucht den CryptoStore zu laden,
	// sollte kein Crypto Store vorhanden sein,
	// wird versucht einer zu erstellen
	cryptoStore, cryptoStoreError := crypto.TryToLoad(hostCryptoStoreDirPath)
	if cryptoStoreError != nil {
		fmt.Println(cryptoStoreError.Error())
		os.Exit(1)
	}

	// Der Host Netzwerk Controller wird erstellt
	//ipnetcon := ipnetwork.NewHostNetworkManagmentUnit()

	// Der Core wird erzeugt
	coreInstanceError := core.Init(cryptoStore)
	if coreInstanceError != nil {
		panic(coreInstanceError)
	}

	// Das Hauptprogramm wird offen gehalten
	cmd.RunCoreConsoleOrBackgroundService()
}
