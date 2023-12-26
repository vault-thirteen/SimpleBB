package gwm

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	s "github.com/vault-thirteen/SimpleBB/pkg/GWM/settings"
	ch "github.com/vault-thirteen/SimpleBB/pkg/common/http"
	netty "github.com/vault-thirteen/SimpleBB/pkg/net"
	"github.com/vault-thirteen/auxie/MIME"
	"github.com/vault-thirteen/auxie/header"
	hh "github.com/vault-thirteen/auxie/http-helper"
	jc "github.com/ybbus/jsonrpc/v3"
)

// Auxiliary functions used in service functions.

const (
	ErrNotEnoughDataInAddress = "not enough data in address"
	ErrFMultipleHeaders       = "multiple headers: %s"
	ErrFHeaderIsNotFound      = "header is not found: %s"
	ErrFUnknownRpcErrorCode   = "unknown RPC error code: %v"
	ErrTypeCast               = "type cast error"
)

func (srv *Server) isIPAddressAllowed(req *http.Request) (ok bool, clientIPA string, err error) {
	clientIPA, err = srv.getClientIPAddress(req)
	if err != nil {
		return false, "", err
	}

	var ipa net.IP
	ipa, err = netty.ParseIPA(clientIPA)
	if err != nil {
		return false, "", err
	}

	var n int
	n, err = srv.dbo.CountBlocksByIPAddress(ipa)
	if err != nil {
		return false, "", srv.databaseError(err)
	}

	if n == 0 {
		return true, clientIPA, nil
	}

	return false, clientIPA, nil
}

func (srv *Server) getClientIPAddress(req *http.Request) (cipa string, err error) {
	var host string

	switch srv.settings.SystemSettings.ClientIPAddressSource {
	case s.ClientIPAddressSource_Direct:
		host, _, err = splitHostPort(req.RemoteAddr)
		if err != nil {
			return "", err
		}

		return host, nil

	case s.ClientIPAddressSource_CustomHeader:
		host, err = getSingleHeader(req, srv.settings.SystemSettings.ClientIPAddressHeader)
		if err != nil {
			return "", err
		}

		return host, nil

	default:
		return "", errors.New(s.ErrUnknownClientIPAddressSource)
	}
}

// While Go language does not have this very basic function, we are
// re-inventing the wheel again and again.
func splitHostPort(addr string) (host, port string, err error) {
	parts := strings.Split(addr, ":")

	if len(parts) != 2 {
		return "", "", errors.New(ErrNotEnoughDataInAddress)
	}

	return parts[0], parts[1], nil
}

// getSingleHeader reads exactly one, single HTTP header. If multiple headers
// are found, an error is returned. Unfortunately, Go language does not do this
// by default and returns only a first header value even when there are
// multiple values available. Such a behaviour may lead to unexpected errors.
func getSingleHeader(req *http.Request, headerName string) (h string, err error) {
	headers := req.Header[headerName]

	if len(headers) > 1 {
		return "", fmt.Errorf(ErrFMultipleHeaders, headerName)
	}

	if len(headers) == 0 {
		return "", fmt.Errorf(ErrFHeaderIsNotFound, headerName)
	}

	return headers[0], nil
}

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

func (srv *Server) checkClientSupportForJson(req *http.Request) (ok bool, err error) {
	var httpAcceptHeader string
	httpAcceptHeader, err = getSingleHeader(req, header.HttpHeaderAccept)
	if err != nil {
		return false, err
	}

	var amts *hh.AcceptedMimeTypes
	amts, err = hh.NewAcceptedMimeTypesFromHeader(httpAcceptHeader)
	if err != nil {
		return false, err
	}

	var amt *hh.AcceptedMimeType
	for {
		amt, err = amts.Next()
		if err != nil {
			break
		}

		switch amt.MimeType {
		case mime.TypeApplicationJson,
			mime.TypeApplicationAny,
			mime.TypeAny:
			return true, nil
		}
	}

	return false, nil
}

func (srv *Server) getAcmHttpStatusCodeByRpcErrorCode(rpcErrorCode int) (httpStatusCode int, err error) {
	var ok bool
	httpStatusCode, ok = srv.acmHttpStatusCodesByRpcErrorCode[rpcErrorCode]
	if !ok {
		return 0, fmt.Errorf(ErrFUnknownRpcErrorCode, rpcErrorCode)
	}

	return httpStatusCode, nil
}

func (srv *Server) processRpcError(rw http.ResponseWriter, jerr *jc.RPCError) {
	var httpStatusCode int
	var err error
	httpStatusCode, err = srv.getAcmHttpStatusCodeByRpcErrorCode(jerr.Code)
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
	srv.respondWithPlainText(rw, jerr.Message)
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
