package c

import (
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/client"
)

// List of supported functions.
const (
	// Ping.
	FuncPing = cc.FuncPing

	// IP address list.
	FuncBlockIPAddress     = "BlockIPAddress"
	FuncIsIPAddressBlocked = "IsIPAddressBlocked"

	// Other.
	FuncShowDiagnosticData = cc.FuncShowDiagnosticData
)
