package gwm

import (
	"encoding/json"
	"io"
	"net/http"

	gm "github.com/vault-thirteen/SimpleBB/pkg/GWM/models"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

const (
	ApiFunctionName_GetProductVersion = "getProductVersion"
)

func (srv *Server) handlePublicApi(rw http.ResponseWriter, req *http.Request, clientIPA string) {
	if req.Method != http.MethodPost {
		srv.processMethodNotAllowed(rw)
		return
	}

	// Check accepted MIME types.
	ok, err := srv.checkClientSupportForJson(req)
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
		srv.processInternalServerError(rw)
		return
	}

	// Check the action.
	var arwoa gm.ApiRequestWithOnlyAction
	err = json.Unmarshal(reqBody, &arwoa)
	if err != nil {
		srv.processBadRequest(rw)
		return
	}

	var handler gm.ApiRequestHandler
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

	var ar = &gm.ApiRequest{
		Action:     arwoa.Action,
		Parameters: arwoa.Parameters,
		Authorisation: &cm.Auth{
			UserIPA: clientIPA,
		},
	}

	if token != nil {
		ar.Authorisation.Token = *token
	}

	handler(ar, rw, req)
	return
}
