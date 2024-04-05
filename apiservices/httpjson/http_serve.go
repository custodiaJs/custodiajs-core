package httpjson

func (o *HttpApiService) Serve(closeSignal chan struct{}) error {
	// Die Basis Urls werden hinzugefügt
	o.serverMux.HandleFunc("/", o.indexHandler)

	// Gibt die einzelnenen VM Informationen aus
	o.serverMux.HandleFunc("/vm", o.vmInfo)

	// Der VM-RPC Handler wird erstellt
	o.serverMux.HandleFunc("/rpc", o.httpCallFunction)

	// Der Websocket Console Stream wird hinzugefügt
	// der Console stream ist nur auf dem Localhost verfügbar
	if o.isLocalhost {
		o.serverMux.HandleFunc("/vm/console", o.handleConsoleStreamWebsocket)
	}

	// Der Webserver wird gestartet
	if err := o.serverObj.ListenAndServeTLS("", ""); err != nil {
		panic("Serve: " + err.Error())
	}

	// Der Vorgagn wurde ohne Fehler durchgeführt
	return nil
}
