package extmodules

import cgowrapper "vnh1/extmodules/cgo_wrapper"

type ExternModuleFunction struct {
	*cgowrapper.CGOWrappedLibModuleFunction
}

type ExternalModule struct {
	*cgowrapper.CGOWrappedLibModule
}
