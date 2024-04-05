package localgrpcservice

import (
	"context"
	"fmt"
	"vnh1/localgrpcproto"
)

func (s *CliGrpcServer) GetVMDetails(ctx context.Context, vmDetailParms *localgrpcproto.VmDetailsParms) (*localgrpcproto.VmDetailsResponse, error) {
	switch vmDetailParms.Value.(type) {
	case localgrpcproto.VmDetailsParms_Id:
	case localgrpcproto.VmDetailsParms_Name:
	default:
		return nil, fmt.Errorf("invalid 'get vm details' parameter vm id/name")
	}

	// Die Daten werden zurückgesendet
	return &localgrpcproto.VmDetailsResponse{}, nil
}
