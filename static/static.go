package static

import "github.com/CustodiaJS/custodiajs-core/types"

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
	C_REPO types.REPO = "https://github.com/CustodiaJS/custodiajs-core"

	// Legt die Dateipfade für z.b Unix Sockets fest
	NONE_ROOT_UNIX_SOCKET            types.SOCKET_PATH              = "/tmp/cusjs_none_root_sock"
	ROOT_UNIX_SOCKET                 types.SOCKET_PATH              = "/tmp/cusjs_root_sock"
	UNIX_ALTERNATIVE_SERVICES        types.ALTERNATIVE_SERVICE_PATH = "/var/lib/cusjs/alternativeservices"
	UNIX_LINUX_LOGGING_DIR           types.LOG_DIR                  = "/var/log/cusjs"
	UNIX_LINUX_LOGGING_DIR_NONE_ROOT types.LOG_DIR                  = "/tmp/cusjs"

	// Gibt die Verfügabren Quellen eines Funktionsaufrufes an
	LOCAL  types.RPCCallSource = 0
	REMOTE types.RPCCallSource = 1

	// Gibt die Verfügabren Typen für einen RPC Request an
	HTTP_REQUEST      types.RpcRequestMethode = 0
	IPC_REQUEST       types.RpcRequestMethode = 1
	WEBSOCKET_REQUEST types.RpcRequestMethode = 2

	// Gibt die Verfügabren EventLoop Methoden an
	KERNEL_EVENT_LOOP_FUNCTION    types.KernelEventLoopOperationMethode = 0
	KERNEL_EVENT_LOOP_SOURCE_CODE types.KernelEventLoopOperationMethode = 1
)
