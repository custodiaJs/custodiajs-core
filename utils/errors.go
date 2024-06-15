package utils

import "vnh1/types"

func MakeV8Error(errorByFunction string, Error error) *types.SpecificError {
	return &types.SpecificError{}
}

func MakeUnkownMethodeError(errorByFunction string) *types.SpecificError {
	return &types.SpecificError{}
}

func MakeRequestTypeIsNotHttpRequest(errorByFunction string) *types.SpecificError {
	return &types.SpecificError{}
}

func MakeHttpRequestIsClosed() *types.SpecificError {
	return &types.SpecificError{}
}

func MakeHttpConnectionIsClosed() *types.SpecificError {
	return &types.SpecificError{}
}

func MakeAlreadyAnsweredRPCRequestError() *types.SpecificError {
	return &types.SpecificError{}
}
