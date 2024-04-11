package extmodules

func (o *ExternModuleFunction) Call() (string, interface{}, error) {
	return o.CGOWrappedLibModuleFunction.Call()
}
