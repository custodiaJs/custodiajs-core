package webgrpc

import (
	"fmt"
	"io"
	"log"
	"vnh1/grpc/publicgrpc"
)

func (s *GrpcServer) Chat(stream publicgrpc.ChatService_ChatServer) error {
	fmt.Println("STREAM")
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// Ende des Streams
			return nil
		}
		if err != nil {
			log.Fatalf("Fehler beim Empfangen der Chat-Nachricht: %v", err)
			return err
		}
		log.Printf("Nachricht von %s: %s", in.GetUser(), in.GetMessage())
		if err := stream.Send(in); err != nil {
			log.Fatalf("Fehler beim Senden der Chat-Nachricht: %v", err)
			return err
		}
	}
}
