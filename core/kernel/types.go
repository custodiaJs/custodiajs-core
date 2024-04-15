package kernel

import (
	"sync"
	"vnh1/core/consolecache"
	"vnh1/types"

	v8 "rogchap.com/v8go"
)

type KernelConfig struct {
	Modules []types.KernelModuleInterface
}

type Kernel struct {
	*v8.Context
	id        string
	config    *KernelConfig
	mutex     *sync.Mutex
	console   *consolecache.ConsoleOutputCache
	register  map[string]interface{}
	vmImports map[string]*v8.Value
}
