package types

type VmState int

const (
	Closed    VmState = 1
	Running   VmState = 2
	Starting  VmState = 3
	StillWait VmState = 0
)
