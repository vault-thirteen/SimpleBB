package gwm

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	ac "github.com/vault-thirteen/SimpleBB/pkg/ACM/client"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	"github.com/vault-thirteen/SimpleBB/pkg/GWM/models/api"
	"github.com/vault-thirteen/SimpleBB/pkg/common/app"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

// Service functions.

func (srv *Server) GetProductVersion(_ *api.Request, _ *http.Request, hrw http.ResponseWriter) {
	srv.respondWithPlainText(hrw, srv.settings.VersionInfo.ProgramVersionString())
}

// ACM.

func (srv *Server) RegisterUser(ar *api.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	rawParameters, ok := ar.Parameters.(json.RawMessage)
	if !ok {
		err = errors.New(ErrTypeCast)
		srv.processInternalServerError(hrw, err)
		return
	}

	var params am.RegisterUserParams
	err = json.Unmarshal(rawParameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.RegisterUserResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncRegisterUser, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ApproveAndRegisterUser(ar *api.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	rawParameters, ok := ar.Parameters.(json.RawMessage)
	if !ok {
		err = errors.New(ErrTypeCast)
		srv.processInternalServerError(hrw, err)
		return
	}

	var params am.ApproveAndRegisterUserParams
	err = json.Unmarshal(rawParameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.ApproveAndRegisterUserResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncApproveAndRegisterUser, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) LogUserIn(ar *api.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	rawParameters, ok := ar.Parameters.(json.RawMessage)
	if !ok {
		err = errors.New(ErrTypeCast)
		srv.processInternalServerError(hrw, err)
		return
	}

	var params am.LogUserInParams
	err = json.Unmarshal(rawParameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.LogUserInResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncLogUserIn, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) LogUserOut(ar *api.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	rawParameters, ok := ar.Parameters.(json.RawMessage)
	if !ok {
		err = errors.New(ErrTypeCast)
		srv.processInternalServerError(hrw, err)
		return
	}

	var params am.LogUserOutParams
	err = json.Unmarshal(rawParameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.LogUserOutResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncLogUserOut, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetListOfLoggedUsers(ar *api.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	rawParameters, ok := ar.Parameters.(json.RawMessage)
	if !ok {
		err = errors.New(ErrTypeCast)
		srv.processInternalServerError(hrw, err)
		return
	}

	var params am.GetListOfLoggedUsersParams
	err = json.Unmarshal(rawParameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.GetListOfLoggedUsersResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncGetListOfLoggedUsers, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) IsUserLoggedIn(ar *api.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	rawParameters, ok := ar.Parameters.(json.RawMessage)
	if !ok {
		err = errors.New(ErrTypeCast)
		srv.processInternalServerError(hrw, err)
		return
	}

	var params am.IsUserLoggedInParams
	err = json.Unmarshal(rawParameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.IsUserLoggedInResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncIsUserLoggedIn, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ChangePassword(ar *api.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	rawParameters, ok := ar.Parameters.(json.RawMessage)
	if !ok {
		err = errors.New(ErrTypeCast)
		srv.processInternalServerError(hrw, err)
		return
	}

	var params am.ChangePasswordParams
	err = json.Unmarshal(rawParameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.ChangePasswordResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncChangePassword, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ChangeEmail(ar *api.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	rawParameters, ok := ar.Parameters.(json.RawMessage)
	if !ok {
		err = errors.New(ErrTypeCast)
		srv.processInternalServerError(hrw, err)
		return
	}

	var params am.ChangeEmailParams
	err = json.Unmarshal(rawParameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.ChangeEmailResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncChangeEmail, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetUserRoles(ar *api.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	rawParameters, ok := ar.Parameters.(json.RawMessage)
	if !ok {
		err = errors.New(ErrTypeCast)
		srv.processInternalServerError(hrw, err)
		return
	}

	var params am.GetUserRolesParams
	err = json.Unmarshal(rawParameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.GetUserRolesResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncGetUserRoles, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ViewUserParameters(ar *api.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	rawParameters, ok := ar.Parameters.(json.RawMessage)
	if !ok {
		err = errors.New(ErrTypeCast)
		srv.processInternalServerError(hrw, err)
		return
	}

	var params am.ViewUserParametersParams
	err = json.Unmarshal(rawParameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.ViewUserParametersResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncViewUserParameters, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) SetUserRoleAuthor(ar *api.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	rawParameters, ok := ar.Parameters.(json.RawMessage)
	if !ok {
		err = errors.New(ErrTypeCast)
		srv.processInternalServerError(hrw, err)
		return
	}

	var params am.SetUserRoleAuthorParams
	err = json.Unmarshal(rawParameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.SetUserRoleAuthorResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncSetUserRoleAuthor, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) SetUserRoleWriter(ar *api.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	rawParameters, ok := ar.Parameters.(json.RawMessage)
	if !ok {
		err = errors.New(ErrTypeCast)
		srv.processInternalServerError(hrw, err)
		return
	}

	var params am.SetUserRoleWriterParams
	err = json.Unmarshal(rawParameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.SetUserRoleWriterResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncSetUserRoleWriter, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) SetUserRoleReader(ar *api.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	rawParameters, ok := ar.Parameters.(json.RawMessage)
	if !ok {
		err = errors.New(ErrTypeCast)
		srv.processInternalServerError(hrw, err)
		return
	}

	var params am.SetUserRoleReaderParams
	err = json.Unmarshal(rawParameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.SetUserRoleReaderResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncSetUserRoleReader, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetSelfRoles(ar *api.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	rawParameters, ok := ar.Parameters.(json.RawMessage)
	if !ok {
		err = errors.New(ErrTypeCast)
		srv.processInternalServerError(hrw, err)
		return
	}

	var params am.GetSelfRolesParams
	err = json.Unmarshal(rawParameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.GetSelfRolesResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncGetSelfRoles, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) BanUser(ar *api.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	rawParameters, ok := ar.Parameters.(json.RawMessage)
	if !ok {
		err = errors.New(ErrTypeCast)
		srv.processInternalServerError(hrw, err)
		return
	}

	var params am.BanUserParams
	err = json.Unmarshal(rawParameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.BanUserResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncBanUser, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) UnbanUser(ar *api.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	rawParameters, ok := ar.Parameters.(json.RawMessage)
	if !ok {
		err = errors.New(ErrTypeCast)
		srv.processInternalServerError(hrw, err)
		return
	}

	var params am.UnbanUserParams
	err = json.Unmarshal(rawParameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.UnbanUserResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncUnbanUser, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

//TODO
