package procslog

import (
	"fmt"
	"log"
	"strings"

	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/utils"
	"github.com/fatih/color"
)

var (
	formatGreen  = color.New(color.FgGreen).SprintFunc()
	formatBold   = color.New(color.Bold).SprintFunc()
	foldEnd      = color.New(color.Bold).SprintfFunc()(":-$")
	foldedOpen   = color.New(color.Bold).SprintFunc()("(")
	foldedClosed = color.New(color.Bold).SprintFunc()(")")
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
	LogPrint(fmt.Sprintf("%s%s%s%s%s %s", foldedHeader, foldedOpen, o.sessionColorFunc(strings.ToUpper(o.id)), foldedClosed, foldEnd, userinput))
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

func LogPrint(text string) {
	log.Print(text)
}

func NewProcLogSession() *ProcLogSession {
	randHex, _ := utils.RandomHex(4)
	sessioncollorFormater := utils.DetermineColorFromHex(randHex)
	val := &ProcLogSession{id: randHex, sessionColorFunc: sessioncollorFormater}
	val.Log(formatBold(formatGreen("New Process Log Session created")))
	return val
}

func NewProcLogSessionWithHeader(header string) *ProcLogSession {
	randHex, _ := utils.RandomHex(4)
	sessioncollorFormater := utils.DetermineColorFromHex(randHex)
	val := &ProcLogSession{id: randHex, header: header, sessionColorFunc: sessioncollorFormater}
	val.Log(formatBold(formatGreen("New Process Log Session created")))
	return val
}

func NewProcLogForCore() *ProcLogSession {
	randHex, _ := utils.RandomHex(4)
	sessioncollorFormater := utils.DetermineColorFromHex(randHex)
	val := &ProcLogSession{id: randHex, header: "Core", sessionColorFunc: sessioncollorFormater}
	fmt.Printf("New Core Process Log Session created '%s'\n", randHex)
	return val
}
