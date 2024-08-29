package context

import (
	"github.com/CustodiaJS/custodiajs-core/global/types"
)

func (o *Context) IsConnected() bool {
	return o.isConnected.Bool()
}

func (o *Context) Close() {
	o.proc.Debug("Close request")
}

func (o *Context) Done() {
	o.proc.Debug("Request completed")
}

func (o *Context) GetChildProcessLog(header string) types.ProcessLogSessionInterface {
	return o.proc.GetChildLog(header)
}

func (o *Context) GetProcessLog() types.ProcessLogSessionInterface {
	return o.proc
}
