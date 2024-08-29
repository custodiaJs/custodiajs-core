package procslog

import (
	"fmt"
	"log"
	"strings"

	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/fatih/color"
)

var (
	foldEnd = color.New(color.Bold).SprintfFunc()(":-$")
)

func (o *ProcLogSession) LogPrint(header, format string, value ...interface{}) {
	// Die Eingabe wird formatiert
	userinput := fmt.Sprintf(format, value...)

	// Der Header wird ermittelt
	var foldedHeader string
	if header != "" {
		foldedHeader = color.New(color.Bold).SprintFunc()(header)
	} else {
		foldedHeader = color.New(color.Bold).SprintFunc()(o.header)
	}

	// Der Text wird angezeigt
	logPrint(fmt.Sprintf("[%s] %s%s %s", o.sessionColorFunc(strings.ToUpper(o.id)), foldedHeader, foldEnd, userinput))
}

func (o *ProcLogSession) Log(format string, value ...interface{}) {
	o.LogPrint("", format, value...)
}

func (o *ProcLogSession) Debug(format string, value ...interface{}) {
	o.LogPrint("", format, value...)
}

func (o *ProcLogSession) LogPrintSuccs(format string, value ...interface{}) {
	o.LogPrint("", format, value...)
}

func (o *ProcLogSession) LogPrintError(format string, value ...interface{}) {
	o.LogPrint("", format, value...)
}

func (o *ProcLogSession) GetChildLog(header string) types.ProcessLogSessionInterface {
	return &ProcLogChildSession{mother: o, header: header}
}

func (o *ProcLogSession) GetID() string {
	return o.id
}

func logPrint(text string) {
	log.Print(text)
}
