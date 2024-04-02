package types

type VmState int
type CLIUserRight int

const (
	Closed    VmState = 1
	Running   VmState = 2
	Starting  VmState = 3
	StillWait VmState = 0

	NONE_ROOT_ADMIN           CLIUserRight = 1
	ROOT_ADMIN                CLIUserRight = 2
	CONTAINER_NONE_ROOT_ADMIN CLIUserRight = 3
	CONAINER_ROOT_ADMIN       CLIUserRight = 4
)
