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

package main

import (
	"os"
	"runtime"

	cmd "github.com/custodia-cenv/cenvx-core/src/cmd"
	"github.com/custodia-cenv/cenvx-core/src/core"
	"github.com/custodia-cenv/cenvx-core/src/crypto"
	"github.com/custodia-cenv/cenvx-core/src/log"
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
		log.LogError(cryptoStoreError.Error())
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
