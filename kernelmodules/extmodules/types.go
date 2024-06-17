package extmodules

import cgowrapper "github.com/CustodiaJS/custodiajs-core/kernelmodules/extmodules/cgo_wrapper"

type ExternModuleImport struct {
}

type ExternModuleObject struct {
}

type ExternModuleEvent struct {
}

type ExternalModule struct {
	*cgowrapper.CGOWrappedLibModule
}
