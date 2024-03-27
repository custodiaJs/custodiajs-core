package jsvm

import (
	"strings"

	"github.com/dop251/goja"
)

func console_base(runtime *goja.Runtime, call goja.FunctionCall, vm *JsVM) goja.Value {
	_ = call
	return runtime.ToValue(func(parms goja.FunctionCall) goja.Value {
		var args []string
		for _, arg := range parms.Arguments[1:] {
			args = append(args, arg.String())
		}
		output := strings.Join(args, " ")

		switch parms.Arguments[1].String() {
		case "info":
			vm.consoleCache.InfoLog(output)
		case "error":
			vm.consoleCache.ErrorLog(output)
		default:
			vm.consoleCache.Log(output)
		}

		return goja.Undefined()
	})
}
