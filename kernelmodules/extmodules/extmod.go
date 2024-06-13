package extmodules

import cgowrapper "vnh1/kernelmodules/extmodules/cgo_wrapper"

func (o *ExternalModule) GetGlobalFunctions() []*cgowrapper.CGOWrappedLibModuleFunction {
	return o.CGOWrappedLibModule.GetGlobalFunctions()
}

func (o *ExternalModule) GetImports() []*ExternModuleImport {
	return []*ExternModuleImport{}
}

func (o *ExternalModule) GetGlobalObjects() []*ExternModuleObject {
	return []*ExternModuleObject{}
}

func (o *ExternalModule) GetEventTriggers() []*ExternModuleEvent {
	return []*ExternModuleEvent{}
}

func (o *ExternalModule) GetName() string {
	return o.CGOWrappedLibModule.GetName()
}

func (o *ExternalModule) GetVersion() uint64 {
	return uint64(o.CGOWrappedLibModule.GetVersion())
}
