package errormsgs

import (
	"github.com/CustodiaJS/custodiajs-core/global/types"
)

func HTTP_API_RPC_VM_NOT_RUNNING(funcname string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_CORE_CONTEXT_EXTRACTION_ERROR(funcname string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_UNKOWN_METHODE(funcname string, hasMethode string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_INVALID_METHODE(funcname string, hasMethode string, needMethode string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_HAS_NO_TLS_ENCRYPTION(funcname string, requestSource string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_INVALID_CODEC(funcname string, codec string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_INVALID_QUERY_PARM_SIZE(funcname string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_INVALID_QUERY_PARMS(funcname string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_INVALID_QUERY_PARM_ID_SIZE(funcname string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_INVALID_QUERY_PARM_HEX_DECODING(funcname string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_INVALID_QUERY_PARM_VMNAME(funcname string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_INVALID_QUERY_PARM_HAS_NOT_ID_AND_NAME(funcname string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_REQUEST_HAS_UNKOWN_VM_IDENT_METHODE(funcname string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_REQUEST_VM_NOT_FOUND(funcname string, vmname string, rm types.RPCRequestVMIdentificationMethode) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_REQUEST_READING_ERROR(funcname string, errormsg string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_REQUEST_READING_CBOR_ERROR(funcname string, errormsg string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_REQUEST_READING_JSON_ERROR(funcname string, errormsg string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_REQUEST_READING_UNKOWN_ERROR(funcname string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_REQUEST_NOT_AUTHORIZED_X_SOURCE(funcname string, source *types.XRequestedWithData) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_REQUEST_INVALID_BODY_DATA(funcname string, source string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_REQUEST_INVALID_BODY_DATA_CALLED_FUNCTION_NAME(funcname string, tryCalledFunctionName string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_REQUEST_FUNCTION_CALL_PANIC(funcname string, panicMsg string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_REQUEST_INVALID_PARAMETER_TYPE_AT_POSITION_X(funcname string, x int, need_dtype string, has_stype string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_REQUEST_INVALID_PARAMETER_SLICE_SIZE(funcname string, x int, nx int) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_REQUEST_INVALID_RPC_FUNCTION_DATATYPES(funcname string, vals []*types.RPCParmeterReadingError) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_REQUEST_NOT_ALLOWED_SOURCE_IP(funcname string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_REQUEST_VM_FUNCTION_NOT_FOUND_ERROR(funcname string, vmname string, rm types.RPCRequestVMIdentificationMethode, functionsig *types.FunctionSignature) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_RESPONSE_WRITING_CONNECTION_CLOSED_ERROR(funcname string, datasize int) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_RESPONSE_WRITING_UNKOWN_ERROR(funcname string, gerr error) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_RESPONSE_WRITING_UNKOWN_ENCODING_ERROR(funcname string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_RESPONSE_WRITING_ENCODING_ERROR(funcname string, encoding types.HttpRequestContentType, err error) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_REQUEST_INTERNAL_ERROR_BY_READING_RETURN_CHAN(funcname string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_INTERNAL_ERROR_BY_EMITTING_CAPSLE_SIZE(funcname string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_REFERER_READING_ERROR(funcname string, refererurl string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_ORIGIN_READING_ERROR(funcname string, refererurl string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_LOCAL_IP_READING_ERROR(funcname string, ipadr string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func HTTP_API_SERVICE_REQUESR_HAS_NOT_ALLOWED_PROTOCOL_SHEME_ERROR(funcname string, used_sheme string, needed_sheme string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}
