package localgrpc

import (
	"net"
	"vnh1/types"

	"google.golang.org/grpc"
)

type HostCliService struct {
	grpcServer *grpc.Server
	netListner net.Listener
	core       types.CoreInterface
}
