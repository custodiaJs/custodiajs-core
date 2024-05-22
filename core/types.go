package core

import (
	"crypto/tls"
	"sync"
	"vnh1/databaseservices"
	"vnh1/databaseservices/services"
	"vnh1/extmodules"
	"vnh1/identkeydatabase"
	"vnh1/kernel"
	"vnh1/types"
	"vnh1/vmdb"
)

type Core struct {
	hostIdentKeyDatabase *identkeydatabase.IdenKeyDatabase
	databaseService      *databaseservices.DbService
	apiSockets           []types.APISocketInterface
	hostTlsCert          *tls.Certificate
	vmsByID              map[string]*CoreVM
	vmsByName            map[string]*CoreVM
	vmKernelPtr          map[types.KernelID]*CoreVM
	vmSyncWaitGroup      sync.WaitGroup
	apiSyncWaitGroup     sync.WaitGroup
	state                types.CoreState
	extModules           map[string]*extmodules.ExternalModule
	serviceSignaling     chan struct{}
	holdOpenChan         chan struct{}
	logDIR               types.LOG_DIR
	objectMutex          *sync.Mutex
	vms                  []*CoreVM
}

type CoreVM struct {
	*kernel.Kernel
	core            *Core
	scriptLoaded    bool
	startTimeUnix   uint64
	objectMutex     *sync.Mutex
	vmState         types.VmState
	vmDbEntry       *vmdb.VmDBEntry
	externalModules []*extmodules.ExternalModule
	dbServiceLinks  []services.DbServiceLinkinterface
}
