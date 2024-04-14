package kernel

import (
	kmodules "vnh1/core/kernel/modules"
	"vnh1/types"
)

var DEFAULT_CONFIG = KernelConfig{
	Modules: []types.KernelModuleInterface{
		kmodules.NewConsoleModule(),
		kmodules.NewRPCModule(),
	},
}
