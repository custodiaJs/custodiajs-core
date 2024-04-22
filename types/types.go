package types

// VM und Core Status Typen sowie Repo Datentypen
type ALTERNATIVE_SERVICE_PATH string // Alternativer Socket Path
type VmState uint8                   // VM Status
type CoreState uint8                 // Core Status
type CLIUserRight uint8              // CLI Benutzerrecht
type VERSION uint32                  // Version des Hauptpgrogrammes
type REPO string                     // URL der Sourccode Qeulle
type SOCKET_PATH string              // Gibt einen Socket Path an
type LOG_DIR string                  // Gibt den Path des Log Dir's unter

// RPC Transport & Call Typen
type RpcCallTransportProtocol uint8 // RPC Transport Protokoll
type HttpRequestContentType uint8   // HTTP Request Content Type

// ID Typen
type KernelID string                     // Gibt die ID eines Kernels an
type KernelFingerprint string            // Gibt die Kernel VM-ID an
type CoreVMFingerprint KernelFingerprint // Gibt die ID einer CoreVM zurück
type RPCCallSource uint8                 // Gibt an ob es sich um eine Lokale Anfrage oder eine Remote Anfrage handelt

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

type FunctionCallReturnData struct {
	CType string
	Value interface{}
}

type FunctionCallState struct {
	State  string
	Error  string
	Return []*FunctionCallReturnData
}

type ExportedV8Value struct {
	Type  string
	Value interface{}
}

type FunctionSignature struct {
	VMID         string
	VMName       string
	FunctionName string
	Params       []string
	ReturnType   string
}

type FunctionParameterCapsle struct {
	Value interface{}
	CType string
}

type RpcRequest struct {
	Parms      []*FunctionParameterCapsle
	RpcRequest HttpJsonRequestData
}
