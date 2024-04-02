package clirpc

import (
	"vnh1/types"
)

func NewCliGrpcServer(core types.CoreInterface) *CliGrpcServer {
	return &CliGrpcServer{core: core}
}
