package core

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"vnh1/core/databaseservices"
	"vnh1/core/identkeydatabase"
	"vnh1/core/vmdb"
	extmodules "vnh1/extmodules"
	"vnh1/static"
	"vnh1/types"
	"vnh1/utils"
)

// Fügt eine Externe Modul Lib dem Core hinzu
func (o *Core) AddExternalModuleLibrary(modLib *extmodules.ExternalModule) error {
	// Es wird geprüft ob es sich um einen Zulässigen Module namen handelt
	if val := utils.ValidateExternalModuleName(modLib.GetName()); !val {
		return fmt.Errorf("Core->AddExternalModuleLibrary: Invalid module name, cant added module")
	}

	// Der Mutex wird angewendet
	o.objectMutex.Lock()

	// Es wird ermittelt ob es bereits ein Externes Module mit dem gleichen Namen gibt
	if _, found := o.extModules[modLib.GetName()]; found {
		o.objectMutex.Unlock()
		return fmt.Errorf("Core->AddExternalModuleLibrary: module always added")
	}

	// Die Module Lib wird zwischengspeichert
	o.extModules[modLib.GetName()] = modLib

	// Der Mutex wird freigegeben
	o.objectMutex.Unlock()

	// Es ist kein Fehler aufgetreten
	return nil
}

