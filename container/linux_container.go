package container

func NewLinuxContainer() (*VmContainer, error) {
	return &VmContainer{}, nil
}

func NewWindowsDockerWSLContainer() (*VmContainer, error) {
	return nil, nil
}

func NewMacOSContainer() (*VmContainer, error) {
	return nil, nil
}

func NewFreeBSDContainer() (*VmContainer, error) {
	return nil, nil
}

func NewOpenBSDContainer() (*VmContainer, error) {
	return nil, nil
}

func NewNetBSDContainer() (*VmContainer, error) {
	return nil, nil
}

func CheckWindowsHasDockerOrWSL() bool {
	return false
}
