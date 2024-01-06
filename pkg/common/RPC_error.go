package common

import "net/http"

// Common error codes start with 1000.

// Codes.
const (
	RpcErrorCode_FunctionIsNotImplemented = 1000
	RpcErrorCode_IPAddressError           = 1001
	RpcErrorCode_MalformedRequest         = 1002
	RpcErrorCode_DatabaseError            = 1003
	RpcErrorCode_InsufficientPermission   = 1004
	RpcErrorCode_GetUserDataByAuthToken   = 1005
	RpcErrorCode_UidList                  = 1006
	RpcErrorCode_InvalidToken             = 1007
)

// Messages.
const (
	RpcErrorMsg_FunctionIsNotImplemented = "function is not implemented"
	RpcErrorMsgF_IPAddressError          = "IP address error: %s" // Template.
	RpcErrorMsg_MalformedRequest         = "malformed request"
	RpcErrorMsg_DatabaseError            = "database error"
	RpcErrorMsg_InsufficientPermission   = "insufficient permission"
	RpcErrorMsgF_GetUserDataByAuthToken  = "GetUserDataByAuthToken: %s" // Template.
	RpcErrorMsgF_UidList                 = "UidList error: %s"          // Template.
	RpcErrorMsg_InvalidToken             = "invalid token"
)

// Unique HTTP status codes used in the map:
// - 400 (Bad request);
// - 403 (Forbidden);
// - 404 (Not found);
// - 500 (Internal server error).
func GetMapOfHttpStatusCodesByRpcErrorCodes() map[int]int {
	return map[int]int{
		RpcErrorCode_FunctionIsNotImplemented: http.StatusNotFound,
		RpcErrorCode_IPAddressError:           http.StatusInternalServerError,
		RpcErrorCode_MalformedRequest:         http.StatusBadRequest,
		RpcErrorCode_DatabaseError:            http.StatusInternalServerError,
		RpcErrorCode_InsufficientPermission:   http.StatusForbidden,
		RpcErrorCode_GetUserDataByAuthToken:   http.StatusForbidden,
		RpcErrorCode_UidList:                  http.StatusInternalServerError,
		RpcErrorCode_InvalidToken:             http.StatusForbidden,
	}
}
