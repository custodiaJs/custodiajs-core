package clirpc

import (
	grpccli "vnh1/grpc/cligrpc"
	"vnh1/types"
)

type CliGrpcServer struct {
	grpccli.UnimplementedMyServiceServer
	core types.CoreInterface
}

func NewCliGrpcServer(core types.CoreInterface) *CliGrpcServer {
	return &CliGrpcServer{core: core}
}
