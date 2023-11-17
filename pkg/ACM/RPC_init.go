package acm

import (
	ac "github.com/vault-thirteen/SimpleBB/pkg/ACM/client"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
)

func (srv *Server) initJsonRpcHandlers() (err error) {
	// Ping.
	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncPing, PingHandler{Server: srv}, am.PingParams{}, am.PingResult{})
	if err != nil {
		return err
	}

	// User registration.
	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncRegisterUser, RegisterUserHandler{Server: srv}, am.RegisterUserParams{}, am.RegisterUserResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncApproveAndRegisterUser, ApproveAndRegisterUserHandler{Server: srv}, am.ApproveAndRegisterUserParams{}, am.ApproveAndRegisterUserResult{})
	if err != nil {
		return err
	}

	// Logging in and out.
	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncLogUserIn, LogUserInHandler{Server: srv}, am.LogUserInParams{}, am.LogUserInResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncLogUserOut, LogUserOutHandler{Server: srv}, am.LogUserOutParams{}, am.LogUserOutResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncGetListOfLoggedUsers, GetListOfLoggedUsersHandler{Server: srv}, am.GetListOfLoggedUsersParams{}, am.GetListOfLoggedUsersResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncIsUserLoggedIn, IsUserLoggedInHandler{Server: srv}, am.IsUserLoggedInParams{}, am.IsUserLoggedInResult{})
	if err != nil {
		return err
	}

	// User properties.
	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncGetUserRoles, GetUserRolesHandler{Server: srv}, am.GetUserRolesParams{}, am.GetUserRolesResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncViewUserParameters, ViewUserParametersHandler{Server: srv}, am.ViewUserParametersParams{}, am.ViewUserParametersResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncSetUserRoleAuthor, SetUserRoleAuthorHandler{Server: srv}, am.SetUserRoleAuthorParams{}, am.SetUserRoleAuthorResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncSetUserRoleWriter, SetUserRoleWriterHandler{Server: srv}, am.SetUserRoleWriterParams{}, am.SetUserRoleWriterResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncSetUserRoleReader, SetUserRoleReaderHandler{Server: srv}, am.SetUserRoleReaderParams{}, am.SetUserRoleReaderResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncGetSelfRoles, GetSelfRolesHandler{Server: srv}, am.GetSelfRolesParams{}, am.GetSelfRolesResult{})
	if err != nil {
		return err
	}

	// User banning.
	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncBanUser, BanUserHandler{Server: srv}, am.BanUserParams{}, am.BanUserResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncUnbanUser, UnbanUserHandler{Server: srv}, am.UnbanUserParams{}, am.UnbanUserResult{})
	if err != nil {
		return err
	}

	// Other.
	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncShowDiagnosticData, ShowDiagnosticDataHandler{Server: srv}, am.ShowDiagnosticDataParams{}, am.ShowDiagnosticDataResult{})
	if err != nil {
		return err
	}

	if srv.settings.SystemSettings.IsDebugMode {
		err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncTest, TestHandler{Server: srv}, am.TestParams{}, am.TestResult{})
		if err != nil {
			return err
		}
	}

	// Template.
	//err = srv.jsonRpcHandlers.RegisterMethod(FuncAbc, AbcHandler{Server: srv}, am.AbcParams{}, am.AbcResult{})
	//if err != nil {
	//	return err
	//}

	return nil
}
