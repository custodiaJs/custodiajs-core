package core

import (
	"sync"

	"github.com/CustodiaJS/custodiajs-core/apiservices/http/context"
	"github.com/CustodiaJS/custodiajs-core/core/ipnetwork"
	"github.com/CustodiaJS/custodiajs-core/crypto"
	"github.com/CustodiaJS/custodiajs-core/types"
)

type Core struct {
	coreLog          types.ProcessLogSessionInterface
	cpmu             *context.ContextManagmentUnit
	apiSockets       []types.APISocketInterface
	cryptoStore      *crypto.CryptoStore
	vmsByID          map[string]types.VmInterface
	vmsByName        map[string]types.VmInterface
	vmKernelPtr      map[types.KernelID]types.VmInterface
	vmSyncWaitGroup  sync.WaitGroup
	apiSyncWaitGroup sync.WaitGroup
	state            types.CoreState
	serviceSignaling chan struct{}
	holdOpenChan     chan struct{}
	logDIR           types.LOG_DIR
	objectMutex      *sync.Mutex
	vms              []types.VmInterface
	hostnetmanager   *ipnetwork.HostNetworkManagmentUnit
}

type CoreFirewall struct {
}
