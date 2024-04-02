package core

import (
	"sync"
	"vnh1/core/identkeydatabase"
	"vnh1/core/jsvm"
	"vnh1/core/vmdb"
	"vnh1/types"
)

type Core struct {
	hostIdentKeyDatabase *identkeydatabase.IdenKeyDatabase
	apiSockets           []types.APISocketInterface
	vmSyncWaitGroup      sync.WaitGroup
	apiSyncWaitGroup     sync.WaitGroup
	state                types.CoreState
	vms                  []*CoreVM
	vmsByID              map[string]*CoreVM
	vmsByName            map[string]*CoreVM
	serviceSignaling     chan struct{}
	holdOpenChan         chan struct{}
	objectMutex          *sync.Mutex
}

type CoreVM struct {
	*jsvm.JsVM
	vmDbEntry *vmdb.VmDBEntry
	vmState   types.VmState
}

type TransportWhitelistVmEntry struct {
	url   string
	alias string
}
