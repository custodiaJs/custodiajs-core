package core

import (
	"sync"
	"vnh1/core/identkeydatabase"
	"vnh1/types"
)

type CoreState int

const (
	NEW CoreState = iota
	SERVING
	SHUTDOWN
	CLOSED
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
	state                CoreState
}
