package c

import (
	"fmt"

	cc "github.com/vault-thirteen/SimpleBB/pkg/common/client"
)

// List of supported functions.
const (
	// Ping.
	FuncPing = "Ping"

	// IP address list.
	FuncBlockIPAddress     = "BlockIPAddress"
	FuncIsIPAddressBlocked = "IsIPAddressBlocked"

	// Other.
	FuncShowDiagnosticData = "ShowDiagnosticData"
)

func NewClient(host string, port uint16, path string) (c *cc.Client, err error) {
	dsn := fmt.Sprintf("http://%s:%d%s", host, port, path)
	return cc.NewClient(dsn, false)
}
