package kmodulerpc

import (
	"vnh1/types"

	v8 "rogchap.com/v8go"
)

type SharedFunction struct {
	v8VM         *v8.Context
	callFunction *v8.Function
	name         string
	parmTypes    []string
	returnType   string
}

type SharedLocalFunction struct {
	*SharedFunction
}

type SharedPublicFunction struct {
	*SharedFunction
}

type SharedFunctionRequest struct {
	resolveChan chan *types.FunctionCallState
	parms       types.RpcRequestInterface
}
