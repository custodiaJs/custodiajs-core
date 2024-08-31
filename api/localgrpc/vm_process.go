package localgrpc

import (
	"sync"

	"github.com/CustodiaJS/custodiajs-core/global/types"
)

func (o *APIProcessVm) GetQVMID() types.QVMID {
	return ""
}

func (o *APIProcessVm) GetManifest() *types.Manifest {
	return o.manifest
}

func (o *APIProcessVm) GetProcessId() types.ProcessId {
	return o.context.procUUID
}

func (o *APIProcessVm) GetScriptHash() types.VmScriptHash {
	return o.scriptHash
}

func (o *APIProcessVm) GetConsoleOutputWatcher() types.WatcherInterface {
	return nil
}

func (o *APIProcessVm) GetAllSharedFunctions() []types.SharedFunctionInterface {
	return o.shared_functions
}

func (o *APIProcessVm) Serve(*sync.WaitGroup) error {
	return nil
}

func (o *APIProcessVm) GetSharedFunctionBySignature(types.RPCCallSource, *types.FunctionSignature) (types.SharedFunctionInterface, bool, *types.SpecificError) {
	return nil, false, nil
}

func (o *APIProcessVm) GetStartingTimestamp() uint64 {
	return 0
}

func (o *APIProcessVm) GetKId() types.KernelID {
	return o.kid
}

func (o *APIProcessVm) SignalShutdown() {
}

func (o *APIProcessVm) GetState() types.VmState {
	return 0
}
