package localgrpc

import (
	"net"

	"github.com/CustodiaJS/custodiajs-core/global/types"
)

type APIContext struct {
	procUUID types.ProcessId

	openvm *APIProcessVm
	Log    types.ProcessLogSessionInterface
}

type APIProcessVm struct {
	manifest         *types.Manifest
	shared_functions []types.SharedFunctionInterface
	scriptHash       types.VmScriptHash
	kid              types.KernelID
	context          *APIContext
}

type HostAPIService struct {
	procLog    types.ProcessLogSessionInterface
	processes  map[string]*APIContext
	netListner net.Listener
	core       types.CoreInterface
}
