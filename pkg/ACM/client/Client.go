package c

import (
	"fmt"

	cc "github.com/vault-thirteen/SimpleBB/pkg/common/client"
)

// List of supported functions.
const (
	FuncPing                   = "Ping"
	FuncRegisterUser           = "RegisterUser"
	FuncApproveAndRegisterUser = "ApproveAndRegisterUser"
	FuncLogUserIn              = "LogUserIn"
	FuncLogUserOut             = "LogUserOut"
	FuncGetListOfLoggedUsers   = "GetListOfLoggedUsers"
	FuncIsUserLoggedIn         = "IsUserLoggedIn"
	FuncGetUserRoles           = "GetUserRoles"
	FuncViewUserParameters     = "ViewUserParameters"
	FuncSetUserRoleAuthor      = "SetUserRoleAuthor"
	FuncSetUserRoleWriter      = "SetUserRoleWriter"
	FuncSetUserRoleReader      = "SetUserRoleReader"
	FuncBanUser                = "BanUser"
	FuncUnbanUser              = "UnbanUser"
	FuncGetSelfRoles           = "GetSelfRoles"
	FuncShowDiagnosticData     = "ShowDiagnosticData"
	FuncTest                   = "Test"
)

func NewClient(host string, port uint16, path string, enableSelfSignedCertificate bool) (c *cc.Client, err error) {
	dsn := fmt.Sprintf("https://%s:%d%s", host, port, path)
	return cc.NewClient(dsn, enableSelfSignedCertificate)
}
