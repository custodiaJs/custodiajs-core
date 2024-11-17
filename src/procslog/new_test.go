package procslog

import (
	"testing"

	cenvxcore "github.com/custodia-cenv/cenvx-core/src"
)

func NewProcLogForCoreTest(t *testing.T) *ProcLogSession {
	new := NewProcLogForCore()
	new.printFunction = func(text string) { t.Log(text) }
	return new
}

func NewProcLogForHostAPISocketTest(t *testing.T) *ProcLogSession {
	new := NewProcLogForHostAPISocket()
	new.printFunction = func(text string) { t.Log(text) }
	return new
}

func NewProcLogSessionWithHeaderTest(header string, t *testing.T) *ProcLogSession {
	new := NewProcLogSessionWithHeader(header)
	new.printFunction = func(text string) { t.Log(text) }
	return new
}

func NewChainMergedProcLogTest(t *testing.T, vat ...cenvxcore.ProcessLogSessionInterface) *ProcLogSession {
	new := NewChainMergedProcLog(vat...)
	new.printFunction = func(text string) { t.Log(text) }
	return new
}
