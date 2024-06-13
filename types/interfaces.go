// Author: fluffelpuff
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package types

import (
	"sync"
	"vnh1/consolecache"
	"vnh1/databaseservices/services"

	v8 "rogchap.com/v8go"
)

type CoreInterface interface {
	GetAllVMs() []VmInterface
	GetAllActiveScriptContainerIDs() []string
	GetScriptContainerVMByID(vmid string) (VmInterface, bool, error)
	GetScriptContainerByVMName(string) (VmInterface, error)
}

type VmInterface interface {
	GetVMName() string
	GetFingerprint() CoreVMFingerprint
	GetConsoleOutputWatcher() WatcherInterface
	GetAllSharedFunctions() []SharedFunctionInterface
	Serve(*sync.WaitGroup) error
	GetSharedFunctionBySignature(RPCCallSource, *FunctionSignature) (SharedFunctionInterface, bool, error)
	GetWhitelist() []*TransportWhitelistVmEntryData
	ValidateRPCRequestSource(soruce string) bool
	GetDatabaseServices() []*VMEntryBaseData
	AddDatabaseServiceLink(dbserviceLink services.DbServiceLinkinterface) error
	GetRootMemberIDS() []*CAMemberData
	GetStartingTimestamp() uint64
	GetKId() KernelID
	SignalShutdown()
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
	GetReturnDatatype() string
	EnterFunctionCall(*RpcRequest) error
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
	AddToEventLoop(KernelEventLoopOperationInterface) error
	GetFingerprint() KernelFingerprint
	AsCoreVM() VmInterface
	GetCAMembershipIDs() []string
	GetCore() CoreInterface
	GetKId() KernelID
}

type KernelModuleInterface interface {
	GetName() string
	Init(KernelInterface, *v8.Isolate, *v8.Context) error
	OnlyForMain() bool
}

type ProcessLogSessionInterface interface {
	LogPrint(string, string, ...interface{})
	LogPrintSuccs(string, ...interface{})
	LogPrintError(string, ...interface{})
	GetID() string
}

type KernelEventLoopContextInterface interface {
	SetError(error)
	SetResult(*v8.Value)
}

type KernelEventLoopOperationInterface interface {
	GetType() KernelEventLoopOperationMethode
	GetFunction() func(*v8.Context, KernelEventLoopContextInterface)
	WaitOfResponse() (*v8.Value, error)
	GetOperation() KernelEventLoopContextInterface
	GetSourceCode() string
	SetResult(*v8.Value)
	SetError(error)
}
