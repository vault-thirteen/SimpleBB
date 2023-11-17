package client

import (
	"fmt"

	cc "github.com/vault-thirteen/SimpleBB/pkg/common/client"
)

// List of supported functions.
const (
	// Ping.
	FuncPing = "Ping"

	// Captcha.
	FuncCreateCaptcha = "CreateCaptcha"
	FuncCheckCaptcha  = "CheckCaptcha"

	// Other.
	FuncShowDiagnosticData = "ShowDiagnosticData"
)

func NewClient(host string, port uint16, path string) (c *cc.Client, err error) {
	dsn := fmt.Sprintf("http://%s:%d%s", host, port, path)
	return cc.NewClient(dsn, false)
}
