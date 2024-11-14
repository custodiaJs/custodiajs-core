package log

import "fmt"

func LogError(format string, value ...any) {
	fmt.Println(fmt.Sprintf(format, value...))
}

func InfoLogPrint(format string, value ...any) {
	fmt.Println(fmt.Sprintf(format, value...))
}

func DebugLogPrint(format string, value ...any) {
	fmt.Println(fmt.Sprintf(format, value...))
}
