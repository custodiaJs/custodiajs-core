package localgrpcservice

import (
	"context"
	"fmt"

	"github.com/CustodiaJS/custodiajs-core/localgrpcproto"
	"github.com/CustodiaJS/custodiajs-core/static"
)

func (s *CliGrpcServer) WelcomeClient(ctx context.Context, req *localgrpcproto.ClientWelcomeRequest) (*localgrpcproto.ClientWelcomeResponse, error) {
	// Es wird geprüft ob es sich um eine VM handelt
	if req.ClientType == localgrpcproto.ClientType_VM {

	}

	fmt.Println("WELCOME CLIENT", req.ClientMethode, req.VmClient.Manifest)
	return &localgrpcproto.ClientWelcomeResponse{Version: uint64(static.C_VESION)}, nil
}
