package utils

import v8 "rogchap.com/v8go"

func V8ContextThrow(context *v8.Context, msg string) {
	errMsg, err := v8.NewValue(context.Isolate(), msg)
	if err != nil {
		panic(err)
	}
	context.Isolate().ThrowException(errMsg)
}
