package kernel

import (
	"sync"

	"github.com/CustodiaJS/custodiajs-core/core/consolecache"
	"github.com/CustodiaJS/custodiajs-core/global/types"

	v8 "rogchap.com/v8go"
)

type KernelConfig struct {
	Modules []types.KernelModuleInterface
}

type Kernel struct {
	*v8.Context
	id        types.KernelID
	config    *KernelConfig
	mutex     *sync.Mutex
	core      types.CoreInterface
	console   *consolecache.ConsoleOutputCache
	vmLink    types.VmInterface
	register  map[string]interface{}
	vmImports map[string]*v8.Value
	//dbEntry           *vmdb.VmDBEntry
	eventLoopStack    []types.KernelEventLoopOperationInterface
	eventLoopLockCond *sync.Cond
	hasCloseSignal    bool
}
