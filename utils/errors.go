package utils

import "github.com/CustodiaJS/custodiajs-core/types"

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

func MakeConnectionIsClosed() *types.SpecificError {
	return &types.SpecificError{}
}

func MakeHttpConnectionIsClosed() *types.SpecificError {
	return &types.SpecificError{}
}

func MakeAlreadyAnsweredRPCRequestError() *types.SpecificError {
	return &types.SpecificError{}
}

func MakeHttpRequestIsClosedBeforeException() *types.SpecificError {
	return &types.SpecificError{}
}

func TimeoutFunctionCallError() *types.SpecificError {
	return &types.SpecificError{}
}

func V8ObjectWritingError() *types.SpecificError {
	return &types.SpecificError{}
}

func V8ObjectInstanceCreatingError() *types.SpecificError {
	return &types.SpecificError{}
}

func V8ValueConvertingError() *types.SpecificError {
	return &types.SpecificError{}
}

func RPCFunctionCallNullSharedFunctionObject() *types.SpecificError {
	return &types.SpecificError{}
}

func RPCFunctionCallNullRequest() *types.SpecificError {
	return &types.SpecificError{}
}

func MakeRPCFunctionCallParametersNumberUnequal(required uint, have uint) *types.SpecificError {
	return &types.SpecificError{}
}
