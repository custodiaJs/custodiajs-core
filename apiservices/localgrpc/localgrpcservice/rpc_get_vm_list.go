package localgrpcservice

import (
	"context"
	"fmt"
	"strings"
	"vnh1/localgrpcproto"
	"vnh1/utils"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *CliGrpcServer) ListVMs(ctx context.Context, _ *emptypb.Empty) (*localgrpcproto.VmListResponse, error) {
	// Die Werte werden abgeabreitet
	entry := []*localgrpcproto.VmListEntry{}
	for _, item := range s.core.GetAllVMs() {
		// Die geteilten Funktionen werden abgerufen
		sharf := make([]string, 0)
		for _, sharfnc := range item.GetAllSharedFunctions() {
			sharf = append(sharf, sharfnc.GetName())
		}

		// Die Erlaubten Domains werden abgerufen
		allowedDomains := []string{}
		for _, item := range item.GetWhitelist() {
			allowedDomains = append(allowedDomains, item.URL())
		}

		// Der Eintrag wird hinzugefügt
		entry = append(entry, &localgrpcproto.VmListEntry{
			Name:            item.GetVMName(),
			Id:              strings.ToUpper(item.GetFingerprint()),
			State:           uint32(item.GetState()),
			StartTime:       item.GetStartingTimestamp(),
			NodeJsModules:   item.GetVMModuleNames(),
			DomainWhiteList: allowedDomains,
			UsedHostKeyIds:  item.GetMemberCertKeyIds(),
			SharedFunctions: sharf,
		})
	}

	// Der Rückgabewert wird erzeugt
	returnValue := &localgrpcproto.VmListResponse{Vms: entry}

	// Log
	utils.LogPrint(fmt.Sprintf("CLI: Retrieve VmList with %d items\n", len(entry)))

	// Die Daten werden zurückgesendet
	return returnValue, nil
}
