package httpjson

type RPCFunctionParameter struct {
	Type  string      `json:"type" cbor:"type"`
	Value interface{} `json:"value" cbor:"value"`
}

type RPCFunctionCall struct {
	FunctionName string                 `json:"name" cbor:"name"`
	Parms        []RPCFunctionParameter `json:"parms" cbor:"parms"`
}

type RPCResponseData struct {
	DType string      `json:"type" cbor:"type"`
	Value interface{} `json:"value" cbor:"value"`
}

type RPCResponse struct {
	Result string           `json:"result" cbor:"result"`
	Data   *RPCResponseData `json:"data,omitempty" cbor:"data,omitempty"`
	Error  *string          `json:"error" cbor:"error"`
}

type Response struct {
	Version          uint32   `json:"version" cbor:"version"`
	RemoteConsole    bool     `json:"remoteconsole" cbor:"remoteconsole"`
	ScriptContainers []string `json:"scriptcontainers" cbor:"scriptcontainers"`
}

type SharedFunction struct {
	Name      string   `json:"name" cbor:"name"`
	ParmTypes []string `json:"parmtypes" cbor:"parmtypes"`
}

type SharedFunctions struct {
	Public []SharedFunction `json:"public" cbor:"public"`
	Local  []SharedFunction `json:"local" cbor:"local"`
}

type vmInfoResponse struct {
	Name            string          `json:"name" cbor:"name"`
	Id              string          `json:"hash" cbor:"hash"`
	Modules         []string        `json:"modules" cbor:"modules"`
	State           string          `json:"state" cbor:"state"`
	SharedFunctions SharedFunctions `json:"sharedfunctions" cbor:"sharedfunctions"`
}
