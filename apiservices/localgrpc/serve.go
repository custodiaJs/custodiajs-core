package localgrpc

import (
	"log"
	"vnh1/apiservices/localgrpc/localgrpcservice"
	"vnh1/localgrpcproto"
)

func (o *HostCliService) Serve(closeSignal chan struct{}) error {
	// Das CLI gRPC Serverobjekt wird erstellt
	localgrpc := localgrpcservice.NewCliGrpcServer(o.core)
	localgrpcproto.RegisterCLIServiceServer(o.grpcServer, localgrpc)

	// Der grpc Server wird gestartet
	if err := o.grpcServer.Serve(o.netListner); err != nil {
		log.Fatalf("Fehler beim Starten des gRPC-Servers: %v", err)
	}

	// Der Vorgagn wurde ohne Fehler durchgeführt
	return nil
}
