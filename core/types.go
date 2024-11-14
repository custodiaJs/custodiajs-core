package core

// Allgemeine Typen die nur im Core benötigt werden
type _VmIpcServerState uint8

const (
	NEW     _VmIpcServerState = 1
	INITED  _VmIpcServerState = 2
	SERVING _VmIpcServerState = 3
	CLOSING _VmIpcServerState = 4
	CLOSED  _VmIpcServerState = 5
)

// Gibt einen
