package jsvm

import (
	"sync"
	"vnh1/core/consolecache"

	"github.com/dop251/goja"
)

type JsVM struct {
	sharedLocalFunctions  map[string]*SharedLocalFunction
	sharedPublicFunctions map[string]*SharedPublicFunction
	cache                 map[string]interface{}
	consoleCache          *consolecache.ConsoleOutputCache
	allowedBuckets        []string
	config                *JsVMConfig
	gojaVM                *goja.Runtime
	exports               *goja.Object
	loadRootLib           bool
	scriptLoaded          bool
	startTimeUnix         uint64
	mutex                 *sync.Mutex
}

type S3MetaData struct {
}

type S3Bucket interface {
	UploadObject(string, interface{}, interface{}) error
	DownloadObject(string, interface{}) (interface{}, error)
	DeleteObject(string, interface{}) error
}

type LocalVMS3Bucket struct {
}

type SharedLocalFunction struct {
	gojaVM       *goja.Runtime
	callFunction goja.Callable
	name         string
	parmTypes    []string
}

type SharedPublicFunction struct {
	gojaVM       *goja.Runtime
	callFunction goja.Callable
	name         string
	parmTypes    []string
}

type FunctionCallReturn struct {
	CType string
	Value interface{}
}

type JsVMConfig struct {
	EnableWebsockets      bool
	EnableFunctionSharing bool
	EnableCache           bool
	EnableS3              bool
}
