package core

import (
	"crypto/tls"
	"sync"
	"vnh1/core/databaseservices"
	"vnh1/core/databaseservices/services"
	"vnh1/core/identkeydatabase"
	"vnh1/core/jsvm"
	"vnh1/core/vmdb"
	"vnh1/types"
)

type Core struct {
	hostIdentKeyDatabase *identkeydatabase.IdenKeyDatabase
	databaseService      *databaseservices.DbService
	apiSockets           []types.APISocketInterface
	hostTlsCert          *tls.Certificate
	vmsByID              map[string]*CoreVM
	vmsByName            map[string]*CoreVM
	vmSyncWaitGroup      sync.WaitGroup
	apiSyncWaitGroup     sync.WaitGroup
	state                types.CoreState
	serviceSignaling     chan struct{}
	holdOpenChan         chan struct{}
	objectMutex          *sync.Mutex
	vms                  []*CoreVM
}

type CoreVM struct {
	*jsvm.JsVM
	dbServiceLinks []services.DbServiceLinkinterface
	vmDbEntry      *vmdb.VmDBEntry
	vmState        types.VmState
}
