package webservice

import (
	"context"
	"vnh1/grpc/publicgrpc"
)

func (s *GrpcServer) CallFunction(ctx context.Context, in *publicgrpc.RPCFunctionCall) (*publicgrpc.RPCResponse, error) {
	// Beispiel: Implementiere die Logik, um die Funktion basierend auf in.FunctionName aufzurufen
	// und die Parameter in.Parms zu verarbeiten.

	// Beispielantwort
	return &publicgrpc.RPCResponse{
		Result: "Erfolg",
		Data: &publicgrpc.RPCResponseData{
			Type:  "string",
			Value: &publicgrpc.RPCResponseData_StringValue{StringValue: "Beispielantwort"},
		},
		Error: "",
	}, nil
}
