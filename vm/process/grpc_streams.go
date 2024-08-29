package process

import (
	"context"
	"log"
	"time"

	"github.com/CustodiaJS/custodiajs-core/localgrpcproto"
)

func grpc_stream_process_control_serve(vmproc *VmProcess) error {
	// DEBUG
	vmproc.Log.Debug("Starting Process Control GRPC-Stream")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	stream, err := vmproc.rpcInstance.CoreToProcessControl(ctx)
	if err != nil {
		log.Fatalf("error while opening stream: %v", err)
	}

	// Dem Server wird ein HelloProcess Paket mit der ID des Prozesses geschickt
	// Sende eine TCC Nachricht
	vmproc.Log.Debug("Write Process Register paket")
	if err := stream.Send(&localgrpcproto.ProcessControlPackage{Req: &localgrpcproto.ProcessControlPackage_RegisterProcess{RegisterProcess: string(vmproc.processId)}}); err != nil {
		log.Fatalf("error while sending TCC: %v", err)
	}

	// DEBUG
	vmproc.Log.Debug("Process Control stream registered")

	cancel()

	return nil
}
