package cligrpc

import (
	"log"
	"vnh1/apiservices/cligrpc/clirpc"
	grpccli "vnh1/grpc/cligrpc"
)

func (o *HostCliService) Serve(closeSignal chan struct{}) error {
	// Das CLI gRPC Serverobjekt wird erstellt
	cliGrpc := clirpc.NewCliGrpcServer(o.core)
	grpccli.RegisterCLIServiceServer(o.grpcServer, cliGrpc)

	// Der grpc Server wird gestartet
	if err := o.grpcServer.Serve(o.netListner); err != nil {
		log.Fatalf("Fehler beim Starten des gRPC-Servers: %v", err)
	}

	// Der Vorgagn wurde ohne Fehler durchgeführt
	return nil
}
