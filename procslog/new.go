package procslog

import (
	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/utils"
)

func NewProcLogSession() *ProcLogSession {
	randHex, _ := utils.RandomHex(4)
	sessioncollorFormater := utils.DetermineColorFromHex(randHex)
	val := &ProcLogSession{id: randHex, sessionColorFunc: sessioncollorFormater}
	return val
}

func NewProcLogSessionWithHeader(header string) *ProcLogSession {
	randHex, _ := utils.RandomHex(4)
	sessioncollorFormater := utils.DetermineColorFromHex(randHex)
	val := &ProcLogSession{id: randHex, header: header, sessionColorFunc: sessioncollorFormater}
	return val
}

func NewProcLogForCore() *ProcLogSession {
	return glogger.NewProcLogForCore()
}

func NewProcLogForHostAPIService() *ProcLogSession {
	return glogger.NewProcLogForHostAPISocket()
}

func NewProcLogForHttpAPIService() *ProcLogSession {
	return glogger.NewProcLogForHttpAPISocket()
}

func NewChainMergedProcLog(vat ...types.ProcessLogSessionInterface) *ProcLogSession {
	return nil
}
