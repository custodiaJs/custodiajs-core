package kernel

import (
	"sync"
	"vnh1/core/consolecache"

	v8 "rogchap.com/v8go"
)

type Kernel struct {
	*v8.Context
	mutex                 *sync.Mutex
	Console               *consolecache.ConsoleOutputCache
	sharedLocalFunctions  map[string]*SharedLocalFunction
	sharedPublicFunctions map[string]*SharedPublicFunction
}

type SharedLocalFunction struct {
	v8VM         *v8.Context
	callFunction *v8.Function
	name         string
	parmTypes    []string
	returnType   string
}

type SharedPublicFunction struct {
	v8VM         *v8.Context
	callFunction *v8.Function
	name         string
	parmTypes    []string
}

type FunctionCallReturn struct {
	CType string
	Value interface{}
}
