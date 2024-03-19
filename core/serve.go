package core

import (
	"fmt"
)

func (o *Core) Serve() error {
	// Es wird signalisiert das der Core ausgeführt wird
	o.state = SERVING

	// Es werden alle Socketservices gestartet
	for _, item := range o.apiSockets {
		o.apiSyncWaitGroup.Add(1)
		go func(pitem APISocketInterface) {
			if err := pitem.Serve(o.serviceSignaling); err != nil {
				fmt.Println("Core:Serve: " + err.Error())
			}
			o.apiSyncWaitGroup.Done()
		}(item)
	}

	// Es werden alle Virtual Machines gestartet
	for _, item := range o.vms {
		o.vmSyncWaitGroup.Add(1)
		go func(item *CoreVM) {
			item.RunScript(item.jsCode)
			o.vmSyncWaitGroup.Done()
		}(item)
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
