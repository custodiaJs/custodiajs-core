package localgrpcservice

import (
	"github.com/CustodiaJS/custodiajs-core/localgrpcproto"
	"github.com/CustodiaJS/custodiajs-core/types"
)

type CliGrpcServer struct {
	localgrpcproto.UnimplementedLocalhostAPIServiceServer
	core types.CoreInterface
}
