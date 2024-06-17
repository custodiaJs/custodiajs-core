package localgrpc

import (
	"net"

	"github.com/CustodiaJS/custodiajs-core/types"

	"google.golang.org/grpc"
)

type HostCliService struct {
	grpcServer *grpc.Server
	netListner net.Listener
	core       types.CoreInterface
}
