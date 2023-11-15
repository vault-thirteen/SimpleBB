package acm

// RPC handlers.

import (
	"context"
	"encoding/json"
	"time"

	js "github.com/osamingo/jsonrpc/v2"
	ac "github.com/vault-thirteen/SimpleBB/pkg/ACM/client"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
)

func (srv *Server) initJsonRpcHandlers() (err error) {
	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncPing, PingHandler{Server: srv}, am.PingParams{}, am.PingResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncRegisterUser, RegisterUserHandler{Server: srv}, am.RegisterUserParams{}, am.RegisterUserResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncApproveAndRegisterUser, ApproveAndRegisterUserHandler{Server: srv}, am.ApproveAndRegisterUserParams{}, am.ApproveAndRegisterUserResult{})
	if err != nil {
		return err
	}

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

	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncBanUser, BanUserHandler{Server: srv}, am.BanUserParams{}, am.BanUserResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncUnbanUser, UnbanUserHandler{Server: srv}, am.UnbanUserParams{}, am.UnbanUserResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(ac.FuncGetSelfRoles, GetSelfRolesHandler{Server: srv}, am.GetSelfRolesParams{}, am.GetSelfRolesResult{})
	if err != nil {
		return err
	}

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

type PingHandler struct {
	Server *Server
}

func (h PingHandler) ServeJSONRPC(_ context.Context, _ *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	result := am.PingResult{OK: true}
	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type RegisterUserHandler struct {
	Server *Server
}

func (h RegisterUserHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.RegisterUserParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *am.RegisterUserResult
	result, jerr = h.Server.registerUser(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type ApproveAndRegisterUserHandler struct {
	Server *Server
}

func (h ApproveAndRegisterUserHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.ApproveAndRegisterUserParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *am.ApproveAndRegisterUserResult
	result, jerr = h.Server.approveAndRegisterUser(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type LogUserInHandler struct {
	Server *Server
}

func (h LogUserInHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.LogUserInParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *am.LogUserInResult
	result, jerr = h.Server.logUserIn(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type LogUserOutHandler struct {
	Server *Server
}

func (h LogUserOutHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.LogUserOutParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *am.LogUserOutResult
	result, jerr = h.Server.logUserOut(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type GetListOfLoggedUsersHandler struct {
	Server *Server
}

func (h GetListOfLoggedUsersHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.GetListOfLoggedUsersParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *am.GetListOfLoggedUsersResult
	result, jerr = h.Server.getListOfLoggedUsers(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type IsUserLoggedInHandler struct {
	Server *Server
}

func (h IsUserLoggedInHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.IsUserLoggedInParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *am.IsUserLoggedInResult
	result, jerr = h.Server.isUserLoggedIn(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type GetUserRolesHandler struct {
	Server *Server
}

func (h GetUserRolesHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.GetUserRolesParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *am.GetUserRolesResult
	result, jerr = h.Server.getUserRoles(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type ViewUserParametersHandler struct {
	Server *Server
}

func (h ViewUserParametersHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.ViewUserParametersParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *am.ViewUserParametersResult
	result, jerr = h.Server.viewUserParameters(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type SetUserRoleAuthorHandler struct {
	Server *Server
}

func (h SetUserRoleAuthorHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.SetUserRoleAuthorParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *am.SetUserRoleAuthorResult
	result, jerr = h.Server.setUserRoleAuthor(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type SetUserRoleWriterHandler struct {
	Server *Server
}

func (h SetUserRoleWriterHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.SetUserRoleWriterParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *am.SetUserRoleWriterResult
	result, jerr = h.Server.setUserRoleWriter(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type SetUserRoleReaderHandler struct {
	Server *Server
}

func (h SetUserRoleReaderHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.SetUserRoleReaderParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *am.SetUserRoleReaderResult
	result, jerr = h.Server.setUserRoleReader(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type BanUserHandler struct {
	Server *Server
}

func (h BanUserHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.BanUserParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *am.BanUserResult
	result, jerr = h.Server.banUser(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type UnbanUserHandler struct {
	Server *Server
}

func (h UnbanUserHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.UnbanUserParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *am.UnbanUserResult
	result, jerr = h.Server.unbanUser(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type GetSelfRolesHandler struct {
	Server *Server
}

func (h GetSelfRolesHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.GetSelfRolesParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *am.GetSelfRolesResult
	result, jerr = h.Server.getSelfRoles(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type ShowDiagnosticDataHandler struct {
	Server *Server
}

func (h ShowDiagnosticDataHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	result, jerr := h.Server.showDiagnosticData()
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type TestHandler struct {
	Server *Server
}

func (h TestHandler) ServeJSONRPC(ctx context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	//requestId, jerr := c.GetRequestId(ctx)
	//if jerr != nil {
	//	return nil, jerr
	//}
	//log.Println("Test. Request start. ID:", requestId) //TODO

	var p am.TestParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *am.TestResult
	result, jerr = h.Server.doTest(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	//log.Println("Test. Request end. ID:", requestId) //TODO

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}
