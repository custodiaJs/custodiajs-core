package filesystem

import (
	"fmt"
	"path"
	"vnh1/types"
)

func MakeLogDirForVM(rootdir types.LOG_DIR, vmName string) (types.LOG_DIR, error) {
	np := path.Join(string(rootdir), vmName)
	if !FolderExists(np) {
		if err := CreateDirectory(np); err != nil {
			return "", fmt.Errorf("MakeLogDirForVM: " + err.Error())
		}
	}
	return types.LOG_DIR(np), nil
}
