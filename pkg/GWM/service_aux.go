package gwm

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	s "github.com/vault-thirteen/SimpleBB/pkg/GWM/settings"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	netty "github.com/vault-thirteen/SimpleBB/pkg/net"
)

// Auxiliary functions used in service functions.

const (
	HttpHeaderXForwardedFor = "X-Forwarded-For"
)

const (
	ErrNotEnoughDataInAddress = "not enough data in address"
	ErrFMultipleHeaders       = "multiple headers: %s"
	ErrFHeaderIsNotFound      = "header is not found: %s"
)

func (srv *Server) isIPAddressAllowed(req *http.Request) (ok bool, err error) {
	var cipaText string
	cipaText, err = srv.getClientIPAddress(req)
	if err != nil {
		return false, err
	}

	var ipa net.IP
	ipa, err = netty.ParseIPA(cipaText)
	if err != nil {
		return false, err
	}

	var n int
	n, err = srv.dbo.CountBlocksByIPAddress(ipa)
	if err != nil {
		return false, c.DatabaseError(err, srv.dbErrors)
	}

	if n == 0 {
		return true, nil
	}

	return false, nil
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

	case s.ClientIPAddressSource_XFF:
		host, err = getSingleHeader(req, HttpHeaderXForwardedFor)
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

func (srv *Server) processInternalError(rw http.ResponseWriter, err error) {
	log.Println(err.Error())
	rw.WriteHeader(http.StatusInternalServerError)
}

func (srv *Server) processBlockedClient(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusForbidden)
}
