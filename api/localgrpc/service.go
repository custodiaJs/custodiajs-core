package localgrpc

import (
	"fmt"
	"net"
	"runtime"

	"github.com/CustodiaJS/custodiajs-core/global/procslog"
	"github.com/CustodiaJS/custodiajs-core/global/types"
	"github.com/CustodiaJS/custodiajs-core/ipc"
	"github.com/google/uuid"
)

// Hält den Server am leben
func (o *HostAPIService) Serve(closeSignal chan struct{}) error {

	// Der Vorgagn wurde ohne Fehler durchgeführt
	return nil
}

// Verknüpft den Core mit dem Service
func (o *HostAPIService) LinkCore(coreObj types.CoreInterface) error {
	// Der Core wird abgespeichert
	o.core = coreObj

	o.procLog.Debug("Core linked")

	// Der Vorgang ist ohne fehler durchgeführt wurden
	return nil
}

// Erstellt einen neuen API Context
func (o *HostAPIService) NewAPIContext() *APIContext {
	// Es wird eine neue UUID erzeugt
	prod_uuid := uuid.New().String()

	// Es wird eine neue Child Session erstellt
	child_session := o.procLog.GetChildLog("ProcessInstance")

	// Es wird ein neuer Context erzeugt
	new_context := &APIContext{procUUID: types.ProcessId(prod_uuid), Log: child_session}

	// Der Context wird zwischengespeichert
	o.processes[prod_uuid] = new_context

	// Die UID wird zurückgegeben
	return new_context
}

// Wird verwendet um einen API Context abzurufen
func (o *HostAPIService) GetContextByProcessId(procId types.ProcessId) *APIContext {
	f := o.processes[string(procId)]
	f.Log.Debug("Passed from memory")
	return f
}

// Wird verwendet um eine Lokale Kommunikation über UnixSockets oder Named Windows Pipes zu ermöglichen
func New(unixOrWinNamedPipeAddr types.SOCKET_PATH, userRightState types.IPCRight) (*HostAPIService, error) {
	// Es wird passend zum Hostos der Richtige Listener erzeugt
	var cliSocket net.Listener
	var err error
	switch runtime.GOOS {
	case "windows":
		err = fmt.Errorf("not supported os")
	case "darwin", "linux":
		cliSocket, err = ipc.CreateNewUnixSocket(string(unixOrWinNamedPipeAddr), userRightState)
	default:
		err = fmt.Errorf("unkown os")
	}

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if err != nil {
		return nil, fmt.Errorf("New: " + err.Error())
	}

	// Das HostCLI Objekt wird erstellt
	hcs := &HostAPIService{
		netListner: cliSocket,
		processes:  make(map[string]*APIContext),
		procLog:    procslog.NewProcLogForHostAPISocket(),
	}

	// LOG
	hcs.procLog.Log("Created on '%s'", unixOrWinNamedPipeAddr)

	// Das Objekt wird zurückgegeben
	return hcs, nil
}
