package errormsgs

import (
	"github.com/CustodiaJS/custodiajs-core/types"
)

func KMDOULE_RPC_SHARED_FUNCTION_FUNCTION_CALL_EVENT_LOOP_ENTERING_ERROR(funcname string, pos uint, err error) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func KMDOULE_RPC_SHARED_FUNCTION_FUNCTION_REQUEST_CONVERTING_PARM_INVALID_DTYPE_ERROR(funcname string, hight int, need string, has string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func KMDOULE_RPC_SHARED_FUNCTION_REQUEST_CONVERTING_ERROR(funcname string, hight int, reason error, hasdtype string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func KMDOULE_RPC_SHARED_FUNCTION_REQUEST_UNKOWN_DATATYPE(funcname string, hight int) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}
