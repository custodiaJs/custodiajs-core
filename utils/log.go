package utils

import (
	"fmt"
	"log"
	"path"
	"vnh1/types"
)

func LogPrint(text string) {
	log.Print(text)
}

func MakeLogDirForVM(rootdir types.LOG_DIR, vmName string) (types.LOG_DIR, error) {
	np := path.Join(string(rootdir), vmName)
	if !FolderExists(np) {
		if err := CreateDirectory(np); err != nil {
			return "", fmt.Errorf("MakeLogDirForVM: " + err.Error())
		}
	}
	return types.LOG_DIR(np), nil
}
