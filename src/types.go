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

type VmWorkingDir string
type VmProcessId uint64
type QualifiedVmID string

// VM und Core Status Typen sowie Repo Datentypen
type ALTERNATIVE_SERVICE_PATH string        // Alternativer Socket Path
type VmState uint8                          // VM Status
type CoreState uint8                        // Core Status
type IPCRight uint8                         // CLI Benutzerrecht
type VERSION uint32                         // Version des Hauptpgrogrammes
type REPO string                            // URL der Sourccode Qeulle
type SOCKET_PATH string                     // Gibt einen Socket Path an
type LOG_DIR string                         // Gibt den Path des Log Dir's unter
type HOST_CRYPTOSTORE_WATCH_DIR_PATH string // Gibt den Ordner an, in dem sich alle Zertifikate und Schlüssel des Hosts befinden
type HOST_CONFIG_FILE_PATH string           // Gibt den Pfad der Config Datei an
type HOST_CONFIG_PATH string
type CHN_CORE_SOCKET_PATH string

// Gibt die QUID einer VM an
type QVMID string

// Gibt den Hash eines Scriptes an
type VmScriptHash string

// Gibt die ProcessId an
type ProcessId string

// Gibt die VmID an
type VmId string

// Gibt an ob es sich bei einer IP um eine TOR IP-handelt
type TorIpState bool
