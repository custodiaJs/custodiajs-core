package procslog

import (
	"fmt"

	"github.com/CustodiaJS/custodiajs-core/global/types"
)

func ProcFormatConsoleText(procLog types.ProcessLogSessionInterface, header string, consoleValue types.CONSOLE_TEXT, value ...string) {
	anyValues := make([]any, 0)
	for _, item := range value {
		anyValues = append(anyValues, item)
	}
	procLog.LogPrint(header, fmt.Sprintf(string(consoleValue), anyValues...))
}
