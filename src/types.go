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

// Gibt den Namen der Anwendung an
type CoreApplicationName string

// Gibt den Status des Cores an
type CoreState uint8

// Gibt die Version des Cores an
type Vesion uint32

// Gibt die Aktuelle Repo des Cores an
type CoreRepoUrl string

// Gibt den Aktuellen Preifx des Sockets an
type CoreIpcVmSocketIdentifierPrefix string

// Gibt den Path für die Generelle Core Config an
type CoreGeneralConfigPath string

// Gibt den Path eines IpcVm Sockets an
type CoreVmIpcSocketPath string
type CoreVmIpcSocketPathTemplate string

// Gibt den Path für das Logging dir an
type LogDirPath string

// Gibt die QUID einer VM an
type QVMID string

// Gibt den Hash eines Scriptes an
type VmScriptHash string

// Gibt die ProcessId an
type ProcessId string

// Gibt die VmID an
type VmId string

// Gibt den Status einer VM an
type VmState uint8

// Gibt die Vm Process ID an
type VmProcessId uint

// Gibt den Path des Crypto Stores an
type CoreCryptoStorePath string
