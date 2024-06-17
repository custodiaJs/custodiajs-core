package saftychan

import (
	"github.com/CustodiaJS/custodiajs-core/types"
)

func (o *FunctionCallReturnChan) Read() (*types.FunctionCallReturn, bool) {
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
	convertedValue := val.(*types.FunctionCallReturn)
	return convertedValue, ol
}

func (o *FunctionCallReturnChan) WriteAndClose(value *types.FunctionCallReturn) {
	o.baseSecureChan.WriteAndClose(value)
}

func NewFunctionCallReturnChan() *FunctionCallReturnChan {
	return &FunctionCallReturnChan{newBaseSecureChan()}
}
