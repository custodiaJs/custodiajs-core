package procslog

import (
	"fmt"
	"strings"
	"vnh1/utils"
)

type ProcLogSession struct {
	id string
}

func (o *ProcLogSession) LogPrint(header, format string, value ...interface{}) {
	userinput := fmt.Sprintf(format, value...)
	utils.LogPrint(fmt.Sprintf("%s(%s):-$ %s", header, strings.ToUpper(o.id), userinput))
}

func (o *ProcLogSession) LogPrintSuccs(format string, value ...interface{}) {
	userinput := fmt.Sprintf(format, value...)
	utils.LogPrint(fmt.Sprintf("%s:-$ %s", o.id, userinput))
}

func (o *ProcLogSession) LogPrintError(format string, value ...interface{}) {
	userinput := fmt.Sprintf(format, value...)
	utils.LogPrint(fmt.Sprintf("%s:-$ %s", o.id, userinput))
}

func (o *ProcLogSession) GetID() string {
	return o.id
}

func NewProcLogSession() *ProcLogSession {
	randHex, _ := utils.RandomHex(4)
	return &ProcLogSession{id: randHex}
}
