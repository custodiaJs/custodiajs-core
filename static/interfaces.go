package static

type JsVmInterface interface {
	GetPublicShareddFunctions() []SharedPublicFunctionInterface
	GetLocalShareddFunctions() []SharedLocalFunctionInterface
	GetState() VmState
}

type CoreInterface interface {
	GetAllActiveScriptContainerIDs() []string
	GetScriptContainerVMByID(string) (CoreVMInterface, error)
}

type CoreVMInterface interface {
	GetVMName() string
	GetFingerprint() string
	GetVMModuleNames() []string
	GetLocalShareddFunctions() []SharedLocalFunctionInterface
	GetPublicShareddFunctions() []SharedPublicFunctionInterface
	GetState() VmState
}

type APISocketInterface interface {
	Serve(chan struct{}) error
	SetupCore(CoreInterface) error
}

type SharedLocalFunctionInterface interface {
	GetName() string
	GetParmTypes() []string
}

type SharedPublicFunctionInterface interface {
	GetName() string
	GetParmTypes() []string
}
