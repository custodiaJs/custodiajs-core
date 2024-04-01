package webservice

import "vnh1/types"

func (o *RpcRequest) GetParms() []types.FunctionParameterBundleInterface {
	return o.parms
}
