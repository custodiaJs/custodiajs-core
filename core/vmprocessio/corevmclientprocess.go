package vmprocessio

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/CustodiaJS/custodiajs-core/crypto"
	"github.com/CustodiaJS/custodiajs-core/ipc"
	"github.com/CustodiaJS/custodiajs-core/localgrpcproto"
	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/vmimage"
)

func (o *CoreVmClientProcess) GetAllVMs() []types.VmInterface {
	return nil
}

func (o *CoreVmClientProcess) GetAllActiveScriptContainerIDs(processLog types.ProcessLogSessionInterface) []string {
	return nil
}

func (o *CoreVmClientProcess) GetScriptContainerVMByID(vmid string) (types.VmInterface, bool, *types.SpecificError) {
	return nil, false, nil
}

func (o *CoreVmClientProcess) GetScriptContainerByVMName(string) (types.VmInterface, bool, *types.SpecificError) {
	return nil, false, nil
}

func (o *CoreVmClientProcess) GetCoreSessionManagmentUnit() types.ContextManagmentUnitInterface {
	return nil
}

func NewCoreVmClientProcess(isRoot bool, socketPath types.SOCKET_PATH, cryptoStore *crypto.VmCryptoStore, manifest *vmimage.Manifest) (*CoreVmClientProcess, error) {
	// Es wird versucht die Socket verbindung herzustellen
	conn, err := ipc.CreateNewUnixSocketGRPC(string(socketPath), isRoot)
	if err != nil {
		return nil, fmt.Errorf("NewCoreVmClientProcess: " + err.Error())
	}

	// Erstellen eines Kontextes mit Timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)

	// Die RPC Instanz wird erstellt
	rpcInstance := localgrpcproto.NewLocalhostAPIServiceClient(conn)

	// Die Methode wird ermittelt
	var methode localgrpcproto.ClientMethode
	if isRoot {
		methode = localgrpcproto.ClientMethode_SYSTEM
	} else {
		methode = localgrpcproto.ClientMethode_USER
	}

	// Das Manifest der Anwendung wird ausgewertet
	jsonData, jsonError := json.Marshal(manifest)
	if jsonError != nil {
		// Der Vorgang wird abgebrochen
		cancel()

		// Der Fehler wird zurückgegeben
		return nil, fmt.Errorf("NewCoreVmClientProcess: " + jsonError.Error())
	}

	// Das Request Objekt wird erzeugt
	requestObject := &localgrpcproto.ClientWelcomeRequest{
		Version:       uint64(static.C_VESION),
		ClientType:    localgrpcproto.ClientType_VM,
		ClientMethode: methode,
		UnixProcInfo: &localgrpcproto.ProcessInfo{
			RunInContainer: false,
			UserId:         uint32(os.Getuid()),
			ProcessId:      uint32(os.Getpid()),
		},
		VmClient: &localgrpcproto.VMClient{
			Manifest: string(jsonData),
		},
	}

	// Der Prozess wird registriert
	rpcInstance.WelcomeClient(ctx, requestObject)

	// Das VM Objekt wird erzeugt
	vmObject := &CoreVmClientProcess{
		cancel: cancel,
	}

	// Das Objekt wird zurückgegeben
	return vmObject, nil
}
