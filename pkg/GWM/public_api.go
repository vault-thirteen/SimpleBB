package gwm

import (
	"encoding/json"
	"io"
	"net/http"

	gm "github.com/vault-thirteen/SimpleBB/pkg/GWM/models"
	"github.com/vault-thirteen/SimpleBB/pkg/GWM/models/api"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	hh "github.com/vault-thirteen/auxie/http-helper"
)

const (
	ApiFunctionName_GetProductVersion      = "getProductVersion"
	ApiFunctionName_RegisterUser           = "registerUser"
	ApiFunctionName_ApproveAndRegisterUser = "approveAndRegisterUser"
	//TODO
)

func (srv *Server) handlePublicApi(rw http.ResponseWriter, req *http.Request, clientIPA string) {
	if req.Method != http.MethodPost {
		srv.respondMethodNotAllowed(rw)
		return
	}

	// Check accepted MIME types.
	ok, err := hh.CheckBrowserSupportForJson(req)
	if err != nil {
		srv.respondBadRequest(rw)
		return
	}

	if !ok {
		srv.respondNotAcceptable(rw)
		return
	}

	var reqBody []byte
	reqBody, err = io.ReadAll(req.Body)
	if err != nil {
		srv.processInternalServerError(rw, err)
		return
	}

	// Check the action.
	var arwoa api.RequestWithOnlyAction
	err = json.Unmarshal(reqBody, &arwoa)
	if err != nil {
		srv.respondBadRequest(rw)
		return
	}

	if arwoa.Action == nil {
		srv.respondBadRequest(rw)
		return
	}

	var handler api.RequestHandler
	handler, ok = srv.apiHandlers[*arwoa.Action]
	if !ok {
		srv.respondNotFound(rw)
		return
	}

	var token *string
	token, err = gm.GetToken(req)
	if err != nil {
		srv.respondBadRequest(rw)
		return
	}

	var ar = &api.Request{
		Action:     arwoa.Action,
		Parameters: arwoa.Parameters,
		Authorisation: &cmr.Auth{
			UserIPA: clientIPA,
		},
	}

	if token != nil {
		ar.Authorisation.Token = *token
	}

	handler(ar, req, rw)
	return
}
