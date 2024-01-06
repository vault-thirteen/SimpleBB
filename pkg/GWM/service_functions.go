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

func (srv *Server) GetProductVersion(_ *api.Request, rw http.ResponseWriter, _ *http.Request) {
	srv.respondWithPlainText(rw, srv.settings.VersionInfo.ProgramVersionString())
}

func (srv *Server) RegisterUser(ar *api.Request, rw http.ResponseWriter, _ *http.Request) {
	var err error
	rawParameters, ok := ar.Parameters.(json.RawMessage)
	if !ok {
		err = errors.New(ErrTypeCast)
		srv.processInternalServerError(rw, err)
		return
	}

	var p am.RegisterUserParams
	err = json.Unmarshal(rawParameters, &p)
	if err != nil {
		srv.processBadRequest(rw)
		return
	}

	p.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.RegisterUserResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncRegisterUser, p, result)
	if err != nil {
		srv.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &api.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(rw, response)
	return
}

func (srv *Server) ApproveAndRegisterUser(ar *api.Request, rw http.ResponseWriter, _ *http.Request) {
	var err error
	rawParameters, ok := ar.Parameters.(json.RawMessage)
	if !ok {
		err = errors.New(ErrTypeCast)
		srv.processInternalServerError(rw, err)
		return
	}

	var p am.ApproveAndRegisterUserParams
	err = json.Unmarshal(rawParameters, &p)
	if err != nil {
		srv.processBadRequest(rw)
		return
	}

	p.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.ApproveAndRegisterUserResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncApproveAndRegisterUser, p, result)
	if err != nil {
		srv.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &api.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(rw, response)
	return
}

//TODO
