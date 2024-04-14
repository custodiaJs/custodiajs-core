package kernel

import (
	v8 "rogchap.com/v8go"
)

func (o *Kernel) _kernel_throw(context *v8.Context, msg string) {
	errMsg, err := v8.NewValue(o.Isolate(), msg)
	if err != nil {
		panic(err)
	}
	context.Isolate().ThrowException(errMsg)
}
