package core

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"sync"

	"github.com/CustodiaJS/custodiajs-core/container"
	"github.com/CustodiaJS/custodiajs-core/databaseservices"
	"github.com/CustodiaJS/custodiajs-core/filesystem"
	"github.com/CustodiaJS/custodiajs-core/identkeydatabase"
	"github.com/CustodiaJS/custodiajs-core/ipnetwork"
	"github.com/CustodiaJS/custodiajs-core/saftychan"
	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/static/errormsgs"
	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/utils"
	"github.com/CustodiaJS/custodiajs-core/utils/grsbool"
	"github.com/CustodiaJS/custodiajs-core/utils/procslog"
	"github.com/CustodiaJS/custodiajs-core/vm"
	"github.com/CustodiaJS/custodiajs-core/vmdb"
)

/* Fügt eine Externe Modul Lib dem Core hinzu
func (o *Core) AddExternalModuleLibrary(modLib *external_module.ExternalModule) error {
	// Es wird geprüft ob es sich um einen Zulässigen Module namen handelt
	if val := utils.ValidateExternalModuleName(modLib.GetName()); !val {
		return fmt.Errorf("Core->AddExternalModuleLibrary: Invalid module name, cant added module")
	}

	// Der Mutex wird angewendet
	o.objectMutex.Lock()

	// Es wird ermittelt ob es bereits ein Externes Module mit dem gleichen Namen gibt
	if _, found := o.extModules[modLib.GetName()]; found {
		o.objectMutex.Unlock()
		return fmt.Errorf("Core->AddExternalModuleLibrary: module always added")
	}

	// Die Module Lib wird zwischengspeichert
	o.extModules[modLib.GetName()] = modLib

	// Der Mutex wird freigegeben
	o.objectMutex.Unlock()

	// Es ist kein Fehler aufgetreten
	return nil
}
*/

