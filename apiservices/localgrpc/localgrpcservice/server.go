package localgrpcservice

import (
	"vnh1/localgrpcproto"
	"vnh1/types"
)

type CliGrpcServer struct {
	localgrpcproto.UnimplementedLocalhostAPIServiceServer
	core types.CoreInterface
}
