// Author: fluffelpuff
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package core

import (
	"fmt"
	"sync"

	cenvxcore "github.com/custodia-cenv/cenvx-core/src"
)

// Wird verwendet um den Core geöffnet zu halten
func Serve() {
	// Der Status wird auf Serving gesetzt
	coreSetState(cenvxcore.SERVING, true)
	defer coreSetState(cenvxcore.SHUTDOWN, true)

	// Es wird ein neuer Waiter erzeugt
	waiter := &sync.WaitGroup{}

	// Es wird gewartet bis das Hold open geschlossen wird
	<-holdOpenChan

	// Der Objekt Mutex wird angewendet
	coremutex.Lock()

	// Der Status wird auf Shutdown gesetzt
	coreSetState(cenvxcore.SHUTDOWN, false)

	// Es wird allen Virtuellen CJS Vm's mitgeteilt dass der Core beendet wird,
	// die Funktion trennt nach dem Übermitteln des Signales alle IPC Verbindungen zu den VM's.
	signalCoreIsClosingAndCloseAllIpcConnections(waiter)

	// Es werden alle VM-IPC Server geschlossen
	closeVMIpcServer()

	// Der Status wird auf geschlossen gesetzt
	coreSetState(cenvxcore.CLOSED, false)

	// Der Objekt Mutex wird freigegeben
	coremutex.Unlock()

	// Es wird gewartet bis alle VM's geschlossen wurden
	waiter.Wait()

	// Es wird gewartet dass alle VM's geschlossen wurden
	vmSyncWaitGroup.Wait()

	// Log
	fmt.Println("Core closed, by.")
}
