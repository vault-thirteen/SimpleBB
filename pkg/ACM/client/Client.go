package c

import (
	"fmt"

	cc "github.com/vault-thirteen/SimpleBB/pkg/common/client"
)

// List of supported functions.
const (
	// Ping.
	FuncPing = "Ping"

	// User registration.
	FuncRegisterUser           = "RegisterUser"
	FuncApproveAndRegisterUser = "ApproveAndRegisterUser"

	// Logging in and out.
	FuncLogUserIn            = "LogUserIn"
	FuncLogUserOut           = "LogUserOut"
	FuncGetListOfLoggedUsers = "GetListOfLoggedUsers"
	FuncIsUserLoggedIn       = "IsUserLoggedIn"

	// Various actions.
	FuncChangePassword = "ChangePassword"
	FuncChangeEmail    = "ChangeEmail"

	// User properties.
	FuncGetUserRoles       = "GetUserRoles"
	FuncViewUserParameters = "ViewUserParameters"
	FuncSetUserRoleAuthor  = "SetUserRoleAuthor"
	FuncSetUserRoleWriter  = "SetUserRoleWriter"
	FuncSetUserRoleReader  = "SetUserRoleReader"
	FuncGetSelfRoles       = "GetSelfRoles"

	// User banning.
	FuncBanUser   = "BanUser"
	FuncUnbanUser = "UnbanUser"

	// Other.
	FuncShowDiagnosticData = "ShowDiagnosticData"
	FuncTest               = "Test"
)

func NewClient(host string, port uint16, path string, enableSelfSignedCertificate bool) (c *cc.Client, err error) {
	dsn := fmt.Sprintf("https://%s:%d%s", host, port, path)
	return cc.NewClient(dsn, enableSelfSignedCertificate)
}
