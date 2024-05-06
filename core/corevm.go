package core

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"vnh1/core/consolecache"
	"vnh1/core/databaseservices/services"
	"vnh1/core/kernel"
	"vnh1/core/vmdb"
	extmodules "vnh1/extmodules"
	"vnh1/static"
	"vnh1/types"
	"vnh1/utils"
)

func (o *CoreVM) GetVMName() string {
	return o.vmDbEntry.GetVMName()
}

func (o *CoreVM) GetFingerprint() types.CoreVMFingerprint {
	return types.CoreVMFingerprint(o.Kernel.GetFingerprint())
}

func (o *CoreVM) GetOwner() string {
	return o.vmDbEntry.GetOwner()
}

func (o *CoreVM) GetRepoURL() string {
	return o.vmDbEntry.GetRepoURL()
}

func (o *CoreVM) GetMode() string {
	return o.vmDbEntry.GetMode()
}

func (o *CoreVM) _routine(scriptContent []byte, syncWaitGroup *sync.WaitGroup) {
	// Der Mutex wird verwendet
	o.objectMutex.Lock()

	// Es wird geptüft ob der Stauts des Objektes STILL_WAIT ist
	if o.vmState != static.Starting {
		// Der Mutext wird freigegeben
		o.objectMutex.Unlock()

		// Rückgabe
		return
	}

	// Die Startzeit wird festgelegt
	o.startTimeUnix = uint64(time.Now().Unix())

	// Der Mutex wird freigegeben
	o.objectMutex.Unlock()

	// Das Script wird ausgeführt
	o.runScript(string(scriptContent))

	// Das Script wird als Running Markiert
	o.objectMutex.Lock()

	// Es wird geprüft wie der Aktuele Status des Scriptes ist
	if o.vmState != static.Starting {
		// Der Mutext wird freigegeben
		o.objectMutex.Unlock()

		// Rückgabe
		return
	}

	// Der Status wird auf Running gesetzt
	o.vmState = static.Running

	// Der Mutext wird freigegeben
	o.objectMutex.Unlock()

	// Log
	o.LogPrint("", "Eventloop started")

	// Die Schleife wird solange ausgeführt, solange der Status, running ist.
	// Die Schleife für den Eventloop des Kernels auf
	for !o.Kernel.HasCloseSignal() {
		if err := o.Kernel.ServeEventLoop(); err != nil {
			panic(err)
		}
	}

	// Der Status wird geupdated
	o.objectMutex.Lock()
	o.vmState = static.Closed
	o.objectMutex.Unlock()

	// Log
	o.LogPrint("", "Eventloop stoped")

	// Es wird Signalisiert das die VM nicht mehr ausgeführt wird
	syncWaitGroup.Done()
}

func (o *CoreVM) serveGorutine(syncWaitGroup *sync.WaitGroup) error {
	// Es wird geprüft ob der Server bereits gestartet wurde
	if o.GetState() != static.StillWait && o.GetState() != static.Closed {
		return fmt.Errorf("serveGorutine: vm always running")
	}

	// Es wird der SyncWaitGroup Signalisiert dass eine weitere Routine ausgeführt wird
	syncWaitGroup.Add(1)

	// Die VM wird als Startend Markiert
	o.vmState = static.Starting

	// Es wird versucht den MainCode einzulesen
	mainCode := o.vmDbEntry.GetMainCodeFile()

	// Es wird versucht den Inhalt der Datei zu laden
	scriptContent, err := mainCode.GetContent()
	if err != nil {
		return fmt.Errorf("CoreVM->serveGorutine: " + err.Error())
	}

	// Diese Funktion wird als Goroutine ausgeführt
	go o._routine([]byte(scriptContent), syncWaitGroup)

	// Es ist kein Fehler aufgetreten
	return nil
}

func (o *CoreVM) GetState() types.VmState {
	return o.vmState
}

func (o *CoreVM) GetWhitelist() []*types.TransportWhitelistVmEntryData {
	returnList := make([]*types.TransportWhitelistVmEntryData, 0)
	for _, item := range o.vmDbEntry.GetWhitelist() {
		returnList = append(returnList, &types.TransportWhitelistVmEntryData{
			WildCardDomains: item.Endpoint.Domain.Wildcards,
			ExactDomains:    item.Endpoint.Domain.Exact,
			Methods:         item.Methods,
			IPv4List:        item.Endpoint.IPv4List,
			Ipv6List:        item.Endpoint.IPv6List,
		})
	}
	return returnList
}

func (o *CoreVM) GetRootMemberIDS() []*types.CAMemberData {
	ret := make([]*types.CAMemberData, 0)
	for _, item := range o.vmDbEntry.GetRootMemberIDS() {
		ret = append(ret, &types.CAMemberData{
			Fingerprint: item.Fingerprint,
			Type:        item.Type,
			ID:          item.ID,
		})
	}
	return ret
}

func (o *CoreVM) GetDatabaseServices() []*types.VMDatabaseData {
	vmdlist := make([]*types.VMDatabaseData, 0)
	for _, item := range o.vmDbEntry.GetAllDatabaseServices() {
		vmdlist = append(vmdlist, &types.VMDatabaseData{
			Type:     item.Type,
			Host:     item.Host,
			Port:     item.Port,
			Username: item.Username,
			Password: item.Password,
			Database: item.Database,
			Alias:    item.Alias,
		})
	}
	return vmdlist
}

