package vm

import (
	"sync"
	"vnh1/databaseservices/services"
	"vnh1/kernel"
	extmodules "vnh1/kernelmodules/extmodules"
	"vnh1/types"
	"vnh1/vmdb"
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
