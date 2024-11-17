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
