package core

import (
	"fmt"
	"sync"
	"vnh1/core/jsvm"
	"vnh1/core/vmdb"
	"vnh1/types"
)

type CoreVM struct {
	*jsvm.JsVM
	vmDbEntry      *vmdb.VmDBEntry
	jsMainFilePath string
	jsCode         string
	state          types.VmState
}

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

func (o *CoreVM) GetState() types.VmState {
	return o.state
}

func (o *CoreVM) serveGorutine(syncWaitGroup *sync.WaitGroup) error {
	// Es wird geprüft ob der Server bereits gestartet wurde
	if o.state != types.StillWait && o.state != types.Closed {
		return fmt.Errorf("serveGorutine: vm always running")
	}

	// Es wird der SyncWaitGroup Signalisiert dass eine weitere Routine ausgeführt wird
	syncWaitGroup.Add(1)

	// Der Aktuelle Status wird festgelegt
	o.state = types.Starting

	// Diese Funktion wird als Goroutine ausgeführt
	go func(item *CoreVM) {
		o.state = types.Running
		item.RunScript(item.jsCode)
		syncWaitGroup.Done()
	}(o)

	// Es ist kein Fehler aufgetreten
	return nil
}

func (o *CoreVM) GetConsoleOutputWatcher() types.WatcherInterface {
	return o.JsVM.GetConsoleOutputWatcher()
}
