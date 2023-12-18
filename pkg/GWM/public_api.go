package gwm

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/vault-thirteen/SimpleBB/pkg/GWM/models"
)

func (srv *Server) handlePublicApi(rw http.ResponseWriter, req *http.Request) {
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
	var arwoa models.ApiRequestWithOnlyAction
	err = json.Unmarshal(reqBody, &arwoa)
	if err != nil {
		srv.processBadRequest(rw)
		return
	}

	// switch arwoa.Action {}

	//TODO
	_, err = rw.Write([]byte("TODO:API"))
	if err != nil {
		log.Println(err)
	}
}
