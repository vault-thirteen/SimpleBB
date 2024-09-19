package c

import (
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/client"
)

// List of supported functions.
const (
	// Ping.
	FuncPing = cc.FuncPing

	// User registration.
	FuncRegisterUser                           = "RegisterUser"
	FuncGetListOfRegistrationsReadyForApproval = "GetListOfRegistrationsReadyForApproval"
	FuncApproveAndRegisterUser                 = "ApproveAndRegisterUser"

	// Logging in and out.
	FuncLogUserIn            = "LogUserIn"
	FuncLogUserOut           = "LogUserOut"
	FuncGetListOfLoggedUsers = "GetListOfLoggedUsers"
	FuncIsUserLoggedIn       = "IsUserLoggedIn"

	// Various actions.
	FuncChangePassword    = "ChangePassword"
	FuncChangeEmail       = "ChangeEmail"
	FuncGetListOfAllUsers = "GetListOfAllUsers"

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
	FuncShowDiagnosticData = cc.FuncShowDiagnosticData
	FuncTest               = "Test"
)
