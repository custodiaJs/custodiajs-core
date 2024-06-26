package core

import (
	"crypto/tls"
	"sync"

	"github.com/CustodiaJS/custodiajs-core/databaseservices"
	"github.com/CustodiaJS/custodiajs-core/identkeydatabase"
	"github.com/CustodiaJS/custodiajs-core/kernelmodules/extmodules"
	"github.com/CustodiaJS/custodiajs-core/types"
)

type CoreSessionManagmentUnit struct {
}

type Core struct {
	cpmu                 *CoreSessionManagmentUnit
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
