package vmprocessio

import (
	"context"
	"fmt"
	"time"

	"github.com/CustodiaJS/custodiajs-core/crypto"
	"github.com/CustodiaJS/custodiajs-core/ipc"
	"github.com/CustodiaJS/custodiajs-core/localgrpcproto"
	"github.com/CustodiaJS/custodiajs-core/types"
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

func NewCoreVmClientProcess(isRoot bool, socketPath types.SOCKET_PATH, cryptoStore *crypto.VmCryptoStore) (*CoreVmClientProcess, error) {
	// Es wird versucht die Socket verbindung herzustellen
	conn, err := ipc.CreateNewUnixSocketGRPC(string(socketPath), isRoot)
	if err != nil {
		return nil, fmt.Errorf("NewCoreVmClientProcess: " + err.Error())
	}

	// Erstellen eines Kontextes mit Timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Die RPC Instanz wird erstellt
	rpcInstance := localgrpcproto.NewLocalhostAPIServiceClient(conn)

	// Der Prozess wird registriert
	rpcInstance.WelcomeClient(ctx, &localgrpcproto.ClientWelcomeRequest{Version: 10000000000000000000})

	return nil, nil
}
