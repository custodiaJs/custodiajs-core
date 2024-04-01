package webgrpc

import (
	"fmt"
	"log"
	"net"
)

func (o *WebGrpcService) Serve(closeSignal chan struct{}) error {
	// Starte den gRPC-Server auf Port 50051
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", o.localAddress.LocalIP, o.localAddress.LocalPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := o.serverObj.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	// Der Vorgagn wurde ohne Fehler durchgeführt
	return nil
}
