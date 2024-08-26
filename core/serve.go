package core

import (
	"fmt"
	"sync"

	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/types"
)

// Wird verwendet um den Core geöffnet zu halten
func (o *Core) Serve() {
	// Der Mutex wird angewendet
	o.objectMutex.Lock()

	// Der Status wird auf Serving gesetzt
	setState(o, static.SERVING, false)
	defer setState(o, static.SHUTDOWN, true)

	// Es werden alle Socketservices gestartet
	for hight, item := range o.apiSockets {
		// Die API Sockets werden ausgeführt
		o.apiSyncWaitGroup.Add(1)
		go func(hight int, pitem types.APISocketInterface) {
			if err := pitem.Serve(o.serviceSignaling); err != nil {
				fmt.Printf("API error:: %s\n", err.Error())
			}
			o.apiSyncWaitGroup.Done()
		}(hight, item)
	}

	// Es werden alle Virtual Machines gestartet
	for _, item := range o.vms {
		// Es wird Signalisiert dass eine VM Instanz mehr ausgeführt wird
		o.vmSyncWaitGroup.Add(1)

		// Die VM wird ausgeführt
		item.Serve(&o.vmSyncWaitGroup)
		/*if err := item.Serve(&o.vmSyncWaitGroup); err != nil {

		}*/
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
	closeAllVirtualMachines(o, waiter)

	// Der Objekt Mutex wird freigegeben
	o.objectMutex.Unlock()

	// Es wird gewartet bis alle VM's geschlossen wurden
	waiter.Wait()

	// Es wird gewartet dass alle VM's geschlossen wurden
	o.vmSyncWaitGroup.Wait()

	// Log
	fmt.Println("Core closed, by.")
}
