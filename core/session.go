package core

import "github.com/CustodiaJS/custodiajs-core/types"

func (o *CoreWebRequestRPCSession) GetProcLogSession() types.ProcessLogSessionInterface {
	return o.proc
}

func (o *CoreWebRequestRPCSession) IsConnected() bool {
	return o.isConnected.Bool()
}

func (o *CoreWebRequestRPCSession) GetReturnChan() types.FunctionCallReturnChanInterface {
	return o.saftyResponseChan
}

func (o *CoreWebRequestRPCSession) SignalsThatAnErrorHasOccurredWhenTheErrorIsSent(size int, error *types.SpecificError) {

}

func (o *CoreWebRequestRPCSession) SignalThatTheErrorWasSuccessfullyTransmitted(size int) {

}

func (o *CoreWebRequestRPCSession) SignalTheErrorSignalCouldNotBeTransmittedTheConnectionWasLost(size int, error *types.SpecificError) {

}

func (o *CoreWebRequestRPCSession) SignalTheResponseCouldNotBeSent(size int, error *types.SpecificError) {

}

func (o *CoreWebRequestRPCSession) SignalTheResponseWasTransmittedSuccessfully(size int, packageHash string) {

}

func (o *CoreWebRequestRPCSession) SignalHasWritedResponseData() {

}

func (o *CoreWebRequestRPCSession) CloseBecauseFunctionReturned() {

}

func (o *CoreWebRequestRPCSession) Done() {

}
