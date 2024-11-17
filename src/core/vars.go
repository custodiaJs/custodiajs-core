package core

import (
	"sync"

	"github.com/CustodiaJS/bngsocket"
	cenvxcore "github.com/custodia-cenv/cenvx-core/src"
	"github.com/custodia-cenv/cenvx-core/src/crypto"
	"github.com/custodia-cenv/cenvx-core/src/static"
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
	vmsByID         map[string]cenvxcore.VmInterface
	vmsByName       map[string]cenvxcore.VmInterface
	vmSyncWaitGroup sync.WaitGroup
	vms             []cenvxcore.VmInterface
)

// Core Variablen
var (
	coreLog      cenvxcore.ProcessLogSessionInterface
	coreState    cenvxcore.CoreState = static.NEW
	cryptoStore  *crypto.CryptoStore
	holdOpenChan chan struct{}
)
