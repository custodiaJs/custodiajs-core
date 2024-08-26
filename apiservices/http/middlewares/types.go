package middlewares

import "github.com/CustodiaJS/custodiajs-core/types"

type HttpResponseCapsle struct {
	Data  []*types.RPCResponseData
	Error string
}
