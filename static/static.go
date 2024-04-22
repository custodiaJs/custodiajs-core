package static

import "vnh1/types"

const (
	// Gibt den Status einer VM an
	Closed    types.VmState = 1
	Running   types.VmState = 2
	Starting  types.VmState = 3
	StillWait types.VmState = 0

	// Gibt die Notwendigen Nutzerrechte an
	NONE_ROOT_ADMIN           types.CLIUserRight = 1
	ROOT_ADMIN                types.CLIUserRight = 2
	CONTAINER_NONE_ROOT_ADMIN types.CLIUserRight = 3
	CONAINER_ROOT_ADMIN       types.CLIUserRight = 4

	// Gibt den Status des Core Osbjektes an
	NEW      types.CoreState = 1
	SERVING  types.CoreState = 2
	SHUTDOWN types.CoreState = 3
	CLOSED   types.CoreState = 4

	// Gibt an ob es sich um CBOR oder JSON Daten handelt
	HTTP_CONTENT_CBOR types.HttpRequestContentType = 1
	HTTP_CONTENT_JSON types.HttpRequestContentType = 2

	// Gibt an, mit welchem Protokoll der Funktionsaufruf durchgeführt wurde
	HTTP_JSON types.RpcCallTransportProtocol = 1
	GRPC      types.RpcCallTransportProtocol = 2

	// Legt die Aktuelle Version fest
	C_VESION types.VERSION = 1000000000

	// Die Repo wird festgelegt
	C_REPO types.REPO = "https://github.com/DGPCorpSoftwares/vnh1"

	// Legt die Dateipfade für z.b Unix Sockets fest
	NONE_ROOT_UNIX_SOCKET            types.SOCKET_PATH              = "/tmp/vnh1_none_root_sock"
	ROOT_UNIX_SOCKET                 types.SOCKET_PATH              = "/tmp/vnh1_root_sock"
	UNIX_ALTERNATIVE_SERVICES        types.ALTERNATIVE_SERVICE_PATH = "/var/lib/vnh1/alternativeservices"
	UNIX_LINUX_LOGGING_DIR           types.LOG_DIR                  = "/var/log/vnh1"
	UNIX_LINUX_LOGGING_DIR_NONE_ROOT types.LOG_DIR                  = "/tmp/vnh1"

	// Gibt die Verfügabren Quellen eines Funktionsaufrufes an
	LOCAL  types.RPCCallSource = 0
	REMOTE types.RPCCallSource = 1
)