// Erstellt eine neue VM Instanz
func (o *Core) AddNewVMInstance(vmDbEntry *vmdb.VmDBEntry) (types.VmInterface, error) {
	// Die Virtuelle Maschine wird geprüft
	if !vmDbEntry.ValidateVM() {
		return nil, fmt.Errorf("AddNewVMInstance: Broken Virtual Machine")
	}

	// Es werden alle benötigten Host CA Extrahiert
	// sollte kein Passendes Gefunden werden, wird der Vorgang abgebrochen
	for _, item := range vmDbEntry.GetRootMemberIDS() {
		// Es wird ermittelt ob es sich um ein SSL Eintrag handelt, wenn nicht wird dieser Ignoriert
		if item.Type != "ssl" {
			continue
		}

		// Es wird ermittelt ob der Host ein SSL-Cert/Privatekey paar besitzt welches von dem Aktuellen RootCA Signiert wurde
		// wenn nicht wird geprüft ob der Fingerabdruck des CERTS mit dem des Lokalen Certs übereinstimmt
		if !o.hostIdentKeyDatabase.ValidateRootCAMembershipByFingerprint(item.Fingerprint) {
			if !strings.EqualFold(hex.EncodeToString(utils.ComputeTlsCertFingerprint(o.hostTlsCert)), item.Fingerprint) {
				return nil, fmt.Errorf("Core->AddNewVMInstance: unkown host ca membership '%s'", strings.ToUpper(item.Fingerprint))
			}
		}
	}

	// Es wird eine Liste mit allen Benötigten externen Libs abgerufen
	neededExternalModulesNameSlice := make([]string, 0)
	for _, item := range vmDbEntry.GetAllExternalServices() {
		neededExternalModulesNameSlice = append(neededExternalModulesNameSlice, item.Name)
	}

	/*
		// Es werden alle Externen Module herausgefiltertet
		modList := make([]*external_modules.ExternalModule, 0)
		for _, item := range neededExternalModulesNameSlice {
			for _, mitem := range o.extModules {
				if item == mitem.GetName() {
					modList = append(modList, mitem)
				}
			}
		}

		// Es wird geprüft ob die benötigten Module gefunden wurden
		notFoundExtModules := make([]string, 0)
		for _, item := range vmDbEntry.GetAllExternalServices() {
			if item.Required {
				found := false
				for _, xtem := range modList {
					if xtem.GetName() == item.Name {
						if xtem.GetVersion() >= uint64(item.MinVersion) {
							found = true
							break
						}
					}
				}
				if !found {
					notFoundExtModules = append(notFoundExtModules, item.Name)
				}
			}
		}


		// Es wird ein Fehler ausgelöst wenn ein benötigtes Modul nicht gefunden wurde
		if len(notFoundExtModules) != 0 {
			return nil, fmt.Errorf("Core->AddNewVMInstance: external modules '%s' not found", strings.Join(neededExternalModulesNameSlice, ","))
		}
	*/

	// Das Logging Verzeichniss wird erstellt
	logPath, err := filesystem.MakeLogDirForVM(o.logDIR, vmDbEntry.GetVMName())
	if err != nil {
		return nil, fmt.Errorf("Core->AddNewVMInstance: " + err.Error())
	}

	// Der Mutex wird angewendet
	o.objectMutex.Lock()

	// Es wird geprüft ob bereits eiein VM Link hinzugefügtne VM mit der Selben ID vorhanden ist
	if _, foundVM := o.vmsByID[strings.ToLower(vmDbEntry.GetVMContainerMerkleHash())]; foundVM {
		o.objectMutex.Unlock()
		return nil, fmt.Errorf("Core->AddNewVMInstance: You cannot add a VM container '%s' multiple times", vmDbEntry.GetVMContainerMerkleHash())
	}

	// WARNUNG!
	// Handelt es sich um ein Linux-System, werden die CustodiaJS-Container vorbereitet.
	// Alle Container und ihre Konfigurationen werden aus der Config-Datenbank ausgelesen.
	// Handelt es sich um ein Windows- oder macOS-System, werden die Container mithilfe von Firewall und Benutzerrechten umgesetzt
	var vmContainer *container.VmContainer
	var vmContainerErr error
	runAsProcess := false
	switch runtime.GOOS {
	case "linux":
		vmContainer, vmContainerErr = container.NewLinuxContainer()
	case "windows":
		// Es wird geprüft ob Docker oder WSL auf Windows Installiert ist,
		// sollte WSL oder Docker Installiert sein, wird der ein neuer Docker/WSL Container erzeugt
		if container.CheckWindowsHasDockerOrWSL() {
			vmContainer, vmContainerErr = container.NewWindowsDockerWSLContainer()
		} else {
			runAsProcess = true
		}
	case "darwin":
		vmContainer, vmContainerErr = container.NewMacOSContainer()
	case "freebsd":
		vmContainer, vmContainerErr = container.NewFreeBSDContainer()
	case "openbsd":
		vmContainer, vmContainerErr = container.NewOpenBSDContainer()
	case "netbsd":
		vmContainer, vmContainerErr = container.NewNetBSDContainer()
	}

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if vmContainerErr != nil {
		return nil, fmt.Errorf("Core->AddNewVMInstance: " + vmContainerErr.Error())
	}

	// Der Mutex wird freigegeben
	o.objectMutex.Unlock()

	// Die VM wird erzeugt, entwender wird eine InProcess VM erzeugt, ein VM Prozess oder eine Container VM Prozess.
	// Es wird vorher geprüft ob der Tag 'ForceInProcessVM' gesetzt wurde, wenn ja werden ausschlißlich InProcess VM's verwendet.
	var vmInstance types.VmInterface
	var vmInstanceErr error
	if !static.FORCE_INPROCESS_VM {
		if vmContainer != nil {
			panic("not implemented")
		} else if runAsProcess {
			panic("not implemented")
		} else {
			vmInstance, vmInstanceErr = vm.NewCoreVM(o, vmDbEntry, logPath)
		}
	} else {
		vmInstance, vmInstanceErr = vm.NewCoreVM(o, vmDbEntry, logPath)
	}

	// Es wird geprüft ob ein Fehler beim erstellen der VM aufgetreten ist
	if vmInstanceErr != nil {
		return nil, fmt.Errorf("AddNewVMInstance: " + vmInstanceErr.Error())
	}

	// Der Mutex wird angewendet
	o.objectMutex.Lock()

	// Das VMObjekt wird zwischengespeichert
	o.vmsByID[strings.ToLower(vmDbEntry.GetVMContainerMerkleHash())] = vmInstance // Merklehash
	o.vmsByName[strings.ToLower(vmDbEntry.GetVMName())] = vmInstance              // VM-Name
	o.vmKernelPtr[vmInstance.GetKId()] = vmInstance                               // Speichert die VM ab, diese wird verwendet um die VM durch den Kernel der VM auffindbar zu machen
	o.vms = append(o.vms, vmInstance)                                             // Die VM wird abgespeichert

	// Der Mutex wird freigegeben
	o.objectMutex.Unlock()

	// Die VM wird mit allen Datenbankdiensten Verknüpft
	for _, item := range vmDbEntry.GetAllDatabaseServices() {
		// Es wird ein neuer Link für die VM erzeugt
		link, err := o.databaseService.GetDBServiceLink(item.GetDatabaseFingerprint())
		if err != nil {
			return nil, fmt.Errorf("Core->AddNewVMInstance: " + err.Error())
		}

		// Der Link für den Datenbank Dienst wird abgespeichert
		if err := vmInstance.AddDatabaseServiceLink(link); err != nil {
			return nil, fmt.Errorf("Core->AddNewVMInstance: " + err.Error())
		}
	}

	// Das VM Objekt wird zwischengespeichert
	return vmInstance, nil
}

