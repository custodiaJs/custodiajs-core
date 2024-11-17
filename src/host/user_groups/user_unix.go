//go:build !windows
// +build !windows

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
