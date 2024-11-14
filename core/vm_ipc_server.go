package core

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/CustodiaJS/bngsocket"
	"github.com/CustodiaJS/custodiajs-core/global/types"
	"github.com/CustodiaJS/custodiajs-core/log"
)

// initVmIpcServer erstellt Sockets für Root, spezifische Gruppen, alle Benutzer und spezifische Benutzer, wenn der Prozess als Root ausgeführt wird.
// Ist der Prozess nicht Root, wird nur ein Socket für den aktuellen Benutzer erstellt.
func coreInitVmIpcServer(basePath string, groupNames, userNames []string) error {
	// Es wird geprüft ob der VM-IPC Server initalisiert wurde
	if vmipcState != NEW {
		return fmt.Errorf("vm ipc always initalized")
	}

	// Log
	log.DebugLogPrint("The VM-IPC interface is prepared")

	// Das Open Connections Array wird erzeugt
	vmipcOpenConnections = make([]*bngsocket.BngConn, 0)
	vmipcSpecificListeners = make(map[string]net.Listener)

	// Der Aktuelle Benutzer wird ermittelt
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("konnte aktuellen benutzer nicht ermitteln: %w", err)
	}

	// Ermitteln, ob der Prozess als Root ausgeführt wird
	isRoot := currentUser.Uid == "0"

	// Die VM-IPC wird als Initalisiert Makiert
	vmipcState = INITED

	// Sollte es sich um den Root Benutzer handeln, so werden 3 Sockets erzeugt,
	// ansonsten wird nur ein Socket für den Aktuellen Benutezr erzeugt
	totalIfaces := 0
	if isRoot {
		// Sockets für Root und allgemeine Benutzer erstellen
		vmipcRootListener, err = createSocketWithHelper(basePath, "root_socket.sock", 0, 0, 0600)
		if err != nil {
			return fmt.Errorf("fehler beim erstellen des root-sockets: %w", err)
		}

		vmipcSpecificListeners, err = createGroupAndUserSockets(basePath, groupNames, userNames)
		if err != nil {
			return err
		}

		// Die Einzelnen Listener werden in einer Acceptor Goroutine verwendet
		go processListenerGoroutine(vmipcRootListener)
		for _, item := range vmipcSpecificListeners {
			go processListenerGoroutine(item)
		}

		// Die VM-IPC wird als Serving Makiert
		vmipcState = SERVING
	} else {
		// Der Listener für den Aktuellen Benutzer wird erstellt
		userListener, err := createSocketWithHelper(basePath, fmt.Sprintf("user_%s_socket.sock", currentUser.Username), atoi(currentUser.Uid), -1, 0600)
		if err != nil {
			return err
		}

		// LOG
		log.DebugLogPrint("VM-IPC created for the current user %s", currentUser.Username)

		// Der Listener wird als Spefic Listener zwischnegespeichert
		vmipcSpecificListeners[fmt.Sprintf("user:%s", currentUser.Username)] = userListener

		// Der Einzelnen Listener wird in einer Acceptor Goroutine verwendet
		go processListenerGoroutine(userListener)

		// Die VM-IPC wird als Serving Makiert
		vmipcState = SERVING
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
func processListenerGoroutine(nlist net.Listener) {
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

// createSocketForUser erstellt einen UNIX-Socket mit Berechtigungen für einen bestimmten Benutzer oder eine Gruppe
func createSocketForUser(socketPath string, uid, gid int, permissions os.FileMode) (net.Listener, error) {
	if _, err := os.Stat(socketPath); err == nil {
		if err := os.Remove(socketPath); err != nil {
			return nil, fmt.Errorf("konnte bestehenden socket nicht entfernen: %w", err)
		}
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("konnte unix-socket nicht erstellen: %w", err)
	}

	if err := os.Chown(socketPath, uid, gid); err != nil {
		listener.Close()
		return nil, fmt.Errorf("konnte eigentümer des sockets nicht setzen: %w", err)
	}
	if err := os.Chmod(socketPath, permissions); err != nil {
		listener.Close()
		return nil, fmt.Errorf("konnte berechtigungen des sockets nicht setzen: %w", err)
	}

	return listener, nil
}

// createSocketWithHelper erstellt einen Socket mit den angegebenen Berechtigungen
func createSocketWithHelper(basePath, name string, uid, gid int, perms os.FileMode) (net.Listener, error) {
	socketPath := filepath.Join(basePath, name)
	return createSocketForUser(socketPath, uid, gid, perms)
}

// createGroupAndUserSockets erstellt Sockets für die angegebenen Gruppen und Benutzer und gibt sie in einer Liste zurück
func createGroupAndUserSockets(basePath string, groupNames, userNames []string) (map[string]net.Listener, error) {
	specificListeners := make(map[string]net.Listener)

	for _, groupName := range groupNames {
		if _, found := specificListeners[fmt.Sprintf("group:%s", groupName)]; found {
			continue
		}

		group, err := user.LookupGroup(groupName)
		if err != nil {
			return nil, fmt.Errorf("konnte gruppe %s nicht finden: %w", groupName, err)
		}

		gid, _ := strconv.Atoi(group.Gid)
		groupSocketPath := filepath.Join(basePath, fmt.Sprintf("group_%s_socket.sock", groupName))
		listener, err := createSocketForUser(groupSocketPath, 0, gid, 0660)
		if err != nil {
			return nil, fmt.Errorf("fehler beim erstellen des sockets für gruppe %s: %w", groupName, err)
		}
		specificListeners[fmt.Sprintf("group:%s", groupName)] = listener
		log.DebugLogPrint("Created VM-IPC listener for user group: %s", groupName)
	}

	for _, userName := range userNames {
		if _, found := specificListeners[fmt.Sprintf("user:%s", userName)]; found {
			continue
		}

		userInfo, err := user.Lookup(userName)
		if err != nil {
			return nil, fmt.Errorf("konnte benutzer %s nicht finden: %w", userName, err)
		}

		uid, _ := strconv.Atoi(userInfo.Uid)
		userSocketPath := filepath.Join(basePath, fmt.Sprintf("user_%s_socket.sock", userName))
		listener, err := createSocketForUser(userSocketPath, uid, -1, 0600)
		if err != nil {
			return nil, fmt.Errorf("fehler beim erstellen des sockets für benutzer %s: %w", userName, err)
		}
		specificListeners[fmt.Sprintf("user:%s", userName)] = listener
		log.DebugLogPrint("Created VM-IPC listener for user: %s", userName)
	}

	return specificListeners, nil
}

// Wird verwendet um allen Clients zu Signalisieren das der Core beendet wird und trennt die Verbindung zu allen Vorhandenen IPC-VMS
// Signalisiert allen VM's dass sie beendet werden
func signalCoreIsClosingAndCloseAllIpcConnections(wg *sync.WaitGroup) {
	// Es werden alle VM's abgearbeitet und geschlossen
	for _, item := range vms {
		wg.Add(1)
		go func(cvm types.VmInterface) {
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

	if vmipcRootListener != nil {
		vmipcRootListener.Close()
	}

	// Die Spiziellen Werden aufgelistet
	extracted := []string{}
	for t, item := range vmipcSpecificListeners {
		item.Close()
		extracted = append(extracted, t)
		log.DebugLogPrint("VM-IPC listener '%s' are closed", t)
	}
	for _, item := range extracted {
		delete(vmipcSpecificListeners, item)
	}

	// Der Status wird auf Closed gesetzt
	vmipcState = CLOSED

	// LOG
	log.DebugLogPrint("All VM-IPC listeners are closed")
}

// Hilfsfunktion zur Umwandlung von UID-Strings in int, um Fehlerbehandlung zu verkürzen
func atoi(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}
