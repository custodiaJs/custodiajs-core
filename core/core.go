package core

import (
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"vnh1/core/identkeydatabase"
	"vnh1/core/vmdb"

	"vnh1/core/jsvm"

	"github.com/dop251/goja"
)

type CoreState int

const (
	NEW CoreState = iota
	SERVING
	SHUTDOWN
	CLOSED
)

type APISocketInterface interface {
	Serve(chan struct{}) error
}

type Core struct {
	hostIdentKeyDatabase *identkeydatabase.IdenKeyDatabase
	vmsByID              map[string]*CoreVM
	vmsByName            map[string]*CoreVM
	vms                  []*CoreVM
	vmSyncWaitGroup      sync.WaitGroup
	apiSyncWaitGroup     sync.WaitGroup
	apiSockets           []APISocketInterface
	serviceSignaling     chan struct{}
	holdOpenChan         chan struct{}
	state                CoreState
}

func (o *Core) AddNewVM(vmDbEntry *vmdb.VmDBEntry) (*CoreVM, error) {
	// Die Datei wird zusammengefasst
	fullPath := filepath.Join(vmDbEntry.Path, "main.js")

	// Die Virtuelle Maschine wird geprüft
	if !vmDbEntry.ValidateVM() {
		return nil, fmt.Errorf("AddNewVM: Broken Virtual Machine")
	}

	// Es wird versucht die Manifestdatei zuladen
	fileData, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("AddNewVM: " + err.Error())
	}

	// Es wird eine neue VM erzeugt
	tvmobj, err := jsvm.NewVM(o, nil)
	if err != nil {
		return nil, fmt.Errorf("AddNewVM: " + err.Error())
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

func (o *Core) RegisterSharedLocalFunction(vm *jsvm.JsVM, funcName string, totalParms []string, function goja.Callable) error {
	fmt.Println("CORE:SHARE_LOCAL_FUNCTION:", funcName, totalParms)
	return nil
}

func (o *Core) AddAPISocket(apiSocket APISocketInterface) error {
	o.apiSockets = append(o.apiSockets, apiSocket)
	return nil
}

func NewCore(hostTlsCert *tls.Certificate, hostIdenKeyDatabase *identkeydatabase.IdenKeyDatabase) (*Core, error) {
	// Das Coreobjekt wird erstellt
	coreObj := &Core{
		vmsByID:    make(map[string]*CoreVM),
		vmsByName:  make(map[string]*CoreVM),
		vms:        make([]*CoreVM, 0),
		apiSockets: make([]APISocketInterface, 0),
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
