package core

import (
	"crypto/tls"
	"sync"
	"vnh1/databaseservices"
	"vnh1/identkeydatabase"
	"vnh1/kernelmodules/extmodules"
	"vnh1/types"
)

type Core struct {
	hostIdentKeyDatabase *identkeydatabase.IdenKeyDatabase
	databaseService      *databaseservices.DbService
	apiSockets           []types.APISocketInterface
	hostTlsCert          *tls.Certificate
	vmsByID              map[string]types.VmInterface
	vmsByName            map[string]types.VmInterface
	vmKernelPtr          map[types.KernelID]types.VmInterface
	vmSyncWaitGroup      sync.WaitGroup
	apiSyncWaitGroup     sync.WaitGroup
	state                types.CoreState
	extModules           map[string]*extmodules.ExternalModule
	serviceSignaling     chan struct{}
	holdOpenChan         chan struct{}
	logDIR               types.LOG_DIR
	objectMutex          *sync.Mutex
	vms                  []types.VmInterface
}
