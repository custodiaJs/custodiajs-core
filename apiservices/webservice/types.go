package webservice

import (
	"crypto/x509"
	"net"
	"net/http"
	"vnh1/types"

	"github.com/soheilhy/cmux"
)

type RpcRequest struct {
	parms []types.FunctionParameterBundleInterface
}

type Webservice struct {
	core         types.CoreInterface
	cert         *x509.Certificate
	localAddress *LocalAddress
	serverObj    *http.Server
	serverMux    *http.ServeMux
	httpSocket   net.Listener
	grpcSocket   net.Listener
	tcpMux       cmux.CMux
	isLocalhost  bool
}

type LocalAddress struct {
	LocalIP   string
	LocalPort uint32
}

// JSON

type RPCFunctionParameter struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type RPCFunctionCall struct {
	FunctionName string                 `json:"name"`
	Parms        []RPCFunctionParameter `json:"parms"`
}

type RPCResponseData struct {
	DType string      `json:"type"`
	Value interface{} `json:"value"`
}

type RPCResponse struct {
	Result string          `json:"result"`
	Data   RPCResponseData `json:"data"`
	Error  *string         `json:"error"`
}

type Response struct {
	Version          uint32   `json:"version"`
	RemoteConsole    bool     `json:"remoteconsole"`
	ScriptContainers []string `json:"scriptcontainers"`
}

type SharedFunction struct {
	Name      string   `json:"name"`
	ParmTypes []string `json:"parmtypes"`
}

type SharedFunctions struct {
	Public []SharedFunction `json:"public"`
	Local  []SharedFunction `json:"local"`
}

type vmInfoResponse struct {
	Name            string          `json:"name"`
	Hash            string          `json:"hash"`
	Modules         []string        `json:"modules"`
	State           string          `json:"state"`
	SharedFunctions SharedFunctions `json:"sharedfunctions"`
}
