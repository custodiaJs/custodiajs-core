package errormsgs

import "github.com/CustodiaJS/custodiajs-core/types"

func VM_GET_FUNCTION_BY_SIGNATURE_TABLE_NULL_ERROR(funcname string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func VM_GET_FUNCTION_RPC_REIGSTER_ERROR(funcname string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}