// Fügt einen API Socket hinzu
func (o *Core) AddAPISocket(apiSocket types.APISocketInterface) error {
	// Es wird geprüft das kein Null Wert übergeben wurde
	if apiSocket == nil {
		return fmt.Errorf("Core->AddAPISocket: null api socket not allowed")
	}

	// Der Core wird in dem  Registriert
	err := apiSocket.SetupCore(o)
	if err != nil {
		return fmt.Errorf("AddAPISocket: ")
	}

	// Der Mutex wird angewendet
	// und nach beenden der Funktion freigegeben
	o.objectMutex.Lock()
	defer o.objectMutex.Unlock()

	// Der API Socket wird zwischengespeichert
	o.apiSockets = append(o.apiSockets, apiSocket)

	// Es ist kein Fehler aufgetreten
	return nil
}

// Gibt eine Spezifisichen Container VM anhand ihrer ID zurück
func (o *Core) GetScriptContainerVMByID(vmid string) (types.VmInterface, bool, *types.SpecificError) {
	// Es wird geprüft ob es sich um einen zulässigen vm Namen handelt
	if !utils.ValidateVMIdString(vmid) {
		return nil, false, nil //fmt.Errorf("Core->GetScriptContainerVMByID: invalid vm container id")
	}

	// Die ID wird lowercast
	lowerCaseId := strings.ToLower(vmid)

	// Der Mutex wird angewendet
	// und nach beenden der Funktion freigegeben
	o.objectMutex.Lock()
	defer o.objectMutex.Unlock()

	// Es wird geprüft ob die VM exestiert
	vmObj, found := o.vmsByID[lowerCaseId]
	if !found {
		return nil, false, nil
	}

	// Das Objekt wird zurückgegeben
	return vmObj, true, nil
}

