package rpcrequest

import "github.com/CustodiaJS/custodiajs-core/types"

func IsHttpRequest(rpcreq *types.RpcRequest) bool {
	if rpcreq == nil {
		return false
	}
	if rpcreq.HttpRequest == nil {
		return false
	}
	return true
}

func ConnectionIsOpen(rpcreq *types.RpcRequest) bool {
	switch {
	case IsHttpRequest(rpcreq):
		return rpcreq.HttpRequest.IsConnected()
	default:
		return false
	}
}

func IsRemoteConnection(rpcreq *types.RpcRequest) bool {
	return false
}
