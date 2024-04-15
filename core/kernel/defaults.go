package kernel

import (
	kmoduleconsole "vnh1/core/kernel/modules/console"
	kmodulecrypto "vnh1/core/kernel/modules/crypto"
	kmoduledatabase "vnh1/core/kernel/modules/database"
	kmodulehttp "vnh1/core/kernel/modules/http"
	kmodulenet "vnh1/core/kernel/modules/network"
	kmodulerpc "vnh1/core/kernel/modules/rpc"
	"vnh1/types"
)

var DEFAULT_CONFIG = KernelConfig{
	Modules: []types.KernelModuleInterface{
		kmoduleconsole.NewConsoleModule(),
		kmodulerpc.NewRPCModule(),
		kmoduledatabase.NewDatabaseModule(),
		kmodulecrypto.NewCryptoModule(),
		kmodulenet.NewNetworkModule(),
		kmodulehttp.NewHttpModule(),
	},
}
