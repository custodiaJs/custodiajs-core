package clirpc

import (
	grpccli "vnh1/grpc/cligrpc"
	"vnh1/types"
)

type CliGrpcServer struct {
	grpccli.UnimplementedCLIServiceServer
	core types.CoreInterface
}
