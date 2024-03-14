package jsvm

import (
	"fmt"
	"log"
	"strings"

	"github.com/dop251/goja"
)

type consoleMethode int
type WebscoketMethode int

const (
	console_log   consoleMethode = 1
	console_info  consoleMethode = 2
	console_error consoleMethode = 3

	ws_send  WebscoketMethode = 4
	ws_open  WebscoketMethode = 5
	ws_close WebscoketMethode = 6
	ws_error WebscoketMethode = 7
)

func console_base(methode consoleMethode, runtime *goja.Runtime, call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 1 {
		fmt.Println(call.Arguments)
		return runtime.ToValue("error")
	}

	return runtime.ToValue(func(parms goja.FunctionCall) goja.Value {
		var args []string
		for _, arg := range parms.Arguments {
			args = append(args, arg.String())
		}
		output := strings.Join(args, " ")

		switch methode {
		case console_info:
			log.Println(output)
		case console_error:
			log.Fatalln(output)
		default:
			log.Println(output)
		}

		return goja.Undefined()
	})
}

func sharefunction_base(runtime *goja.Runtime, call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 1 {
		fmt.Println(call.Arguments)
		return runtime.ToValue("error")
	}

	return runtime.ToValue(func(parms goja.FunctionCall) goja.Value {
		fmt.Println(parms.Arguments[0].Export())
		fmt.Println(parms.Arguments[1].Export())
		return goja.Undefined()
	})
}

func websocket_base(methode WebscoketMethode, runtime *goja.Runtime, call goja.FunctionCall) goja.Value {
	return nil
}

func (o *JsVM) vnhCOMAPIEP(call goja.FunctionCall) goja.Value {
	// Es wird ermittelt um welchen vorgang es sich handelt
	if len(call.Arguments) < 1 {
		return o.gojaVM.ToValue("invalid")
	}

	// Die jeweilige Funktion wird ermittelt
	switch call.Arguments[0].String() {
	// Konsolen funktionen
	case "console/log":
		return console_base(console_log, o.gojaVM, call)
	case "console/info":
		return console_base(console_info, o.gojaVM, call)
	case "console/error":
		return console_base(console_error, o.gojaVM, call)
	// Websocket
	case "webscoket/send":
		return websocket_base(ws_send, o.gojaVM, call)
	case "webscoket/open":
		return websocket_base(ws_open, o.gojaVM, call)
	case "webscoket/close":
		return websocket_base(ws_close, o.gojaVM, call)
	case "webscoket/error":
		return websocket_base(ws_error, o.gojaVM, call)
	// Share Functions
	case "root/sharefunction":
		return sharefunction_base(o.gojaVM, call)
	default:
		return o.gojaVM.ToValue("invalid")
	}
}
