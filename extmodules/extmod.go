package extmodules

func (o *ExternalModule) GetGlobalFunctions() []*ExternModuleFunction {
	vat := make([]*ExternModuleFunction, 0)
	for _, item := range o.CGOWrappedLibModule.GetGlobalFunctions() {
		vat = append(vat, &ExternModuleFunction{item})
	}
	return vat
}

func (o *ExternalModule) GetName() string {
	return o.CGOWrappedLibModule.GetName()
}

func (o *ExternalModule) GetVersion() uint64 {
	return uint64(o.CGOWrappedLibModule.GetVersion())
}
