package types

type VmState uint8
type CoreState uint8
type CLIUserRight uint8
type RpcCallTransportProtocol uint8
type HttpRequestContentType uint8

const (
	Closed    VmState = 1
	Running   VmState = 2
	Starting  VmState = 3
	StillWait VmState = 0

	NONE_ROOT_ADMIN           CLIUserRight = 1
	ROOT_ADMIN                CLIUserRight = 2
	CONTAINER_NONE_ROOT_ADMIN CLIUserRight = 3
	CONAINER_ROOT_ADMIN       CLIUserRight = 4

	NEW      CoreState = 1
	SERVING  CoreState = 2
	SHUTDOWN CoreState = 3
	CLOSED   CoreState = 4

	HTTP_CONTENT_CBOR HttpRequestContentType = 1
	HTTP_CONTENT_JSON HttpRequestContentType = 2

	HTTP_JSON RpcCallTransportProtocol = 1
	GRPC      RpcCallTransportProtocol = 2
)
