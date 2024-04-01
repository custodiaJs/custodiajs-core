package httpapi

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"vnh1/types"
)

type RpcRequest struct {
	parms []types.FunctionParameterBundleInterface
}

type HttpApiService struct {
	core         types.CoreInterface
	cert         *x509.Certificate
	localAddress *LocalAddress
	serverObj    *http.Server
	serverMux    *http.ServeMux
	tlsConfig    *tls.Config
	isLocalhost  bool
}

type LocalAddress struct {
	LocalIP   string
	LocalPort uint32
}

type FunctionParameterCapsle struct {
	Value interface{}
	CType string
}
