package localgrpc

import (
	"fmt"

	"github.com/CustodiaJS/custodiajs-core/global/static"
	"github.com/CustodiaJS/custodiajs-core/localgrpcproto"
)

// Wird verwendet um neue Prozesse zu Registrieren
func (s *HostAPIService) SpawnProcess(stream localgrpcproto.LocalhostAPIService_SpawnProcessServer) error {
	// Es wird auf die Register Nachricht gewartet
	registerRecv, registerRecvErr := stream.Recv()
	if registerRecvErr != nil {
		return fmt.Errorf("HostAPIService->CoreToProcessControl: " + registerRecvErr.Error())
	}

	// Überprüfen, ob req.Req vom Typ *localgrpcproto.Process_RegisterProcessRequest ist
	newProcessRegister, ok := registerRecv.Request.(*localgrpcproto.ProcessTransport_WelcomeRequest)
	if !ok {
		return fmt.Errorf("HostAPIService->SpawnProcess: invalid process type, expected RegisterProcessRequest")
	}

	// Abrufen des Contexts aus dem Kontext
	hacontext := s.NewAPIContext()

	// Es wird dem Context der Typ des Prozesses mitgeteilt
	switch newProcessRegister.WelcomeRequest.ClientType {
	// Es wird geprüft ob es sich um eine VM handelt
	case localgrpcproto.ClientType_VM:
		// Der Type wird hinzugefügt
		hacontext.SetType(localgrpcproto.ClientType_VM)
	// Es handelt sich um ein unbekannten client typen
	default:
		return fmt.Errorf("unkown client type")
	}

	// Das Paket wird gebaut
	response := &localgrpcproto.ProcessTransport{
		Respone: &localgrpcproto.ProcessTransport_WelcomeResponse{
			WelcomeResponse: &localgrpcproto.ServerWelcomeResponse{
				Version:   uint64(static.C_VESION),
				ProcessId: string(hacontext.procUUID),
				Accepted:  true,
			},
		},
	}

	// Dem Client wird die Context ID zurückgesendet
	if err := stream.Send(response); err != nil {
		return err
	}

	// DEBUG
	s.procLog.Debug("New Process inited")

	// Es ist kein Fehler Fehler aufgetreten
	return nil
}
