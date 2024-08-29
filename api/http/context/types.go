package context

import (
	"crypto/x509"
	"net"
	"net/url"

	"github.com/CustodiaJS/custodiajs-core/api/http/context/saftychan"
	"github.com/CustodiaJS/custodiajs-core/global/procslog"
	"github.com/CustodiaJS/custodiajs-core/global/types"
	"github.com/CustodiaJS/custodiajs-core/global/utils/grsbool"
)

type Context struct {
	isConnected        *grsbool.Grsbool
	proc               *procslog.ProcLogSession
	method             types.HTTP_METHOD
	contentType        types.HttpRequestContentType
	xRequestedWithData *types.XRequestedWithData
	refererURL         *url.URL
	originURL          *url.URL
	tlsCert            []*x509.Certificate
	fncs               *types.FunctionSignature
}

type HttpContext struct {
	*Context
	saftyResponseChan *saftychan.FunctionCallReturnChan
	localIp           net.IP
	remoteIp          net.IP
}

type ContextManagmentUnit struct {
	core types.CoreInterface
}
