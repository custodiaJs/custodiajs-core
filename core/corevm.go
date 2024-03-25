package core

import (
	"vnh1/core/jsvm"
	"vnh1/core/vmdb"
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