// Gibt eine Spezifisichen Container VM anhand ihrer ID zurück
func (o *Core) GetScriptContainerByVMName(vmName string) (types.VmInterface, bool, *types.SpecificError) {
	// Es wird geprüft ob es sich um einen zulässigen Namen handelt
	if !utils.ValidateVMName(vmName) {
		return nil, false, nil //fmt.Errorf("Core->GetScriptContainerByVMName: invalid vm container name")
	}

	// Der Name wird lowercast
	lowerCaseVmName := strings.ToLower(vmName)

	// Der Mutex wird angewendet
	// und nach beenden der Funktion freigegeben
	o.objectMutex.Lock()
	defer o.objectMutex.Unlock()

	// Es wird geprüft ob die VM exestiert
	vmObj, found := o.vmsByName[lowerCaseVmName]
	if !found {
		return nil, false, nil // fmt.Errorf("Core->GetScriptContainerByVMName: unkown vm '%s'", lowerCaseVmName)
	}

	// Das Objekt wird zurückgegeben
	return vmObj, true, nil
}

// Gibt die ID's der Aktiven VM-Container zurück
func (o *Core) GetAllActiveScriptContainerIDs() []string {
	// Der Mutex wird angewendet
	// und nach beenden der Funktion freigegeben
	o.objectMutex.Lock()
	defer o.objectMutex.Unlock()

	// Es wird eine Liste mit allen Verfügbaren VM's erstellt
	extr := make([]string, 0)
	for _, item := range o.vmsByID {
		extr = append(extr, string(item.GetFingerprint()))
	}

	// Die Liste wird zurückgegeben
	return extr
}

// Gibt alle VM-Container zurück
func (o *Core) GetAllVMs() []types.VmInterface {
	// Der Mutex wird angewendet
	// und nach beenden der Funktion freigegeben
	o.objectMutex.Lock()
	defer o.objectMutex.Unlock()

	// Es wird eine Liste mit allen Verfügbaren VM-Containern erstellt
	extr := make([]types.VmInterface, 0)
	for _, item := range o.vmsByID {
		extr = append(extr, item)
	}

	// Die Liste wird zurückgegeben
	return extr
}

// Signalisiert dem Core, dass er beendet werden soll
func (o *Core) SignalShutdown() {
	// Der Mutex wird angewendet
	o.objectMutex.Lock()
	defer o.objectMutex.Unlock()

	// Die Chan wird geschlossen
	close(o.holdOpenChan)
}

// Gibt die Prozess Managment Unit zurück
func (o *Core) GetCoreSessionManagmentUnit() types.CoreSessionManagmentUnitInterface {
	return o.cpmu
}

// Erzeugt aus einer IP-Adresse eine Verified Core IP Address (LRSAP)
func (o *Core) ConvertLagacyIPAddressToLRSAP(lagacyRemoteIPAddress string, lagacyLocalIPAddress string) (types.VerifiedCoreIPAddressInterface, *types.SpecificError) {
	// Die Quell IP-Adresse wird eingelesen
	sourceIp, errSourceIp := o.hostnetmanager.TryParseIp(lagacyRemoteIPAddress)
	if errSourceIp != nil {
		// Die Aktuelle Funktion wird dem Fehler hinzugefügt
		errSourceIp.AddCallerFunctionToHistory("Core->ConvertLagacyIPAddressToLRSAP")

		// Der Fehler wird zurückgegeben
		return nil, errSourceIp
	}

	// Die Ziel IP-Adresse wird eingelesen
	destinationIp, errDestIp := o.hostnetmanager.TryParseIp(lagacyLocalIPAddress)
	if errDestIp != nil {
		// Die Aktuelle Funktion wird dem Fehler hinzugefügt
		errDestIp.AddCallerFunctionToHistory("Core->ConvertLagacyIPAddressToLRSAP")

		// Der Fehler wird zurückgegeben
		return nil, errDestIp
	}

	// Es wird versucht anhand der IP-Adresse das Zugehörige Lokale Netzwerk Interface zu ermitteln
	localNetIfaceByDestinationIp := o.hostnetmanager.GetNetworkInterfaceByLocalIp(destinationIp)
	if localNetIfaceByDestinationIp != nil {
		return nil, errormsgs.CORE_CANT_RETRIVE_NETWORK_INTERFAC_BY_IP("Core->ConvertLagacyIPAddressToLRSAP", destinationIp)
	}

	// Es wird ein neues CoreLRSAP Struct erzeugt
	coreLRSAP := &CoreLRSAP{
		SourceAddress:           sourceIp,
		LocalAddress:            destinationIp,
		ConnectionOverInterface: localNetIfaceByDestinationIp,
	}

	// Das Objekt wird zurückgegeben, ohen Fehler
	return coreLRSAP, nil
}

