//go:build linux

package container

import "fmt"

func NewLinuxContainer() (*VmContainer, error) {
	return &VmContainer{}, nil
}

func NewWindowsDockerWSLContainer() (*VmContainer, error) {
	return nil, fmt.Errorf("NewWindowsDockerWSLContainer: Only allowed on Windows")
}

func NewMacOSContainer() (*VmContainer, error) {
	return nil, fmt.Errorf("NewWindowsDockerWSLContainer: Only allowed on macOS")
}

func NewFreeBSDContainer() (*VmContainer, error) {
	return nil, fmt.Errorf("NewWindowsDockerWSLContainer: Only allowed on FreeBSD")
}

func NewOpenBSDContainer() (*VmContainer, error) {
	return nil, fmt.Errorf("NewWindowsDockerWSLContainer: Only allowed on OpenBSD")
}

func NewNetBSDContainer() (*VmContainer, error) {
	return nil, fmt.Errorf("NewWindowsDockerWSLContainer: Only allowed on NetBSD")
}

func CheckWindowsHasDockerOrWSL() bool {
	return false
}
