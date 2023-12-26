package rcs

// RPC errors.

// Error codes must not exceed 999.

// Codes.
const (
	RpcErrorCode_Exception      = 1
	RpcErrorCode_TaskIdIsNotSet = 2
	RpcErrorCode_AnswerIsNotSet = 3
	RpcErrorCode_CheckError     = 4
	RpcErrorCode_CreateError    = 5
)

// Messages.
const (
	RpcErrorMsg_Exception      = "exception"
	RpcErrorMsg_TaskIdIsNotSet = "task ID is not set"
	RpcErrorMsg_AnswerIsNotSet = "answer is not set"
	RpcErrorMsgF_CheckError    = "check error: %s"
	RpcErrorMsgF_CreateError   = "captcha creation error: %s"
)
