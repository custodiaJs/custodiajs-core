package webservice

import (
	"fmt"
	"log"
)

func (o *Webservice) Serve(closeSignal chan struct{}) error {
	// Die Basis Urls werden hinzugefügt
	o.serverMux.HandleFunc("/", o.indexHandler)

	// Gibt die einzelnenen VM Informationen aus
	o.serverMux.HandleFunc("/vm", o.vmInfo)

	// Der VM-RPC Handler wird erstellt
	o.serverMux.HandleFunc("/rpc", o.httpRPCHandler)

	// Der Websocket Console Stream wird hinzugefügt
	// der Console stream ist nur auf dem Localhost verfügbar
	if o.isLocalhost {
		o.serverMux.HandleFunc("/vm/console", o.handleConsoleStreamWebsocket)
	}

	// Hier fehlt die Registrierung Ihrer gRPC-Services
	go func() {
		if err := o.grpcServer.Serve(o.grpcSocket); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Der HTTP Server wird gestartet
	go func() {
		if err := o.serverObj.Serve(o.httpSocket); err != nil {
			panic("Serve: " + err.Error())
		}
	}()

	// Der Mux Server wird gestartet
	if err := o.tcpMux.Serve(); err != nil {
		return fmt.Errorf("Serve: " + err.Error())
	}

	// Der Vorgagn wurde ohne Fehler durchgeführt
	return nil
}
