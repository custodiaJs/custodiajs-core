package vm

import (
	"sync"

	"github.com/CustodiaJS/custodiajs-core/global/types"
	kernel "github.com/CustodiaJS/custodiajs-core/vm/context"
	"github.com/CustodiaJS/custodiajs-core/vm/image"
)

type VmInstance struct {
	*kernel.VmContext
	core          types.CoreInterface
	scriptLoaded  bool
	startTimeUnix uint64
	objectMutex   *sync.Mutex
	vmState       types.VmState
	vmImage       *image.VmImage
	_signal_CLOSE bool
}
