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

package kmodulerpc

import (
	"sync"
	"vnh1/types"

	v8 "rogchap.com/v8go"
)

type FunctionType string

const (
	Local  FunctionType = "local"
	Public FunctionType = "public"
)

type SharedFunction struct {
	cond               *sync.Cond            // Speichert den Objekt Cond ab
	mutex              *sync.Mutex           // Speichert den Objektmutex ab
	kernel             types.KernelInterface // Speichert den Verwendeten Kernel ab
	v8Function         *v8.Function          // Speichert die eigentliche Funktion des V8 Codes ab
	name               string                // Speichert den Namen der Funktion ab
	parmTypes          []string              // Speichert die Parameterdatentypen welche die Funktion erwartet ab
	returnType         string                // Speichert den Rückgabetypen welche die Funktion zurück gibt ab
	eventOnRequest     []*v8.Function        // Speichert die Funktionen ab, welche aufgerufen werden sobald das Event eintritt
	eventOnRequestFail []*v8.Function        // Speichert die Funktionen ab, welche aufgerufen werden sobald das Event eintritt
}

type SharedFunctionRequest struct {
	resolveChan chan *types.FunctionCallState
	kernel      types.KernelInterface
	mutex       *sync.Mutex
	_destroyed  bool
}

type CallFunctionSignature struct {
	*types.FunctionSignature
}
