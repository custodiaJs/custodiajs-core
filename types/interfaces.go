package types

import "github.com/dop251/goja"

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
	GetConsoleOutputWatcher() WatcherInterface
	GetState() VmState
}

type APISocketInterface interface {
	Serve(chan struct{}) error
	SetupCore(CoreInterface) error
}

type SharedLocalFunctionInterface interface {
	GetName() string
	GetParmTypes() []string
	EnterFunctionCall(RpcRequestInterface) (goja.Value, error)
}

type SharedPublicFunctionInterface interface {
	GetName() string
	GetParmTypes() []string
	EnterFunctionCall(RpcRequestInterface) (goja.Value, error)
}

type SharedFunctionInterface interface {
	GetName() string
	GetParmTypes() []string
	EnterFunctionCall(RpcRequestInterface) (goja.Value, error)
}

type WatcherInterface interface {
	Read() string
}

type RpcRequestInterface interface {
	GetParms() []FunctionParameterBundleInterface
}

type FunctionParameterBundleInterface interface {
	GetType() string
	GetValue() interface{}
}
