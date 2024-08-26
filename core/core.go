package core

import (
	"fmt"
	"strings"
	"sync"

	"github.com/CustodiaJS/custodiajs-core/context"
	"github.com/CustodiaJS/custodiajs-core/core/ipnetwork"

	"github.com/CustodiaJS/custodiajs-core/crypto"
	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/utils"
	"github.com/CustodiaJS/custodiajs-core/utils/procslog"
)

// Fügt einen API Socket hinzu
func (o *Core) AddAPISocket(apiSocket types.APISocketInterface) error {
	// Es wird geprüft das kein Null Wert übergeben wurde
	if apiSocket == nil {
		return fmt.Errorf("Core->AddAPISocket: null api socket not allowed")
	}

	// Der Core wird in dem  Registriert
	err := apiSocket.SetupCore(o)
	if err != nil {
		return fmt.Errorf("AddAPISocket: ")
	}

	// Der Mutex wird angewendet
	// und nach beenden der Funktion freigegeben
	o.objectMutex.Lock()
	defer o.objectMutex.Unlock()

	// Der API Socket wird zwischengespeichert
	o.apiSockets = append(o.apiSockets, apiSocket)

	// Es ist kein Fehler aufgetreten
	return nil
}

// Gibt eine Spezifisichen Container VM anhand ihrer ID zurück
func (o *Core) GetScriptContainerVMByID(vmid string) (types.VmInterface, bool, *types.SpecificError) {
	// Es wird geprüft ob es sich um einen zulässigen vm Namen handelt
	if !utils.ValidateVMIdString(vmid) {
		return nil, false, nil //fmt.Errorf("Core->GetScriptContainerVMByID: invalid vm container id")
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
func (o *Core) GetScriptContainerByVMName(vmName string) (types.VmInterface, bool, *types.SpecificError) {
	// Es wird geprüft ob es sich um einen zulässigen Namen handelt
	if !utils.ValidateVMName(vmName) {
		return nil, false, nil //fmt.Errorf("Core->GetScriptContainerByVMName: invalid vm container name")
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
		return nil, false, nil // fmt.Errorf("Core->GetScriptContainerByVMName: unkown vm '%s'", lowerCaseVmName)
	}

	// Das Objekt wird zurückgegeben
	return vmObj, true, nil
}

// Gibt die ID's der Aktiven VM-Container zurück
func (o *Core) GetAllActiveScriptContainerIDs(processLog types.ProcessLogSessionInterface) []string {
	// Es wird eine neue Debug einheit erzeugt
	var plog types.ProcessLogSessionInterface
	if processLog != nil {
		plog = processLog.GetChildLog("Core")
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
		extr = append(extr, string(item.GetFingerprint()))
	}

	// DEBUG
	plog.Debug("%d Active VMs were retrieved", len(extr))

	// Die Liste wird zurückgegeben
	return extr
}

// Gibt alle VM-Container zurück
func (o *Core) GetAllVMs() []types.VmInterface {
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
func (o *Core) GetCoreSessionManagmentUnit() types.ContextManagmentUnitInterface {
	return o.cpmu
}

// Gibt das Aktuelle Primäre Host Cert für API Verbindungen zurück
func (o *Core) GetLocalhostCryptoStore() *crypto.CryptoStore {
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
		//extModules:      make(map[string]*external_modules.ExternalModule),
		// Chans
		holdOpenChan:     make(chan struct{}),
		serviceSignaling: make(chan struct{}),
		vmSyncWaitGroup:  sync.WaitGroup{},
		apiSyncWaitGroup: sync.WaitGroup{},
		// Datenbanken
		//hostIdentKeyDatabase: hostIdenKeyDatabase,
		// Mutexes
		objectMutex: &sync.Mutex{},
		// Log
		logDIR: logDIRPath,
		// IP-Info Einheit
		hostnetmanager: ipnet,
	}

	// Der Core sowie der Context Manager werden miteinander gepaart
	context.PairCoreToContextManager(coreObj.cpmu, coreObj)

	// Das Objekt wird zurückgegeben
	return coreObj, nil
}
