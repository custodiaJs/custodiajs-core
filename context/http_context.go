package context

import (
	"crypto/x509"
	"encoding/json"
	"log"
	"net/url"

	"github.com/CustodiaJS/custodiajs-core/types"
)

func (o *HttpContext) SetMethod(method types.HTTP_METHOD) {
	o.method = method
	o.proc.Debug("Set Method '%s'", method)
}

func (o *HttpContext) SetContentType(contentType types.HttpRequestContentType) {
	o.contentType = contentType
	o.proc.Debug("Set Content Type '%d'", contentType)
}

func (o *HttpContext) SetXRequestedWith(xRequestedWithData *types.XRequestedWithData) {
	o.xRequestedWithData = xRequestedWithData
	jsonData, err := json.MarshalIndent(xRequestedWithData, "", "    ")
	if err != nil {
		log.Fatalf("Error occurred during JSON marshaling. Error: %s", err)
	}
	o.proc.Debug("Set XRequestedWith '%s'", jsonData)
}

func (o *HttpContext) SetReferer(refererURL *url.URL) {
	o.refererURL = refererURL
	o.proc.Debug("Set RefererURL '%s'", refererURL)
}

func (o *HttpContext) SetOrigin(originURL *url.URL) {
	o.originURL = originURL
	o.proc.Debug("Set OriginURL '%s'", originURL)
}

func (o *HttpContext) SetTLSCertificate(tlsCert []*x509.Certificate) {
	o.tlsCert = tlsCert
	o.proc.Debug("Set TLS Certificate '%s'", tlsCert)
}

func (o *HttpContext) AddSearchedFunctionSignature(fncs *types.FunctionSignature) {
	o.fncs = fncs
	o.proc.Debug("Set Method '%s'")
}

func (o *HttpContext) GetSearchedFunctionSignature() *types.FunctionSignature {
	return o.fncs
}

func (o *HttpContext) GetMethod() types.HTTP_METHOD {
	return o.method
}

func (o *HttpContext) GetContentType() types.HttpRequestContentType {
	return o.contentType
}

func (o *HttpContext) GetXRequestedWith() *types.XRequestedWithData {
	return o.xRequestedWithData
}

func (o *HttpContext) GetReferer() *url.URL {
	return o.refererURL
}

func (o *HttpContext) GetOrigin() *url.URL {
	return o.originURL
}

func (o *HttpContext) GetTLSCertificate() []*x509.Certificate {
	return o.tlsCert
}

func (o *HttpContext) SignalsThatAnErrorHasOccurredWhenTheErrorIsSent(size int, error *types.SpecificError) {

}

func (o *HttpContext) SignalTheErrorSignalCouldNotBeTransmittedTheConnectionWasLost(size int, error *types.SpecificError) {

}

func (o *HttpContext) SignalTheResponseWasTransmittedSuccessfully(size int, packageHash string) {

}

func (o *HttpContext) SignalTheResponseCouldNotBeSent(size int, error *types.SpecificError) {

}

func (o *HttpContext) SignalThatTheErrorWasSuccessfullyTransmitted(size int) {

}

func (o *HttpContext) GetReturnChan() types.FunctionCallReturnChanInterface {
	return o.saftyResponseChan
}

func (o *HttpContext) CloseBecauseFunctionReturned() {

}
