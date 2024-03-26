package core

import (
	"vnh1/core/jsvm"
	"vnh1/core/vmdb"
	"vnh1/static"
)

type CoreVM struct {
	*jsvm.JsVM
	vmDbEntry      *vmdb.VmDBEntry
	jsMainFilePath string
	jsCode         string
}

func (o *CoreVM) GetVMName() string {
	return o.vmDbEntry.GetVMName()
}

func (o *CoreVM) GetFingerprint() string {
	return o.vmDbEntry.GetVMContainerMerkleHash()
}

func (o *CoreVM) GetVMModuleNames() []string {
	modNames := make([]string, 0)
	for _, item := range o.vmDbEntry.GetNodeJsModules() {
		modNames = append(modNames, item.GetName())
	}
	return modNames
}

func (o *CoreVM) GetLocalShareddFunctions() []static.SharedLocalFunctionInterface {
	return o.JsVM.GetLocalShareddFunctions()
}
