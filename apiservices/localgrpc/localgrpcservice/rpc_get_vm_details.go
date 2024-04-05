package localgrpcservice

import (
	"context"
	"vnh1/localgrpcproto"
)

func (s *CliGrpcServer) GetVMDetails(ctx context.Context, vmDetailParms *localgrpcproto.VmDetailsParms) (*localgrpcproto.VmDetailsResponse, error) {

	// Die Daten werden zurückgesendet
	return &localgrpcproto.VmDetailsResponse{}, nil
}
