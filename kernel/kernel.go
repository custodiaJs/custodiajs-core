package kernel

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/CustodiaJS/custodiajs-core/consolecache"
	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/utils"
	"github.com/CustodiaJS/custodiajs-core/vmdb"

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

func (o *Kernel) LinkKernelWithCoreVM(coreObj types.VmInterface) error {
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

func (o *Kernel) AsCoreVM() types.VmInterface {
	return o.vmLink
}

func (o *Kernel) IsClosed() bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	cstate := bool(o.hasCloseSignal)
	return cstate
}

func (o *Kernel) ServeEventLoop() error {
	// Der Mutex wird verwendet
	o.eventLoopLockCond.L.Lock()

	// Es wird geprüft ob ein Eintrag vorhanden ist, wenn nicht wird gewartet
	if len(o.eventLoopStack) == 0 {
		o.eventLoopLockCond.Wait()
	}

	// Es wird ermittelt ob der Kernel beendet werden soll
	if o.hasCloseSignal {
		o.eventLoopLockCond.L.Unlock()
		return nil
	}

	// Die Funktion wird aus dem Stack entfertn
	eventLoopOperation := o.eventLoopStack[0]
	o.eventLoopStack = o.eventLoopStack[1:]

	// Der Mutex wird freigegeben
	o.eventLoopLockCond.L.Unlock()

	// Es wird geprüft ob es sich um eine Funktion oder einen Sourcecode Call handelt
	switch eventLoopOperation.GetType() {
	case static.KERNEL_EVENT_LOOP_FUNCTION:
		// Die Funktion wird abgerufen
		funct := eventLoopOperation.GetFunction()

		// Die Funktion wird ausgeführt
		funct(o.Context, eventLoopOperation.GetOperation())

		// Es ist kein Fehler aufgetreten
		return nil
	case static.KERNEL_EVENT_LOOP_SOURCE_CODE:
		// Der Code wird ausgeführt
		result, err := o.Context.RunScript(eventLoopOperation.GetSourceCode(), "eventloop.js")
		if err != nil {
			// Der Fehler wird zurückgegeben
			eventLoopOperation.SetError(err)

			// Der Fehler wird zurückgegeben
			return fmt.Errorf("Kernel->call_eventloop_function: " + err.Error())
		}

		// Die Rückgabe wird zurückgegeben
		eventLoopOperation.SetResult(result)

		// Es ist kein Fehler aufgetreten
		return nil
	default:
		return fmt.Errorf("Kernel->ServeEventLoop: unkown operation methode")
	}
}

func (o *Kernel) AddToEventLoop(operation types.KernelEventLoopOperationInterface) *types.SpecificError {
	// Mittels Goroutine wird ein neues Event hinzugefügt
	go func() {
		// Der Mutex wird verwendet
		o.eventLoopLockCond.L.Lock()

		// Die Eventfunktion wird hinzugefügt
		o.eventLoopStack = append(o.eventLoopStack, operation)

		// Es wird Signalisiert, dass ein neuer Eintrag vorhanden ist
		o.eventLoopLockCond.Broadcast()

		// Der Cond wird freigegeben
		o.eventLoopLockCond.L.Unlock()
	}()

	// Rückgabe
	return nil
}

func (o *Kernel) Close() {
	// Der Mutex wird angewendet
	o.eventLoopLockCond.L.Lock()

	// Es wird geprüft ob bereits ein Close Signal vorhanden ist
	if o.hasCloseSignal {
		o.eventLoopLockCond.L.Unlock()
		return
	}

	// Es wird Signalisiert dass die VN beendet werden soll
	o.hasCloseSignal = true

	// Die Eventloop wird angehalten
	o.eventLoopLockCond.Broadcast()

	// Der Mutex wird freigegeben
	o.eventLoopLockCond.L.Unlock()

	// Der V8 Context wird geschlossen
	o.ContextV8().Close()
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
		eventLoopStack:    make([]types.KernelEventLoopOperationInterface, 0),
		eventLoopLockCond: sync.NewCond(mutex),
		hasCloseSignal:    false,
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