// Fügt einen neune Script Container hinzu
func (o *Core) AddScriptContainer(vmDbEntry *vmdb.VmDBEntry) (*CoreVM, error) {
	// Die Virtuelle Maschine wird geprüft
	if !vmDbEntry.ValidateVM() {
		return nil, fmt.Errorf("AddScriptContainer: Broken Virtual Machine")
	}

	// Es werden alle benötigten Host CA Extrahiert
	// sollte kein Passendes Gefunden werden, wird der Vorgang abgebrochen
	for _, item := range vmDbEntry.GetRootMemberIDS() {
		// Es wird ermittelt ob es sich um ein SSL Eintrag handelt, wenn nicht wird dieser Ignoriert
		if item.Type != "ssl" {
			continue
		}

		// Es wird ermittelt ob der Host ein SSL-Cert/Privatekey paar besitzt welches von dem Aktuellen RootCA Signiert wurde
		// wenn nicht wird geprüft ob der Fingerabdruck des CERTS mit dem des Lokalen Certs übereinstimmt
		if !o.hostIdentKeyDatabase.ValidateRootCAMembershipByFingerprint(item.Fingerprint) {
			if !strings.EqualFold(hex.EncodeToString(utils.ComputeTlsCertFingerprint(o.hostTlsCert)), item.Fingerprint) {
				return nil, fmt.Errorf("Core->AddScriptContainer: unkown host ca membership '%s'", strings.ToUpper(item.Fingerprint))
			}
		}
	}

	// Es wird eine Liste mit allen Benötigten externen Libs abgerufen
	neededExternalModulesNameSlice := make([]string, 0)
	for _, item := range vmDbEntry.GetAllExternalServices() {
		neededExternalModulesNameSlice = append(neededExternalModulesNameSlice, item.Name)
	}

	// Es werden alle Module welche benötigt werden abgerufen
	modList := o._core_util_get_list_of_extmods_by_name(neededExternalModulesNameSlice...)

	// Es wird geprüft ob die benötigten Module gefunden wurden
	notFoundExtModules := make([]string, 0)
	for _, item := range vmDbEntry.GetAllExternalServices() {
		if item.Required {
			found := false
			for _, xtem := range modList {
				if xtem.GetName() == item.Name {
					if xtem.GetVersion() >= uint64(item.MinVersion) {
						found = true
						break
					}
				}
			}
			if !found {
				notFoundExtModules = append(notFoundExtModules, item.Name)
			}
		}
	}

	// Es wird ein Fehler ausgelöst wenn ein benötigtes Modul nicht gefunden wurde
	if len(notFoundExtModules) != 0 {
		return nil, fmt.Errorf("Core->AddScriptContainer: external modules '%s' not found", strings.Join(neededExternalModulesNameSlice, ","))
	}

	// Das Logging Verzeichniss wird erstellt
	logPath, err := utils.MakeLogDirForVM(o.logDIR, vmDbEntry.GetVMName())
	if err != nil {
		return nil, fmt.Errorf("Core->AddScriptContainer: " + err.Error())
	}

	// Der Mutex wird angewendet
	o.objectMutex.Lock()

	// Es wird geprüft ob bereits eiein VM Link hinzugefügtne VM mit der Selben ID vorhanden ist
	if _, foundVM := o.vmsByID[strings.ToLower(vmDbEntry.GetVMContainerMerkleHash())]; foundVM {
		o.objectMutex.Unlock()
		return nil, fmt.Errorf("Core->AddScriptContainer: You cannot add a VM container '%s' multiple times", vmDbEntry.GetVMContainerMerkleHash())
	}

	// Das Detailspaket wird erzeugt
	vmobject, err := newCoreVM(o, vmDbEntry, modList, logPath)
	if err != nil {
		return nil, fmt.Errorf("AddScriptContainer: " + err.Error())
	}

	// Das VMObjekt wird zwischengespeichert
	o.vmsByID[strings.ToLower(vmDbEntry.GetVMContainerMerkleHash())] = vmobject // Merklehash
	o.vmsByName[strings.ToLower(vmDbEntry.GetVMName())] = vmobject              // VM-Name
	o.vmKernelPtr[vmobject.GetKId()] = vmobject                                 // Speichert die VM ab, diese wird verwendet um die VM durch den Kernel der VM auffindbar zu machen
	o.vms = append(o.vms, vmobject)                                             // Die VM wird abgespeichert

	// Der Mutex wird freigegeben
	o.objectMutex.Unlock()

	// Die VM wird mit allen Datenbankdiensten Verknüpft
	for _, item := range vmDbEntry.GetAllDatabaseServices() {
		// Es wird ein neuer Link für die VM erzeugt
		link, err := o.databaseService.GetDBServiceLink(item.GetDatabaseFingerprint())
		if err != nil {
			return nil, fmt.Errorf("Core->AddScriptContainer: " + err.Error())
		}

		// Der Link für den Datenbank Dienst wird abgespeichert
		if err := vmobject.addDatabaseServiceLink(link); err != nil {
			return nil, fmt.Errorf("Core->AddScriptContainer: " + err.Error())
		}
	}

	// Das VM Objekt wird zwischengespeichert
	return vmobject, nil
}

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
func (o *Core) GetScriptContainerVMByID(vmid string) (types.CoreVMInterface, bool, error) {
	// Es wird geprüft ob es sich um einen zulässigen vm Namen handelt
	if !utils.ValidateVMIdString(vmid) {
		return nil, false, fmt.Errorf("Core->GetScriptContainerVMByID: invalid vm container id")
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
func (o *Core) GetScriptContainerByVMName(vmName string) (types.CoreVMInterface, error) {
	// Es wird geprüft ob es sich um einen zulässigen Namen handelt
	if !utils.ValidateVMName(vmName) {
		return nil, fmt.Errorf("Core->GetScriptContainerByVMName: invalid vm container name")
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
		return nil, fmt.Errorf("Core->GetScriptContainerByVMName: unkown vm '%s'", lowerCaseVmName)
	}

	// Das Objekt wird zurückgegeben
	return vmObj, nil
}

// Gibt die ID's der Aktiven VM-Container zurück
func (o *Core) GetAllActiveScriptContainerIDs() []string {
	// Der Mutex wird angewendet
	// und nach beenden der Funktion freigegeben
	o.objectMutex.Lock()
	defer o.objectMutex.Unlock()

	// Es wird eine Liste mit allen Verfügbaren VM's erstellt
	extr := make([]string, 0)
	for _, item := range o.vmsByID {
		extr = append(extr, string(item.GetFingerprint()))
	}

	// Die Liste wird zurückgegeben
	return extr
}

// Gibt alle VM-Container zurück
func (o *Core) GetAllVMs() []types.CoreVMInterface {
	// Der Mutex wird angewendet
	// und nach beenden der Funktion freigegeben
	o.objectMutex.Lock()
	defer o.objectMutex.Unlock()

	// Es wird eine Liste mit allen Verfügbaren VM-Containern erstellt
	extr := make([]types.CoreVMInterface, 0)
	for _, item := range o.vmsByID {
		extr = append(extr, item)
	}

	// Die Liste wird zurückgegeben
	return extr
}

// Erstellt einen neuen vnh1 Core
func NewCore(hostTlsCert *tls.Certificate, hostIdenKeyDatabase *identkeydatabase.IdenKeyDatabase, dbService *databaseservices.DbService, logDIRPath types.LOG_DIR) (*Core, error) {
	// Das Coreobjekt wird erstellt
	coreObj := &Core{
		vmsByID:         make(map[string]*CoreVM),
		vmsByName:       make(map[string]*CoreVM),
		vmKernelPtr:     make(map[types.KernelID]*CoreVM),
		vms:             make([]*CoreVM, 0),
		apiSockets:      make([]types.APISocketInterface, 0),
		hostTlsCert:     hostTlsCert,
		databaseService: dbService,
		state:           static.NEW,
		extModules:      make(map[string]*extmodules.ExternalModule),
		// Chans
		holdOpenChan:     make(chan struct{}),
		serviceSignaling: make(chan struct{}),
		vmSyncWaitGroup:  sync.WaitGroup{},
		apiSyncWaitGroup: sync.WaitGroup{},
		// Datenbanken
		hostIdentKeyDatabase: hostIdenKeyDatabase,
		// Mutexes
		objectMutex: &sync.Mutex{},
		// Log
		logDIR: logDIRPath,
	}

	// Das Objekt wird zurückgegeben
	return coreObj, nil
}
