package gwm

// RPC errors.

// Error codes must not exceed 999.

// Codes.
const (
	RpcErrorCode_Exception          = 1
	RpcErrorCode_FirewallIsDisabled = 2
	RpcErrorCode_IPAddressIsNotSet  = 3
	RpcErrorCode_BlockTimeIsNotSet  = 4
)

// Messages.
const (
	RpcErrorMsg_Exception          = "exception"
	RpcErrorMsg_FirewallIsDisabled = "Firewall is disabled"
	RpcErrorMsg_IPAddressIsNotSet  = "IP address is not set"
	RpcErrorMsg_BlockTimeIsNotSet  = "Block time is not set"
)
