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
	"crypto/x509"
	"net/http"
	"net/url"
	"sync"

	"github.com/CustodiaJS/custodiajs-core/core/consolecache"

	v8 "rogchap.com/v8go"
)

type CoreInterface interface {
	GetAllVMs() []VmInterface
	GetAllActiveScriptContainerIDs(processLog ProcessLogSessionInterface) []string
	GetScriptContainerVMByID(vmid string) (VmInterface, bool, *SpecificError)
	GetScriptContainerByVMName(string) (VmInterface, bool, *SpecificError)
	GetCoreSessionManagmentUnit() ContextManagmentUnitInterface
}

type VmInterface interface {
	GetVMName() string
	IsAllowedXRequested(xrd *XRequestedWithData) bool
	GetFingerprint() CoreVMFingerprint
	GetConsoleOutputWatcher() WatcherInterface
	GetAllSharedFunctions() []SharedFunctionInterface
	Serve(*sync.WaitGroup) error
	GetSharedFunctionBySignature(RPCCallSource, *FunctionSignature) (SharedFunctionInterface, bool, *SpecificError)
	GetStartingTimestamp() uint64
	GetKId() KernelID
	SignalShutdown()
	GetState() VmState
	GetOwner() string
	GetRepoURL() string
}

type APISocketInterface interface {
	Serve(chan struct{}) error
	SetupCore(CoreInterface) error
}

type SharedFunctionInterface interface {
	GetName() string
	GetParmTypes() []string
	GetReturnDatatype() string
	EnterFunctionCall(*RpcRequest) *SpecificError
	GetScriptVM() VmInterface
}

type WatcherInterface interface {
	Read() string
}

type RpcRequestInterface interface {
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
	AddToEventLoop(KernelEventLoopOperationInterface) *SpecificError
	GetFingerprint() KernelFingerprint
	Signal(id string, value interface{})
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
	GetChildLog(header string) ProcessLogSessionInterface
	Log(format string, value ...interface{})
	Debug(format string, value ...interface{})
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

type FunctionCallReturnChanInterface interface {
	WriteAndClose(value *FunctionCallReturn)
	Read() (*FunctionCallReturn, bool)
	IsClosed() bool
	Close()
}

type FunctionCallStateChanInterface interface {
	Read() (*FunctionCallReturn, bool)
	WriteAndClose(*FunctionCallReturn)
}

type GrsboolInterface interface {
	Set(bval bool)
	Bool() bool
	WaitOfChange(waitOfState bool)
}

type CoreContextInterface interface {
	GetChildProcessLog(header string) ProcessLogSessionInterface
	GetProcessLog() ProcessLogSessionInterface
	IsConnected() bool
	Close()
	Done()
}

type CoreHttpContextInterface interface {
	CoreContextInterface
	SetMethod(method HTTP_METHOD)
	SetContentType(HttpRequestContentType)
	SetXRequestedWith(*XRequestedWithData)
	SetReferer(refererURL *url.URL)
	SetOrigin(originURL *url.URL)
	SetTLSCertificate(tlsCert []*x509.Certificate)
	AddSearchedFunctionSignature(fncs *FunctionSignature)
	GetSearchedFunctionSignature() *FunctionSignature
	GetMethod() HTTP_METHOD
	GetContentType() HttpRequestContentType
	GetXRequestedWith() *XRequestedWithData
	GetReferer() *url.URL
	GetOrigin() *url.URL
	GetTLSCertificate() []*x509.Certificate

	SignalTheErrorSignalCouldNotBeTransmittedTheConnectionWasLost(size int, error *SpecificError)
	SignalsThatAnErrorHasOccurredWhenTheErrorIsSent(size int, err *SpecificError)
	SignalTheResponseWasTransmittedSuccessfully(size int, packageHash string)
	SignalTheResponseCouldNotBeSent(size int, error *SpecificError)
	SignalThatTheErrorWasSuccessfullyTransmitted(size int)
	GetReturnChan() FunctionCallReturnChanInterface
	CloseBecauseFunctionReturned()
}

type ContextManagmentUnitInterface interface {
	NewHTTPBasesSession(r *http.Request) (CoreHttpContextInterface, *SpecificError)
}

type CustodiaJSNetworkHypervisorInterface interface {
	Start() *SpecificError
}
