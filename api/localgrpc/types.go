package localgrpc

import (
	"net"

	"github.com/CustodiaJS/custodiajs-core/global/types"
	"github.com/CustodiaJS/custodiajs-core/localgrpcproto"
	"google.golang.org/grpc"
)

type APIContext struct {
	procUUID types.VmProcessId
	tpe      localgrpcproto.ClientType
	openvm   *APIProcessVm
	Log      types.ProcessLogSessionInterface
}

type APIProcessVm struct {
	manifest         *types.Manifest
	shared_functions []types.SharedFunctionInterface
	scriptHash       types.VmScriptHash
	kid              types.KernelID
	context          *APIContext
}

type HostAPIService struct {
	localgrpcproto.UnimplementedLocalhostapierver
	procLog    types.ProcessLogSessionInterface
	grpcServer *grpc.Server
	processes  map[string]*APIContext
	netListner net.Listener
	core       types.CoreInterface
}
