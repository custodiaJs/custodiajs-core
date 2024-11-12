package core

import (
	"fmt"
	"sync"

	"github.com/CustodiaJS/custodiajs-core/global/static"
	"github.com/CustodiaJS/custodiajs-core/global/types"
)

// Wird verwendet um den Core geöffnet zu halten
func Serve() {
	// Der Mutex wird angewendet
	coremutex.Lock()

	// Der Status wird auf Serving gesetzt
	setState(static.SERVING, false)
	defer setState(static.SHUTDOWN, true)

	// Es werden alle Socketservices gestartet
	for hight, item := range apiSockets {
		// Die API Sockets werden ausgeführt
		apiSyncWaitGroup.Add(1)
		go func(hight int, pitem types.APISocketInterface) {
			if err := pitem.Serve(serviceSignaling); err != nil {
				fmt.Printf("API error:: %s\n", err.Error())
			}
			apiSyncWaitGroup.Done()
		}(hight, item)
	}

	// Es werden alle Virtual Machines gestartet
	for _, item := range vms {
		// Es wird Signalisiert dass eine VM Instanz mehr ausgeführt wird
		vmSyncWaitGroup.Add(1)

		// Die VM wird ausgeführt
		item.Serve(&vmSyncWaitGroup)
		/*if err := item.Serve(&vmSyncWaitGroup); err != nil {

		}*/
	}

	// Der Mutex wird freigegeben
	coremutex.Unlock()

	// Es wird ein neuer Waiter erzeugt
	waiter := &sync.WaitGroup{}

	// Es wird gewartet bis das Hold open geschlossen wird
	<-holdOpenChan

	// Der Objekt Mutex wird angewendet
	coremutex.Lock()

	// Der Status wird auf Shutdown gesetzt
	setState(static.SHUTDOWN, false)

	// Der Beenden wird vorbereitet
	closeAllVirtualMachines(waiter)

	// Der Objekt Mutex wird freigegeben
	coremutex.Unlock()

	// Es wird gewartet bis alle VM's geschlossen wurden
	waiter.Wait()

	// Es wird gewartet dass alle VM's geschlossen wurden
	vmSyncWaitGroup.Wait()

	// Log
	fmt.Println("Core closed, by.")
}
