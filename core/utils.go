package core

import (
	"sync"

	"github.com/CustodiaJS/custodiajs-core/global/types"
)

// Legt den Core Status fest
func setState(tstate types.CoreState, useMutex bool) {
	// Es wird geprüft ob Mutex verwendet werden sollen
	if useMutex {
		objectMutex.Lock()
		defer objectMutex.Unlock()
	}

	// Der Neue Status wird gesetzt
	cstate = tstate
}

// Signalisiert allen VM's dass sie beendet werden
func closeAllVirtualMachines(wg *sync.WaitGroup) {
	// Es werden alle VM's abgearbeitet und geschlossen
	for _, item := range vms {
		wg.Add(1)
		go func(cvm types.VmInterface) {
			cvm.SignalShutdown()
			wg.Done()
		}(item)
	}
}
