package core

import (
	"fmt"
	"strings"
	"sync"
	"vnh1/types"
	"vnh1/utils"
)

func (o *CoreVM) GetVMName() string {
	return o.vmDbEntry.GetVMName()
}

func (o *CoreVM) GetFingerprint() string {
	return o.vmDbEntry.GetVMContainerMerkleHash()
}

func (o *CoreVM) GetVMModuleNames() []string {
	if o.vmDbEntry.GetTotalNodeJsModules() < 1 {
		return make([]string, 0)
	}
	modNames := make([]string, 0)
	for _, item := range o.vmDbEntry.GetNodeJsModules() {
		modNames = append(modNames, item.GetName())
	}
	return modNames
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

func (o *CoreVM) GetWhitelist() []types.TransportWhitelistVmEntryInterface {
	returnList := make([]types.TransportWhitelistVmEntryInterface, 0)
	for _, item := range o.vmDbEntry.GetWhitelist() {
		returnList = append(returnList, &TransportWhitelistVmEntry{url: item.URL, alias: item.Alias})
	}
	return returnList
}

func (o *CoreVM) GetMemberCertKeyIds() []string {
	return o.vmDbEntry.GetMemberCertKeyIds()
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
