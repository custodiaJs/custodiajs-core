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
	vmState        types.VmState
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

func (o *CoreVM) serveGorutine(syncWaitGroup *sync.WaitGroup) error {
	// Es wird geprüft ob der Server bereits gestartet wurde
	if o.GetState() != types.StillWait && o.GetState() != types.Closed {
		return fmt.Errorf("serveGorutine: vm always running")
	}

	// Es wird der SyncWaitGroup Signalisiert dass eine weitere Routine ausgeführt wird
	syncWaitGroup.Add(1)

	// Die VM wird als Startend Markiert
	o.vmState = types.Starting

	// Diese Funktion wird als Goroutine ausgeführt
	go func(item *CoreVM) {
		o.vmState = types.Running
		item.RunScript(item.jsCode)
		syncWaitGroup.Done()
	}(o)

	// Es ist kein Fehler aufgetreten
	return nil
}

func (o *CoreVM) GetState() types.VmState {
	return o.vmState
}

func (o *CoreVM) GetConsoleOutputWatcher() types.WatcherInterface {
	return o.JsVM.GetConsoleOutputWatcher()
}
