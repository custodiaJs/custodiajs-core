package utils

import "fmt"

type ProcLogSession struct {
	id string
}

func (o *ProcLogSession) LogPrint(format string, value ...interface{}) {
	userinput := fmt.Sprintf(format, value...)
	fmt.Printf("[LOG]:%s:-$ %s", o.id, userinput)
}

func (o *ProcLogSession) LogPrintSuccs(format string, value ...interface{}) {
	userinput := fmt.Sprintf(format, value...)
	fmt.Printf("[LOG]:%s:-$ %s", o.id, userinput)
}

func (o *ProcLogSession) LogPrintError(format string, value ...interface{}) {
	userinput := fmt.Sprintf(format, value...)
	fmt.Printf("[LOG]:%s:-$%s", o.id, userinput)
}

func (o *ProcLogSession) GetID() string {
	return o.id
}

func NewProcLogSession() *ProcLogSession {
	randHex, _ := RandomHex(4)
	return &ProcLogSession{id: randHex}
}
