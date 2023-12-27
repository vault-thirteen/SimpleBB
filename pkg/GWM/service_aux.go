package gwm

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	s "github.com/vault-thirteen/SimpleBB/pkg/GWM/settings"
	ch "github.com/vault-thirteen/SimpleBB/pkg/common/http"
	cn "github.com/vault-thirteen/SimpleBB/pkg/common/net"
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
		host, _, err = cn.SplitHostPort(req.RemoteAddr)
		if err != nil {
			return "", err
		}

		return host, nil

	case s.ClientIPAddressSource_CustomHeader:
		host, err = ch.GetSingleHeader(req, srv.settings.SystemSettings.ClientIPAddressHeader)
		if err != nil {
			return "", err
		}

		return host, nil

	default:
		return "", errors.New(s.ErrUnknownClientIPAddressSource)
	}
}

func (srv *Server) getAcmHttpStatusCodeByRpcErrorCode(rpcErrorCode int) (httpStatusCode int, err error) {
	var ok bool
	httpStatusCode, ok = srv.acmHttpStatusCodesByRpcErrorCode[rpcErrorCode]
	if !ok {
		return 0, fmt.Errorf(ErrFUnknownRpcErrorCode, rpcErrorCode)
	}

	return httpStatusCode, nil
}
