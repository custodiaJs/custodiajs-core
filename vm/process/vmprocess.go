package process

import (
	"context"
	"fmt"
	"time"

	"github.com/CustodiaJS/custodiajs-core/crypto"
	"github.com/CustodiaJS/custodiajs-core/global/procslog"
	"github.com/CustodiaJS/custodiajs-core/global/types"
	"github.com/CustodiaJS/custodiajs-core/localgrpcproto"
)

// Wird verwendet um eine neue VM Instanz hinzuzufügen
func (o *VmProcess) AddVMInstance(vmInstance types.VmInterface, plog_a types.ProcessLogSessionInterface) error {
	// Es wird geprüft das nicht mehr als eine VM ausgeführt wird
	if len(o.vmInstances) > 0 {
		return fmt.Errorf("VmProcess->AddVMInstance: only one vm allowed")
	}

	// Erstellen eines Kontextes mit Timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Der Prozess Stream wird gestartet
	stream, err := o.rpcInstance.SpawnProcess(ctx)
	if err != nil {
		return fmt.Errorf("VmProcess->AddVMInstance: " + err.Error())
	}

	// Es wird eien neuer VM Instanz Stream Inistalisiert
	vmid, err := init_vm_instance_stream(ctx, o, stream, vmInstance.GetManifest())
	if err != nil {
		return fmt.Errorf("VmProcess->AddVMInstance: " + err.Error())
	}

	// Die Verbindung wird am leben erhalten
	if err := vm_instance_alive(ctx, o, stream, vmid, vmInstance); err != nil {
		return fmt.Errorf("VmProcess->AddVMInstance: " + err.Error())
	}

	// Die VM Instanz wird zwischengespeichert
	o.vmInstances = append(o.vmInstances, vmInstance)

	// DEBUG
	o.Log.Debug("New VM Instance added, name = '%s', shash = '%s'", vmInstance.GetManifest().Name, vmInstance.GetScriptHash())

	// Es ist kein Fehler aufgetreten
	return nil
}

// Gibt die Aktuelle Prozess ID zurück, diese ID wird vom Core zugewiesen
func (o *VmProcess) GetProcessId(plog_a types.ProcessLogSessionInterface) types.ProcessId {
	return o.processId
}

// Gibt eine Liste mit allen Verfügbaren VM's zurück
func (o *VmProcess) GetAllVMs(plog_a types.ProcessLogSessionInterface) []types.VmInterface {
	// Erstellen eines Kontextes mit Timeout
	_, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Es wurden keine Vms abgerufen
	return nil
}

// Gibt eine Liste mit allen Verfügbaren VM ID's zurück
func (o *VmProcess) GetAllActiveVmIDs(plog_a types.ProcessLogSessionInterface) []string {
	// Erstellen eines Kontextes mit Timeout
	_, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	return nil
}

// Gibt eine Spizielle VM anhand ihrer ID zurück
func (o *VmProcess) GetVmByID(vmid string, plog_a types.ProcessLogSessionInterface) (types.VmInterface, bool, *types.SpecificError) {
	// Erstellen eines Kontextes mit Timeout
	_, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	return nil, false, nil
}

// Gibt eine Spizielle VM anhand ihres Namens zurück
func (o *VmProcess) GetVmByName(vmid string, plog_a types.ProcessLogSessionInterface) (types.VmInterface, bool, *types.SpecificError) {
	// Erstellen eines Kontextes mit Timeout
	_, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	return nil, false, nil
}

// Gibt den Core Session Manager zurück
func (o *VmProcess) GetCoreSessionManagmentUnit(plog_a types.ProcessLogSessionInterface) types.ContextManagmentUnitInterface {
	// Erstellen eines Kontextes mit Timeout
	_, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	return nil
}

// Wird verwendet um eine neue VM Prozess Instanz zu erzeugen
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

	// Die RPC Instanz wird erstellt
	rpcInstance := localgrpcproto.NewLocalhostAPIServiceClient(conn)

	// DEBUG
	plog.Debug("Register Process as VM-Process")

	// Erstellen eines Kontextes mit Timeout
	ctx, _ := context.WithTimeout(context.Background(), 50*time.Millisecond)

	// Der Prozess Stream wird gestartet
	processStream, err := rpcInstance.SpawnProcess(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewVmInstanceClientProcess: " + err.Error())
	}

	// Der Prozess wird registriert
	procId, coreVersion, err := init_process(processStream, isRoot)
	if err != nil {
		return nil, fmt.Errorf("NewVmInstanceClientProcess: " + err.Error())
	}

	// Das VM Objekt wird erzeugt
	vmObject := &VmProcess{
		Log:         plog,
		rpcInstance: rpcInstance,
		processId:   procId,
		coreVersion: coreVersion,
		vmInstances: make([]types.VmInterface, 0),
	}

	// DEBUG
	vmObject.Log.Debug("Process '%s' initiated", procId)

	// Der Prozess wird am leben erhalten
	if err := alive_process(processStream, vmObject); err != nil {
		return nil, fmt.Errorf("NewVmInstanceClientProcess: " + err.Error())
	}

	// Das Objekt wird zurückgegeben
	return vmObject, nil
}
