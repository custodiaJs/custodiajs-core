package core

import (
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"vnh1/core/identkeydatabase"
	"vnh1/core/vmdb"
	"vnh1/types"

	"vnh1/core/jsvm"
)

// Fügt einen neune Script Container hinzu
func (o *Core) AddScriptContainer(vmDbEntry *vmdb.VmDBEntry) (*CoreVM, error) {
	// Die Datei wird zusammengefasst
	fullPath := filepath.Join(vmDbEntry.Path, "main.js")

	// Die Virtuelle Maschine wird geprüft
	if !vmDbEntry.ValidateVM() {
		return nil, fmt.Errorf("AddScriptContainer: Broken Virtual Machine")
	}

	// Es wird versucht die Manifestdatei zuladen
	fileData, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("AddScriptContainer: " + err.Error())
	}

	// Es wird eine neue VM erzeugt
	tvmobj, err := jsvm.NewVM(o, nil)
	if err != nil {
		return nil, fmt.Errorf("AddScriptContainer: " + err.Error())
	}

	// Das Detailspaket wird erzeugt
	vmobject := &CoreVM{JsVM: tvmobj, jsCode: string(fileData), vmDbEntry: vmDbEntry, jsMainFilePath: fullPath, vmState: types.StillWait}

	// Das VMObjekt wird zwischengespeichert
	o.vmsByID[vmDbEntry.GetVMContainerMerkleHash()] = vmobject
	o.vmsByName[vmDbEntry.GetVMName()] = vmobject
	o.vms = append(o.vms, vmobject)

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

	// Der API Socket wird zwischengespeichert
	o.apiSockets = append(o.apiSockets, apiSocket)

	// Es ist kein Fehler aufgetreten
	return nil
}

// Gibt eine Spezifisichen Container VM anhand ihrer ID zurück
func (o *Core) GetScriptContainerVMByID(vmid string) (types.CoreVMInterface, error) {
	// Es wird geprüft ob es sich um einen 64 Zeichen langen String handelt
	if len(vmid) != 64 {
		return nil, fmt.Errorf("Core->GetScriptContainerVMByID: invalid vm container id")
	}

	// Die ID wird lowercast
	lowerCaseID := strings.ToLower(vmid)

	// Es wird geprüft ob die VM exestiert
	vmObj, found := o.vmsByID[vmid]
	if !found {
		return nil, fmt.Errorf("GetScriptContainerVMByID: unkown vm '%s'", lowerCaseID)
	}

	// Das Objekt wird zurückgegeben
	return vmObj, nil
}

// Gibt die ID's der Aktiven VM-Container zurück
func (o *Core) GetAllActiveScriptContainerIDs() []string {
	// Es wird eine Liste mit allen Verfügbaren VM's erstellt
	extr := make([]string, 0)
	for _, item := range o.vmsByID {
		extr = append(extr, item.GetFingerprint())
	}

	// Die Liste wird zurückgegeben
	return extr
}

// Gibt alle VM-Container zurück
func (o *Core) GetAllVMs() []types.CoreVMInterface {
	// Es wird eine Liste mit allen Verfügbaren VM-Containern erstellt
	extr := make([]types.CoreVMInterface, 0)
	for _, item := range o.vmsByID {
		extr = append(extr, item)
	}

	// Die Liste wird zurückgegeben
	return extr
}

// Erstellt einen neuen vnh1 Core
func NewCore(hostTlsCert *tls.Certificate, hostIdenKeyDatabase *identkeydatabase.IdenKeyDatabase) (*Core, error) {
	// Das Coreobjekt wird erstellt
	coreObj := &Core{
		vmsByID:    make(map[string]*CoreVM),
		vmsByName:  make(map[string]*CoreVM),
		vms:        make([]*CoreVM, 0),
		apiSockets: make([]types.APISocketInterface, 0),
		state:      types.NEW,
		// Chans
		holdOpenChan:     make(chan struct{}),
		serviceSignaling: make(chan struct{}),
		vmSyncWaitGroup:  sync.WaitGroup{},
		apiSyncWaitGroup: sync.WaitGroup{},
		// Datenbanken
		hostIdentKeyDatabase: hostIdenKeyDatabase,
	}

	// Das Objekt wird zurückgegeben
	return coreObj, nil
}
