package saftychan

import (
	"github.com/CustodiaJS/custodiajs-core/types"
)

func (o *FunctionCallStateChan) Read() (*types.FunctionCallState, bool) {
	o.lock.Lock()
	if o.isClosed {
		o.lock.Unlock()
		return nil, false
	}
	o.lock.Unlock()

	val, ol := <-o.chanValue
	if val == nil {
		return nil, ol
	}
	convertedValue := val.(*types.FunctionCallState)
	return convertedValue, ol
}

func (o *FunctionCallStateChan) WriteAndClose(value *types.FunctionCallState) {
	o.baseSecureChan.WriteAndClose(value)
}

func NewFunctionCallStateChan() *FunctionCallStateChan {
	return &FunctionCallStateChan{newBaseSecureChan()}
}
