package webservice

import (
	"fmt"
	"vnh1/types"
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

	// Der Websocket gRPC Stream wird erzeugt
	o.serverMux.HandleFunc("/grpc", o.handleGRPC)

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

func (o *Webservice) SetupCore(coreObj types.CoreInterface) error {
	// Es wird geprüft ob der Core festgelegt wurde
	if o.core != nil {
		return fmt.Errorf("SetupCore: always linked with core")
	}

	// Das Objekt wird zwischengespeichert
	o.core = coreObj

	// Der Vorgang ist ohne fehler durchgeführt wurden
	return nil
}
