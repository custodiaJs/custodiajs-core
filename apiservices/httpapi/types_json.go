package httpapi

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

type Response struct {
	Version          uint32   `json:"version"`
	RemoteConsole    bool     `json:"remoteconsole"`
	ScriptContainers []string `json:"scriptcontainers"`
}

type SharedFunction struct {
	Name      string   `json:"name"`
	ParmTypes []string `json:"parmtypes"`
}

type SharedFunctions struct {
	Public []SharedFunction `json:"public"`
	Local  []SharedFunction `json:"local"`
}

type vmInfoResponse struct {
	Name            string          `json:"name"`
	Hash            string          `json:"hash"`
	Modules         []string        `json:"modules"`
	State           string          `json:"state"`
	SharedFunctions SharedFunctions `json:"sharedfunctions"`
}
