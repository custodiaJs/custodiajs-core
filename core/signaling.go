package core

import "fmt"

// Signalisiert dem Core, dass er beendet werden soll
func (o *Core) SignalShutdown() {
	// Log
	fmt.Println("Closing CustodiaJS...")

	// Der Mutex wird angewendet
	o.objectMutex.Lock()
	defer o.objectMutex.Unlock()

	// Die Chan wird geschlossen
	close(o.holdOpenChan)
}
