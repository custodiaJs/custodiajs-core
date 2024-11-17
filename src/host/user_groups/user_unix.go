//go:build !windows
// +build !windows

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

package usergroups

import (
	"os/user"
)

// ... Gibt an ob es sich auf Unix Systemen um einen Root Benutzer,
// oder auf WindowsNT Systemen um einen Administrator handelt
func UserHasPrivilegedSystemRights() bool {
	// Der Aktuelle Benutzer wird ermittelt
	currentUser, err := user.Current()
	if err != nil {
		return false
	}
	return currentUser.Uid == "0"
}