func (o *CoreVM) ValidateRPCRequestSource(soruce string) bool {
	// Es wird geprüft ob es es einen Global Wildcard eintrag gib
	if _, hasWildCard := o.vmDbEntry.GetAllowedHttpSources()["*"]; hasWildCard {
		return true
	}

	// Es wird geprüft ob es für diesen Host einen Eintrag gibt
	if _, checkresult := o.vmDbEntry.GetAllowedHttpSources()[strings.ToLower(soruce)]; checkresult {
		return true
	}

	// Es wird eine mögliche Whitelist erstellt
	whitelist := make([]string, 0)
	for item := range o.vmDbEntry.GetAllowedHttpSources() {
		whitelist = append(whitelist, item)
	}

	// Es wird geprüft ob die Quelldomain sich durch die Whitelist bestätigen lässt, das ergebniss wird zurückgegeben
	return utils.CheckHostInWhitelist(strings.ToLower(soruce), whitelist)
}

func (o *CoreVM) GetConsoleOutputWatcher() types.WatcherInterface {
	return o.Kernel.Console().GetOutputStream()
}

func (o *CoreVM) addDatabaseServiceLink(dbserviceLink services.DbServiceLinkinterface) error {
	o.dbServiceLinks = append(o.dbServiceLinks, dbserviceLink)
	return nil
}

func (o *CoreVM) GetStartingTimestamp() uint64 {
	return o.startTimeUnix
}

func (o *CoreVM) runScript(script string) error {
	// Es wird geprüft ob das Script beretis geladen wurden
	if o.scriptLoaded {
		return fmt.Errorf("LoadScript: always script loaded")
	}

	// Es wird markiert dass das Script geladen wurde
	o.scriptLoaded = true

	// Das Script wird ausgeführt
	_, err := o.Kernel.RunScript(script, "main.js")
	if err != nil {
		panic(err)
	}

	// Es ist kein Fehler aufgetreten
	return nil
}

func (o *CoreVM) GetAllSharedFunctions() []types.SharedFunctionInterface {
	extracted := make([]types.SharedFunctionInterface, 0)
	table := o.GloablRegisterRead("rpc")
	if table == nil {
		return extracted
	}

	ctable, istable := table.(map[string]types.SharedFunctionInterface)
	if !istable {
		return extracted
	}

	for _, item := range ctable {
		extracted = append(extracted, item)
	}

	return extracted
}

func (o *CoreVM) GetSharedFunctionBySignature(sourceType types.RPCCallSource, funcSignature *types.FunctionSignature) (types.SharedFunctionInterface, bool, error) {
	// Es wird versucht die RPC Tabelle zu lesen
	var table interface{}
	if sourceType == static.LOCAL {
		table = o.GloablRegisterRead("rpc")
	} else {
		table = o.GloablRegisterRead("rpc_public")
	}

	// Es wird ermittelt ob die Tabelle gefunden wurde
	if table == nil {
		return nil, false, fmt.Errorf("rpc register reading error")
	}

	// Es wird versucht die Tabelle richtig einzulesen
	ctable, istable := table.(map[string]types.SharedFunctionInterface)
	if !istable {
		return nil, false, fmt.Errorf("rpc register reading error")
	}

	// Es wird geprüft ob in der Tabelle eine Eintrag für die Funktion vorhanden ist
	result, fResult := ctable[utils.FunctionOnlySignatureString(funcSignature)]
	if !fResult {
		return nil, false, nil
	}

	// Das Ergebniss wird zurückgegeben
	return result, true, nil
}

func newCoreVM(core *Core, vmDb *vmdb.VmDBEntry, extModules []*extmodules.ExternalModule, loggingPath types.LOG_DIR) (*CoreVM, error) {
	// Es wird ein neuer Konsolen Stream erzeugt
	consoleStream, err := consolecache.NewConsoleOutputCache(string(loggingPath))
	if err != nil {
		return nil, fmt.Errorf("CoreVM->newCoreVM: " + err.Error())
	}

	// Es werden alle Externen Module geladen
	extMods := make([]types.KernelModuleInterface, 0)
	for _, item := range extModules {
		// Es wird versucht das Modul zu bauen
		extMod, err := kernel.LinkWithExternalModule(item)
		if err != nil {
			return nil, fmt.Errorf("newCoreVM: " + err.Error())
		}

		// Die Daten werden abgespeichert
		extMods = append(extMods, extMod)
	}

	// Die KernelModule werden Initalisiert
	var kernelConfig *kernel.KernelConfig
	if len(extMods) > 0 {
		kernelConfig = kernel.NewFromExist(&kernel.DEFAULT_CONFIG, extMods...)
	} else {
		kernelConfig = &kernel.DEFAULT_CONFIG
	}

	// Es wird ein neuer Kernel erzeugt
	vmKernel, err := kernel.NewKernel(consoleStream, kernelConfig, vmDb, core)
	if err != nil {
		return nil, fmt.Errorf("newCoreVM: " + err.Error())
	}

	// Das Core Objekt wird erstellt
	coreObject := &CoreVM{
		Kernel:          vmKernel,
		core:            core,
		vmDbEntry:       vmDb,
		externalModules: extModules,
		objectMutex:     &sync.Mutex{},
		vmState:         static.StillWait,
		dbServiceLinks:  make([]services.DbServiceLinkinterface, 0),
	}

	// Es wird versucht die VM mit dem Kernel zu verlinken
	if err := vmKernel.LinkKernelWithCoreVM(coreObject); err != nil {
		return nil, fmt.Errorf("newCoreVM: " + err.Error())
	}

	// Das Objekt wird zurückgegeben
	return coreObject, nil
}
