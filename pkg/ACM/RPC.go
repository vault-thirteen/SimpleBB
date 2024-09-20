package acm

// RPC handlers.

import (
	"encoding/json"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	cs "github.com/vault-thirteen/SimpleBB/pkg/common/settings"
)

func (srv *Server) initRpc() (err error) {
	rpcDurationFieldName := cs.RpcDurationFieldName
	rpcRequestIdFieldName := cs.RpcRequestIdFieldName

	ps := &jrm1.ProcessorSettings{
		CatchExceptions:    true,
		LogExceptions:      true,
		CountRequests:      true,
		DurationFieldName:  &rpcDurationFieldName,
		RequestIdFieldName: &rpcRequestIdFieldName,
	}

	srv.js, err = jrm1.NewProcessor(ps)
	if err != nil {
		return err
	}

	fns := []jrm1.RpcFunction{
		srv.Ping,
		srv.RegisterUser,
		srv.GetListOfRegistrationsReadyForApproval,
		srv.RejectRegistrationRequest,
		srv.ApproveAndRegisterUser,
		srv.LogUserIn,
		srv.LogUserOut,
		srv.GetListOfLoggedUsers,
		srv.IsUserLoggedIn,
		srv.ChangePassword,
		srv.ChangeEmail,
		srv.GetListOfAllUsers,
		srv.GetUserRoles,
		srv.ViewUserParameters,
		srv.SetUserRoleAuthor,
		srv.SetUserRoleWriter,
		srv.SetUserRoleReader,
		srv.GetSelfRoles,
		srv.BanUser,
		srv.UnbanUser,
		srv.ShowDiagnosticData,
		srv.Test,
	}

	for _, fn := range fns {
		err = srv.js.AddFunc(fn)
		if err != nil {
			return err
		}
	}

	return nil
}

// Ping.

func (srv *Server) Ping(_ *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	return am.PingResult{OK: true}, nil
}

// User registration.

func (srv *Server) RegisterUser(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.RegisterUserParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.RegisterUserResult
	r, re = srv.registerUser(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetListOfRegistrationsReadyForApproval(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.GetListOfRegistrationsReadyForApprovalParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.GetListOfRegistrationsReadyForApprovalResult
	r, re = srv.getListOfRegistrationsReadyForApproval(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) RejectRegistrationRequest(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.RejectRegistrationRequestParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.RejectRegistrationRequestResult
	r, re = srv.rejectRegistrationRequest(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ApproveAndRegisterUser(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.ApproveAndRegisterUserParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.ApproveAndRegisterUserResult
	r, re = srv.approveAndRegisterUser(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

// Logging in and out.

func (srv *Server) LogUserIn(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.LogUserInParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.LogUserInResult
	r, re = srv.logUserIn(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) LogUserOut(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.LogUserOutParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.LogUserOutResult
	r, re = srv.logUserOut(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetListOfLoggedUsers(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.GetListOfLoggedUsersParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.GetListOfLoggedUsersResult
	r, re = srv.getListOfLoggedUsers(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) IsUserLoggedIn(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.IsUserLoggedInParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.IsUserLoggedInResult
	r, re = srv.isUserLoggedIn(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

// Various actions.

func (srv *Server) ChangePassword(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.ChangePasswordParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.ChangePasswordResult
	r, re = srv.changePassword(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ChangeEmail(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.ChangeEmailParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.ChangeEmailResult
	r, re = srv.changeEmail(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetListOfAllUsers(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.GetListOfAllUsersParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.GetListOfAllUsersResult
	r, re = srv.getListOfAllUsers(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

// User properties.

func (srv *Server) GetUserRoles(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.GetUserRolesParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.GetUserRolesResult
	r, re = srv.getUserRoles(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ViewUserParameters(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.ViewUserParametersParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.ViewUserParametersResult
	r, re = srv.viewUserParameters(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) SetUserRoleAuthor(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.SetUserRoleAuthorParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.SetUserRoleAuthorResult
	r, re = srv.setUserRoleAuthor(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) SetUserRoleWriter(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.SetUserRoleWriterParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.SetUserRoleWriterResult
	r, re = srv.setUserRoleWriter(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) SetUserRoleReader(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.SetUserRoleReaderParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.SetUserRoleReaderResult
	r, re = srv.setUserRoleReader(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetSelfRoles(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.GetSelfRolesParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.GetSelfRolesResult
	r, re = srv.getSelfRoles(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

// User banning.

func (srv *Server) BanUser(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.BanUserParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.BanUserResult
	r, re = srv.banUser(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) UnbanUser(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.UnbanUserParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.UnbanUserResult
	r, re = srv.unbanUser(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

// Other.

func (srv *Server) ShowDiagnosticData(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.ShowDiagnosticDataParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.ShowDiagnosticDataResult
	r, re = srv.showDiagnosticData()
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) Test(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *am.TestParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *am.TestResult
	r, re = srv.test(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}
