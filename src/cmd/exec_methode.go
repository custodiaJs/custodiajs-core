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

package cmd

import (
	"os"
	"os/user"
	"runtime"

	"github.com/custodia-cenv/cenvx-core/src/host"
	"github.com/custodia-cenv/cenvx-core/src/log"
)

// Wird verwendet um zu ermitteln ob es sich um ein unterstützes OS handelt
func OSSupportCheck() {
	// Es wird geprüft ob es sich um Unterstützes OS handelt
	switch runtime.GOOS {
	case "linux":
		if err := host.VerifyLinuxSystem(); err != nil {
			log.InfoLogPrint(err.Error())
			os.Exit(1)
		}
	case "windows":
		if err := host.VerifyWindowsSystem(); err != nil {
			log.InfoLogPrint(err.Error())
			os.Exit(1)
		}
	case "darwin":
		if err := host.VerifyAppleMacOSSystem(); err != nil {
			log.InfoLogPrint(err.Error())
			os.Exit(1)
		}
	case "freebsd", "openbsd", "netbsd":
		if err := host.VerifyBSDSystem(); err != nil {
			log.InfoLogPrint(err.Error())
			os.Exit(1)
		}
	default:
		log.InfoLogPrint("It is an unsupported operating system.")
		os.Exit(1)
	}
}

// Zeigt die Aktuellen Host Informationen an
func PrintHostInformations() {
	// Die Linux Informationen werden ausgelesen
	hostInfo, err := host.DetectLinuxDist()
	if err != nil {
		panic(err)
	}

	// Die Host Informationen werden angezigt
	log.InfoLogPrint("Host OS: %s", hostInfo)

	// Es wird ermittelt ob das Programm in einem Container ausgeführt wird
	isRunningInLinuxContainer := host.IsRunningInContainer()

	// Die Info wird angezeigt
	if isRunningInLinuxContainer {
		log.InfoLogPrint("Running in container: yes")
	} else {
		log.InfoLogPrint("Running in container: no")
	}
}

// Wird verwendet um zu ermittelt ob der Host derzeit im Root Modus ausgeführt wird
func IsRunningAsRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		return false
	}
	uid := currentUser.Uid
	return uid == "0"
}
