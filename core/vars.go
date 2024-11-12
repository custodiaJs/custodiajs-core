package core

import (
	"net"
	"sync"

	"github.com/CustodiaJS/bngsocket"
	"github.com/CustodiaJS/custodiajs-core/core/ipnetwork"
	"github.com/CustodiaJS/custodiajs-core/crypto"
	"github.com/CustodiaJS/custodiajs-core/global/types"
)

// Core Mutex
var coremutex *sync.Mutex = new(sync.Mutex)

// VM-IPC Sockets
var (
	vmipcRootListener      net.Listener
	vmipcAllUsersListener  net.Listener
	vmipcSpecificListeners map[string]net.Listener
	vmipcOpenConnections   []*bngsocket.BngSocket
	vmipcInited            bool
)

// Core Variablen
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
	vms              []types.VmInterface
	hostnetmanager   *ipnetwork.HostNetworkManagmentUnit
)
