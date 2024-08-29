package errormsgs

import (
	"github.com/CustodiaJS/custodiajs-core/core/ipnetwork"
	"github.com/CustodiaJS/custodiajs-core/global/types"
)

func CORE_INVALID_IP_ADDRESS_READING_ERROR(funcname string, ipadr string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func CORE_INVALID_IP_ADDRESS_UNKOWN_ADDRESS_FAMILY_REMOTE(funcname string, ipadr string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func CORE_INVALID_IP_ADDRESS_UNKOWN_ADDRESS_FAMILY_LOCAL(funcname string, ipadr string) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func CORE_CANT_RETRIVE_NETWORK_INTERFAC_BY_IP(funcname string, wlocalip *ipnetwork.IpAddress) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}

func CORE_INVALID_ADRESS_DOSENT_LOCALHOST_IPADR_ERROR(funcname string, wlocalip *ipnetwork.IpAddress) *types.SpecificError {
	return &types.SpecificError{CliError: nil, LocalJSVMError: nil, GoProcessError: nil, LocalApiOrRpcError: types.ApiError{}, RemoteApiOrRpcError: types.ApiError{}, History: []string{funcname}}
}
