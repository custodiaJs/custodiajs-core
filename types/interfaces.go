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
	GetScriptContainerByVMName(string) (CoreVMInterface, error)
}

type CoreVMInterface interface {
	GetVMName() string
	GetFingerprint() CoreVMFingerprint
	GetLocalSharedFunctions() []SharedLocalFunctionInterface
	GetPublicSharedFunctions() []SharedPublicFunctionInterface
	GetConsoleOutputWatcher() WatcherInterface
	GetAllSharedFunctions() []SharedFunctionInterface
	GetWhitelist() []*TransportWhitelistVmEntryData
	ValidateRPCRequestSource(soruce string) bool
	GetDatabaseServices() []*VMDatabaseData
	GetMemberCertsPkeys() []*CAMemberData
	GetStartingTimestamp() uint64
	GetState() VmState
	GetOwner() string
	GetRepoURL() string
	GetMode() string
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

type AlternativeServiceInterface interface {
}
