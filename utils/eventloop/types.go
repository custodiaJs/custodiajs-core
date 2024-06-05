package eventloop

import (
	"sync"
	"vnh1/types"

	"rogchap.com/v8go"
)

type KernelEventLoopOperation struct {
	Type             types.KernelEventLoopOperationMethode
	DirectV8Function types.KernelLoopV8Function
	EngineSourceCode string
	_hasNullReturn   bool
	_cond            *sync.Cond
	_mutex           *sync.Mutex
	_returnError     error
	_returnValue     *v8go.Value
}
