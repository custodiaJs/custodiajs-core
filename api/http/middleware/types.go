package middleware

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"

	"github.com/CustodiaJS/custodiajs-core/global/types"
)

type MiddlewareFunction func(core types.CoreInterface, w http.ResponseWriter, r *http.Request) *types.SpecificError
type MiddlewareFunctionList []MiddlewareFunction

type httpRequest func(ResponseWriter http.ResponseWriter, req *http.Request)

type RequestData struct {
	VmIdentificationMethode types.RPCRequestVMIdentificationMethode
	TransportProtocol       types.RpcCallTransportProtocol
	ContentType             types.HttpRequestContentType
	Source                  net.Addr
	XRequestedWith          *types.XRequestedWithData
	Referer                 *url.URL
	Origin                  *url.URL
	VmNameOrID              string
	TLS                     *tls.ConnectionState
	Cookies                 []*http.Cookie
	Headers                 map[string][]string
}

type HttpResponseCapsle struct {
	Data  []*types.RPCResponseData
	Error string
}
