package clirpc

import (
	"context"
	"vnh1/grpc/cligrpc"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *CliGrpcServer) ListVMs(ctx context.Context, _ *emptypb.Empty) (*cligrpc.VmListResponse, error) {
	// Die Werte werden abgeabreitet
	entry := []*cligrpc.VmListEntry{}
	for _, item := range s.core.GetAllVMs() {
		sharf := make([]string, 0)
		for _, sharfnc := range item.GetAllSharedFunctions() {
			sharf = append(sharf, sharfnc.GetName())
		}
		entry = append(entry, &cligrpc.VmListEntry{
			Name:            item.GetVMName(),
			Id:              item.GetFingerprint(),
			State:           uint32(item.GetState()),
			StartTime:       item.GetStartingTimestamp(),
			SharedFunctions: sharf,
			NodeJsModules:   item.GetVMModuleNames(),
			DomainWhiteList: []string{},
			SslMemberHashes: []string{},
		})
	}

	// Der Rückgabewert wird erzeugt
	returnValue := &cligrpc.VmListResponse{Vms: entry}

	// Die Daten werden zurückgesendet
	return returnValue, nil
}
