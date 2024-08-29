package http

import (
	"github.com/CustodiaJS/custodiajs-core/api/http/middleware"
	"github.com/CustodiaJS/custodiajs-core/api/http/middlewares"
)

var (
	// Diese Inspektoren werden bei /rpc Anfragen verwendet
	RpcMiddlewares middleware.MiddlewareFunctionList = middleware.MiddlewareFunctionList{
		middlewares.ForceTLS,
		middlewares.ParseAndPassVmRpcUrlFunctionSignature,
		middlewares.ProxyOrBrowserRequestValidation,
		middlewares.ValidateRPCRequest,
	}

	// Diese Inspektoren werden bei / Anfragen verwendet
	IndexMiddlewares middleware.MiddlewareFunctionList = middleware.MiddlewareFunctionList{
		middlewares.ForceTLS,
	}

	// Diese Inspektoren werden bei /vm Anfragen verwendet
	VmMiddlewares middleware.MiddlewareFunctionList = middleware.MiddlewareFunctionList{}

	// Diese Inspektoren werden bei /ws anfragen verwendet
	//WebConsoleMiddleware MiddlewareFunctionList = MiddlewareFunctionList{IsWebsocketRequest}
)
