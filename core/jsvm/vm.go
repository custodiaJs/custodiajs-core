package jsvm

import (
	"fmt"
	"time"
	"vnh1/core/consolecache"
	"vnh1/types"

	"github.com/dop251/goja"
)

type JsVM struct {
	sharedLocalFunctions  map[string]*SharedLocalFunction
	sharedPublicFunctions map[string]*SharedPublicFunction
	coreService           types.CoreInterface
	cache                 map[string]interface{}
	consoleCache          *consolecache.ConsoleOutputCache
	allowedBuckets        []string
	config                *JsVMConfig
	gojaVM                *goja.Runtime
	exports               *goja.Object
	loadRootLib           bool
	scriptLoaded          bool
	startTimeUnix         uint64
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

func (o *JsVM) functionIsSharing(functionName string) bool {
	// Es wird geprüft ob die Funktion bereits Registriert wurde
	_, found := o.sharedLocalFunctions[functionName]

	// Das Ergebniss wird zurückgegeben
	return found
}

func (o *JsVM) shareLocalFunction(funcName string, parmTypes []string, function goja.Callable) error {
	// Es wird geprüft ob diese Funktion bereits registriert wurde
	if _, found := o.sharedLocalFunctions[funcName]; found {
		return fmt.Errorf("function always registrated")
	}

	// Die Funktion wird zwischengespeichert
	o.sharedLocalFunctions[funcName] = &SharedLocalFunction{
		callFunction: function,
		name:         funcName,
		parmTypes:    parmTypes,
		gojaVM:       o.gojaVM,
	}

	// Die Funktion wird im Core registriert
	fmt.Println("VM:SHARE_LOCAL_FUNCTION:", funcName, parmTypes)

	// Der Vorgang wurde ohne Fehler durchgeführt
	return nil
}

func (o *JsVM) sharePublicFunction(funcName string, parmTypes []string, function goja.Callable) error {
	// Es wird geprüft ob diese Funktion bereits registriert wurde
	if _, found := o.sharedPublicFunctions[funcName]; found {
		return fmt.Errorf("function always registrated")
	}

	// Die Funktion wird zwischengespeichert
	o.sharedPublicFunctions[funcName] = &SharedPublicFunction{
		callFunction: function,
		name:         funcName,
		parmTypes:    parmTypes,
		gojaVM:       o.gojaVM,
	}

	// Die Funktion wird im Core registriert
	fmt.Println("VM:SHARE_PUBLIC_FUNCTION:", funcName, parmTypes)

	// Der Vorgang wurde ohne Fehler durchgeführt
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

	// Die Aktuelle Uhrzeit wird ermittelt
	o.startTimeUnix = uint64(time.Now().Unix())

	// Es ist kein Fehler aufgetreten
	return nil
}

func (o *JsVM) GetLocalSharedFunctions() []types.SharedLocalFunctionInterface {
	extracted := make([]types.SharedLocalFunctionInterface, 0)
	for _, item := range o.sharedLocalFunctions {
		extracted = append(extracted, item)
	}
	return extracted
}

func (o *JsVM) GetPublicSharedFunctions() []types.SharedPublicFunctionInterface {
	extracted := make([]types.SharedPublicFunctionInterface, 0)
	for _, item := range o.sharedPublicFunctions {
		extracted = append(extracted, item)
	}
	return extracted
}

func (o *JsVM) GetConsoleOutputWatcher() types.WatcherInterface {
	return o.consoleCache.GetOutputStream()
}

func (o *JsVM) GetStartingTimestamp() uint64 {
	return o.startTimeUnix
}

func (o *JsVM) GetAllSharedFunctions() []types.SharedFunctionInterface {
	vat := make([]types.SharedFunctionInterface, 0)
	for _, item := range o.GetLocalSharedFunctions() {
		vat = append(vat, item)
	}
	for _, item := range o.GetPublicSharedFunctions() {
		vat = append(vat, item)
	}
	return vat
}

func NewVM(core types.CoreInterface, config *JsVMConfig) (*JsVM, error) {
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
