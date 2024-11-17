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

package filesystem

import (
	"fmt"
	"path"

	cenvxcore "github.com/custodia-cenv/cenvx-core/src"
)

func MakeLogDirForVM(rootdir cenvxcore.LOG_DIR, vmName string) (cenvxcore.LOG_DIR, error) {
	np := path.Join(string(rootdir), vmName)
	if !FolderExists(np) {
		if err := CreateDirectory(np); err != nil {
			return "", fmt.Errorf("MakeLogDirForVM: " + err.Error())
		}
	}
	return cenvxcore.LOG_DIR(np), nil
}
