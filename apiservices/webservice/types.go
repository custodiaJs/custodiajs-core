package webservice

import "vnh1/types"

type RPCFunctionParameter struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type RPCFunctionCall struct {
	FunctionName string                 `json:"name"`
	Parms        []RPCFunctionParameter `json:"parms"`
}

type RPCResponseData struct {
	DType string      `json:"type"`
	Value interface{} `json:"value"`
}

type RPCResponse struct {
	Result string          `json:"result"`
	Data   RPCResponseData `json:"data"`
	Error  *string         `json:"error"`
}

type RpcRequest struct {
	parms []types.FunctionParameterBundleInterface
}

type Webservice struct {
	core types.CoreInterface
}

type Response struct {
	Version          uint32   `json:"version"`
	RemoteConsole    bool     `json:"remoteconsole"`
	ScriptContainers []string `json:"scriptcontainers"`
}
