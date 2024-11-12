package core

import "fmt"

// Signalisiert dem Core, dass er beendet werden soll
func SignalShutdown() {
	// Log
	fmt.Println("Closing CustodiaJS...")

	// Der Mutex wird angewendet
	objectMutex.Lock()
	defer objectMutex.Unlock()

	// Die Chan wird geschlossen
	close(holdOpenChan)
}
