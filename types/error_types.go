package types

type ExtModCGOPanic struct {
	ErrorValue error
}

type ExtModFunctionCallError struct {
	ErrorValue error
}

func (e *ExtModCGOPanic) Error() string {
	return e.ErrorValue.Error()
}

func (e *ExtModFunctionCallError) Error() string {
	return e.ErrorValue.Error()
}
