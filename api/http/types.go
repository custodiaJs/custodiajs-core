package http

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"github.com/CustodiaJS/custodiajs-core/api/http/middleware"
	"github.com/CustodiaJS/custodiajs-core/global/procslog"
	"github.com/CustodiaJS/custodiajs-core/global/types"
)

type HttpApiService struct {
	middlewareHandlers []middleware.MiddlewareFunction
	plog               *procslog.ProcLogSession
	core               types.CoreInterface
	cert               *x509.Certificate
	localAddress       *LocalAddress
	serverObj          *http.Server
	serverMux          *http.ServeMux
	tlsConfig          *tls.Config
	isLocalhost        bool
}

type LocalAddress struct {
	LocalIP   string
	LocalPort uint32
}
