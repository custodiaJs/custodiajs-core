package procslog

type ProcLogSession struct {
	mergedContainer  bool
	merged           []*ProcLogSession
	id               string
	header           string
	sessionColorFunc func(a ...interface{}) string
	printFunction    func(text string)
}
