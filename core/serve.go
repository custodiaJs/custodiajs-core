package core

import (
	"fmt"
	"vnh1/static"
)

func (o *Core) Serve() error {
	// Es wird signalisiert das der Core ausgeführt wird
	o.state = SERVING

	// Es werden alle Socketservices gestartet
	for _, item := range o.apiSockets {
		o.apiSyncWaitGroup.Add(1)
		go func(pitem static.APISocketInterface) {
			if err := pitem.Serve(o.serviceSignaling); err != nil {
				fmt.Println("Core:Serve: " + err.Error())
			}
			o.apiSyncWaitGroup.Done()
		}(item)
	}

	// Es werden alle Virtual Machines gestartet
	for _, item := range o.vms {
		if err := item.serveGorutine(&o.vmSyncWaitGroup); err != nil {
			return fmt.Errorf("Serve: " + err.Error())
		}
	}

	// Es wird gewartet bis das Hold open geschlossen wird
	<-o.holdOpenChan

	// Der Beenden wird vorbereitet
	o.prepareForShutdown()

	// Der Vorgang wurde ohne Fehöer durchgeführt
	return nil
}

func (o *Core) SignalShutdown() {
	close(o.holdOpenChan)
}

func (o *Core) prepareForShutdown() {
	fmt.Println("By!")
}
