package types

const (
	// Gibt den Status einer VM an
	Closed    VmState = 1
	Running   VmState = 2
	Starting  VmState = 3
	StillWait VmState = 0

	// Gibt die Notwendigen Nutzerrechte an
	NONE_ROOT_ADMIN           CLIUserRight = 1
	ROOT_ADMIN                CLIUserRight = 2
	CONTAINER_NONE_ROOT_ADMIN CLIUserRight = 3
	CONAINER_ROOT_ADMIN       CLIUserRight = 4

	// Gibt den Status des Core Osbjektes an
	NEW      CoreState = 1
	SERVING  CoreState = 2
	SHUTDOWN CoreState = 3
	CLOSED   CoreState = 4

	// Gibt an ob es sich um CBOR oder JSON Daten handelt
	HTTP_CONTENT_CBOR HttpRequestContentType = 1
	HTTP_CONTENT_JSON HttpRequestContentType = 2

	// Gibt an, mit welchem Protokoll der Funktionsaufruf durchgeführt wurde
	HTTP_JSON RpcCallTransportProtocol = 1
	GRPC      RpcCallTransportProtocol = 2

	// Legt die Aktuelle Version fest
	C_VESION VERSION = 1000000000

	// Die Repo wird festgelegt
	C_REPO REPO = "https://github.com/DGPCorpSoftwares/vnh1"

	// Legt die Dateipfade für z.b Unix Sockets fest
	NONE_ROOT_UNIX_SOCKET SOCKET_PATH = "/tmp/vnh1_none_root_sock"
	ROOT_UNIX_SOCKET      SOCKET_PATH = "/tmp/vnh1_root_sock"
)
