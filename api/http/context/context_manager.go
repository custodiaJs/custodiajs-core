package context

import (
	"net"
	"net/http"

	"github.com/CustodiaJS/custodiajs-core/api/http/context/saftychan"
	"github.com/CustodiaJS/custodiajs-core/global/procslog"
	"github.com/CustodiaJS/custodiajs-core/global/types"
	"github.com/CustodiaJS/custodiajs-core/global/utils/grsbool"
)

func (o *ContextManagmentUnit) NewHTTPBasesSession(r *http.Request, proclogMother types.ProcessLogSessionInterface) (types.CoreHttpContextInterface, *types.SpecificError) {
	// Die Basis Variabeln werden erzeugt
	proc := procslog.NewChainMergedProcLog(proclogMother, procslog.NewProcLogSessionWithHeader("HttpContext"))
	isConnected, saftyResponseChan := grsbool.NewGrsbool(true), saftychan.NewFunctionCallReturnChan()

	// Die RemoteIP wird eingelesen
	remoteIp := net.ParseIP(r.RemoteAddr)

	// Es wird versucht die Lokale IP Einzulesen
	addr := r.Context().Value(http.LocalAddrContextKey).(net.Addr)
	tcpAddr, ok := addr.(*net.TCPAddr)
	var localIp net.IP
	if ok {
		localIp = net.ParseIP(tcpAddr.IP.String())
	} else {
		localIp = net.ParseIP("0.0.0.0")
	}

	// Es wird eine Goroutine gestartet, welche prüft ob die Verbindung getrennt wurde
	go func() {
		// Es wird darauf gewartet dass die Verbindung geschlossen wird
		<-r.Context().Done()

		// Es wird Signalisiert dass die Verbindung geschlossen wurde
		isConnected.Set(false)

		// Das SaftyChan wird geschlossen
		saftyResponseChan.Close()
	}()

	// Es wird ein neues Rückgabe Objekt erstellt
	returnObject := &HttpContext{
		Context:           &Context{isConnected: isConnected, proc: proc},
		saftyResponseChan: saftyResponseChan,
		localIp:           localIp,
		remoteIp:          remoteIp,
	}

	// Das HttpContext Objekt wird ohne Fehler zurückgegeben
	return returnObject, nil
}

func NewContextManager() *ContextManagmentUnit {
	return &ContextManagmentUnit{}
}

func PairCoreToContextManager(cntxm *ContextManagmentUnit, core types.CoreInterface) {
	cntxm.core = core
}
