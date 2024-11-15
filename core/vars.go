package core

import (
	"sync"

	"github.com/CustodiaJS/bngsocket"
	"github.com/CustodiaJS/custodiajs-core/crypto"
	"github.com/CustodiaJS/custodiajs-core/global/static"
	"github.com/CustodiaJS/custodiajs-core/global/types"
)

// Core Mutex
var coremutex *sync.Mutex = new(sync.Mutex)

// VM-IPC Sockets
var (
	vmipcState           _VmIpcServerState = NEW
	vmipcListeners       []*_AclListener
	vmipcOpenConnections []*bngsocket.BngConn
)

// Speichert alle VM's ab, welche dem Core bekannt sind
var (
	vmsByID         map[string]types.VmInterface
	vmsByName       map[string]types.VmInterface
	vmSyncWaitGroup sync.WaitGroup
	vms             []types.VmInterface
)

// Core Variablen
var (
	coreLog      types.ProcessLogSessionInterface
	coreState    types.CoreState = static.NEW
	cryptoStore  *crypto.CryptoStore
	holdOpenChan chan struct{}
)
