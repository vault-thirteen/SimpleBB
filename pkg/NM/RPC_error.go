package nm

import "net/http"

// RPC errors.

// Error codes must not exceed 999.

// Codes.
const (
	RpcErrorCode_UserIdIsNotSet            = 1
	RpcErrorCode_TextIsNotSet              = 2
	RpcErrorCode_NotificationIdIsNotSet    = 3
	RpcErrorCode_NotificationIsNotFound    = 4
	RpcErrorCode_NotificationIsAlreadyRead = 5
	RpcErrorCode_TestError                 = 6
)

// Messages.
const (
	RpcErrorMsg_UserIdIsNotSet            = "user ID is not set"
	RpcErrorMsg_TextIsNotSet              = "text is not set"
	RpcErrorMsg_NotificationIdIsNotSet    = "notification ID is not set"
	RpcErrorMsg_NotificationIsNotFound    = "notification is not found"
	RpcErrorMsg_NotificationIsAlreadyRead = "notification is already read"
	RpcErrorMsgF_TestError                = "test error: %s"
)

// Unique HTTP status codes used in the map:
// - 400 (Bad request);
// - 404 (Not found);
// - 409 (Conflict);
// - 500 (Internal server error).
func GetMapOfHttpStatusCodesByRpcErrorCodes() map[int]int {
	return map[int]int{
		RpcErrorCode_UserIdIsNotSet:            http.StatusBadRequest,
		RpcErrorCode_TextIsNotSet:              http.StatusBadRequest,
		RpcErrorCode_NotificationIdIsNotSet:    http.StatusBadRequest,
		RpcErrorCode_NotificationIsNotFound:    http.StatusNotFound,
		RpcErrorCode_NotificationIsAlreadyRead: http.StatusConflict,
		RpcErrorCode_TestError:                 http.StatusInternalServerError,
	}
}
