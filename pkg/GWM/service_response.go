package gwm

import (
	"encoding/json"
	"log"
	"net/http"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	ch "github.com/vault-thirteen/SimpleBB/pkg/common/http"
	"github.com/vault-thirteen/auxie/header"
)

func (srv *Server) processBlockedClient(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusForbidden)
}

func (srv *Server) processPageNotFound(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusNotFound)
}

func (srv *Server) processBadRequest(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusBadRequest)
}

func (srv *Server) processMethodNotAllowed(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (srv *Server) processNotAcceptable(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusNotAcceptable)
}

func (srv *Server) processInternalServerError(rw http.ResponseWriter, err error) {
	srv.logError(err)
	rw.WriteHeader(http.StatusInternalServerError)
}

// processRpcError turns RPC error codes of other services (which are
// non-gateway modules) into an HTTP response with an error.
func (srv *Server) processRpcError(rw http.ResponseWriter, re *jrm1.RpcError) {
	var httpStatusCode int
	var err error
	httpStatusCode, err = srv.getHttpStatusCodeByRpcErrorCode(re.Code.Int())
	if err != nil {
		srv.processInternalServerError(rw, err)
		return
	}

	switch httpStatusCode {
	case http.StatusInternalServerError:
		srv.processInternalServerError(rw, err)
		return
	}

	rw.WriteHeader(httpStatusCode)
	srv.respondWithPlainText(rw, re.AsError().Error())
	return
}

func (srv *Server) respondWithPlainText(rw http.ResponseWriter, text string) {
	rw.Header().Set(header.HttpHeaderContentType, ch.ContentType_TextPlain)

	_, err := rw.Write([]byte(text))
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (srv *Server) respondWithJsonObject(rw http.ResponseWriter, obj any) {
	rw.Header().Set(header.HttpHeaderContentType, ch.ContentType_ApplicationJson)

	err := json.NewEncoder(rw).Encode(obj)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
