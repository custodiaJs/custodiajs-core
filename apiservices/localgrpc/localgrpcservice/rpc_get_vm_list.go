package localgrpcservice

import (
	"context"
	"fmt"
	"strings"
	"vnh1/localgrpcproto"
	"vnh1/utils/procslog"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *CliGrpcServer) ListVMs(ctx context.Context, _ *emptypb.Empty) (*localgrpcproto.VmListResponse, error) {
	// Die Werte werden abgeabreitet
	entry := []*localgrpcproto.VmListEntry{}
	for _, item := range s.core.GetAllVMs() {
		// Der Eintrag wird hinzugefügt
		entry = append(entry, &localgrpcproto.VmListEntry{
			Name:      item.GetVMName(),
			Id:        strings.ToUpper(string(item.GetFingerprint())),
			State:     uint32(item.GetState()),
			StartTime: item.GetStartingTimestamp(),
		})
	}

	// Der Rückgabewert wird erzeugt
	returnValue := &localgrpcproto.VmListResponse{Vms: entry}

	// Log
	procslog.LogPrint(fmt.Sprintf("CLI: Retrieve VmList with %d items\n", len(entry)))

	// Die Daten werden zurückgesendet
	return returnValue, nil
}
