package core

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/CustodiaJS/bngsocket"
)

// initVmIpcServer erstellt Sockets für Root, spezifische Gruppen, alle Benutzer und spezifische Benutzer, wenn der Prozess als Root ausgeführt wird.
// Ist der Prozess nicht Root, wird nur ein Socket für den aktuellen Benutzer erstellt.
func InitVmIpcServer(basePath string, groupNames, userNames []string) error {
	// Es wird geprüft ob der Core Initalisiert wurde
	if !CoreIsInited() {
		return fmt.Errorf("core ist not initalized")
	}

	coremutex.Lock()
	defer coremutex.Unlock()

	// Es wird geprüft ob der VM-IPC Server initalisiert wurde
	if vmipcInited {
		return fmt.Errorf("vm ipc always initalized")
	}

	// Das Open Connections Array wird erzeugt
	vmipcOpenConnections = make([]*bngsocket.BngSocket, 0)

	// Der Aktuelle Benutzer wird ermittelt
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("konnte aktuellen benutzer nicht ermitteln: %w", err)
	}

	// Ermitteln, ob der Prozess als Root ausgeführt wird
	isRoot := currentUser.Uid == "0"

	// Socket für den aktuellen Benutzer, falls nicht Root
	if !isRoot {
		// Der Listener für den Aktuellen Benutzer wird erstellt
		userListener, err := createSocketWithHelper(basePath, fmt.Sprintf("user_%s_socket.sock", currentUser.Username), atoi(currentUser.Uid), -1, 0600)
		if err != nil {
			return err
		}

		// Der Listener wird als Spefic Listener zwischnegespeichert
		vmipcSpecificListeners := map[string]net.Listener{}
		vmipcSpecificListeners[fmt.Sprintf("user:%s", currentUser)] = userListener

		// Der Einzelnen Listener wird in einer Acceptor Goroutine verwendet
		go processListenerGoroutine(userListener)

		// Es ist kein Fehler vorhanden
		return nil
	}

	// Sockets für Root und allgemeine Benutzer erstellen
	vmipcRootListener, err = createSocketWithHelper(basePath, "root_socket.sock", 0, 0, 0600)
	if err != nil {
		return fmt.Errorf("fehler beim erstellen des root-sockets: %w", err)
	}

	vmipcSpecificListeners, err = createGroupAndUserSockets(basePath, groupNames, userNames)
	if err != nil {
		return err
	}

	vmipcAllUsersListener, err = createSocketWithHelper(basePath, "public_socket.sock", 0, 0, 0666)
	if err != nil {
		return fmt.Errorf("fehler beim erstellen des allgemeinen sockets: %w", err)
	}

	// Die Einzelnen Listener werden in einer Acceptor Goroutine verwendet
	go processListenerGoroutine(vmipcAllUsersListener)
	go processListenerGoroutine(vmipcRootListener)
	for _, item := range vmipcSpecificListeners {
		go processListenerGoroutine(item)
	}

	// Die VM-IPC wird als Initalisiert Makiert
	vmipcInited = true

	return nil
}

// Wird als Goroutine ausgeführt um eintreffende Socketanfragen zu verarbeiten
func processListenerGoroutine(nlist net.Listener) {
	for {
		conn, err := nlist.Accept()
		if err != nil {
		}

		upgraddedConn, err := bngsocket.UpgradeSocketToBngSocket(conn)
		if err != nil {
		}

		if addVmIpcConnection(upgraddedConn); err != nil {
		}
	}
}

// Speichert eintreffende Verbindungen von VM Prozessen ab
func addVmIpcConnection(conn *bngsocket.BngSocket) error {
	return nil
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
	}

	return specificListeners, nil
}

// Hilfsfunktion zur Umwandlung von UID-Strings in int, um Fehlerbehandlung zu verkürzen
func atoi(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}
