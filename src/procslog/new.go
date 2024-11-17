package procslog

import (
	cenvxcore "github.com/custodia-cenv/cenvx-core/src"
	"github.com/custodia-cenv/cenvx-core/src/utils"
)

func extractMergedSessions(p *ProcLogSession) []*ProcLogSession {
	var allMergedSessions []*ProcLogSession

	var extract func(sessions []*ProcLogSession)
	extract = func(sessions []*ProcLogSession) {
		for _, session := range sessions {
			allMergedSessions = append(allMergedSessions, session)
			if len(session.merged) > 0 {
				extract(session.merged)
			}
		}
	}

	extract(p.merged)

	return allMergedSessions
}

func NewChainMergedProcLog(vat ...cenvxcore.ProcessLogSessionInterface) *ProcLogSession {
	newProcLog := NewProcLog()
	newProcLog.mergedContainer = true
	newProcLog.merged = []*ProcLogSession{}
	for _, item := range vat {
		x := item.(*ProcLogSession)
		if len(x.merged) != 0 {
			newProcLog.merged = extractMergedSessions(x)
		}
		if !x.mergedContainer {
			newProcLog.merged = append(newProcLog.merged, x)
		}
	}

	return newProcLog
}

func NewProcLogForCore() *ProcLogSession {
	newProcLog := NewProcLog()
	newProcLog.header = "Core"
	return newProcLog
}

func NewProcLogForHostAPISocket() *ProcLogSession {
	newProcLog := NewProcLog()
	newProcLog.header = "Host-API-Socket"
	return newProcLog
}

func NewProcLogForHttpAPISocket() *ProcLogSession {
	newProcLog := NewProcLog()
	newProcLog.header = "Http-API-Socket"
	return newProcLog
}

func NewProcLogForVmProcess() *ProcLogSession {
	newProcLog := NewProcLog()
	newProcLog.header = "VM-Process"
	return newProcLog
}

func NewProcLog() *ProcLogSession {
	randHex, _ := utils.RandomHex(4)
	sessioncollorFormater := utils.DetermineColorFromHex(randHex)
	val := &ProcLogSession{id: randHex, sessionColorFunc: sessioncollorFormater}
	return val
}

func NewProcLogSessionWithHeader(header string) *ProcLogSession {
	newProcLog := NewProcLog()
	newProcLog.header = header
	return newProcLog
}
