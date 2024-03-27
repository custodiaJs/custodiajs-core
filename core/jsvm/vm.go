package jsvm

import (
	"fmt"
	"vnh1/core/consolecache"
	"vnh1/static"

	"github.com/dop251/goja"
)

type JsVM struct {
	sharedLocalFunctions  map[string]*SharedLocalFunction
	sharedPublicFunctions map[string]*SharedPublicFunction
	coreService           static.CoreInterface
	cache                 map[string]interface{}
	consoleCache          *consolecache.ConsoleOutputCache
	allowedBuckets        []string
	config                *JsVMConfig
	gojaVM                *goja.Runtime
	exports               *goja.Object
	loadRootLib           bool
	scriptLoaded          bool
}

func (o *JsVM) prepareVM() error {
	// Die Standardobjekte werden erzeugt
	vnh1Obj := o.gojaVM.NewObject()

	// Die VNH1 Funktionen werden bereitgestellt
	vnh1Obj.Set("com", o.gojaCOMFunctionModule)
	o.gojaVM.Set("vnh1", vnh1Obj)

	// Die JS Exports werden bereitgestellt
	o.gojaVM.Set("exports", o.exports)

	/* Es wird geprüft ob das API Root Script durch die VM bereitgestellt werden soll
	if o.loadRootLib {

	}*/

	// Der Vorgang ist ohne Fehler durchgeführt wurden
	return nil
}

func (o *JsVM) RunScript(script string) error {
	// Es wird geprüft ob das Script beretis geladen wurden
	if o.scriptLoaded {
		return fmt.Errorf("LoadScript: always script loaded")
	}

	// Es wird markiert dass das Script geladen wurde
	o.scriptLoaded = true

	// Das Script wird ausgeführt
	_, err := o.gojaVM.RunString(script)
	if err != nil {
		panic(err)
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

func (o *JsVM) gojaCOMFunctionModule(call goja.FunctionCall) goja.Value {
	// Es wird ermittelt um welchen vorgang es sich handelt
	if len(call.Arguments) < 1 {
		return o.gojaVM.ToValue("invalid")
	}

	// Die jeweilige Funktion wird ermittelt
	switch call.Arguments[0].String() {
	// Konsolen funktionen
	case "console":
		return console_base(o.gojaVM, call, o)
	// Share Functions
	case "root":
		return root_base(o.gojaVM, call, o)
	// S3 Funktionen
	case "s3":
		// Es wird geprüft ob die S3 Funktionen verfügbar sind
		if !o.config.EnableS3 {
			return goja.Undefined()
		}

		// Die S3 Funktionen werden bereitgestellt
		return sthreeb_base(o.gojaVM, call, o)
	// Die Cache Funktionen werden bereitgesllt
	case "cache":
		// Es wird geprüft ob Cache Funktion verfügbar sind
		if !o.config.EnableCache {
			return goja.Undefined()
		}

		// Die Cache Funktionen werden bereitgestellt
		return cache_base(o.gojaVM, call, o)
	// Es handelt sich um ein Unbekanntes Modul
	default:
		return goja.Undefined()
	}
}

func (o *JsVM) GetLocalShareddFunctions() []static.SharedLocalFunctionInterface {
	extracted := make([]static.SharedLocalFunctionInterface, 0)
	for _, item := range o.sharedLocalFunctions {
		extracted = append(extracted, item)
	}
	return extracted
}

func (o *JsVM) GetPublicShareddFunctions() []static.SharedPublicFunctionInterface {
	extracted := make([]static.SharedPublicFunctionInterface, 0)
	for _, item := range o.sharedPublicFunctions {
		extracted = append(extracted, item)
	}
	return extracted
}

func (o *JsVM) GetConsoleOutputWatcher() static.WatcherInterface {
	return o.consoleCache.GetOutputStream()
}

func NewVM(core static.CoreInterface, config *JsVMConfig) (*JsVM, error) {
	// Die GoJA VM wird erstellt
	gojaVM := goja.New()

	// Das Basisobjekt wird erzeugt
	var vmObject *JsVM
	if config == nil {
		vmObject = &JsVM{
			config:                &defaultConfig,
			gojaVM:                gojaVM,
			scriptLoaded:          false,
			exports:               gojaVM.NewObject(),
			sharedLocalFunctions:  make(map[string]*SharedLocalFunction),
			sharedPublicFunctions: make(map[string]*SharedPublicFunction),
			consoleCache:          consolecache.NewConsoleOutputCache(),
			allowedBuckets:        make([]string, 0),
			cache:                 make(map[string]interface{}),
			loadRootLib:           false,
			coreService:           core,
		}
	} else {
		vmObject = &JsVM{
			config:                config,
			gojaVM:                gojaVM,
			scriptLoaded:          false,
			exports:               gojaVM.NewObject(),
			sharedLocalFunctions:  make(map[string]*SharedLocalFunction),
			sharedPublicFunctions: make(map[string]*SharedPublicFunction),
			consoleCache:          consolecache.NewConsoleOutputCache(),
			allowedBuckets:        make([]string, 0),
			cache:                 make(map[string]interface{}),
			loadRootLib:           false,
			coreService:           core,
		}
	}

	// Die Funktionen werden hinzugefügt
	if err := vmObject.prepareVM(); err != nil {
		return nil, fmt.Errorf("NewVM: " + err.Error())
	}

	// Das VM Objekt wird zurückgegeben
	return vmObject, nil
}
