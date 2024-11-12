package core

import (
	"fmt"
	"strings"
	"sync"

	"github.com/CustodiaJS/custodiajs-core/core/ipnetwork"
	"github.com/CustodiaJS/custodiajs-core/crypto"
	"github.com/CustodiaJS/custodiajs-core/global/procslog"
	"github.com/CustodiaJS/custodiajs-core/global/static"
	"github.com/CustodiaJS/custodiajs-core/global/types"
	"github.com/CustodiaJS/custodiajs-core/global/utils"
)

// Erstellt einen neuen CustodiaJS Core
func Init(localHostCryptoStore *crypto.CryptoStore, logDIRPath types.LOG_DIR, ipnet *ipnetwork.HostNetworkManagmentUnit) error {
	coremutex.Lock()

	coreLog = procslog.NewProcLogForCore()
	vmKernelPtr = make(map[types.KernelID]types.VmInterface)
	apiSockets = make([]types.APISocketInterface, 0)
	vmsByName = make(map[string]types.VmInterface)
	vmsByID = make(map[string]types.VmInterface)
	cryptoStore = localHostCryptoStore
	vms = make([]types.VmInterface, 0)
	cstate = static.NEW

	// Chans
	holdOpenChan = make(chan struct{})
	serviceSignaling = make(chan struct{})
	vmSyncWaitGroup = sync.WaitGroup{}
	apiSyncWaitGroup = sync.WaitGroup{}

	// Log
	logDIR = logDIRPath

	// IP-Info Einheit
	hostnetmanager = ipnet

	// Log
	coreLog.Debug("Created")

	// Der VMIPC-Service wird gestartet
	if err := InitVmIpcServer("/tmp", nil, nil); err != nil {
		coremutex.Unlock()
		return err
	}

	coremutex.Unlock()

	// Das Objekt wird zurückgegeben
	return nil
}

// Gibt das Aktuelle Primäre Host Cert für API Verbindungen zurück
func GetLocalhostCryptoStore(plog_a types.ProcessLogSessionInterface) *crypto.CryptoStore {
	return cryptoStore
}

// Gibt alle VM-Container zurück
func GetAllVMs(plog_a types.ProcessLogSessionInterface) []types.VmInterface {
	// Der Mutex wird angewendet
	// und nach beenden der Funktion freigegeben
	coremutex.Lock()
	defer coremutex.Unlock()

	// Es wird eine Liste mit allen Verfügbaren VM-Containern erstellt
	extr := make([]types.VmInterface, 0)
	for _, item := range vmsByID {
		extr = append(extr, item)
	}

	// Die Liste wird zurückgegeben
	return extr
}

// Gibt die ID's der Aktiven VM-Container zurück
func GetAllActiveVmIDs(plog_a types.ProcessLogSessionInterface) []string {
	// Es wird eine neue Debug einheit erzeugt
	var plog types.ProcessLogSessionInterface
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
func GetVmByName(vmName string, plog_a types.ProcessLogSessionInterface) (types.VmInterface, bool, *types.SpecificError) {
	// Es wird geprüft ob es sich um einen zulässigen Namen handelt
	if !utils.ValidateVMName(vmName) {
		return nil, false, nil //fmt.Errorf("Core->GetVmByName: invalid vm container name")
	}

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
func GetVmByID(vmid string, plog_a types.ProcessLogSessionInterface) (types.VmInterface, bool, *types.SpecificError) {
	// Es wird geprüft ob es sich um einen zulässigen vm Namen handelt
	if !utils.ValidateVMIdString(vmid) {
		return nil, false, nil //fmt.Errorf("Core->GetVmByID: invalid vm container id")
	}

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

// Fügt einen API Socket hinzu
func AddAPISocket(apiSocket types.APISocketInterface, plog_a types.ProcessLogSessionInterface) error {
	// Es wird geprüft das kein Null Wert übergeben wurde
	if apiSocket == nil {
		return fmt.Errorf("Core->AddAPISocket: null api socket not allowed")
	}

	/* Der Core wird in dem  Registriert
	err := apiSocket.LinkCore(o)
	if err != nil {
		return fmt.Errorf("AddAPISocket: ")
	}
	*/

	// Der Mutex wird angewendet
	// und nach beenden der Funktion freigegeben
	coremutex.Lock()
	defer coremutex.Unlock()

	// Der API Socket wird zwischengespeichert
	apiSockets = append(apiSockets, apiSocket)
	coreLog.Debug("New API Socket added")

	// Es ist kein Fehler aufgetreten
	return nil
}

// Fügt eine neue API hinzu
func AddVMInstance(vmInstance types.VmInterface, plog_a types.ProcessLogSessionInterface) error {
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
	vmKernelPtr[vmInstance.GetKId()] = vmInstance                          // Speichert die VM ab, diese wird verwendet um die VM durch den Kernel der VM auffindbar zu machen
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

// Gibt an ob der Core Initialisiert wurde
func CoreIsInited() bool {
	return false
}
