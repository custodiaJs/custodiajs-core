package core

import (
	"fmt"
	"strings"
	"sync"

	"github.com/CustodiaJS/custodiajs-core/api/http/context"
	"github.com/CustodiaJS/custodiajs-core/core/ipnetwork"

	"github.com/CustodiaJS/custodiajs-core/crypto"
	"github.com/CustodiaJS/custodiajs-core/global/procslog"
	"github.com/CustodiaJS/custodiajs-core/global/static"
	"github.com/CustodiaJS/custodiajs-core/global/types"
	"github.com/CustodiaJS/custodiajs-core/global/utils"
)

// Fügt eine neue API hinzu
func (o *Core) AddVMInstance(vmInstance types.VmInterface, plog_a types.ProcessLogSessionInterface) error {
	// Es wird geprüft das kein Nill Wert übergeben wurde
	if vmInstance == nil {
		return fmt.Errorf("Core->AddVMInstance: null vm instance not allowed")
	}

	// Der Mutex wird angewendet
	o.objectMutex.Lock()

	// Es wird geprüft ob bereits eiein VM Link hinzugefügtne VM mit der Selben ID vorhanden ist
	if _, foundVM := o.vmsByID[string(vmInstance.GetQVMID())]; foundVM {
		o.objectMutex.Unlock()
		return fmt.Errorf("Core->AddNewVMInstance: You cannot add a VM container '%s' multiple times", vmInstance.GetQVMID())
	}

	// Der Mutex wird freigegeben
	o.objectMutex.Unlock()

	// Der Mutex wird angewendet
	o.objectMutex.Lock()

	// Das VMObjekt wird zwischengespeichert
	o.vmsByID[string(vmInstance.GetQVMID())] = vmInstance                    // Merklehash
	o.vmsByName[strings.ToLower(vmInstance.GetManifest().Name)] = vmInstance // VM-Name
	o.vmKernelPtr[vmInstance.GetKId()] = vmInstance                          // Speichert die VM ab, diese wird verwendet um die VM durch den Kernel der VM auffindbar zu machen
	o.vms = append(o.vms, vmInstance)                                        // Die VM wird abgespeichert

	// Der Mutex wird freigegeben
	o.objectMutex.Unlock()

	o.coreLog.Log("New VM Instance added, name = '%s', shash = '%s'", vmInstance.GetManifest().Name, vmInstance.GetScriptHash())

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

// Fügt einen API Socket hinzu
func (o *Core) AddAPISocket(apiSocket types.APISocketInterface, plog_a types.ProcessLogSessionInterface) error {
	// Es wird geprüft das kein Null Wert übergeben wurde
	if apiSocket == nil {
		return fmt.Errorf("Core->AddAPISocket: null api socket not allowed")
	}

	// Der Core wird in dem  Registriert
	err := apiSocket.LinkCore(o)
	if err != nil {
		return fmt.Errorf("AddAPISocket: ")
	}

	// Der Mutex wird angewendet
	// und nach beenden der Funktion freigegeben
	o.objectMutex.Lock()
	defer o.objectMutex.Unlock()

	// Der API Socket wird zwischengespeichert
	o.apiSockets = append(o.apiSockets, apiSocket)
	o.coreLog.Debug("New API Socket added")

	// Es ist kein Fehler aufgetreten
	return nil
}

// Gibt eine Spezifisichen Container VM anhand ihrer ID zurück
func (o *Core) GetVmByID(vmid string, plog_a types.ProcessLogSessionInterface) (types.VmInterface, bool, *types.SpecificError) {
	// Es wird geprüft ob es sich um einen zulässigen vm Namen handelt
	if !utils.ValidateVMIdString(vmid) {
		return nil, false, nil //fmt.Errorf("Core->GetVmByID: invalid vm container id")
	}

	// Die ID wird lowercast
	lowerCaseId := strings.ToLower(vmid)

	// Der Mutex wird angewendet
	// und nach beenden der Funktion freigegeben
	o.objectMutex.Lock()
	defer o.objectMutex.Unlock()

	// Es wird geprüft ob die VM exestiert
	vmObj, found := o.vmsByID[lowerCaseId]
	if !found {
		return nil, false, nil
	}

	// Das Objekt wird zurückgegeben
	return vmObj, true, nil
}

// Gibt eine Spezifisichen Container VM anhand ihrer ID zurück
func (o *Core) GetVmByName(vmName string, plog_a types.ProcessLogSessionInterface) (types.VmInterface, bool, *types.SpecificError) {
	// Es wird geprüft ob es sich um einen zulässigen Namen handelt
	if !utils.ValidateVMName(vmName) {
		return nil, false, nil //fmt.Errorf("Core->GetVmByName: invalid vm container name")
	}

	// Der Name wird lowercast
	lowerCaseVmName := strings.ToLower(vmName)

	// Der Mutex wird angewendet
	// und nach beenden der Funktion freigegeben
	o.objectMutex.Lock()
	defer o.objectMutex.Unlock()

	// Es wird geprüft ob die VM exestiert
	vmObj, found := o.vmsByName[lowerCaseVmName]
	if !found {
		return nil, false, nil // fmt.Errorf("Core->GetVmByName: unkown vm '%s'", lowerCaseVmName)
	}

	// Das Objekt wird zurückgegeben
	return vmObj, true, nil
}

// Gibt die ID's der Aktiven VM-Container zurück
func (o *Core) GetAllActiveVmIDs(plog_a types.ProcessLogSessionInterface) []string {
	// Es wird eine neue Debug einheit erzeugt
	var plog types.ProcessLogSessionInterface
	if plog_a != nil {
		plog = procslog.NewChainMergedProcLog(plog_a, o.coreLog)
	} else {
		plog = o.coreLog
	}

	// DEBUG
	plog.Debug("All active VMs are retrieved")

	// Der Mutex wird angewendet
	// und nach beenden der Funktion freigegeben
	o.objectMutex.Lock()
	defer o.objectMutex.Unlock()

	// Es wird eine Liste mit allen Verfügbaren VM's erstellt
	extr := make([]string, 0)
	for _, item := range o.vmsByID {
		extr = append(extr, string(item.GetQVMID()))
	}

	// DEBUG
	plog.Debug("%d Active VMs were retrieved", len(extr))

	// Die Liste wird zurückgegeben
	return extr
}

// Gibt alle VM-Container zurück
func (o *Core) GetAllVMs(plog_a types.ProcessLogSessionInterface) []types.VmInterface {
	// Der Mutex wird angewendet
	// und nach beenden der Funktion freigegeben
	o.objectMutex.Lock()
	defer o.objectMutex.Unlock()

	// Es wird eine Liste mit allen Verfügbaren VM-Containern erstellt
	extr := make([]types.VmInterface, 0)
	for _, item := range o.vmsByID {
		extr = append(extr, item)
	}

	// Die Liste wird zurückgegeben
	return extr
}

// Gibt die Prozess Managment Unit zurück
func (o *Core) GetCoreSessionManagmentUnit(plog_a types.ProcessLogSessionInterface) types.ContextManagmentUnitInterface {
	return o.cpmu
}

// Gibt das Aktuelle Primäre Host Cert für API Verbindungen zurück
func (o *Core) GetLocalhostCryptoStore(plog_a types.ProcessLogSessionInterface) *crypto.CryptoStore {
	return o.cryptoStore
}

// Erstellt einen neuen CustodiaJS Core
func NewCore(localHostCryptoStore *crypto.CryptoStore, logDIRPath types.LOG_DIR, ipnet *ipnetwork.HostNetworkManagmentUnit) (*Core, error) {
	// Das Coreobjekt wird erstellt
	coreObj := &Core{
		coreLog:     procslog.NewProcLogForCore(),
		cpmu:        context.NewContextManager(),
		vmsByID:     make(map[string]types.VmInterface),
		vmsByName:   make(map[string]types.VmInterface),
		vmKernelPtr: make(map[types.KernelID]types.VmInterface),
		vms:         make([]types.VmInterface, 0),
		apiSockets:  make([]types.APISocketInterface, 0),
		//hostTlsCert:     hostTlsCert,
		cryptoStore: localHostCryptoStore,
		state:       static.NEW,
		// Chans
		holdOpenChan:     make(chan struct{}),
		serviceSignaling: make(chan struct{}),
		vmSyncWaitGroup:  sync.WaitGroup{},
		apiSyncWaitGroup: sync.WaitGroup{},
		// Mutexes
		objectMutex: &sync.Mutex{},
		// Log
		logDIR: logDIRPath,
		// IP-Info Einheit
		hostnetmanager: ipnet,
	}

	// Der Core sowie der Context Manager werden miteinander gepaart
	context.PairCoreToContextManager(coreObj.cpmu, coreObj)

	coreObj.coreLog.Debug("Created")

	// Das Objekt wird zurückgegeben
	return coreObj, nil
}
