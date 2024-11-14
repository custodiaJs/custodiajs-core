package static

import (
	"github.com/CustodiaJS/custodiajs-core/global/types"
)

const (
	// Gibt den Status einer VM an
	Closed    types.VmState = 1
	Running   types.VmState = 2
	Starting  types.VmState = 3
	StillWait types.VmState = 0

	// Gibt die Notwendigen Nutzerrechte an
	NONE_ROOT_ADMIN           types.IPCRight = 1
	ROOT_ADMIN                types.IPCRight = 2
	CONTAINER_NONE_ROOT_ADMIN types.IPCRight = 3
	CONAINER_ROOT_ADMIN       types.IPCRight = 4

	// Gibt den Status des Core Osbjektes an
	NEW      types.CoreState = 1
	INITED   types.CoreState = 2
	SERVING  types.CoreState = 3
	SHUTDOWN types.CoreState = 4
	CLOSED   types.CoreState = 5

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

	// Legt die Verfügbaren RPC VM Identifikations Methoden fest
	RPC_REQUEST_METHODE_VM_IDENT_ID   types.RPCRequestVMIdentificationMethode = 0
	RPC_REQUEST_METHODE_VM_IDENT_NAME types.RPCRequestVMIdentificationMethode = 1

	// Gibt an, dass nicht ermittelt werden konnte, ob es sich um eine Tor IP handelt
	UNKOWN_TOR_IP_STATE types.TorIpState = false

	// Definiert alle HTTP Methoden
	GET       types.HTTP_METHOD = "GET"
	POST      types.HTTP_METHOD = "POST"
	PUT       types.HTTP_METHOD = "PUT"
	DELETE    types.HTTP_METHOD = "DELETE"
	PATCH     types.HTTP_METHOD = "PATCH"
	HEAD      types.HTTP_METHOD = "HEAD"
	OPTIONS   types.HTTP_METHOD = "OPTIONS"
	CONNECT   types.HTTP_METHOD = "CONNECT"
	TRACE     types.HTTP_METHOD = "TRACE"
	WEBSOCKET types.HTTP_METHOD = "WEBSOCKET"

	// Gibt den Key für das Core Objekt an
	CORE_SESSION_CONTEXT_KEY types.CONTEXT_KEY = types.CONTEXT_KEY("CoreSession")
)
