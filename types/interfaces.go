package types

import "github.com/dop251/goja"

type JsVmInterface interface {
	GetPublicSharedFunctions() []SharedPublicFunctionInterface
	GetLocalSharedFunctions() []SharedLocalFunctionInterface
	GetState() VmState
}

type CoreInterface interface {
	GetAllVMs() []CoreVMInterface
	GetAllActiveScriptContainerIDs() []string
	GetScriptContainerVMByID(string) (CoreVMInterface, error)
}

type CoreVMInterface interface {
	GetVMName() string
	GetFingerprint() string
	GetVMModuleNames() []string
	GetLocalSharedFunctions() []SharedLocalFunctionInterface
	GetPublicSharedFunctions() []SharedPublicFunctionInterface
	GetConsoleOutputWatcher() WatcherInterface
	GetAllSharedFunctions() []SharedFunctionInterface
	GetWhitelist() []TransportWhitelistVmEntryInterface
	GetMemberCertKeyIds() []string
	GetStartingTimestamp() uint64
	GetState() VmState
}

type APISocketInterface interface {
	Serve(chan struct{}) error
	SetupCore(CoreInterface) error
}

type SharedLocalFunctionInterface interface {
	GetName() string
	GetParmTypes() []string
	EnterFunctionCall(RpcRequestData, RpcRequestInterface) (goja.Value, error)
}

type SharedPublicFunctionInterface interface {
	GetName() string
	GetParmTypes() []string
	EnterFunctionCall(RpcRequestData, RpcRequestInterface) (goja.Value, error)
}

type SharedFunctionInterface interface {
	GetName() string
	GetParmTypes() []string
	EnterFunctionCall(RpcRequestData, RpcRequestInterface) (goja.Value, error)
}

type WatcherInterface interface {
	Read() string
}

type RpcRequestInterface interface {
	GetParms() []FunctionParameterBundleInterface
}

type RpcRequestData interface {
}

type FunctionParameterBundleInterface interface {
	GetType() string
	GetValue() interface{}
}

type TransportWhitelistVmEntryInterface interface {
	URL() string
	Alias() string
}
