package eventloop

import (
	"github.com/CustodiaJS/custodiajs-core/global/types"

	v8 "rogchap.com/v8go"
)

func (o *KernelEventLoopOperation) GetType() types.KernelEventLoopOperationMethode {
	return o.Type
}

func (o *KernelEventLoopOperation) GetFunction() func(*v8.Context, types.KernelEventLoopContextInterface) {
	return o.DirectV8Function
}

func (o *KernelEventLoopOperation) GetSourceCode() string {
	return o.EngineSourceCode
}

func (o *KernelEventLoopOperation) SetResult(value *v8.Value) {
	// Der Mutex wird verwendet
	o._cond.L.Lock()

	// Der Wert wird gespeichert
	if value == nil {
		o._hasNullReturn = true
	} else {
		o._returnValue = value
	}

	// Es wird Signalisiert das ein Wert vorhanden ist
	o._cond.Broadcast()

	// Der Mutex wird freigegeben
	o._cond.L.Unlock()
}

func (o *KernelEventLoopOperation) SetError(err error) {
	// Der Mutex wird verwendet
	o._cond.L.Lock()

	// Der Wert wird gespeichert
	o._returnError = err

	// Es wird Signalisiert das ein Wert vorhanden ist
	o._cond.Broadcast()

	// Der Mutex wird freigegeben
	o._cond.L.Unlock()
}

func (o *KernelEventLoopOperation) WaitOfResponse() (*v8.Value, error) {
	// Der Mutex wird verwendet
	o._cond.L.Lock()

	// Es wird geprüft ob ein Rückgabewert vorhanden ist
	for {
		// Es wird geprüft ob ein Fehler vorhanden ist
		if o._returnError != nil {
			// Der Mutex wird freigegeben
			defer o._cond.L.Unlock()

			// Der Wert wird zurückgegeben
			return nil, o._returnError
		}

		// Es wird geprüft ob ein Nullwert vorhanden ist
		if o._hasNullReturn {
			// Der Mutex wird freigegeben
			defer o._cond.L.Unlock()

			// Der Wert wird zurückgegeben
			return nil, nil
		}

		// Es wird auf die Verbindung gewartet
		if o._returnValue != nil {
			// Der Mutex wird freigegeben
			defer o._cond.L.Unlock()

			// Der Wert wird zurückgegeben
			return o._returnValue, nil
		}

		// Es wird gewartet
		o._cond.Wait()
	}
}

func (o *KernelEventLoopOperation) GetOperation() types.KernelEventLoopContextInterface {
	return &KernelEventLoopContext{setError: o.SetError, setResult: o.SetResult}
}
