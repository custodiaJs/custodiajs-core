package procslog

import (
	"fmt"

	cenvxcore "github.com/custodia-cenv/cenvx-core/src"
)

func ProcFormatConsoleText(procLog cenvxcore.ProcessLogSessionInterface, header string, consoleValue cenvxcore.CONSOLE_TEXT, value ...string) {
	anyValues := make([]any, 0)
	for _, item := range value {
		anyValues = append(anyValues, item)
	}
	procLog.LogPrint(header, fmt.Sprintf(string(consoleValue), anyValues...))
}
