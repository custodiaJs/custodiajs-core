package httpjson

import (
	"encoding/json"
	"net/http"
	"vnh1/types"
)

func errorResponse(contentType types.HttpRequestContentType, w http.ResponseWriter, errorMessage string) {
	// Der Passende Typ wird festgelegt
	switch contentType {
	case types.HTTP_CONTENT_CBOR:
		w.Header().Set("Content-Type", "application/cbor")
	default:
		w.Header().Set("Content-Type", "application/json")
	}

	// Die Header werden festgelegt
	w.WriteHeader(http.StatusBadRequest)

	// Die Antwort wird gebaut
	response := &RPCResponse{Result: "failed", Error: &errorMessage}

	// Das Antwortpaket wird ferigestellt
	var bytedResponse []byte
	var err error
	switch contentType {
	case types.HTTP_CONTENT_CBOR:
		bytedResponse, err = json.Marshal(response)
	default:
		bytedResponse, err = json.Marshal(response)
	}

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if err != nil {
		panic(err)
	}

	// Das Paket wird gesendet
	w.Write(bytedResponse)
}

func finalResponse(contentType types.HttpRequestContentType, w http.ResponseWriter, data *RPCResponseData) (uint64, error) {
	// Der Passende Typ wird festgelegt
	switch contentType {
	case types.HTTP_CONTENT_CBOR:
		w.Header().Set("Content-Type", "application/cbor")
	default:
		w.Header().Set("Content-Type", "application/json")
	}

	// Der Status header wird fetgelegt
	w.WriteHeader(http.StatusOK)

	// Die Antwort wird gebaut
	response := &RPCResponse{Result: "success", Data: data}

	// Das Antwortpaket wird ferigestellt
	var bytedResponse []byte
	var err error
	switch contentType {
	case types.HTTP_CONTENT_CBOR:
		bytedResponse, err = json.Marshal(response)
	default:
		bytedResponse, err = json.Marshal(response)
	}

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if err != nil {
		return 0, err
	}

	// Das Paket wird gesendet
	w.Write(bytedResponse)

	// Der Vorgang wurde ohne Fehler durchgeführt
	return uint64(len(bytedResponse)), nil
}
