package saftychan

import (
	"sync"
)

func newBaseSecureChan() *baseSecureChan {
	return &baseSecureChan{&sync.Mutex{}, make(chan interface{}), false}
}

func (o *baseSecureChan) IsClosed() bool {
	o.lock.Lock()
	val := bool(o.isClosed)
	o.lock.Unlock()
	return val
}

func (o *baseSecureChan) Close() {
	o.lock.Lock()
	if o.isClosed {
		o.lock.Unlock()
		return
	}
	close(o.chanValue)
	o.isClosed = true
	o.lock.Unlock()
}

func (o *baseSecureChan) WriteAndClose(value interface{}) {
	// Der Mutex wird verwendet
	o.lock.Lock()

	// Es wird geprüft ob der Wert bereits gesetzt wurde
	if o.isClosed {
		o.lock.Unlock()
		return
	}

	// Es wird Markiert dass der Wert geschreuben wurde
	o.isClosed = true

	// Der Mutex wird freigegeben
	o.lock.Unlock()

	// Es wird auf den Wert gewartet
	o.chanValue <- value
}
