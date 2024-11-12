package core

import (
	"net"
	"sync"

	"github.com/CustodiaJS/custodiajs-core/core/ipnetwork"
	"github.com/CustodiaJS/custodiajs-core/crypto"
	"github.com/CustodiaJS/custodiajs-core/global/types"
)

var (
	coreLog          types.ProcessLogSessionInterface
	apiSockets       []types.APISocketInterface
	cryptoStore      *crypto.CryptoStore
	vmsByID          map[string]types.VmInterface
	vmsByName        map[string]types.VmInterface
	vmKernelPtr      map[types.KernelID]types.VmInterface
	vmSyncWaitGroup  sync.WaitGroup
	apiSyncWaitGroup sync.WaitGroup
	cstate           types.CoreState
	serviceSignaling chan struct{}
	holdOpenChan     chan struct{}
	logDIR           types.LOG_DIR
	objectMutex      *sync.Mutex
	vms              []types.VmInterface
	hostnetmanager   *ipnetwork.HostNetworkManagmentUnit
	rootListener     net.Listener
	coreUserListener net.Listener
	allUsersListener net.Listener
)
