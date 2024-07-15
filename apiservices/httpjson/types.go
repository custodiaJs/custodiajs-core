package httpjson

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"

	"github.com/CustodiaJS/custodiajs-core/types"
)

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

type RequestData struct {
	VmIdentificationMethode types.RPCRequestVMIdentificationMethode
	TransportProtocol       types.RpcCallTransportProtocol
	ContentType             types.HttpRequestContentType
	Source                  types.VerifiedCoreIPAddressInterface
	XRequestedWith          *types.XRequestedWithData
	Referer                 *url.URL
	Origin                  *url.URL
	VmNameOrID              string
	TLS                     *tls.ConnectionState
	Cookies                 []*http.Cookie
}

type ResponseCapsle struct {
	Data  []*RPCResponseData
	Error string
}
