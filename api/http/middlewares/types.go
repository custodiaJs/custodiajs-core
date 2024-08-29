package middlewares

import "github.com/CustodiaJS/custodiajs-core/global/types"

type HttpResponseCapsle struct {
	Data  []*types.RPCResponseData
	Error string
}
