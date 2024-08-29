package eventloop

import (
	"sync"

	"github.com/CustodiaJS/custodiajs-core/global/static"
	"github.com/CustodiaJS/custodiajs-core/global/types"

	v8 "rogchap.com/v8go"
)

func NewKernelEventLoopFunctionOperation(funct func(*v8.Context, types.KernelEventLoopContextInterface)) *KernelEventLoopOperation {
	mutex := &sync.Mutex{}
	return &KernelEventLoopOperation{DirectV8Function: funct, Type: static.KERNEL_EVENT_LOOP_FUNCTION, _cond: sync.NewCond(mutex), _mutex: mutex, _hasNullReturn: false}
}

func NewKernelEventLoopSourceOperation(source string) *KernelEventLoopOperation {
	mutex := &sync.Mutex{}
	return &KernelEventLoopOperation{EngineSourceCode: source, Type: static.KERNEL_EVENT_LOOP_SOURCE_CODE, _cond: sync.NewCond(mutex), _mutex: mutex, _hasNullReturn: false}
}
