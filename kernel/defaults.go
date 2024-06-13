package kernel

import (
	kmoduleconsole "vnh1/kernelmodules/mainmodules/console"
	kmodulecrypto "vnh1/kernelmodules/mainmodules/crypto"
	kmoduledatabase "vnh1/kernelmodules/mainmodules/database"
	kmodulehttp "vnh1/kernelmodules/mainmodules/http"
	kmodulenet "vnh1/kernelmodules/mainmodules/network"
	kmodulerpc "vnh1/kernelmodules/mainmodules/rpc"
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
