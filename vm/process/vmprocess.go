package process

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/CustodiaJS/custodiajs-core/crypto"
	"github.com/CustodiaJS/custodiajs-core/ipc"
	"github.com/CustodiaJS/custodiajs-core/localgrpcproto"
	"github.com/CustodiaJS/custodiajs-core/procslog"
	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/types"
)

func (o *VmProcess) AddVMInstance(vmInstance types.VmInterface, plog_a types.ProcessLogSessionInterface) error {
	// Erstellen eines Kontextes mit Timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Das Manifest wird in JSON umgewandelt
	jsonManifestBytes, err := json.Marshal(vmInstance.GetManifest())
	if err != nil {
		return fmt.Errorf("VmProcess->AddVMInstance: " + err.Error())
	}

	// Der Request wird gebaut
	datafqk := &localgrpcproto.VMClientRegisterRequest{
		ProcessId:  string(o.processId),
		Manifest:   jsonManifestBytes,
		ScriptHash: string(vmInstance.GetScriptHash()),
		Kid:        string(vmInstance.GetKId()),
	}

	// Die VM Instance wird hinzugefügt
	_, err = o.rpcInstance.AddVMInstanceByProcess(ctx, datafqk)
	if err != nil {
		return fmt.Errorf("VmProcess->AddVMInstance: " + err.Error())
	}

	o.Log.Debug("New VM Instance added, name = '%s', shash = '%s'", vmInstance.GetManifest().Name, vmInstance.GetScriptHash())

	return nil
}

func (o *VmProcess) GetVmProcessId(plog_a types.ProcessLogSessionInterface) types.VmProcessId {
	return o.processId
}

func (o *VmProcess) GetAllVMs(plog_a types.ProcessLogSessionInterface) []types.VmInterface {
	// Erstellen eines Kontextes mit Timeout
	_, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	return nil
}

func (o *VmProcess) GetAllActiveScriptContainerIDs(plog_a types.ProcessLogSessionInterface) []string {
	// Erstellen eines Kontextes mit Timeout
	_, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	return nil
}

func (o *VmProcess) GetScriptContainerVMByID(vmid string, plog_a types.ProcessLogSessionInterface) (types.VmInterface, bool, *types.SpecificError) {
	// Erstellen eines Kontextes mit Timeout
	_, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	return nil, false, nil
}

func (o *VmProcess) GetScriptContainerByVMName(vmid string, plog_a types.ProcessLogSessionInterface) (types.VmInterface, bool, *types.SpecificError) {
	// Erstellen eines Kontextes mit Timeout
	_, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	return nil, false, nil
}

func (o *VmProcess) GetCoreSessionManagmentUnit(plog_a types.ProcessLogSessionInterface) types.ContextManagmentUnitInterface {
	// Erstellen eines Kontextes mit Timeout
	_, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	return nil
}

func NewVmInstanceClientProcess(isRoot bool, socketPath types.SOCKET_PATH, cryptoStore *crypto.VmCryptoStore, manifest *types.Manifest) (*VmProcess, error) {
	// Es wird ein neuer LogProc erzeugt
	plog := procslog.NewProcLogForVmProcess()

	// Es wird versucht die Socket verbindung herzustellen
	conn, err := ipc.CreateNewUnixSocketGRPC(string(socketPath), isRoot)
	if err != nil {
		return nil, fmt.Errorf("NewVmInstanceClientProcess: " + err.Error())
	}

	// DEBUG
	plog.Debug("Created")

	// Erstellen eines Kontextes mit Timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Die RPC Instanz wird erstellt
	rpcInstance := localgrpcproto.NewLocalhostAPIServiceClient(conn)

	// Die Methode wird ermittelt
	var methode localgrpcproto.ClientMethode
	if isRoot {
		methode = localgrpcproto.ClientMethode_SYSTEM
	} else {
		methode = localgrpcproto.ClientMethode_USER
	}

	// Das Request Objekt wird erzeugt
	requestObject := &localgrpcproto.ClientWelcomeRequest{
		Version:       uint64(static.C_VESION),
		ClientType:    localgrpcproto.ClientType_VM,
		ClientMethode: methode,
	}

	// DEBUG
	plog.Debug("Register Process as VM-Process")

	// Der Prozess wird registriert
	response, err := rpcInstance.WelcomeClient(ctx, requestObject)
	if err != nil {
		return nil, fmt.Errorf("NewVmInstanceClientProcess: " + err.Error())
	}

	// Das VM Objekt wird erzeugt
	vmObject := &VmProcess{
		Log:         plog,
		rpcInstance: rpcInstance,
		processId:   types.VmProcessId(response.ProcessId),
		coreVersion: requestObject.Version,
	}

	// Der Process Control Stream wird gestartet
	if err := grpc_stream_process_control_serve(vmObject); err != nil {
		return nil, fmt.Errorf("NewVmInstanceClientProcess: " + err.Error())
	}

	// Das Objekt wird zurückgegeben
	return vmObject, nil
}
