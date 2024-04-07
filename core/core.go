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
	"vnh1/types"
	"vnh1/utils"

	"vnh1/core/jsvm"
)

// Fügt einen neune Script Container hinzu
func (o *Core) AddScriptContainer(vmDbEntry *vmdb.VmDBEntry) (*CoreVM, error) {
	// Die Virtuelle Maschine wird geprüft
	if !vmDbEntry.ValidateVM() {
		return nil, fmt.Errorf("AddScriptContainer: Broken Virtual Machine")
	}

	// Es wird geprüft welche HostCAS benötigt werden
	for _, item := range vmDbEntry.GetMemberCertsPkeys() {
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

	// Der Mutex wird angewendet
	// und nach beenden der Funktion freigegeben
	o.objectMutex.Lock()

	// Es wird geprüft ob bereits eiein VM Link hinzugefügtne VM mit der Selben ID vorhanden ist
	if _, foundVM := o.vmsByID[strings.ToLower(vmDbEntry.GetVMContainerMerkleHash())]; foundVM {
		o.objectMutex.Unlock()
		return nil, fmt.Errorf("Core->AddScriptContainer: You cannot add a VM container '%s' multiple times", vmDbEntry.GetVMContainerMerkleHash())
	}

	// Es wird eine neue VM erzeugt
	tvmobj, err := jsvm.NewVM(nil)
	if err != nil {
		o.objectMutex.Unlock()
		return nil, fmt.Errorf("AddScriptContainer: " + err.Error())
	}

	// Das Detailspaket wird erzeugt
	vmobject := newCoreVM(tvmobj, vmDbEntry)

	// Das VMObjekt wird zwischengespeichert
	o.vmsByID[strings.ToLower(vmDbEntry.GetVMContainerMerkleHash())] = vmobject // Merklehash
	o.vmsByName[strings.ToLower(vmDbEntry.GetVMName())] = vmobject              // VM-Name
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

// Fügt einen alternativen Dienst hinzus
func (o *Core) AddAlternativeService(altService types.AlternativeServiceInterface) error {
	panic("FUNCTION_NOT_IMPLEMENTATED")
}

// Gibt eine Spezifisichen Container VM anhand ihrer ID zurück
func (o *Core) GetScriptContainerVMByID(vmid string) (types.CoreVMInterface, error) {
	// Es wird geprüft ob es sich um einen zulässigen vm Namen handelt
	if !utils.ValidateVMIdString(vmid) {
		return nil, fmt.Errorf("Core->GetScriptContainerVMByID: invalid vm container id")
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
		return nil, fmt.Errorf("Core->GetScriptContainerVMByID: unkown vm '%s'", lowerCaseId)
	}

	// Das Objekt wird zurückgegeben
	return vmObj, nil
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
func NewCore(hostTlsCert *tls.Certificate, hostIdenKeyDatabase *identkeydatabase.IdenKeyDatabase, dbService *databaseservices.DbService) (*Core, error) {
	// Das Coreobjekt wird erstellt
	coreObj := &Core{
		vmsByID:         make(map[string]*CoreVM),
		vmsByName:       make(map[string]*CoreVM),
		vms:             make([]*CoreVM, 0),
		apiSockets:      make([]types.APISocketInterface, 0),
		hostTlsCert:     hostTlsCert,
		databaseService: dbService,
		state:           types.NEW,
		// Chans
		holdOpenChan:     make(chan struct{}),
		serviceSignaling: make(chan struct{}),
		vmSyncWaitGroup:  sync.WaitGroup{},
		apiSyncWaitGroup: sync.WaitGroup{},
		// Datenbanken
		hostIdentKeyDatabase: hostIdenKeyDatabase,
		// Mutexes
		objectMutex: &sync.Mutex{},
	}

	// Das Objekt wird zurückgegeben
	return coreObj, nil
}
