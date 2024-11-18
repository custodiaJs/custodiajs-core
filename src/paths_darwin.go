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

import (
	"fmt"
	"path/filepath"
)

var (
	// Gibt das Standard Config Verzeichniss an
	CoreGeneralConfigFilePath CoreGeneralConfigPath = CoreGeneralConfigPath("/Library/Application Support/" + ApplicationName + "/core.conf")

	// Gibt den Speichertort des CryptoStores an
	CoreCryptoStoreDirPath CoreCryptoStorePath = CoreCryptoStorePath("/Library/Application Support/" + ApplicationName + "/crypstore")

	// Log Dir
	CoreLoggingDirPath LogDirPath = LogDirPath(filepath.Join("/", "Library", "Logs", string(ApplicationName)))

	// Legt die Dateipfade für z.b Unix Sockets fest
	CoreVmIpcRootSocketPath            CoreVmIpcSocketPath         = CoreVmIpcSocketPath("/var/run/" + string(ApplicationName) + "/vmipc.sock")
	_CoreVmIpcSocketSpeficUser         CoreVmIpcSocketPathTemplate = CoreVmIpcSocketPathTemplate("/tmp/" + ApplicationName + "_u_%s_vmipc.sock")
	_CoreVmIpcSocketSpeficUserGrpup    CoreVmIpcSocketPathTemplate = CoreVmIpcSocketPathTemplate("/tmp/" + ApplicationName + "_g_%s_vmipc.sock")
	_CoreVmIpcSocketSpeficUserAndGroup CoreVmIpcSocketPathTemplate = CoreVmIpcSocketPathTemplate("/tmp/" + ApplicationName + "_ug_%s_%s_vmipc.sock")
)

func GetCoreSpeficSocketUserPath(username string) CoreVmIpcSocketPath {
	return CoreVmIpcSocketPath(fmt.Sprintf(string(_CoreVmIpcSocketSpeficUser), username))
}

func GetCoreSpeficSocketUserGroupPath(groupName string) CoreVmIpcSocketPath {
	return CoreVmIpcSocketPath(fmt.Sprintf(string(_CoreVmIpcSocketSpeficUserGrpup), groupName))
}

func GetCoreSpeficSocketUserAndGroupPath(username string, groupName string) CoreVmIpcSocketPath {
	return CoreVmIpcSocketPath(fmt.Sprintf(string(_CoreVmIpcSocketSpeficUserAndGroup), username, groupName))
}
