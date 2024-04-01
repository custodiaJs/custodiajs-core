package webservice

import (
	"crypto/x509"
	"net"
	"net/http"
	"vnh1/types"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
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
	grpcServer   *grpc.Server
	tcpMux       cmux.CMux
	isLocalhost  bool
}

type LocalAddress struct {
	LocalIP   string
	LocalPort uint32
}
