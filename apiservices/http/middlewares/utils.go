package middlewares

import (
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/static/errormsgs"
	"github.com/CustodiaJS/custodiajs-core/types"
	"github.com/CustodiaJS/custodiajs-core/utils"
	"github.com/fxamacker/cbor/v2"
)

// Wandelt ein Paket in ein Bytearry um
func convertHttpResponseCapsleToByteSlice(rc *HttpResponseCapsle, requestContentType types.HttpRequestContentType) ([]byte, *types.SpecificError) {
	switch requestContentType {
	case static.HTTP_CONTENT_JSON:
		b, err := json.Marshal(rc)
		if err != nil {
			return nil, errormsgs.HTTP_API_SERVICE_RESPONSE_WRITING_ENCODING_ERROR("convertHttpResponseCapsleToByteSlice", requestContentType, err)
		}
		return b, nil
	case static.HTTP_CONTENT_CBOR:
		b, err := cbor.Marshal(rc)
		if err != nil {
			return nil, errormsgs.HTTP_API_SERVICE_RESPONSE_WRITING_ENCODING_ERROR("convertHttpResponseCapsleToByteSlice", requestContentType, err)
		}
		return b, nil
	default:
		return nil, nil
	}
}

// Ließt den X-Requested-With wert aus
func readXRequestedWithData(xRequestedWith string) (*types.XRequestedWithData, *types.SpecificError) {
	switch strings.ToLower(xRequestedWith) {
	case "xmlhttprequest":
		return &types.XRequestedWithData{IsKnown: true, Value: xRequestedWith}, nil
	case "fetch":
		return &types.XRequestedWithData{IsKnown: true, Value: xRequestedWith}, nil
	case "mobileapp":
		return &types.XRequestedWithData{IsKnown: true, Value: xRequestedWith}, nil
	case "reactnative":
		return &types.XRequestedWithData{IsKnown: true, Value: xRequestedWith}, nil
	default:
		return &types.XRequestedWithData{IsKnown: false, Value: xRequestedWith}, nil
	}
}

// Erzeugt einen SHA3-256 Hash aus einem Response Capsle
func ComputeSHA3_256HashFromHttpResponseCapsle(rc *HttpResponseCapsle, requestContentType types.HttpRequestContentType) (string, *types.SpecificError) {
	byted, err := convertHttpResponseCapsleToByteSlice(rc, requestContentType)
	if err != nil {
		err.AddCallerFunctionToHistory("computeSHA3_256HashFromHttpResponseCapsle")
		return "", err
	}

	hexvalued := hex.EncodeToString(byted)
	hash := utils.HashOfString(hexvalued)

	return hash, nil
}
