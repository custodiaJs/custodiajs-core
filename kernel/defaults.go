package kernel

import (
	kmoduleconsole "github.com/CustodiaJS/custodiajs-core/kernelmodules/mainmodules/console"
	kmodulecrypto "github.com/CustodiaJS/custodiajs-core/kernelmodules/mainmodules/crypto"
	kmoduledatabase "github.com/CustodiaJS/custodiajs-core/kernelmodules/mainmodules/database"
	kmodulehttp "github.com/CustodiaJS/custodiajs-core/kernelmodules/mainmodules/http"
	kmodulenet "github.com/CustodiaJS/custodiajs-core/kernelmodules/mainmodules/network"
	kmodulerpc "github.com/CustodiaJS/custodiajs-core/kernelmodules/mainmodules/rpc"
	"github.com/CustodiaJS/custodiajs-core/types"
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
