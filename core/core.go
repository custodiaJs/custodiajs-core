package core

import (
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"vnh1/core/identkeydatabase"
	"vnh1/core/vmdb"
	"vnh1/types"

	"vnh1/core/jsvm"
)

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
	vmobject := &CoreVM{JsVM: tvmobj, jsCode: string(fileData), vmDbEntry: vmDbEntry, jsMainFilePath: fullPath}

	// Das VMObjekt wird zwischengespeichert
	o.vmsByID[vmDbEntry.GetVMContainerMerkleHash()] = vmobject
	o.vmsByName[vmDbEntry.GetVMName()] = vmobject
	o.vms = append(o.vms, vmobject)

	// Das VM Objekt wird zwischengespeichert
	return vmobject, nil
}

func (o *Core) AddAPISocket(apiSocket types.APISocketInterface) error {
	// Der Core wird in dem API-Socket Registriert
	err := apiSocket.SetupCore(o)
	if err != nil {
		return fmt.Errorf("AddAPISocket: ")
	}

	// Der API Socket wird zwischengespeichert
	o.apiSockets = append(o.apiSockets, apiSocket)

	// Es ist kein Fehler aufgetreten
	return nil
}

func (o *Core) GetScriptContainerVMByID(vmid string) (types.CoreVMInterface, error) {
	// Es wird geprüft ob die VM exestiert
	vmObj, found := o.vmsByID[vmid]
	if !found {
		return nil, fmt.Errorf("GetScriptContainerVMByID: unkown vm")
	}

	// Das Objekt wird zurückgegeben
	return vmObj, nil
}

func (o *Core) GetAllActiveScriptContainerIDs() []string {
	extr := make([]string, 0)
	for _, item := range o.vmsByID {
		extr = append(extr, item.GetFingerprint())
	}
	return extr
}

func NewCore(hostTlsCert *tls.Certificate, hostIdenKeyDatabase *identkeydatabase.IdenKeyDatabase) (*Core, error) {
	// Das Coreobjekt wird erstellt
	coreObj := &Core{
		vmsByID:    make(map[string]*CoreVM),
		vmsByName:  make(map[string]*CoreVM),
		vms:        make([]*CoreVM, 0),
		apiSockets: make([]types.APISocketInterface, 0),
		state:      NEW,
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
