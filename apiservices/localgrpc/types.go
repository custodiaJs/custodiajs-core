package localgrpc

import (
	"net"

	"github.com/CustodiaJS/custodiajs-core/localgrpcproto"
	"github.com/CustodiaJS/custodiajs-core/procslog"
	"github.com/CustodiaJS/custodiajs-core/types"
	"google.golang.org/grpc"
)

type APIContext struct {
	procUUID types.VmProcessId
	tpe      localgrpcproto.ClientType
	openvm   *APIProcessVm
}

type APIProcessVm struct {
	manifest         *types.Manifest
	shared_functions []types.SharedFunctionInterface
	scriptHash       types.VmScriptHash
	kid              types.KernelID
	context          *APIContext
}

type HostAPIService struct {
	localgrpcproto.UnimplementedLocalhostAPIServiceServer
	procLog    *procslog.ProcLogSession
	grpcServer *grpc.Server
	processes  map[string]*APIContext
	netListner net.Listener
	core       types.CoreInterface
}
