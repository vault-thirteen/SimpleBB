package smtp

// RPC errors.

// Error codes must not exceed 999.

// Codes.
const (
	RpcErrorCode_Exception   = 1
	RpcErrorCode_MailerError = 2
)

// Messages.
const (
	RpcErrorMsg_Exception    = "exception"
	RpcErrorMsgF_MailerError = "mailer error: %s"
)
