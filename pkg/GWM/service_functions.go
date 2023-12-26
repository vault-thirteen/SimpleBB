package gwm

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	ac "github.com/vault-thirteen/SimpleBB/pkg/ACM/client"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	"github.com/vault-thirteen/SimpleBB/pkg/GWM/models/api"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	jc "github.com/ybbus/jsonrpc/v3"
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

	p.CommonParams = cm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.RegisterUserResult)
	var jerr *jc.RPCError
	err = srv.acmServiceClient.MakeRequest(context.Background(), result, ac.FuncRegisterUser, p)
	if err != nil {
		jerr, ok = err.(*jc.RPCError)
		if !ok {
			err = errors.New(ErrTypeCast)
			srv.processInternalServerError(rw, err)
			return
		}

		srv.processRpcError(rw, jerr)
		return
	}

	result.CommonResult.TimeSpent = 0
	var response = &api.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(rw, response)
	return
}

//TODO
