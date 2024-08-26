package vm

import (
	"sync"

	"github.com/CustodiaJS/custodiajs-core/kernel"
	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/vmimage"
)

type CoreVM struct {
	*kernel.Kernel
	core          types.CoreInterface
	scriptLoaded  bool
	startTimeUnix uint64
	objectMutex   *sync.Mutex
	vmState       types.VmState
	vmImage       *vmimage.VmImage
	_signal_CLOSE bool
}
