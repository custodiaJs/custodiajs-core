package localgrpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CustodiaJS/custodiajs-core/global/static"
	"github.com/CustodiaJS/custodiajs-core/global/types"
	"github.com/CustodiaJS/custodiajs-core/localgrpcproto"
)

func (s *HostAPIService) WelcomeClient(ctx context.Context, req *localgrpcproto.ClientWelcomeRequest) (*localgrpcproto.ServerWelcomeResponse, error) {
	// Abrufen des Contexts aus dem Kontext
	hacontext := s.NewAPIContext()

	// Es wird dem Context der Typ des Prozesses mitgeteilt
	switch req.ClientType {
	// Es wird geprüft ob es sich um eine VM handelt
	case localgrpcproto.ClientType_VM:
		// Der Type wird hinzugefügt
		hacontext.SetType(req.ClientType)

		// DEBUG
	// Es handelt sich um ein unbekannten client typen
	default:
		return nil, fmt.Errorf("unkown client type")
	}

	// Die Client ID wird zurückgegebn
	return &localgrpcproto.ServerWelcomeResponse{Version: uint64(static.C_VESION), ProcessId: string(hacontext.procUUID)}, nil
}

func (s *HostAPIService) AddVMInstanceByProcess(ctx context.Context, req *localgrpcproto.VMClientRegisterRequest) (*localgrpcproto.VMClientRegisterResponse, error) {
	// Es wird versucht einen vorhanden Context anhand der ProcessId wieder zu finden
	hacontext := s.GetContextByProcessId(types.VmProcessId(req.ProcessId))

	// Es wird versucht das Manifest einzulesen
	var manifest *types.Manifest
	if err := json.Unmarshal(req.Manifest, &manifest); err != nil {
		return nil, err
	}

	// Es wird eineue Vm Instanz erzeugt
	vmInstance, vmid, err := hacontext.CreateVmInstance(manifest, types.VmScriptHash(req.ScriptHash), types.KernelID(req.Kid), hacontext.procUUID)
	if err != nil {
		return nil, err
	}

	// Die VM wird im Core hinzugefügt
	if err := s.core.AddVMInstance(vmInstance, nil); err != nil {
		return nil, err
	}

	// Das VM Objekt wird zurückgegeben
	return &localgrpcproto.VMClientRegisterResponse{Vmid: vmid}, nil
}
