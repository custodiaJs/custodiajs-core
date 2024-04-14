package kernel

import (
	kmoduleconsole "vnh1/core/kernel/modules/console"
	kmodulerpc "vnh1/core/kernel/modules/rpc"
	"vnh1/types"
)

var DEFAULT_CONFIG = KernelConfig{
	Modules: []types.KernelModuleInterface{
		kmoduleconsole.NewConsoleModule(),
		kmodulerpc.NewRPCModule(),
	},
}
