package core

import (
	"fmt"
	"strings"
	"sync"

	cenvxcore "github.com/custodia-cenv/cenvx-core/src"
	"github.com/custodia-cenv/cenvx-core/src/crypto"
	"github.com/custodia-cenv/cenvx-core/src/log"
	"github.com/custodia-cenv/cenvx-core/src/procslog"
	"github.com/custodia-cenv/cenvx-core/src/static"
)

// Erstellt einen neuen CustodiaJS Core
func Init(localHostCryptoStore *crypto.CryptoStore) error {
	// Der Mutex wird verwendet
	coremutex.Lock()

	// Log
	log.InfoLogPrint("Core initializing...")

	// Es wird geprüft ob der Core bereits Initalisiert wurde
	if coreState > 1 && coreState < 3 {
		return fmt.Errorf("core is always inited")
	}

	// Die Laufzeitvariabeln werden festgelegt
	coreLog = procslog.NewProcLogForCore()
	vmsByName = make(map[string]cenvxcore.VmInterface)
	vmsByID = make(map[string]cenvxcore.VmInterface)
	cryptoStore = localHostCryptoStore
	vms = make([]cenvxcore.VmInterface, 0)

	// Chans und Waitgroups
	holdOpenChan = make(chan struct{})
	vmSyncWaitGroup = sync.WaitGroup{}

	// Der VMIPC-Service wird gestartet
	if err := coreInitVmIpcServer("/tmp", []*ACL{}); err != nil {
		coremutex.Unlock()
		return err
	}

	// Der Core Status wird auf Inited geändert
	coreSetState(static.INITED, false)

	// Der Core Mutex wird freigegeben
	coremutex.Unlock()

	// Log
	log.InfoLogPrint("Core Initialized")

	// Das Objekt wird zurückgegeben
	return nil
}

// Gibt das Aktuelle Primäre Host Cert für API Verbindungen zurück
func GetLocalhostCryptoStore(plog_a cenvxcore.ProcessLogSessionInterface) *crypto.CryptoStore {
	return cryptoStore
}

// Gibt alle VM-Container zurück
func GetAllVMs(plog_a cenvxcore.ProcessLogSessionInterface) []cenvxcore.VmInterface {
	// Der Mutex wird angewendet
	// und nach beenden der Funktion freigegeben
	coremutex.Lock()
	defer coremutex.Unlock()

	// Es wird eine Liste mit allen Verfügbaren VM-Containern erstellt
	extr := make([]cenvxcore.VmInterface, 0)
	for _, item := range vmsByID {
		extr = append(extr, item)
	}

	// Die Liste wird zurückgegeben
	return extr
}

// Gibt die ID's der Aktiven VM-Container zurück
func GetAllActiveVmIDs(plog_a cenvxcore.ProcessLogSessionInterface) []string {
	// Es wird eine neue Debug einheit erzeugt
	var plog cenvxcore.ProcessLogSessionInterface
	if plog_a != nil {
		plog = procslog.NewChainMergedProcLog(plog_a, coreLog)
	} else {
		plog = coreLog
	}

	// DEBUG
	plog.Debug("All active VMs are retrieved")

	// Der Mutex wird angewendet
	// und nach beenden der Funktion freigegeben
	coremutex.Lock()
	defer coremutex.Unlock()

	// Es wird eine Liste mit allen Verfügbaren VM's erstellt
	extr := make([]string, 0)
	for _, item := range vmsByID {
		extr = append(extr, string(item.GetQVMID()))
	}

	// DEBUG
	plog.Debug("%d Active VMs were retrieved", len(extr))

	// Die Liste wird zurückgegeben
	return extr
}

// Gibt eine Spezifisichen Container VM anhand ihrer ID zurück
func GetVmByName(vmName string, plog_a cenvxcore.ProcessLogSessionInterface) (cenvxcore.VmInterface, bool, *cenvxcore.SpecificError) {
	// Der Name wird lowercast
	lowerCaseVmName := strings.ToLower(vmName)

	// Der Mutex wird angewendet
	// und nach beenden der Funktion freigegeben
	coremutex.Lock()
	defer coremutex.Unlock()

	// Es wird geprüft ob die VM exestiert
	vmObj, found := vmsByName[lowerCaseVmName]
	if !found {
		return nil, false, nil // fmt.Errorf("Core->GetVmByName: unkown vm '%s'", lowerCaseVmName)
	}

	// Das Objekt wird zurückgegeben
	return vmObj, true, nil
}

