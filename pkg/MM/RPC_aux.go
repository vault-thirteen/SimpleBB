package mm

import (
	"context"
	"fmt"
	"sync"

	js "github.com/osamingo/jsonrpc/v2"
	ac "github.com/vault-thirteen/SimpleBB/pkg/ACM/client"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	"github.com/vault-thirteen/SimpleBB/pkg/net"
)

// Auxiliary functions used in RPC functions.

// mustBeAnAuthToken ensures that an authorization token is present and is
// valid. If the token is absent or invalid, an error is returned and the caller
// of this function must stop and return this error. User data is returned when
// token is valid.
func (srv *Server) mustBeAnAuthToken(auth *cm.Auth) (userRoles *am.GetSelfRolesResult, jerr *js.Error) {
	jerr = srv.mustBeAuthUserIPA(auth)
	if jerr != nil {
		return nil, jerr
	}

	if len(auth.Token) == 0 {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	var err error
	userRoles, err = srv.getUserSelfRoles(auth)
	if err != nil {
		jerr = &js.Error{Code: c.RpcErrorCode_GetUserDataByAuthToken, Message: fmt.Sprintf(c.RpcErrorMsgF_GetUserDataByAuthToken, err.Error())}
		return nil, jerr
	}

	return userRoles, nil
}

// mustBeAuthUserIPA ensures that user's IP address is set. If it is not set,
// an error is returned and the caller of this function must stop and return
// this error.
func (srv *Server) mustBeAuthUserIPA(auth *cm.Auth) (jerr *js.Error) {
	if auth == nil {
		return &js.Error{Code: c.RpcErrorCode_MalformedRequest, Message: c.RpcErrorMsg_MalformedRequest}
	}

	if len(auth.UserIPA) == 0 {
		return &js.Error{Code: c.RpcErrorCode_MalformedRequest, Message: c.RpcErrorMsg_MalformedRequest}
	}

	var err error
	auth.UserIPAB, err = net.ParseIPA(auth.UserIPA)
	if err != nil {
		return &js.Error{Code: c.RpcErrorCode_IPAddressError, Message: fmt.Sprintf(c.RpcErrorMsgF_IPAddressError, err.Error())}
	}

	return nil
}

func (srv *Server) getUserSelfRoles(auth *cm.Auth) (userRoles *am.GetSelfRolesResult, err error) {
	var params = am.GetSelfRolesParams{
		CommonParams: cm.CommonParams{
			Auth: auth,
		},
	}

	userRoles = &am.GetSelfRolesResult{}
	err = srv.acmServiceClient.MakeRequest(context.Background(), userRoles, ac.FuncGetSelfRoles, params)
	if err != nil {
		return nil, err
	}

	return userRoles, nil
}

func (srv *Server) doTestA(wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()

	var ap = am.TestParams{}
	var ar = am.TestResult{}
	err := srv.acmServiceClient.MakeRequest(context.Background(), &ar, ac.FuncTest, ap)
	if err != nil {
		errChan <- err
	}
}
