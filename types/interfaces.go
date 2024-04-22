package types

import (
	"vnh1/core/consolecache"

	v8 "rogchap.com/v8go"
)

type CoreInterface interface {
	GetAllVMs() []CoreVMInterface
	GetAllActiveScriptContainerIDs() []string
	GetScriptContainerVMByID(vmid string) (CoreVMInterface, bool, error)
	GetScriptContainerByVMName(string) (CoreVMInterface, error)
}

type CoreVMInterface interface {
	GetVMName() string
	GetFingerprint() CoreVMFingerprint
	GetConsoleOutputWatcher() WatcherInterface
	GetAllSharedFunctions() []SharedFunctionInterface
	GetSharedFunctionBySignature(sourceType RPCCallSource, funcSignature *FunctionSignature) (SharedFunctionInterface, bool, error)
	GetWhitelist() []*TransportWhitelistVmEntryData
	ValidateRPCRequestSource(soruce string) bool
	GetDatabaseServices() []*VMDatabaseData
	GetRootMemberIDS() []*CAMemberData
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

type SharedFunctionInterface interface {
	GetName() string
	GetParmTypes() []string
	GetReturnDType() string
	EnterFunctionCall(*RpcRequest) (*FunctionCallState, error)
}

type WatcherInterface interface {
	Read() string
}

type HttpJsonRequestData interface {
}

type AlternativeServiceInterface interface {
}

type VmCaMembershipCertInterface interface {
}

type KernelInterface interface {
	LogPrint(header string, format string, v ...any)
	GloablRegisterWrite(string, interface{}) error
	Console() *consolecache.ConsoleOutputCache
	AddImportModule(string, *v8.Value) error
	GloablRegisterRead(string) interface{}
	GetNewIsolateContext() (*v8.Isolate, *v8.Context, error)
	GetCAMembershipCerts() []VmCaMembershipCertInterface
	GetFingerprint() KernelFingerprint
	AsCoreVM() CoreVMInterface
	GetCAMembershipIDs() []string
	GetCore() CoreInterface
	GetKId() KernelID
}

type KernelModuleInterface interface {
	GetName() string
	Init(KernelInterface, *v8.Isolate, *v8.Context) error
	OnlyForMain() bool
}
