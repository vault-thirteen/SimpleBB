package acm

// RPC handlers.

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	js "github.com/osamingo/jsonrpc/v2"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
)

// Ping.

type PingHandler struct {
	Server *Server
}

func (h PingHandler) ServeJSONRPC(_ context.Context, _ *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	result := am.PingResult{OK: true}
	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

// User registration.

type RegisterUserHandler struct {
	Server *Server
}

func (h RegisterUserHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.RegisterUserParams
	jerr = js.Unmarshal(params, &p)
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

func (h ApproveAndRegisterUserHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.ApproveAndRegisterUserParams
	jerr = js.Unmarshal(params, &p)
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

// Logging in and out.

type LogUserInHandler struct {
	Server *Server
}

func (h LogUserInHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.LogUserInParams
	jerr = js.Unmarshal(params, &p)
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

func (h LogUserOutHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.LogUserOutParams
	jerr = js.Unmarshal(params, &p)
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

func (h GetListOfLoggedUsersHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.GetListOfLoggedUsersParams
	jerr = js.Unmarshal(params, &p)
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

func (h IsUserLoggedInHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.IsUserLoggedInParams
	jerr = js.Unmarshal(params, &p)
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

// Various actions.

type ChangePasswordHandler struct {
	Server *Server
}

func (h ChangePasswordHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.ChangePasswordParams
	jerr = js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *am.ChangePasswordResult
	result, jerr = h.Server.changePassword(&p)
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

type ChangeEmailHandler struct {
	Server *Server
}

func (h ChangeEmailHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.ChangeEmailParams
	jerr = js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *am.ChangeEmailResult
	result, jerr = h.Server.changeEmail(&p)
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

// User properties.

type GetUserRolesHandler struct {
	Server *Server
}

func (h GetUserRolesHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.GetUserRolesParams
	jerr = js.Unmarshal(params, &p)
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

func (h ViewUserParametersHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.ViewUserParametersParams
	jerr = js.Unmarshal(params, &p)
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

func (h SetUserRoleAuthorHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.SetUserRoleAuthorParams
	jerr = js.Unmarshal(params, &p)
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

func (h SetUserRoleWriterHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.SetUserRoleWriterParams
	jerr = js.Unmarshal(params, &p)
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

func (h SetUserRoleReaderHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.SetUserRoleReaderParams
	jerr = js.Unmarshal(params, &p)
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

type GetSelfRolesHandler struct {
	Server *Server
}

func (h GetSelfRolesHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.GetSelfRolesParams
	jerr = js.Unmarshal(params, &p)
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

// User banning.

type BanUserHandler struct {
	Server *Server
}

func (h BanUserHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.BanUserParams
	jerr = js.Unmarshal(params, &p)
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

func (h UnbanUserHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.UnbanUserParams
	jerr = js.Unmarshal(params, &p)
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

// Other.

type ShowDiagnosticDataHandler struct {
	Server *Server
}

func (h ShowDiagnosticDataHandler) ServeJSONRPC(_ context.Context, _ *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var result *am.ShowDiagnosticDataResult
	result, jerr = h.Server.showDiagnosticData()
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

func (h TestHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p am.TestParams
	jerr = js.Unmarshal(params, &p)
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

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}
