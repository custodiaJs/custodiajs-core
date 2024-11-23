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
	"sync"

	"github.com/custodia-cenv/bngsocket-go"
	cenvxcore "github.com/custodia-cenv/cenvx-core/src"
	"github.com/custodia-cenv/cenvx-core/src/crypto"
)

// Core Mutex
var coremutex *sync.Mutex = new(sync.Mutex)

// VM-IPC Sockets
var (
	vmipcState           _VmIpcServerState = NEW
	vmipcListeners       []*_AclListener
	vmipcOpenConnections []*bngsocket.BngConn
)

// Speichert alle VM's ab, welche dem Core bekannt sind
var (
	vmsByID         map[string]cenvxcore.VmInterface
	vmsByName       map[string]cenvxcore.VmInterface
	vmSyncWaitGroup sync.WaitGroup
	vms             []cenvxcore.VmInterface
)

// Core Variablen
var (
	coreState    cenvxcore.CoreState = cenvxcore.NEW
	cryptoStore  *crypto.CryptoStore
	holdOpenChan chan struct{}
)
