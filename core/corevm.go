package core

import (
	"fmt"
	"strings"
	"sync"
	"vnh1/core/databaseservices/services"
	"vnh1/core/jsvm"
	"vnh1/core/vmdb"
	"vnh1/types"
	"vnh1/utils"
)

func (o *CoreVM) GetVMName() string {
	return o.vmDbEntry.GetVMName()
}

func (o *CoreVM) GetFingerprint() types.CoreVMFingerprint {
	return types.CoreVMFingerprint(strings.ToLower(o.vmDbEntry.GetVMContainerMerkleHash()))
}

func (o *CoreVM) GetVMJSModules() []*types.VmNodeJsModuleDetails {
	if o.vmDbEntry.GetTotalNodeJsModules() < 1 {
		return make([]*types.VmNodeJsModuleDetails, 0)
	}
	modNames := make([]*types.VmNodeJsModuleDetails, 0)
	for _, item := range o.vmDbEntry.GetNodeJsModules() {
		modNames = append(modNames, &types.VmNodeJsModuleDetails{
			Name:  item.GetName(),
			Alias: item.GetAlias(),
		})
	}
	return modNames
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

func (o *CoreVM) serveGorutine(syncWaitGroup *sync.WaitGroup) error {
	// Es wird geprüft ob der Server bereits gestartet wurde
	if o.GetState() != types.StillWait && o.GetState() != types.Closed {
		return fmt.Errorf("serveGorutine: vm always running")
	}

	// Es wird der SyncWaitGroup Signalisiert dass eine weitere Routine ausgeführt wird
	syncWaitGroup.Add(1)

	// Die VM wird als Startend Markiert
	o.vmState = types.Starting

	// Es wird versucht den MainCode einzulesen
	mainCode := o.vmDbEntry.GetMainCodeFile()

	// Es wird versucht den Inhalt der Datei zu laden
	scriptContent, err := mainCode.GetContent()
	if err != nil {
		return fmt.Errorf("CoreVM->serveGorutine: " + err.Error())
	}

	// Diese Funktion wird als Goroutine ausgeführt
	go func(item *CoreVM, scriptContent []byte) {
		o.vmState = types.Running
		item.RunScript(string(scriptContent))
		syncWaitGroup.Done()
	}(o, scriptContent)

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

func (o *CoreVM) GetMemberCertsPkeys() []*types.CAMemberData {
	ret := make([]*types.CAMemberData, 0)
	for _, item := range o.vmDbEntry.GetMemberCertsPkeys() {
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
	return o.JsVM.GetConsoleOutputWatcher()
}

func (o *CoreVM) addDatabaseServiceLink(dbserviceLink services.DbServiceLinkinterface) error {
	o.dbServiceLinks = append(o.dbServiceLinks, dbserviceLink)
	return nil
}

func newCoreVM(jsvm *jsvm.JsVM, vmDb *vmdb.VmDBEntry) *CoreVM {
	return &CoreVM{JsVM: jsvm, vmDbEntry: vmDb, vmState: types.StillWait, dbServiceLinks: make([]services.DbServiceLinkinterface, 0)}
}
