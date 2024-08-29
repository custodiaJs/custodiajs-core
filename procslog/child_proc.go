package procslog

import (
	"github.com/CustodiaJS/custodiajs-core/types"
)

func (o *ProcLogChildSession) LogPrint(header, format string, value ...interface{}) {
	if header != "" {
		o.mother.LogPrint(header, format, value...)
	} else {
		o.mother.LogPrint(o.header, format, value...)
	}
}

func (o *ProcLogChildSession) Log(format string, value ...interface{}) {
	o.LogPrint(o.header, format, value...)
}

func (o *ProcLogChildSession) Debug(format string, value ...interface{}) {
	o.LogPrint(o.header, format, value...)
}

func (o *ProcLogChildSession) LogPrintSuccs(format string, value ...interface{}) {
	o.LogPrint(o.header, format, value...)
}

func (o *ProcLogChildSession) LogPrintError(format string, value ...interface{}) {
	o.LogPrint(o.header, format, value...)
}

func (o *ProcLogChildSession) GetChildLog(header string) types.ProcessLogSessionInterface {
	return &ProcLogChildSession{mother: o.mother, header: header}
}

func (o *ProcLogChildSession) GetID() string {
	return o.mother.id
}
