package eventloop

import v8 "rogchap.com/v8go"

func (o *KernelEventLoopContext) SetError(err error) {
	o.setError(err)
}

func (o *KernelEventLoopContext) SetResult(val *v8.Value) {
	o.setResult(val)
}
