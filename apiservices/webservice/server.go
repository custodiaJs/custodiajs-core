package webservice

import (
	"fmt"
	"net/http"
	"vnh1/types"
)

func (o *Webservice) Serve(closeSignal chan struct{}) error {
	// Die Basis Urls werden hinzugefügt
	http.HandleFunc("/", o.indexHandler)

	// Gibt die einzelnenen VM Informationen aus
	http.HandleFunc("/vm", o.vmInfo)

	// Der VM-RPC Handler wird erstellt
	http.HandleFunc("/rpc", o.vmRPCHandler)

	// Der Websocket Console Stream wird erzeugt
	http.HandleFunc("/consolestream", o.handleConsoleStreamWebsocket)

	// Der HTTP Server wird gestartet
	if err := http.ListenAndServe(":8080", nil); err != nil {
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
