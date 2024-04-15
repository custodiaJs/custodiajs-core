package kernel

import (
	"fmt"
	"log"
	"sync"
	"vnh1/core/consolecache"
	"vnh1/utils"

	v8 "rogchap.com/v8go"
)

func (o *Kernel) Console() *consolecache.ConsoleOutputCache {
	return o.console
}

func (o *Kernel) ContextV8() *v8.Context {
	return o.Context
}

func (o *Kernel) KernelThrow(context *v8.Context, msg string) {
	errMsg, err := v8.NewValue(o.Isolate(), msg)
	if err != nil {
		panic(err)
	}
	context.Isolate().ThrowException(errMsg)
}

func (o *Kernel) GloablRegisterRead(name_id string) interface{} {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	value, no := o.register[name_id]
	if !no {
		return nil
	}
	//fmt.Println("GLOB_REG_READ: " + name_id)

	return value
}

func (o *Kernel) GloablRegisterWrite(name_id string, value interface{}) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.register[name_id] = value
	//fmt.Println("GLOB_REG_WRITE: "+name_id, value)

	return nil
}

func (o *Kernel) AddImportModule(name string, v8Value *v8.Value) error {
	// Der Mutex wird verwendet
	o.mutex.Lock()
	defer o.mutex.Unlock()

	// Der Eintrag wird abgespeichert
	o.vmImports[name] = v8Value

	// Es ist kein Fehler aufgetreten
	return nil
}

func (o *Kernel) LogPrint(header string, format string, v ...any) {
	if header != "" {
		log.Printf("vm@%s: %s:-$ %s", o.id, header, fmt.Sprintf(format, v...))
	} else {
		log.Printf("vm@%s:-$ %s", o.id, fmt.Sprintf(format, v...))
	}
}

func (o *Kernel) GetKId() string {
	return o.id
}

func NewKernel(consoleCache *consolecache.ConsoleOutputCache, kernelConfig *KernelConfig) (*Kernel, error) {
	// DIe VM wird erezugt
	iso := v8.NewIsolate()

	// Die KernelID wird erzeugt
	kid, err := utils.RandomHex(6)
	if err != nil {
		return nil, fmt.Errorf("Kernel->NewKernel:" + err.Error())
	}

	// Das Kernelobjekt wird erzeugt
	kernelObj := &Kernel{
		id:        kid,
		register:  make(map[string]interface{}),
		mutex:     &sync.Mutex{},
		console:   consoleCache,
		Context:   v8.NewContext(iso),
		config:    kernelConfig,
		vmImports: make(map[string]*v8.Value),
	}

	// Die Require Funktionen werden Registriert
	if err := kernelObj._setup_require(); err != nil {
		return nil, fmt.Errorf("Kernel->NewKernel: " + err.Error())
	}

	// Die Einzelnen Kernel Module werden Registriert
	for _, item := range kernelConfig.Modules {
		if err := item.Init(kernelObj); err != nil {
			return nil, fmt.Errorf("Kernel->NewKernel: " + err.Error())
		}
	}

	// Das Kernelobejkt wird zurückgegeben
	return kernelObj, nil
}
