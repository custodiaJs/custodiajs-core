package process

import (
	"github.com/CustodiaJS/custodiajs-core/global/procslog"
	"github.com/CustodiaJS/custodiajs-core/global/types"
	"github.com/CustodiaJS/custodiajs-core/localgrpcproto"
	"google.golang.org/grpc"
)

type processStream grpc.BidiStreamingClient[localgrpcproto.ProcessTransport, localgrpcproto.ProcessTransport]
type vmInstanceStream grpc.BidiStreamingClient[localgrpcproto.ProcessTransport, localgrpcproto.ProcessTransport]

type VmProcess struct {
	Log         *procslog.ProcLogSession
	rpcInstance localgrpcproto.LocalhostAPIServiceClient
	processId   types.ProcessId
	coreVersion uint64
	vmInstances []types.VmInterface
}
