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
	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/utils/grsbool"

	v8 "rogchap.com/v8go"
)

type FunctionType string

const (
	Local  FunctionType = "local"
	Public FunctionType = "public"
)

type SharedFunction struct {
	kernel     types.KernelInterface    // Speichert den Verwendeten Kernel ab
	v8Function *v8.Function             // Speichert die eigentliche Funktion des V8 Codes ab
	name       string                   // Speichert den Namen der Funktion ab
	signature  *types.FunctionSignature // Speichert die Signatur der Fuktion ab
}

type CallFunctionSignature struct {
	*types.FunctionSignature
}

type RequestResponseWaiter struct {
	CallState *types.FunctionCallState
	Error     error
}

type SharedFunctionRequestContext struct {
	kernel          types.KernelInterface
	_rprequest      *types.RpcRequest
	_returnDataType string
	_wasResponded   *grsbool.Grsbool
	_destroyed      *grsbool.Grsbool
}
