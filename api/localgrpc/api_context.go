package localgrpc

import (
	"github.com/CustodiaJS/custodiajs-core/global/types"
	"github.com/CustodiaJS/custodiajs-core/localgrpcproto"
	"github.com/google/uuid"
)

func (o *APIContext) SetType(tpe localgrpcproto.ClientType) {
	o.Log.Debug("Set Type = ", string(tpe))
	o.tpe = tpe
}

func (o *APIContext) CreateVmInstance(manifest *types.Manifest, scriptHash types.VmScriptHash, kid types.KernelID, puuid types.ProcessId) (*APIProcessVm, string, error) {
	uuid := uuid.New().String()
	procvm := &APIProcessVm{manifest: manifest, scriptHash: scriptHash, kid: kid, context: o}
	o.openvm = procvm
	o.Log.Debug("New Process VM created '%s'", uuid)
	return procvm, uuid, nil
}
