package core

import (
	"sync"
	"vnh1/core/identkeydatabase"
	"vnh1/types"
)

type Core struct {
	hostIdentKeyDatabase *identkeydatabase.IdenKeyDatabase
	vmsByID              map[string]*CoreVM
	vmsByName            map[string]*CoreVM
	vms                  []*CoreVM
	vmSyncWaitGroup      sync.WaitGroup
	apiSyncWaitGroup     sync.WaitGroup
	apiSockets           []types.APISocketInterface
	serviceSignaling     chan struct{}
	holdOpenChan         chan struct{}
	state                types.CoreState
}

type TransportWhitelistVmEntry struct {
	url   string
	alias string
}
