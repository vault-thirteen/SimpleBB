package gwm

import (
	"encoding/json"
	"log"
	"net/http"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	ch "github.com/vault-thirteen/SimpleBB/pkg/common/http"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	"github.com/vault-thirteen/auxie/header"
)

// processInternalServerError logs the internal error and responds via HTTP.
func (srv *Server) processInternalServerError(rw http.ResponseWriter, err error) {
	srv.logError(err)
	rw.WriteHeader(http.StatusInternalServerError)
}

// processRpcError finds an appropriate HTTP status code and message text for
// an RPC error and responds via HTTP.
func (srv *Server) processRpcError(moduleId byte, re *jrm1.RpcError, rw http.ResponseWriter) {
	var httpStatusCode int
	var err error
	httpStatusCode, err = srv.getHttpStatusCodeByRpcErrorCode(moduleId, re.Code.Int())
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

// respondWithPlainText responds via HTTP with a simple text message.
func (srv *Server) respondWithPlainText(rw http.ResponseWriter, text string) {
	rw.Header().Set(header.HttpHeaderContentType, ch.ContentType_TextPlain)

	_, err := rw.Write([]byte(text))
	if err != nil {
		log.Println(err.Error())
		return
	}
}

// respondWithJsonObject responds via HTTP with a JSON object.
func (srv *Server) respondWithJsonObject(rw http.ResponseWriter, obj any) {
	rw.Header().Set(header.HttpHeaderContentType, ch.ContentType_ApplicationJson)

	err := json.NewEncoder(rw).Encode(obj)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (srv *Server) respondBadRequest(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusBadRequest)
}

func (srv *Server) respondForbidden(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusForbidden)
}

func (srv *Server) respondNotFound(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusNotFound)
}

func (srv *Server) respondMethodNotAllowed(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (srv *Server) respondNotAcceptable(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusNotAcceptable)
}

func (srv *Server) setTokenCookie(rw http.ResponseWriter, token string) {
	var c = &http.Cookie{
		Name:     cm.CookieName_Token,
		Value:    token,
		MaxAge:   int(srv.settings.SystemSettings.SessionMaxDuration),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	}

	ch.SetCookie(rw, c)
}

func (srv *Server) clearTokenCookie(rw http.ResponseWriter) {
	ch.UnsetCookie(rw, cm.CookieName_Token)
}
