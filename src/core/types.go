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

import "net"

// Allgemeine Typen die nur im Core benötigt werden
type _VmIpcServerState uint8

const (
	NEW     _VmIpcServerState = 1
	INITED  _VmIpcServerState = 2
	SERVING _VmIpcServerState = 3
	CLOSING _VmIpcServerState = 4
	CLOSED  _VmIpcServerState = 5
)

// Gibt die ACL Regeln für einen Benutzer an
type ACL struct {
	Username  *string
	Groupname *string
}

type _AclListener struct {
	net.Listener
	AclRule *ACL
}
