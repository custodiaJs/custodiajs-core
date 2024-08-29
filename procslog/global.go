package procslog

import (
	"github.com/CustodiaJS/custodiajs-core/utils"
)

func newGlobalProcLoc() *gloablProcLog {
	return nil
}

func (o *gloablProcLog) NewProcLogForCore() *ProcLogSession {
	randHex, _ := utils.RandomHex(4)
	sessioncollorFormater := utils.DetermineColorFromHex(randHex)
	val := &ProcLogSession{id: randHex, header: "Core", sessionColorFunc: sessioncollorFormater}
	return val
}

func (o *gloablProcLog) NewProcLogForHostAPISocket() *ProcLogSession {
	randHex, _ := utils.RandomHex(4)
	sessioncollorFormater := utils.DetermineColorFromHex(randHex)
	val := &ProcLogSession{id: randHex, header: "Host-API-Socket", sessionColorFunc: sessioncollorFormater}
	return val
}

func (o *gloablProcLog) NewProcLogForHttpAPISocket() *ProcLogSession {
	randHex, _ := utils.RandomHex(4)
	sessioncollorFormater := utils.DetermineColorFromHex(randHex)
	val := &ProcLogSession{id: randHex, header: "Http-API-Socket", sessionColorFunc: sessioncollorFormater}
	return val
}

func (o *gloablProcLog) NewProcessProcLog() *ProcLogSession {
	return &ProcLogSession{}
}

var glogger *gloablProcLog = newGlobalProcLoc()
