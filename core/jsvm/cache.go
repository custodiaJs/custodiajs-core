package jsvm

import (
	"fmt"

	"github.com/dop251/goja"
)

func cache_write(call goja.FunctionCall, vm *JsVM) goja.Value {
	vm.cache[call.Arguments[1].String()] = call.Arguments[2].Export()
	return nil
}

func cache_read(runtime *goja.Runtime, call goja.FunctionCall, vm *JsVM) goja.Value {
	value, found := vm.cache[call.Arguments[1].String()]
	if !found {
		panic(runtime.NewTypeError("Zweites Argument ist keine Funktion"))
	}
	return runtime.ToValue(value)
}

func cache_base(runtime *goja.Runtime, call goja.FunctionCall, vm *JsVM) goja.Value {
	_ = call
	return runtime.ToValue(func(parms goja.FunctionCall) goja.Value {
		switch parms.Arguments[0].String() {
		case "write":
			fmt.Println("VM_CACHE: WRITE")
			return cache_write(parms, vm)
		case "read":
			fmt.Println("VM_CACHE: READ")
			c := cache_read(runtime, parms, vm)
			return c
		default:
			return goja.Undefined()
		}
	})
}
