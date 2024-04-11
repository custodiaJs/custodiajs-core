package types

type VmState uint8
type CoreState uint8
type CLIUserRight uint8
type RpcCallTransportProtocol uint8
type HttpRequestContentType uint8
type VERSION uint32
type REPO string
type SOCKET_PATH string
type ALTERNATIVE_SERVICE_PATH string
type CoreVMFingerprint string

type TransportWhitelistVmEntryData struct {
	WildCardDomains []string
	ExactDomains    []string
	Methods         []string
	IPv4List        []string
	Ipv6List        []string
}

type CAMemberData struct {
	Fingerprint string
	Type        string
	ID          string
}

type VMDatabaseData struct {
	Type     string
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Alias    string
}
