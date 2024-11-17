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

package cenvxcore

import "sync"

type VmInterface interface {
	// Gibt das Manifest zurück
	GetManifest() *Manifest
	// Gibt den Scripthash zurück
	GetScriptHash() string
	// Gibt den Consolen Output Watcher zurück
	GetConsoleOutputWatcher() interface{}
	// Gibt alle Geteilten RPC Funktionen zurück
	GetAllSharedFunctions() []interface{}
	// Hält die Vm am leben
	Serve(*sync.WaitGroup) error
	// Gibt eine Geteilte Funktion anhand ihrer Signatur zurück
	GetSharedFunctionBySignature(interface{}, *interface{}) (interface{}, bool, *interface{})
	// Gibt den Timestamp zurück der angebit wann die VM gestartet wurde
	GetStartingTimestamp() uint64
	// Signalisiert dass die VM beendet werden soll
	SignalShutdown()
	// Gibt den Aktuellen Status der VM zurück
	GetState() VmState
	// Gibt die ProzessID zurück
	GetProcessId() VmProcessId
	// Gibt die Qualified Full VM ID (QVMID) zurück
	GetQVMID() VmId
}
