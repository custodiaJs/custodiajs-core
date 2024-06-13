package kmodulerpc

import (
	"fmt"
	"vnh1/types"
	rpcrequest "vnh1/utils/rpc_request"
)

func (o *RequestResponseUnit) WaitOfState() (*types.FunctionCallState, error) {
	// Diese Chan wird verwendet um auf das Ergebniss zu warten
	resultChan := make(chan *RequestResponseWaiter)

	// Diese Funktion wartet darauf, das neue Daten vorliegen
	go func() {
		// Es wird darauf gewartet dass Daten eintreffen
		result, ok := <-o.request.resolveChan
		if !ok {
			return
		}

		// Es wird geprüft ob das Ergebniss NULL ist
		if result == nil {
			return
		}

		// Das Ergebniss wird zurückgegegben
		resultChan <- &RequestResponseWaiter{CallState: result, Error: nil}
	}()

	// Diese Funktion wartet darauf, dass die Verbindung getrennt wird
	go func() {
		// Es wird darauf gewartet, dass die Verbindung geschlossen wurde
		rpcrequest.WaitOfConnectionStateChange(o.request._rprequest, false)

		// Die ResultChan wird geschlossen
		close(resultChan)

		// Die *.request.resolveChan wird geschlossen
		close(o.request.resolveChan)
	}()

	// Es wird auf das Ergebniss der beiden Geroutinen gewartets
	result, ok := <-resultChan
	if !ok {
		return nil, fmt.Errorf("")
	}

	// Es wird geprüft ob ein Ergebniss zurückgegeben wurde
	if result == nil {
		// Es wird geprüft ob die Verbindung getrennt wurde, wenn ja wird der Vorgang abgebrochen ohne Fehler
		if !rpcrequest.ConnectionIsOpen(o.request._rprequest) {
			return nil, nil
		}

		// Es wird ein Fehler zurückgegeben
		return nil, fmt.Errorf("unkown, internal error")
	}

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if result.Error != nil {
		return nil, result.Error
	}

	// Es wird geprüft ob Daten vorhanden sind
	if result.CallState == nil {
		return nil, nil
	}

	// Rückgabe
	return result.CallState, nil
}

func newRequestResponseWaiter(request *SharedFunctionRequestContext) (*RequestResponseUnit, error) {
	newRRW := &RequestResponseUnit{request: request}
	return newRRW, nil
}
