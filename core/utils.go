package core

import (
	"sync"

	"github.com/CustodiaJS/custodiajs-core/global/types"
)

// Legt den Core Status fest
func setState(core *Core, state types.CoreState, useMutex bool) {
	// Es wird geprüft ob Mutex verwendet werden sollen
	if useMutex {
		core.objectMutex.Lock()
		defer core.objectMutex.Unlock()
	}

	// Es wird geprüft ob der neue Status, der Aktuelle ist
	if core.state == state {
		return
	}

	// Der Neue Status wird gesetzt
	core.state = state
}

// Signalisiert allen VM's dass sie beendet werden
func closeAllVirtualMachines(o *Core, wg *sync.WaitGroup) {
	// Es werden alle VM's abgearbeitet und geschlossen
	for _, item := range o.vms {
		wg.Add(1)
		go func(cvm types.VmInterface) {
			cvm.SignalShutdown()
			wg.Done()
		}(item)
	}
}
