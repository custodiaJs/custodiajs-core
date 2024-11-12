// Author: fluffelpuff
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package types

import (
	"crypto/tls"
	"net/http"
	"net/url"
)

// VM und Core Status Typen sowie Repo Datentypen
type ALTERNATIVE_SERVICE_PATH string        // Alternativer Socket Path
type VmState uint8                          // VM Status
type CoreState uint8                        // Core Status
type IPCRight uint8                         // CLI Benutzerrecht
type VERSION uint32                         // Version des Hauptpgrogrammes
type REPO string                            // URL der Sourccode Qeulle
type SOCKET_PATH string                     // Gibt einen Socket Path an
type LOG_DIR string                         // Gibt den Path des Log Dir's unter
type HOST_CRYPTOSTORE_WATCH_DIR_PATH string // Gibt den Ordner an, in dem sich alle Zertifikate und Schlüssel des Hosts befinden
type HOST_CONFIG_FILE_PATH string           // Gibt den Pfad der Config Datei an
type HOST_CONFIG_PATH string
type CNH_FIRMWARE_PATH string
type CHN_CORE_SOCKET_PATH string
type CONTEXT_KEY string
type HTTP_METHOD string

// Gibt die QUID einer VM an
type QVMID string

// Gibt den Hash eines Scriptes an
type VmScriptHash string

// Gibt die ProcessId an
type ProcessId string

// Gibt die VmID an
type VmId string

// RPC Transport & Call Typen
type RpcCallTransportProtocol uint8 // RPC Transport Protokoll
type HttpRequestContentType uint8   // HTTP Request Content Type

// ID Typen
type KernelID string          // Gibt die ID eines Kernels an
type KernelFingerprint string // Gibt die Kernel VM-ID an
type RPCCallSource uint8      // Gibt an ob es sich um eine Lokale Anfrage oder eine Remote Anfrage handelt

// RPC Request Methode
type RpcRequestMethode uint8 // Gibt an, über welche Methode der RPC Request Empfangen wurde

// Kernel Loop Operation Methode
type KernelEventLoopOperationMethode uint8

// Vererbte Structs
type FunctionCallReturnData ExportedV8Value

// Gibt die Indentifizierungsmethode an
type RPCRequestVMIdentificationMethode uint8

// Gibt an ob es sich bei einer IP um eine TOR IP-handelt
type TorIpState bool

type TransportWhitelistVmEntryData struct {
	WildCardDomains []string
	ExactDomains    []string
	Methods         []string
	IPv4List        []string
	Ipv6List        []string
}

type FunctionCallState struct {
	State  string
	Error  string
	Return []*FunctionCallReturnData
}

type FunctionCallReturn struct {
	*FunctionCallState
	Resolve func()
	Reject  func()
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

type HttpRpcRequestUserData struct {
	Username string
	Password string
}

type HttpContext struct {
	IsConnected      func() bool
	ContentLength    int64
	PostForm         url.Values
	Header           http.Header
	Host             string
	Form             url.Values
	Proto            string
	RemoteAddr       string
	RequestURI       string
	TLS              *tls.ConnectionState
	TransferEncoding []string
	URL              *url.URL
	Cookies          []*http.Cookie
	BasicAuth        *HttpRpcRequestUserData
	UserAgent        string
}

type RpcRequest struct {
	RequestType   RpcRequestMethode
	ProcessLog    ProcessLogSessionInterface
	Parms         []*FunctionParameterCapsle
	Request       RpcRequestInterface
	WriteResponse func(*FunctionCallReturn) error
	Context       CoreContextInterface
}

type RPCParmeterReadingError struct {
	Pos       int
	Has       string
	Need      string
	SpeficMsg string
}

type XRequestedWithData struct {
	IsKnown bool
	Value   string
}

type HttpRpcResponseCapsle struct {
	Data  []*RPCResponseData
	Error string
}

type HostKeyCert struct {
	Algorithm string
	FilePath  string
	Alias     string
}

type VmInstanceProcessParameters struct {
	VMImageFilePath   string
	VMWorkingDir      string
	HostKeyCerts      []HostKeyCert
	DisableCoreCrypto bool
}

type VmInstanceInstanceConfig struct {
	VmProcessParameters *VmInstanceProcessParameters
}
