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
		return rpcreq.HttpRequest.IsConnected.Bool()
	default:
		return false
	}
}

func WaitOfConnectionStateChange(rpcreq *types.RpcRequest, cvalue bool) {
	switch {
	case IsHttpRequest(rpcreq):
		rpcreq.HttpRequest.IsConnected.WaitOfChange(cvalue)
	default:
		return
	}
}

func IsRemoteConnection(rpcreq *types.RpcRequest) bool {
	return false
}
