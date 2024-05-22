package kernel

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"vnh1/consolecache"
	"vnh1/types"
	"vnh1/utils"
	"vnh1/vmdb"

	v8 "rogchap.com/v8go"
)

func (o *Kernel) Console() *consolecache.ConsoleOutputCache {
	return o.console
}

func (o *Kernel) GetCAMembershipCerts() []types.VmCaMembershipCertInterface {
	return nil
}

func (o *Kernel) GetNewIsolateContext() (*v8.Isolate, *v8.Context, error) {
	// Es wird versucht eine neue ISO und einen neuen Context mit VM Zugehörigkeit zu erzeugen
	iso, context, err := makeIsolationAndContext(o, false)
	if err != nil {
		return nil, nil, fmt.Errorf("Kernel->GetNewIsolateContext: " + err.Error())
	}

	// Die Objekte werden zurückgegeben
	return iso, context, nil
}

func (o *Kernel) ContextV8() *v8.Context {
	return o.Context
}

func (o *Kernel) GetFingerprint() types.KernelFingerprint {
	return types.KernelFingerprint(strings.ToLower(o.dbEntry.GetVMContainerMerkleHash()))
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

func (o *Kernel) HostRegisterRead(name_id string) interface{} {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	value, no := o.register[name_id]
	if !no {
		return nil
	}
	//fmt.Println("GLOB_REG_READ: " + name_id)

	return value
}

func (o *Kernel) HostRegisterWrite(name_id string, value interface{}) error {
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
		log.Printf("VM(%s): %s:-$ %s", o.id, header, fmt.Sprintf(format, v...))
	} else {
		log.Printf("VM(%s):-$ %s", o.id, fmt.Sprintf(format, v...))
	}
}

func (o *Kernel) GetKId() types.KernelID {
	return o.id
}

func (o *Kernel) GetCAMembershipIDs() []string {
	membIds := make([]string, 0)
	for _, item := range o.dbEntry.GetRootMemberIDS() {
		membIds = append(membIds, item.Fingerprint)
	}
	return membIds
}

func (o *Kernel) GetCore() types.CoreInterface {
	return o.core
}

func (o *Kernel) LinkKernelWithCoreVM(coreObj types.CoreVMInterface) error {
	// Der Mutex wird verwendet
	o.mutex.Lock()
	defer o.mutex.Unlock()

	// Es wird geprüft ob bereits eine VM mit dem Kernel verlinkt wurde
	if o.vmLink != nil {
		return fmt.Errorf("vm always linked with kernel")
	}

	// Der Kernel wird mit dem VM Verlinkt
	o.vmLink = coreObj

	// Es ist kein Fehler aufgetreten
	return nil
}

func (o *Kernel) AsCoreVM() types.CoreVMInterface {
	return o.vmLink
}

func (o *Kernel) HasCloseSignal() bool {
	return false
}

func (o *Kernel) ServeEventLoop() error {
	// Der Mutex wird verwendet
	o.eventLoopLockCond.L.Lock()

	// Es wird geprüft ob ein Eintrag vorhanden ist, wenn nicht wird gewartet
	if len(o.eventLoopStack) == 0 {
		o.eventLoopLockCond.Wait()
	}

	// Die Funktion wird aus dem Stack entfertn
	resolvedFunction := o.eventLoopStack[0]
	o.eventLoopStack = o.eventLoopStack[1:]

	// Der Mutex wird freigegeben
	o.eventLoopLockCond.L.Unlock()

	// Die Funktion wird aufgerufen
	o.LogPrint("", "Eventloop processes operation")

	// Die Funktion wird aufgerufen
	if err := resolvedFunction(o.Context); err != nil {
		return fmt.Errorf("ServeEventLoop: " + err.Error())
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

func (o *Kernel) AddFunctionCallToEventLoop(funcv func(*v8.Context) error) error {
	// Der Mutex wird verwendet
	o.eventLoopLockCond.L.Lock()

	// Die Eventfunktion wird hinzugefügt
	o.eventLoopStack = append(o.eventLoopStack, funcv)

	// Es wird Signalisiert, dass ein neuer Eintrag vorhanden ist
	o.eventLoopLockCond.Broadcast()

	// Der Cond wird freigegeben
	o.eventLoopLockCond.L.Unlock()

	// Rückgabe
	return nil
}

func makeIsolationAndContext(kernel *Kernel, isMain bool) (*v8.Isolate, *v8.Context, error) {
	// Die Isolation wird erezrugt
	iso := v8.NewIsolate()

	// Der Context wird erzeugt
	context := v8.NewContext(iso)

	// Es werden alle Standard Module geladen
	for _, item := range kernel.config.Modules {
		// Es wird geprüft ob es sich um den Main Context handelt
		if item.OnlyForMain() {
			if !isMain {
				continue
			}
		}

		// Das Modul wird geladen
		if err := item.Init(kernel, iso, context); err != nil {
			return nil, nil, fmt.Errorf("makeIsolationAndContext: " + err.Error())
		}
	}

	// Es ist kein Fehler aufgetreten
	return iso, context, nil
}

func NewKernel(consoleCache *consolecache.ConsoleOutputCache, kernelConfig *KernelConfig, dbEntry *vmdb.VmDBEntry, coreIface types.CoreInterface) (*Kernel, error) {
	// Die KernelID wird erzeugt
	kid, err := utils.RandomHex(6)
	if err != nil {
		return nil, fmt.Errorf("Kernel->NewKernel:" + err.Error())
	}

	// Der Mutex wird erzeugt
	mutex := &sync.Mutex{}

	// Das Kernelobjekt wird erzeugt
	kernelObj := &Kernel{
		Context:           nil,
		id:                types.KernelID(kid),
		register:          make(map[string]interface{}),
		mutex:             mutex,
		console:           consoleCache,
		config:            kernelConfig,
		core:              coreIface,
		vmImports:         make(map[string]*v8.Value),
		dbEntry:           dbEntry,
		eventLoopStack:    make([]func(*v8.Context) error, 0),
		eventLoopLockCond: sync.NewCond(mutex),
	}

	// Der Context wird im Kernel Objekt gespeichert
	_, context, err := makeIsolationAndContext(kernelObj, true)
	if err != nil {
		return nil, fmt.Errorf("Kernel->NewKernel:" + err.Error())
	}

	// Der Context wird im Kernel abgespeichert
	kernelObj.Context = context

	// Die Require Funktionen werden Registriert
	if err := kernelObj._setup_require(); err != nil {
		return nil, fmt.Errorf("Kernel->NewKernel: " + err.Error())
	}

	// Das Kernelobejkt wird zurückgegeben
	return kernelObj, nil
}
