package grsbool

import (
	"sync"
)

type Grsbool struct {
	_bool  bool
	_mutex *sync.Mutex
	_cond  *sync.Cond
}

func (o *Grsbool) Set(bval bool) {
	o._mutex.Lock()
	defer o._mutex.Unlock()
	o._bool = bval
	o._cond.Broadcast()
}

func (o *Grsbool) Bool() bool {
	o._mutex.Lock()
	defer o._mutex.Unlock()
	return o._bool
}

func (o *Grsbool) WaitOfChange(waitOfState bool) {
	o._cond.L.Lock()
	for o._bool != waitOfState {
		o._cond.Wait()
	}
	o._cond.L.Unlock()
}

func NewGrsbool(v bool) *Grsbool {
	m := &sync.Mutex{}
	c := sync.NewCond(m)
	return &Grsbool{_mutex: m, _bool: v, _cond: c}
}
