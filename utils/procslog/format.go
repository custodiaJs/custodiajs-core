package procslog

import (
	"fmt"
	"vnh1/types"
)

func FormatConsoleText(consoleValue types.CONSOLE_TEXT, value ...string) string {
	anyValues := make([]any, 0)
	for _, item := range value {
		anyValues = append(anyValues, item)
	}
	return fmt.Sprintf(string(consoleValue), anyValues...)
}

func ProcFormatConsoleText(procLog *ProcLogSession, header string, consoleValue types.CONSOLE_TEXT, value ...string) {
	anyValues := make([]any, 0)
	for _, item := range value {
		anyValues = append(anyValues, item)
	}
	procLog.LogPrint(header, fmt.Sprintf(string(consoleValue), anyValues...))
}
