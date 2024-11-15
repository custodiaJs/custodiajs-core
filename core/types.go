package core

import "net"

// Allgemeine Typen die nur im Core benötigt werden
type _VmIpcServerState uint8

const (
	NEW     _VmIpcServerState = 1
	INITED  _VmIpcServerState = 2
	SERVING _VmIpcServerState = 3
	CLOSING _VmIpcServerState = 4
	CLOSED  _VmIpcServerState = 5
)

// Gibt die ACL Regeln für einen Benutzer an
type ACL struct {
	Username  *string
	Groupname *string
}

type _AclListener struct {
	net.Listener
	AclRule *ACL
}
