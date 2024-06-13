package kmodulerpc

import "vnh1/types"

func makeV8Error(errorByFunction string, Error error) *types.SpecificError {
	return &types.SpecificError{}
}

func makeUnkownMethodeError(errorByFunction string) *types.SpecificError {
	return &types.SpecificError{}
}

func makeRequestTypeIsNotHttpRequest(errorByFunction string) *types.SpecificError {
	return &types.SpecificError{}
}
