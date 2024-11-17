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
	"io"
	"net"
	"os/user"
	"sync"

	"github.com/CustodiaJS/bngsocket"
	cenvxcore "github.com/custodia-cenv/cenvx-core/src"
	"github.com/custodia-cenv/cenvx-core/src/log"
)

// initVmIpcServer erstellt Sockets für Root, spezifische Gruppen, alle Benutzer und spezifische Benutzer, wenn der Prozess als Root ausgeführt wird.
// Ist der Prozess nicht Root, wird nur ein Socket für den aktuellen Benutzer erstellt.
func coreInitVmIpcServer(basePath string, acls []*ACL) error {
	// Es wird geprüft ob der VM-IPC Server initalisiert wurde
	if vmipcState != NEW {
		return fmt.Errorf("vm ipc always initalized")
	}

	// Log
	log.DebugLogPrint("The VM-IPC interface is prepared")

	// Das Open Connections Array wird erzeugt
	vmipcOpenConnections = make([]*bngsocket.BngConn, 0)
	vmipcListeners = make([]*_AclListener, 0)

	// Der Aktuelle Benutzer wird ermittelt
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("konnte aktuellen benutzer nicht ermitteln: %w", err)
	}

	// Ermitteln, ob der Prozess als Root ausgeführt wird
	isRoot := currentUser.Uid == "0"

	// Die Sockets werden erzeugt
	var aclListeners []*_AclListener
	var aclerr error
	if isRoot {
		aclListeners, aclerr = createAclListeners(acls, basePath)
	} else {
		// Es wird versucht für den Aktuellen Benutzer ein ACL zu erstellen
		cacl, err := createACLForCurrentUser()
		if err != nil {
			return err
		}
		aclListeners, aclerr = createAclListeners([]*ACL{cacl}, basePath)
	}

	// Es wird geprüft ob ein Fehler beim Erstellen der ACL Sockets aufgetreten ist
	if aclerr != nil {
		return aclerr
	}

	// Die VM-IPC wird als Initalisiert Makiert
	vmipcState = INITED

	// Sollte es sich um den Root Benutzer handeln, so werden 3 Sockets erzeugt,
	// ansonsten wird nur ein Socket für den Aktuellen Benutezr erzeugt
	totalIfaces := 0
	for _, item := range aclListeners {
		go processListenerGoroutine(item)
		totalIfaces = totalIfaces + 1
	}

	// LOG
	if totalIfaces == 1 {
		log.DebugLogPrint("The VM-IPC interface have been successfully prepared")
	} else {
		log.DebugLogPrint("The VM-IPC interfaces have been successfully prepared")
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

// Wird als Goroutine ausgeführt um eintreffende Socketanfragen zu verarbeiten
func processListenerGoroutine(nlist *_AclListener) {
	log.DebugLogPrint("VM-IPC listener started")
	for {
		conn, err := nlist.Accept()
		if err != nil {
			cstate := getVmIpcServerState()
			if cstate != CLOSING && cstate != CLOSED {
				corePanic(err)
				return
			}
		}

		// LOG
		log.DebugLogPrint("New VM-IPC connection accepted")
		go processConnectionGoroutine(conn)
	}
}

// Wird verwendet um zu ermitteln ob die Verbindung geschlossen wurde
func processConnectionGoroutine(conn net.Conn) {
	// Die Verbindung wird geupgradet
	upgraddedConn, err := bngsocket.UpgradeSocketToBngConn(conn)
	if err != nil {
	}

	// Die Verbindung wird zwischengespeichert
	addVmIpcConnection(upgraddedConn)

	// Es wird eine Go Routine gestartet, welche das Monitoring der Verbindung übernimmt
	go func() {
		// Die Verbindung wird nach abschluss der Funktion entfernt
		defer removeVmIpcConnection(upgraddedConn)

		// Es wird darauf gewartet dass die Verbindung geschlossen wird
		mresult := bngsocket.MonitorConnection(upgraddedConn)
		log.DebugLogPrint("VM-IPC connection closed")
		if mresult != nil {
			// Es wird geprüft ob die Verbindung Regulär getrennt wurde,
			// sollte die Verbindung nicht Regulär getrennt wurden sein,
			// so wird der Vorgang in den Error Log geschrieben
			if mresult != io.EOF && mresult != bngsocket.ErrConnectionClosedEOF {
				// Der Fehler wird in den Log geschrieben
				log.LogError("VM-IPC# Monitoring error: %s", mresult.Error())
				return
			}
		}
	}()

	// Die Kernfunktion wird Registriert, über diese Funktion Registriert sich ein Client
	upgraddedConn.RegisterFunction("init", func(req *bngsocket.BngRequest) error {
		return nil
	})
}

// Speichert eintreffende Verbindungen von VM Prozessen ab
func addVmIpcConnection(conn *bngsocket.BngConn) {
	coremutex.Lock()
	vmipcOpenConnections = append(vmipcOpenConnections, conn)
	coremutex.Unlock()
	log.DebugLogPrint("New VM-IPC connection cached")
}

// Entfernt VM Prozesse
func removeVmIpcConnection(conn *bngsocket.BngConn) {
	coremutex.Lock()
	vmipcOpenConnections = append(vmipcOpenConnections, conn)
	coremutex.Unlock()
	log.DebugLogPrint("VM-IPC connection removed from cache")
}

// Wird verwendet um allen Clients zu Signalisieren das der Core beendet wird und trennt die Verbindung zu allen Vorhandenen IPC-VMS
// Signalisiert allen VM's dass sie beendet werden
func signalCoreIsClosingAndCloseAllIpcConnections(wg *sync.WaitGroup) {
	// Es werden alle VM's abgearbeitet und geschlossen
	for _, item := range vms {
		wg.Add(1)
		go func(cvm cenvxcore.VmInterface) {
			cvm.SignalShutdown()
			wg.Done()
		}(item)
	}
}

// Gibt den Status des Aktuellen VM-IPC Servers zurück
func getVmIpcServerState() _VmIpcServerState {
	coremutex.Lock()
	t := vmipcState
	coremutex.Unlock()
	return t
}

// Schließt alle Verfügabren VM-IPC Listener
func closeVMIpcServer() {
	// Es wird geprüft ob der VmIPC Server geschlossen wurde
	if vmipcState != SERVING {
		return
	}

	// Der Status wird auf Closing gesetzt
	vmipcState = CLOSING

	// Es werden alle Sitzungen geschlossen
	for len(vmipcListeners) != 0 {
		item := vmipcListeners[0]
		vmipcListeners = vmipcListeners[1:]
		item.Close()
	}

	// Der Status wird auf Closed gesetzt
	vmipcState = CLOSED

	// LOG
	log.DebugLogPrint("All VM-IPC listeners are closed")
}