// Überprüft ob die Quelle zulässig ist
func (o *Core) LRSAPSourceIsAllowed(LRSAP types.VerifiedCoreIPAddressInterface) bool {
	return false
}

// Signalisiert allen VM's dass sie beendet werden
func (o *Core) signalVmsShutdown(wg *sync.WaitGroup) {
	for _, item := range o.vms {
		wg.Add(1)
		go func(cvm types.VmInterface) {
			cvm.SignalShutdown()
			wg.Done()
		}(item)
	}
}

// Erzeugt eine neue Web bases RPC Sitzung
func (o *CoreSessionManagmentUnit) NewWebRequestBasedRPCSession(httpRequest *http.Request) (types.WebRequestBasedRPCSessionInterface, *types.SpecificError) {
	// Die Basis Variabeln werden erzeugt
	proc, isConnected, saftyResponseChan := procslog.NewProcLogSession(), grsbool.NewGrsbool(true), saftychan.NewFunctionCallReturnChan()

	// Es wird eine Goroutine gestartet, welche prüft ob die Verbindung getrennt wurde
	go func() {
		// Es wird darauf gewartet dass die Verbindung geschlossen wird
		<-httpRequest.Context().Done()

		// Es wird Signalisiert dass die Verbindung geschlossen wurde
		isConnected.Set(false)

		// Das SaftyChan wird geschlossen
		saftyResponseChan.Close()
	}()

	// Es wird ein neues Rückgabe Objekt erstellt
	returnObject := &CoreWebRequestRPCSession{saftyResponseChan: saftyResponseChan, proc: proc, isConnected: isConnected}

	// Das CoreWebRequestRPCSession Objekt wird ohne Fehler zurückgegeben
	return returnObject, nil
}

// Legt den Core Status fest
func setState(core *Core, state types.CoreState, useMutex bool) {
	// Es wird geprüft ob Mutex verwendet werden sollen
	if useMutex {
		core.objectMutex.Lock()
		defer core.objectMutex.Unlock()
	}

	// Es wird geprüft ob der neue Status, der Aktuelle ist
	if core.state == state {
		return
	}

	// Der Neue Status wird gesetzt
	core.state = state
}

// Erstellt einen neuen CustodiaJS Core
func NewCore(hostTlsCert *tls.Certificate, hostIdenKeyDatabase *identkeydatabase.IdenKeyDatabase, dbService *databaseservices.DbService, logDIRPath types.LOG_DIR, ipnet *ipnetwork.HostNetworkManagmentUnit) (*Core, error) {
	// Das Coreobjekt wird erstellt
	coreObj := &Core{
		vmsByID:         make(map[string]types.VmInterface),
		vmsByName:       make(map[string]types.VmInterface),
		vmKernelPtr:     make(map[types.KernelID]types.VmInterface),
		vms:             make([]types.VmInterface, 0),
		apiSockets:      make([]types.APISocketInterface, 0),
		hostTlsCert:     hostTlsCert,
		databaseService: dbService,
		state:           static.NEW,
		//extModules:      make(map[string]*external_modules.ExternalModule),
		// Chans
		holdOpenChan:     make(chan struct{}),
		serviceSignaling: make(chan struct{}),
		vmSyncWaitGroup:  sync.WaitGroup{},
		apiSyncWaitGroup: sync.WaitGroup{},
		// Datenbanken
		hostIdentKeyDatabase: hostIdenKeyDatabase,
		// Mutexes
		objectMutex: &sync.Mutex{},
		// Log
		logDIR: logDIRPath,
		// IP-Info Einheit
		hostnetmanager: ipnet,
	}

	// Das Objekt wird zurückgegeben
	return coreObj, nil
}
