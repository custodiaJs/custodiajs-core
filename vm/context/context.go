package vmcontext

import (
	"fmt"
	"log"
	"sync"

	"github.com/CustodiaJS/custodiajs-core/core/consolecache"
	"github.com/CustodiaJS/custodiajs-core/global/static"
	"github.com/CustodiaJS/custodiajs-core/global/types"
	"github.com/CustodiaJS/custodiajs-core/global/utils"
	contextconsole "github.com/CustodiaJS/custodiajs-core/vm/context/console"
	rpcjsprocessor "github.com/CustodiaJS/custodiajs-core/vm/context/rpc"

	v8 "rogchap.com/v8go"
)

func (o *VmContext) Console() *consolecache.ConsoleOutputCache {
	return o.console
}

func (o *VmContext) GetCAMembershipCerts() []types.VmCaMembershipCertInterface {
	return nil
}

func (o *VmContext) GetNewIsolateContext() (*v8.Isolate, *v8.Context, error) {
	// Es wird versucht eine neue ISO und einen neuen Context mit VM Zugehörigkeit zu erzeugen
	iso, context, err := makeIsolationAndContext(o, map[string]bool{})
	if err != nil {
		return nil, nil, fmt.Errorf("VmContext->GetNewIsolateContext: " + err.Error())
	}

	// Die Objekte werden zurückgegeben
	return iso, context, nil
}

func (o *VmContext) ContextV8() *v8.Context {
	return o.Context
}

func (o *VmContext) GetFingerprint() types.KernelFingerprint {
	return types.KernelFingerprint(o.vmLink.GetManifest().Filehash)
}

func (o *VmContext) GloablRegisterRead(name_id string) interface{} {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	value, no := o.register[name_id]
	if !no {
		return nil
	}
	//fmt.Println("GLOB_REG_READ: " + name_id)

	return value
}

func (o *VmContext) GloablRegisterWrite(name_id string, value interface{}) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.register[name_id] = value
	//fmt.Println("GLOB_REG_WRITE: "+name_id, value)

	return nil
}

func (o *VmContext) HostRegisterRead(name_id string) interface{} {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	value, no := o.register[name_id]
	if !no {
		return nil
	}
	//fmt.Println("GLOB_REG_READ: " + name_id)

	return value
}

func (o *VmContext) HostRegisterWrite(name_id string, value interface{}) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.register[name_id] = value
	//fmt.Println("GLOB_REG_WRITE: "+name_id, value)

	return nil
}

func (o *VmContext) AddImportModule(name string, v8Value *v8.Value) error {
	// Der Mutex wird verwendet
	o.mutex.Lock()
	defer o.mutex.Unlock()

	// Der Eintrag wird abgespeichert
	o.vmImports[name] = v8Value

	// Es ist kein Fehler aufgetreten
	return nil
}

func (o *VmContext) LogPrint(header string, format string, v ...any) {
	if header != "" {
		log.Printf("VM(%s): %s:-$ %s", o.id, header, fmt.Sprintf(format, v...))
	} else {
		log.Printf("VM(%s):-$ %s", o.id, fmt.Sprintf(format, v...))
	}
}

func (o *VmContext) GetKId() types.KernelID {
	return o.id
}

func (o *VmContext) GetCAMembershipIDs() []string {
	/*
		membIds := make([]string, 0)
		for _, item := range o.dbEntry.GetRootMemberIDS() {
			membIds = append(membIds, item.Fingerprint)
		}
		return membIds
	*/
	panic("not implemented")
}

func (o *VmContext) GetCore() types.CoreInterface {
	return o.core
}

func (o *VmContext) LinkKernelWithVmInstance(vmInstance types.VmInterface) error {
	// Der Mutex wird verwendet
	o.mutex.Lock()
	defer o.mutex.Unlock()

	// Es wird geprüft ob bereits eine VM mit dem VmContext verlinkt wurde
	if o.vmLink != nil {
		return fmt.Errorf("vm always linked with VmContext")
	}

	// Der VmContext wird mit dem VM Verlinkt
	o.vmLink = vmInstance

	// Es ist kein Fehler aufgetreten
	return nil
}

