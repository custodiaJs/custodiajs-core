package process

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CustodiaJS/custodiajs-core/global/static"
	"github.com/CustodiaJS/custodiajs-core/global/types"
	"github.com/CustodiaJS/custodiajs-core/localgrpcproto"
)

// Wird verwendet um den neuen Prozess zu Initalisieren
func init_process(stream processStream, isRoot bool) (types.ProcessId, uint64, error) {
	// Die Methode wird ermittelt
	var methode localgrpcproto.ClientMethode
	if isRoot {
		methode = localgrpcproto.ClientMethode_SYSTEM
	} else {
		methode = localgrpcproto.ClientMethode_USER
	}

	// Der Request wird gebaut
	request := &localgrpcproto.ProcessTransport{
		Request: &localgrpcproto.ProcessTransport_WelcomeRequest{
			WelcomeRequest: &localgrpcproto.ClientWelcomeRequest{
				Version:       uint64(static.C_VESION),
				ClientType:    localgrpcproto.ClientType_VM,
				ClientMethode: methode,
			},
		},
	}

	// Der Request wird übermittelt
	if err := stream.Send(request); err != nil {
		return "", 0, fmt.Errorf("init_process: " + err.Error())
	}

	// Es wird auf die Antwort gewartet
	response, err := stream.Recv()
	if err != nil {
		return "", 0, fmt.Errorf("init_process: " + err.Error())
	}

	// Überprüfen, ob req.Req vom Typ *localgrpcproto.Process_RegisterProcessResponse ist
	newProcessRegister, ok := response.Respone.(*localgrpcproto.ProcessTransport_WelcomeResponse)
	if !ok {
		return "", 0, fmt.Errorf("init_process: invalid process type, expected RegisterProcessRequest")
	}

	// Es wird geprüft ob das Paket angenommen wurde
	if !newProcessRegister.WelcomeResponse.Accepted {
		return "", 0, fmt.Errorf("init_process: " + newProcessRegister.WelcomeResponse.Reason)
	}

	// Die ID wird zurück
	return types.ProcessId(newProcessRegister.WelcomeResponse.ProcessId), newProcessRegister.WelcomeResponse.GetVersion(), nil
}

// Wird verwendet um einen Prozess am leben zu halten, sobald diese Funktion zuende ist, gillt der Prozess als geschlossen
func alive_process(_ processStream, _ *VmProcess) error {
	return nil
}

// Wird verwendet um einen neuen Vm Instanz Stream zu erstellen
func init_vm_instance_stream(ctx context.Context, vmProcess *VmProcess, stream vmInstanceStream, manifest *types.Manifest) (types.VmId, error) {
	// Das VmProcess Register Package wird gebaut

	// Das Manifest wird in JSON umgewandelt
	_, err := json.Marshal(manifest)
	if err != nil {
		return "", fmt.Errorf("VmProcess->AddVMInstance: " + err.Error())
	}
	return "", nil
}

// Wird verwendet um die VM Stream Instanz am leben zu erhalten
func vm_instance_alive(ctx context.Context, vmProcess *VmProcess, stream vmInstanceStream, id types.VmId, vmInstance types.VmInterface) error {
	return nil
}
