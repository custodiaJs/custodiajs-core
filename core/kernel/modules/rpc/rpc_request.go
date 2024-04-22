package kmodulerpc

import "vnh1/types"

func (o *RpcRequest) GetParms() []*types.FunctionParameterCapsle {
	return o.parms
}
