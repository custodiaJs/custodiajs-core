package external_modules

import cgowrapper "github.com/CustodiaJS/custodiajs-core/kernel/external_modules/cgo_wrapper"

type ExternModuleImport struct {
}

type ExternModuleObject struct {
}

type ExternModuleEvent struct {
}

type ExternalModule struct {
	*cgowrapper.CGOWrappedLibModule
}
