package webgrpc

import (
	"crypto/tls"
	"crypto/x509"
	"vnh1/grpc/publicgrpc"
	"vnh1/types"

	"google.golang.org/grpc"
)

type LocalAddress struct {
	LocalIP   string
	LocalPort uint32
}

type WebGrpcService struct {
	core         types.CoreInterface
	cert         *x509.Certificate
	localAddress *LocalAddress
	serverObj    *grpc.Server
	tlsConfig    *tls.Config
	isLocalhost  bool
}

type GrpcServer struct {
	publicgrpc.UnsafeRPCServiceServer
	core types.CoreInterface
}

type FunctionParameterCapsle struct {
	Value interface{}
	CType string
}

type RpcRequest struct {
	parms []types.FunctionParameterBundleInterface
}
