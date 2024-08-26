package procslog

type ProcLogSession struct {
	id               string
	header           string
	sessionColorFunc func(a ...interface{}) string
}

type ProcLogChildSession struct {
	header string
	mother *ProcLogSession
}
