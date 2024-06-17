package vm

import (
	"sync"

	"github.com/CustodiaJS/custodiajs-core/databaseservices/services"
	"github.com/CustodiaJS/custodiajs-core/kernel"
	extmodules "github.com/CustodiaJS/custodiajs-core/kernelmodules/extmodules"
	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/vmdb"
)

type CoreVM struct {
	*kernel.Kernel
	core            types.CoreInterface
	scriptLoaded    bool
	startTimeUnix   uint64
	objectMutex     *sync.Mutex
	vmState         types.VmState
	vmDbEntry       *vmdb.VmDBEntry
	externalModules []*extmodules.ExternalModule
	dbServiceLinks  []services.DbServiceLinkinterface
	_signal_CLOSE   bool
}
