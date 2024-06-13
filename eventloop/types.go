package eventloop

import (
	"sync"
	"vnh1/types"

	v8 "rogchap.com/v8go"
)

type KernelEventLoopOperation struct {
	Type             types.KernelEventLoopOperationMethode
	DirectV8Function func(*v8.Context, types.KernelEventLoopContextInterface)
	EngineSourceCode string
	_hasNullReturn   bool
	_cond            *sync.Cond
	_mutex           *sync.Mutex
	_returnError     error
	_returnValue     *v8.Value
}

type KernelEventLoopContext struct {
	setError  func(error)
	setResult func(*v8.Value)
}
