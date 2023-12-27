package mm

import (
	"context"
	"fmt"
	"hash/crc32"
	"log"
	"sync"

	js "github.com/osamingo/jsonrpc/v2"
	ac "github.com/vault-thirteen/SimpleBB/pkg/ACM/client"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	cn "github.com/vault-thirteen/SimpleBB/pkg/common/net"
)

// Auxiliary functions used in RPC functions.

// logError logs error if debug mode is enabled.
func (srv *Server) logError(err error) {
	if srv.settings.SystemSettings.IsDebugMode {
		log.Println(err)
	}
}

// databaseError processes the database error and returns a JSON RPC error.
func (srv *Server) databaseError(err error) (jerr *js.Error) {
	if c.IsNetworkError(err) {
		log.Println(fmt.Sprintf(c.ErrFDatabaseNetwork, err.Error()))
		*(srv.dbErrors) <- err
	} else {
		srv.logError(err)
	}

	return &js.Error{Code: c.RpcErrorCode_DatabaseError, Message: c.RpcErrorMsg_DatabaseError}
}

// Token-related functions.

// mustBeAuthUserIPA ensures that user's IP address is set. If it is not set,
// an error is returned and the caller of this function must stop and return
// this error.
func (srv *Server) mustBeAuthUserIPA(auth *cmr.Auth) (jerr *js.Error) {
	if auth == nil {
		return &js.Error{Code: c.RpcErrorCode_MalformedRequest, Message: c.RpcErrorMsg_MalformedRequest}
	}

	if len(auth.UserIPA) == 0 {
		return &js.Error{Code: c.RpcErrorCode_MalformedRequest, Message: c.RpcErrorMsg_MalformedRequest}
	}

	var err error
	auth.UserIPAB, err = cn.ParseIPA(auth.UserIPA)
	if err != nil {
		srv.logError(err)
		return &js.Error{Code: c.RpcErrorCode_IPAddressError, Message: fmt.Sprintf(c.RpcErrorMsgF_IPAddressError, err.Error())}
	}

	return nil
}

// mustBeAnAuthToken ensures that an authorisation token is present and is
// valid. If the token is absent or invalid, an error is returned and the caller
// of this function must stop and return this error. User data is returned when
// token is valid.
func (srv *Server) mustBeAnAuthToken(auth *cmr.Auth) (userRoles *am.GetSelfRolesResult, jerr *js.Error) {
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
		srv.logError(err)
		return nil, &js.Error{Code: c.RpcErrorCode_GetUserDataByAuthToken, Message: fmt.Sprintf(c.RpcErrorMsgF_GetUserDataByAuthToken, err.Error())}
	}

	return userRoles, nil
}

// Other functions.

func (srv *Server) getUserSelfRoles(auth *cmr.Auth) (userRoles *am.GetSelfRolesResult, err error) {
	var params = am.GetSelfRolesParams{
		CommonParams: cmr.CommonParams{
			Auth: auth,
		},
	}

	userRoles = new(am.GetSelfRolesResult)
	err = srv.acmServiceClient.MakeRequest(context.Background(), userRoles, ac.FuncGetSelfRoles, params)
	if err != nil {
		return nil, err
	}

	return userRoles, nil
}

func (srv *Server) doTestA(wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()

	var ap = am.TestParams{}

	var ar = new(am.TestResult)
	err := srv.acmServiceClient.MakeRequest(context.Background(), ar, ac.FuncTest, ap)
	if err != nil {
		errChan <- err
	}
}

func (srv *Server) getMessageTextChecksum(msgText string) (checksum uint32) {
	return crc32.Checksum([]byte(msgText), srv.crcTable)
}

func (srv *Server) checkMessageTextChecksum(msgText string, checksum uint32) (ok bool) {
	return srv.getMessageTextChecksum(msgText) == checksum
}
