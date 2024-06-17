package localgrpcservice

import (
	"github.com/CustodiaJS/custodiajs-core/types"
)

func NewCliGrpcServer(core types.CoreInterface) *CliGrpcServer {
	return &CliGrpcServer{core: core}
}
