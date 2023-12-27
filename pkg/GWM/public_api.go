package gwm

import (
	"encoding/json"
	"io"
	"net/http"

	gm "github.com/vault-thirteen/SimpleBB/pkg/GWM/models"
	"github.com/vault-thirteen/SimpleBB/pkg/GWM/models/api"
	ch "github.com/vault-thirteen/SimpleBB/pkg/common/http"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

const (
	ApiFunctionName_GetProductVersion      = "getProductVersion"
	ApiFunctionName_RegisterUser           = "registerUser"
	ApiFunctionName_ApproveAndRegisterUser = "approveAndRegisterUser"
	//TODO
)

func (srv *Server) handlePublicApi(rw http.ResponseWriter, req *http.Request, clientIPA string) {
	if req.Method != http.MethodPost {
		srv.processMethodNotAllowed(rw)
		return
	}

	// Check accepted MIME types.
	ok, err := ch.CheckClientSupportForJson(req)
	if err != nil {
		srv.processBadRequest(rw)
		return
	}

	if !ok {
		srv.processNotAcceptable(rw)
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
		srv.processBadRequest(rw)
		return
	}

	var handler api.RequestHandler
	handler, ok = srv.apiHandlers[arwoa.Action]
	if !ok {
		srv.processPageNotFound(rw)
		return
	}

	var token *string
	token, err = gm.GetToken(req)
	if err != nil {
		srv.processBadRequest(rw)
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

	handler(ar, rw, req)
	return
}
