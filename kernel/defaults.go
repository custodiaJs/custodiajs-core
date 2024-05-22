package kernel

import (
	kmoduleconsole "vnh1/kernel/modules/console"
	kmodulecrypto "vnh1/kernel/modules/crypto"
	kmoduledatabase "vnh1/kernel/modules/database"
	kmodulehttp "vnh1/kernel/modules/http"
	kmodulenet "vnh1/kernel/modules/network"
	kmodulerpc "vnh1/kernel/modules/rpc"
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
