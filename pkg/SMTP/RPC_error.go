package smtp

// RPC errors.

// Error codes must not exceed 999.

// Codes.
const (
	RpcErrorCode_MailerError = 1
)

// Messages.
const (
	RpcErrorMsgF_MailerError = "mailer error: %s"
)
