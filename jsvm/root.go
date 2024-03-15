package jsvm

import (
	"strings"

	"github.com/dop251/goja"
)

func root_mavail(moduleName string, runtime *goja.Runtime, vm *JsVM) goja.Value {
	switch moduleName {
	case "function_share":
		return runtime.ToValue(vm.config.EnableFunctionSharing)
	case "s3":
		return runtime.ToValue(vm.config.EnableS3)
	default:
		return runtime.ToValue(false)
	}
}

func root_fshare(methode string, functionName string, function goja.Value, runtime *goja.Runtime, vm *JsVM) goja.Value {
	// Die JS Funktion wird geprüft und abgespeichert
	jsFunc, ok := goja.AssertFunction(function)
	if !ok {
		panic(runtime.NewTypeError("Zweites Argument ist keine Funktion"))
	}

	// Es wird geprüft ob eine Funktion mit dem Namen bereits geteilt wird
	if vm.functionIsSharing(functionName) {
		// Der Vorgang wird abgebrochen
		panic(runtime.NewTypeError("Zweites Argument ist keine Funktion"))
	}

	// Die Funktion wird geteilt
	switch strings.ToLower(methode) {
	case "local":
		if err := vm.sharLocalFunction(functionName, jsFunc); err != nil {
			panic(runtime.NewTypeError(err.Error()))
		}
	case "public":
		if err := vm.sharePublicFunction(functionName, jsFunc); err != nil {
			panic(runtime.NewTypeError(err.Error()))
		}
	default:
		panic(runtime.NewTypeError("Zweites: " + methode))
	}

	// Es wird ein Undefined zurückgegeben
	return goja.Undefined()
}

func root_base(runtime *goja.Runtime, call goja.FunctionCall, vm *JsVM) goja.Value {
	_ = call
	return runtime.ToValue(func(parms goja.FunctionCall) goja.Value {
		switch parms.Arguments[0].String() {
		case "mavail":
			return root_mavail(parms.Arguments[1].String(), runtime, vm)
		case "fshare":
			return root_fshare(parms.Arguments[1].String(), parms.Arguments[2].String(), parms.Arguments[3], runtime, vm)
		case "finsh":
			vm.gojaVM.Set("vnh1", goja.Undefined())
			return runtime.ToValue(true)
		default:
			return goja.Undefined()
		}
	})
}
