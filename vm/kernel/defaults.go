package kernel

import (
	"github.com/CustodiaJS/custodiajs-core/global/types"
	kmoduleconsole "github.com/CustodiaJS/custodiajs-core/vm/kernel/base_modules/console"
	kmodulecrypto "github.com/CustodiaJS/custodiajs-core/vm/kernel/base_modules/crypto"
	kmoduledatabase "github.com/CustodiaJS/custodiajs-core/vm/kernel/base_modules/database"
	kmodulehttp "github.com/CustodiaJS/custodiajs-core/vm/kernel/base_modules/http"
	kmodulenet "github.com/CustodiaJS/custodiajs-core/vm/kernel/base_modules/network"
	kmodulerpc "github.com/CustodiaJS/custodiajs-core/vm/kernel/base_modules/rpc"
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
