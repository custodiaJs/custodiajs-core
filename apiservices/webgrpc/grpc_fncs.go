package webgrpc

import (
	"context"
	"fmt"
	"log"
	"vnh1/grpc/publicgrpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func (s *GrpcServer) CallFunction(ctx context.Context, in *publicgrpc.RPCFunctionCall) (*publicgrpc.RPCResponse, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		log.Println("Keine Peer-Informationen gefunden")
		return nil, status.Error(codes.Internal, "Fehler beim Abrufen der Peer-Informationen")
	}

	tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		log.Println("Keine TLS-Authentifizierungsinformationen gefunden")
		return nil, status.Error(codes.Unauthenticated, "Keine TLS-Authentifizierungsinformationen")
	}

	fmt.Println(tlsInfo)

	return &publicgrpc.RPCResponse{
		Result: "Erfolg",
		Data: &publicgrpc.RPCResponseData{
			Type:  "string",
			Value: &publicgrpc.RPCResponseData_StringValue{StringValue: "Beispielantwort"},
		},
		Error: "",
	}, nil
}
