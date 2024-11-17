package static

import (
	cenvxcore "github.com/custodia-cenv/cenvx-core/src"
)

const (
	// Gibt den Status einer VM an
	Closed    cenvxcore.VmState = 1
	Running   cenvxcore.VmState = 2
	Starting  cenvxcore.VmState = 3
	StillWait cenvxcore.VmState = 0

	// Gibt die Notwendigen Nutzerrechte an
	NONE_ROOT_ADMIN           cenvxcore.IPCRight = 1
	ROOT_ADMIN                cenvxcore.IPCRight = 2
	CONTAINER_NONE_ROOT_ADMIN cenvxcore.IPCRight = 3
	CONAINER_ROOT_ADMIN       cenvxcore.IPCRight = 4

	// Gibt den Status des Core Osbjektes an
	NEW      cenvxcore.CoreState = 1
	INITED   cenvxcore.CoreState = 2
	SERVING  cenvxcore.CoreState = 3
	SHUTDOWN cenvxcore.CoreState = 4
	CLOSED   cenvxcore.CoreState = 5

	// Legt die Aktuelle Version fest
	C_VESION cenvxcore.VERSION = 1000000000

	// Die Repo wird festgelegt
	C_REPO cenvxcore.REPO = "https://github.com/custodia-cenv/cenvx-core"

	// Gibt an, dass nicht ermittelt werden konnte, ob es sich um eine Tor IP handelt
	UNKOWN_TOR_IP_STATE cenvxcore.TorIpState = false
)
