package procslog

import (
	"testing"
)

func Test(t *testing.T) {
	core := NewProcLogForCoreTest(t)
	apisocket := NewProcLogForHostAPISocketTest(t)
	session := NewChainMergedProcLogTest(t, apisocket, NewProcLogSessionWithHeader("Session"))
	testConnection := NewChainMergedProcLogTest(t, session, core)

	core.Debug("Created")
	apisocket.Debug("Created")
	session.Debug("Created")
	testConnection.Debug("Test")
}
