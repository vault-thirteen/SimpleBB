package gwm

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	s "github.com/vault-thirteen/SimpleBB/pkg/GWM/settings"
	"github.com/vault-thirteen/SimpleBB/pkg/common/app"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cn "github.com/vault-thirteen/SimpleBB/pkg/common/net"
	hh "github.com/vault-thirteen/auxie/http-helper"
)

// Auxiliary functions used in service functions.

const (
	ErrFUnknownRpcErrorCode = "unknown RPC error code: %v"
	ErrTypeCast             = "type cast error"
)

func (srv *Server) isIPAddressAllowed(req *http.Request) (ok bool, clientIPA string, err error) {
	clientIPA, err = srv.getClientIPAddress(req)
	if err != nil {
		return false, "", err
	}

	var ipa net.IP
	ipa, err = cn.ParseIPA(clientIPA)
	if err != nil {
		return false, "", err
	}

	var n int
	n, err = srv.dbo.CountBlocksByIPAddress(ipa)
	if err != nil {
		re := srv.databaseError(err)
		return false, "", re.AsError()
	}

	if n != 0 {
		return false, clientIPA, nil
	}

	return true, clientIPA, nil
}

func (srv *Server) getClientIPAddress(req *http.Request) (cipa string, err error) {
	var host string

	switch srv.settings.SystemSettings.ClientIPAddressSource {
	case s.ClientIPAddressSource_Direct:
		host, _, err = cn.SplitHostPort(req.RemoteAddr)
		if err != nil {
			return "", err
		}

		return host, nil

	case s.ClientIPAddressSource_CustomHeader:
		host, err = hh.GetSingleHttpHeader(req, srv.settings.SystemSettings.ClientIPAddressHeader)
		if err != nil {
			return "", err
		}

		return host, nil

	default:
		return "", errors.New(s.ErrUnknownClientIPAddressSource)
	}
}

func (srv *Server) getHttpStatusCodeByRpcErrorCode(moduleId cm.Module, rpcErrorCode int) (httpStatusCode int, err error) {
	var ok bool

	httpStatusCode, ok = srv.commonHttpStatusCodesByRpcErrorCode[rpcErrorCode]
	if ok {
		return httpStatusCode, nil
	}

	switch moduleId {
	case app.ModuleId_ACM:
		httpStatusCode, ok = srv.acmHttpStatusCodesByRpcErrorCode[rpcErrorCode]
		if ok {
			return httpStatusCode, nil
		}

	case app.ModuleId_MM:
		httpStatusCode, ok = srv.mmHttpStatusCodesByRpcErrorCode[rpcErrorCode]
		if ok {
			return httpStatusCode, nil
		}

	case app.ModuleId_NM:
		httpStatusCode, ok = srv.nmHttpStatusCodesByRpcErrorCode[rpcErrorCode]
		if ok {
			return httpStatusCode, nil
		}

	case app.ModuleId_SM:
		httpStatusCode, ok = srv.smHttpStatusCodesByRpcErrorCode[rpcErrorCode]
		if ok {
			return httpStatusCode, nil
		}
	}

	return 0, fmt.Errorf(ErrFUnknownRpcErrorCode, rpcErrorCode)
}
