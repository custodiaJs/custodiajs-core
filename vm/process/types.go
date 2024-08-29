package process

import (
	"github.com/CustodiaJS/custodiajs-core/global/procslog"
	"github.com/CustodiaJS/custodiajs-core/global/types"
	"github.com/CustodiaJS/custodiajs-core/localgrpcproto"
)

type VmProcess struct {
	Log         *procslog.ProcLogSession
	rpcInstance localgrpcproto.LocalhostAPIServiceClient
	processId   types.VmProcessId
	coreVersion uint64
}
