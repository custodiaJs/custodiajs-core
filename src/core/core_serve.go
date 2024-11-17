package core

import (
	"fmt"
	"sync"

	"github.com/custodia-cenv/cenvx-core/src/static"
)

// Wird verwendet um den Core geöffnet zu halten
func Serve() {
	// Der Status wird auf Serving gesetzt
	coreSetState(static.SERVING, true)
	defer coreSetState(static.SHUTDOWN, true)

	// Es wird ein neuer Waiter erzeugt
	waiter := &sync.WaitGroup{}

	// Es wird gewartet bis das Hold open geschlossen wird
	<-holdOpenChan

	// Der Objekt Mutex wird angewendet
	coremutex.Lock()

	// Der Status wird auf Shutdown gesetzt
	coreSetState(static.SHUTDOWN, false)

	// Es wird allen Virtuellen CJS Vm's mitgeteilt dass der Core beendet wird,
	// die Funktion trennt nach dem Übermitteln des Signales alle IPC Verbindungen zu den VM's.
	signalCoreIsClosingAndCloseAllIpcConnections(waiter)

	// Es werden alle VM-IPC Server geschlossen
	closeVMIpcServer()

	// Der Status wird auf geschlossen gesetzt
	coreSetState(static.CLOSED, false)

	// Der Objekt Mutex wird freigegeben
	coremutex.Unlock()

	// Es wird gewartet bis alle VM's geschlossen wurden
	waiter.Wait()

	// Es wird gewartet dass alle VM's geschlossen wurden
	vmSyncWaitGroup.Wait()

	// Log
	fmt.Println("Core closed, by.")
}
