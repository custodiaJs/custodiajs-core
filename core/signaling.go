package core

import "fmt"

// Signalisiert dem Core, dass er beendet werden soll
func SignalShutdown() {
	// Log
	fmt.Println("Closing CustodiaJS...")

	// Der Mutex wird angewendet
	coremutex.Lock()
	defer coremutex.Unlock()

	// Die Chan wird geschlossen
	close(holdOpenChan)
}
