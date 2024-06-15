package saftychan

import (
	"sync"
)

type baseSecureChan struct {
	lock      *sync.Mutex
	chanValue chan interface{}
	isClosed  bool
}

type FunctionCallStateChan struct {
	*baseSecureChan
}

type FunctionCallReturnChan struct {
	*baseSecureChan
}
