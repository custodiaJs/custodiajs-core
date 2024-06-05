package core

import (
	"fmt"
	"sync"
	"vnh1/static"
	"vnh1/types"
)

func (o *Core) Serve() error {
	// Der Mutex wird angewendet
	o.objectMutex.Lock()

	// Der Status wird auf Serving gesetzt
	setState(o, static.SERVING, false)
	defer setState(o, static.SHUTDOWN, true)

	// Es werden alle Socketservices gestartet
	for _, item := range o.apiSockets {
		// Die API Sockets werden ausgeführt
		o.apiSyncWaitGroup.Add(1)
		go func(pitem types.APISocketInterface) {
			if err := pitem.Serve(o.serviceSignaling); err != nil {
				fmt.Println("Core:Serve: " + err.Error())
			}
			o.apiSyncWaitGroup.Done()
		}(item)
	}

	// Es werden alle Virtual Machines gestartet
	for _, item := range o.vms {
		// Es wird Signalisiert dass eine VM Instanz mehr ausgeführt wird
		o.vmSyncWaitGroup.Add(1)

		// Die VM wird ausgeführt
		if err := item.Serve(&o.vmSyncWaitGroup); err != nil {
			return fmt.Errorf("Serve: " + err.Error())
		}
	}

	// Der Mutex wird freigegeben
	o.objectMutex.Unlock()

	// Es wird ein neuer Waiter erzeugt
	waiter := &sync.WaitGroup{}

	// Es wird gewartet bis das Hold open geschlossen wird
	<-o.holdOpenChan

	// Der Objekt Mutex wird angewendet
	o.objectMutex.Lock()

	// Der Status wird auf Shutdown gesetzt
	setState(o, static.SHUTDOWN, false)

	// Der Beenden wird vorbereitet
	o.signalVmsShutdown(waiter)

	// Der Objekt Mutex wird freigegeben
	o.objectMutex.Unlock()

	// Es wird gewartet bis alle VM's geschlossen wurden
	waiter.Wait()

	// Es wird gewartet dass alle VM's geschlossen wurden
	o.vmSyncWaitGroup.Wait()

	// Log
	fmt.Println("Core closed, by.")

	// Der Vorgang wurde ohne Fehöer durchgeführt
	return nil
}
