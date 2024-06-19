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

func MakeConnectionIsClosedError(cfuncname string) *types.SpecificError {
	return &types.SpecificError{}
}

func MakeHttpConnectionIsClosedError() *types.SpecificError {
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

func MakeSharedFunctionRequestContextObjectError(cfuncname string) *types.SpecificError {
	return &types.SpecificError{}
}

func MakeSharedFunctionCallStateError(cfuncname string) *types.SpecificError {
	return &types.SpecificError{}
}

func MakeV8MissingParameters(cfuncname string, need uint, have int) *types.SpecificError {
	return &types.SpecificError{}
}

func MakeV8InvalidParameterDatatype(cfuncname string, hight uint, want string) *types.SpecificError {
	return &types.SpecificError{}
}

func MakeV8PromiseCreatingError(cfuncname string, eError error) *types.SpecificError {
	return &types.SpecificError{}
}

func MakeV8ConvertValueToStringError(cfuncname string, eError error) *types.SpecificError {
	return &types.SpecificError{}
}

func MakeNewRPCSharedFunctionContextKernelIsNullError(cfuncname string) *types.SpecificError {
	return &types.SpecificError{}
}

func MakeNewRPCSharedFunctionContextReturnDatatypeStringIsInvalidError(cfuncname string, hasType string) *types.SpecificError {
	return &types.SpecificError{}
}

func MakeNewRPCSharedFunctionContextRPCRequestIsNullError(cfuncname string) *types.SpecificError {
	return &types.SpecificError{}
}

func MakeNewRPCSharedFunctionInvalidContextObjectError(cfuncname string) *types.SpecificError {
	return &types.SpecificError{}
}

func MakeNewRPCSharedFunctionNewV8ObjectInstanceError(cfuncname string, err error) *types.SpecificError {
	return &types.SpecificError{}
}

func MakeRPCRequestAlwaysResponsedError(cfuncname string) *types.SpecificError {
	return &types.SpecificError{}
}

func MakeV8FunctionCallbackInfoIsNullError(cfuncname string) *types.SpecificError {
	return &types.SpecificError{}
}

func MakeV8ToGoConvertingError(cfuncname string) *types.SpecificError {
	return &types.SpecificError{}
}

func MakeRPCRequestContextIsClosedAndDestroyed(cfuncname string) *types.SpecificError {
	return &types.SpecificError{}
}

func MakeRPCResolvingDataError(cfuncname string) *types.SpecificError {
	return &types.SpecificError{}
}
