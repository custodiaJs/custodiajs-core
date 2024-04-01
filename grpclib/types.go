package grpclib

import "vnh1/grpclib/publicgrpc"

type GrpcServer struct {
	publicgrpc.UnsafeRPCServiceServer
}