// Gibt eine Spezifisichen Container VM anhand ihrer ID zurück
func GetVmByID(vmid string, plog_a cenvxcore.ProcessLogSessionInterface) (cenvxcore.VmInterface, bool, *cenvxcore.SpecificError) {
	// Die ID wird lowercast
	lowerCaseId := strings.ToLower(vmid)

	// Der Mutex wird angewendet
	// und nach beenden der Funktion freigegeben
	coremutex.Lock()
	defer coremutex.Unlock()

	// Es wird geprüft ob die VM exestiert
	vmObj, found := vmsByID[lowerCaseId]
	if !found {
		return nil, false, nil
	}

	// Das Objekt wird zurückgegeben
	return vmObj, true, nil
}

// Fügt eine neue API hinzu
func AddVMInstance(vmInstance cenvxcore.VmInterface, plog_a cenvxcore.ProcessLogSessionInterface) error {
	// Es wird geprüft das kein Nill Wert übergeben wurde
	if vmInstance == nil {
		return fmt.Errorf("Core->AddVMInstance: null vm instance not allowed")
	}

	// Der Mutex wird angewendet
	coremutex.Lock()

	// Es wird geprüft ob bereits eiein VM Link hinzugefügtne VM mit der Selben ID vorhanden ist
	if _, foundVM := vmsByID[string(vmInstance.GetQVMID())]; foundVM {
		coremutex.Unlock()
		return fmt.Errorf("Core->AddNewVMInstance: You cannot add a VM container '%s' multiple times", vmInstance.GetQVMID())
	}

	// Der Mutex wird freigegeben
	coremutex.Unlock()

	// Der Mutex wird angewendet
	coremutex.Lock()

	// Das VMObjekt wird zwischengespeichert
	vmsByID[string(vmInstance.GetQVMID())] = vmInstance                    // Merklehash
	vmsByName[strings.ToLower(vmInstance.GetManifest().Name)] = vmInstance // VM-Name
	vms = append(vms, vmInstance)                                          // Die VM wird abgespeichert

	// Der Mutex wird freigegeben
	coremutex.Unlock()

	coreLog.Log("New VM Instance added, name = '%s', shash = '%s'", vmInstance.GetManifest().Name, vmInstance.GetScriptHash())

	/* Die VM wird mit allen Datenbankdiensten Verknüpft
	for _, item := range vmDbEntry.GetAllDatabaseServices() {
		// Es wird ein neuer Link für die VM erzeugt
		link, err := o.databaseService.GetDBServiceLink(item.GetDatabaseFingerprint())
		if err != nil {
			return nil, fmt.Errorf("Core->AddNewVMInstance: " + err.Error())
		}

		// Der Link für den Datenbank Dienst wird abgespeichert
		if err := vmInstance.AddDatabaseServiceLink(link); err != nil {
			return nil, fmt.Errorf("Core->AddNewVMInstance: " + err.Error())
		}
	}
	*/

	// Das VM Objekt wird zwischengespeichert
	return nil
}

// Signalisiert dem Core, dass er beendet werden soll
func SignalShutdown() {
	// Log
	fmt.Println("Closing CustodiaJS...")

	// Der Mutex wird angewendet
	coremutex.Lock()
	defer coremutex.Unlock()

	// Die Chan wird geschlossen
	close(holdOpenChan)
}

// Gibt an ob der Core Initialisiert wurde
func CoreIsInited() bool {
	coremutex.Lock()
	defer coremutex.Unlock()
	return coreState > 1 && coreState < 3
}

// Legt den Core Status fest
func coreSetState(tstate cenvxcore.CoreState, useMutex bool) {
	// Es wird geprüft ob Mutex verwendet werden sollen
	if useMutex {
		coremutex.Lock()
		defer coremutex.Unlock()
	}

	// Der Neue Status wird gesetzt
	coreState = tstate
}
