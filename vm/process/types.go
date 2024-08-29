package process

import (
	"github.com/CustodiaJS/custodiajs-core/localgrpcproto"
	"github.com/CustodiaJS/custodiajs-core/procslog"
	"github.com/CustodiaJS/custodiajs-core/types"
)

type VmProcess struct {
	Log         *procslog.ProcLogSession
	rpcInstance localgrpcproto.LocalhostAPIServiceClient
	processId   types.VmProcessId
	coreVersion uint64
}
