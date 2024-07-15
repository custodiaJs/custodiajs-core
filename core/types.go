package core

import (
	"crypto/tls"
	"sync"

	"github.com/CustodiaJS/custodiajs-core/databaseservices"
	"github.com/CustodiaJS/custodiajs-core/identkeydatabase"
	"github.com/CustodiaJS/custodiajs-core/ipnetwork"
	"github.com/CustodiaJS/custodiajs-core/kernel/external_modules"
	"github.com/CustodiaJS/custodiajs-core/saftychan"
	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/utils/grsbool"
	"github.com/CustodiaJS/custodiajs-core/utils/procslog"
)

type CoreWebRequestRPCSession struct {
	saftyResponseChan *saftychan.FunctionCallReturnChan
	proc              *procslog.ProcLogSession
	isConnected       *grsbool.Grsbool
}

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
	extModules           map[string]*external_modules.ExternalModule
	serviceSignaling     chan struct{}
	holdOpenChan         chan struct{}
	logDIR               types.LOG_DIR
	objectMutex          *sync.Mutex
	vms                  []types.VmInterface
	hostnetmanager       *ipnetwork.HostNetworkManagmentUnit
}

type CoreFirewall struct {
}

type CoreLRSAP struct {
	SourceAddress           *ipnetwork.IpAddress
	LocalAddress            *ipnetwork.IpAddress
	ConnectionOverInterface *ipnetwork.NetworkInterface
}