func (o *VmContext) AsVmInstance() types.VmInterface {
	return o.vmLink
}

func (o *VmContext) IsClosed() bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	cstate := bool(o.hasCloseSignal)
	return cstate
}

func (o *VmContext) ServeEventLoop() error {
	// Der Mutex wird verwendet
	o.eventLoopLockCond.L.Lock()

	// Es wird geprüft ob ein Eintrag vorhanden ist, wenn nicht wird gewartet
	if len(o.eventLoopStack) == 0 {
		o.eventLoopLockCond.Wait()
	}

	// Es wird ermittelt ob der VmContext beendet werden soll
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
			return fmt.Errorf("VmContext->call_eventloop_function: " + err.Error())
		}

		// Die Rückgabe wird zurückgegeben
		eventLoopOperation.SetResult(result)

		// Es ist kein Fehler aufgetreten
		return nil
	default:
		return fmt.Errorf("VmContext->ServeEventLoop: unkown operation methode")
	}
}

func (o *VmContext) AddToEventLoop(operation types.KernelEventLoopOperationInterface) *types.SpecificError {
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

func (o *VmContext) Close() {
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

func (o *VmContext) Signal(id string, value interface{}) {

}

func makeIsolationAndContext(VmContext *VmContext, enabeldFunctions map[string]bool) (*v8.Isolate, *v8.Context, error) {
	// Die Isolation wird erezrugt
	iso := v8.NewIsolate()

	// Der Context wird erzeugt
	context := v8.NewContext(iso)

	// Die Consolen Funktionen werden bereitgestellt
	if err := contextconsole.NewConsoleModule().Init(VmContext, iso, context); err != nil {
		return nil, nil, fmt.Errorf("makeIsolationAndContext: " + err.Error())
	}

	// Es wird geprüft ob die RPC Funktionen verfügbar sind
	if value, exists := enabeldFunctions["rpc"]; exists && value {
		if err := rpcjsprocessor.NewRPCModule().Init(VmContext, iso, context); err != nil {
			return nil, nil, fmt.Errorf("makeIsolationAndContext: " + err.Error())
		}
	}

	// Es ist kein Fehler aufgetreten
	return iso, context, nil
}

func NewKernel(consoleCache *consolecache.ConsoleOutputCache, coreIface types.CoreInterface, enabeldFunctions map[string]bool) (*VmContext, error) {
	// Die KernelID wird erzeugt
	kid, err := utils.RandomHex(6)
	if err != nil {
		return nil, fmt.Errorf("VmContext->NewKernel:" + err.Error())
	}

	// Der Mutex wird erzeugt
	mutex := &sync.Mutex{}

	// Das Kernelobjekt wird erzeugt
	jsprocessorObj := &VmContext{
		Context:           nil,
		id:                types.KernelID(kid),
		register:          make(map[string]interface{}),
		mutex:             mutex,
		console:           consoleCache,
		core:              coreIface,
		vmImports:         make(map[string]*v8.Value),
		eventLoopStack:    make([]types.KernelEventLoopOperationInterface, 0),
		eventLoopLockCond: sync.NewCond(mutex),
		hasCloseSignal:    false,
	}

	// Der Context wird im VmContext Objekt gespeichert
	_, context, err := makeIsolationAndContext(jsprocessorObj, enabeldFunctions)
	if err != nil {
		return nil, fmt.Errorf("VmContext->NewKernel:" + err.Error())
	}

	// Der Context wird im VmContext abgespeichert
	jsprocessorObj.Context = context

	// Die Require Funktionen werden Registriert
	if err := jsprocessorObj._setup_require(); err != nil {
		return nil, fmt.Errorf("VmContext->NewKernel: " + err.Error())
	}

	// Das Kernelobejkt wird zurückgegeben
	return jsprocessorObj, nil
}
