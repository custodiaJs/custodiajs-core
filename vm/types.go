package vm

import (
	"sync"

	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/vm/image"
	"github.com/CustodiaJS/custodiajs-core/vm/kernel"
)

type VmInstance struct {
	*kernel.Kernel
	core          types.CoreInterface
	scriptLoaded  bool
	startTimeUnix uint64
	objectMutex   *sync.Mutex
	vmState       types.VmState
	vmImage       *image.VmImage
	_signal_CLOSE bool
}
