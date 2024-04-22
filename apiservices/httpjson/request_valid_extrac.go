package httpjson

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"vnh1/static"
	"vnh1/types"

	"github.com/fxamacker/cbor/v2"
)

func validateRequestAndGetRequestData(methode string, r *http.Request) (*RequestData, error) {
	// Es wird geprüft ob es sich um die POST Methode handelt
	if r.Method != methode {
		return nil, fmt.Errorf("invalid allow methode")
	}

	// Es wird geprüft ob eine TLS Verbindung vorhanden ist
	if r.TLS == nil {
		return nil, fmt.Errorf("ssl needed")
	}

	// Der Content Typ wird geprüft
	var contentType types.HttpRequestContentType
	switch r.Header.Get("content-type") {
	case "application/json":
		contentType = static.HTTP_CONTENT_JSON
	case "application/cbor":
		contentType = static.HTTP_CONTENT_CBOR
	default:
		return nil, fmt.Errorf("unsuported content type")
	}

	// Es wird geprüft ob der VM Query angegeben wurde
	queryParams := r.URL.Query()

	// Prüfe, ob mehr als ein Query-Parameter vorhanden ist oder der Parameter 'name' nicht existiert oder mehr als einen Wert hat.
	if len(queryParams) != 1 {
		return nil, fmt.Errorf("invalid query len")
	}

	// Prüfen, ob 'name' existiert und genau einen Wert hat.
	value, ok := queryParams["id"]
	if !ok || len(queryParams["id"]) != 1 {
		return nil, fmt.Errorf("invalid query parm")
	}

	// Die ID wird geprüft
	if len(value) != 1 {
		return nil, fmt.Errorf("invalid query parm")
	}
	if len(value[0]) != 64 {
		return nil, fmt.Errorf("invalid query parm")
	}

	// Es wird geprüft ob es sich um einen Hexwert handelt
	decodedId, err := hex.DecodeString(value[0])
	if err != nil {
		return nil, fmt.Errorf("invalid query parm")
	}

	// Der String wird zurückerstellt
	recodedHexStr := hex.EncodeToString(decodedId)

	// Das Rückgabe Objekt wird erstellt
	returnObj := &RequestData{
		TransportProtocol: static.HTTP_JSON,
		Source:            r.RemoteAddr,
		ContentType:       contentType,
		Cookies:           r.Cookies(),
		XRequestedWith:    r.Header.Get("X-Requested-With"),
		Referer:           r.Header.Get("Referer"),
		Origin:            r.Header.Get("Origin"),
		VmId:              strings.ToLower(recodedHexStr),
		TLS:               r.TLS,
	}

	// Die VM ID wird zurückgegeben
	return returnObj, nil
}

func validatePOSTRequestAndGetRequestData(r *http.Request) (*RequestData, error) {
	return validateRequestAndGetRequestData("POST", r)
}

func validateGETRequestAndGetRequestData(r *http.Request) (*RequestData, error) {
	return validateRequestAndGetRequestData("GET", r)
}

func validateWSRequestAndGetRequestData(r *http.Request) (*RequestData, error) {
	return validateRequestAndGetRequestData("GET", r)
}

func getRefererOrXRequestedWith(r *RequestData) string {
	if r.Referer != "" {
		o, err := url.Parse(r.Referer)
		if err != nil {
			return ""
		}
		return o.Host
	}
	if r.XRequestedWith != "" {
		o, err := url.Parse(r.XRequestedWith)
		if err != nil {
			return ""
		}
		return o.Host
	}
	return ""
}

func hasRefererOrXRequestedWith(r *RequestData) bool {
	return getRefererOrXRequestedWith(r) != ""
}

func extractHttpRpcBody(requestContentType types.HttpRequestContentType, body io.ReadCloser) (*RPCFunctionCall, error) {
	var data *RPCFunctionCall
	switch requestContentType {
	case static.HTTP_CONTENT_CBOR:
		body, err := io.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("extractRpcBody: " + err.Error())
		}
		if err := cbor.Unmarshal(body, &data); err != nil {
			return nil, fmt.Errorf("extractRpcBody: " + err.Error())
		}
	case static.HTTP_CONTENT_JSON:
		if err := json.NewDecoder(body).Decode(&data); err != nil {
			return nil, fmt.Errorf("extractRpcBody: " + err.Error())
		}
	default:
		return nil, fmt.Errorf("extractRpcBody: invalid content data")
	}
	return data, nil
}
