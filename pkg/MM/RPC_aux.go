package mm

import (
	"context"
	"fmt"
	"hash/crc32"
	"log"
	"sync"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	ac "github.com/vault-thirteen/SimpleBB/pkg/ACM/client"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	cn "github.com/vault-thirteen/SimpleBB/pkg/common/net"
)

// Auxiliary functions used in RPC functions.

// logError logs error if debug mode is enabled.
func (srv *Server) logError(err error) {
	if err == nil {
		return
	}

	if srv.settings.SystemSettings.IsDebugMode {
		log.Println(err)
	}
}

// databaseError processes the database error and returns a JSON RPC error.
func (srv *Server) databaseError(err error) (re *jrm1.RpcError) {
	if err == nil {
		return nil
	}

	if c.IsNetworkError(err) {
		log.Println(fmt.Sprintf(c.ErrFDatabaseNetwork, err.Error()))
		*(srv.dbErrors) <- err
	} else {
		srv.logError(err)
	}

	return jrm1.NewRpcErrorByUser(c.RpcErrorCode_DatabaseError, c.RpcErrorMsg_DatabaseError, nil)
}

// Token-related functions.

// mustBeAuthUserIPA ensures that user's IP address is set. If it is not set,
// an error is returned and the caller of this function must stop and return
// this error.
func (srv *Server) mustBeAuthUserIPA(auth *cmr.Auth) (re *jrm1.RpcError) {
	if auth == nil {
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_MalformedRequest, c.RpcErrorMsg_MalformedRequest, nil)
	}

	if len(auth.UserIPA) == 0 {
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_MalformedRequest, c.RpcErrorMsg_MalformedRequest, nil)
	}

	var err error
	auth.UserIPAB, err = cn.ParseIPA(auth.UserIPA)
	if err != nil {
		srv.logError(err)
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_IPAddressError, fmt.Sprintf(c.RpcErrorMsgF_IPAddressError, err.Error()), nil)
	}

	return nil
}

// mustBeAnAuthToken ensures that an authorisation token is present and is
// valid. If the token is absent or invalid, an error is returned and the caller
// of this function must stop and return this error. User data is returned when
// token is valid.
func (srv *Server) mustBeAnAuthToken(auth *cmr.Auth) (userRoles *am.GetSelfRolesResult, re *jrm1.RpcError) {
	re = srv.mustBeAuthUserIPA(auth)
	if re != nil {
		return nil, re
	}

	if len(auth.Token) == 0 {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_InsufficientPermission, c.RpcErrorMsg_InsufficientPermission, nil)
	}

	var err error
	userRoles, err = srv.getUserSelfRoles(auth)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_GetUserDataByAuthToken, fmt.Sprintf(c.RpcErrorMsgF_GetUserDataByAuthToken, err.Error()), nil)
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
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncGetSelfRoles, params, userRoles)
	if err != nil {
		return nil, err
	}
	if re != nil {
		return nil, re.AsError()
	}

	return userRoles, nil
}

func (srv *Server) doTestA(wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()

	var ap = am.TestParams{}

	var ar = new(am.TestResult)
	re, err := srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncTest, ap, ar)
	if err != nil {
		errChan <- err
	}
	if re != nil {
		errChan <- re.AsError()
	}
}

func (srv *Server) getMessageTextChecksum(msgText string) (checksum uint32) {
	return crc32.Checksum([]byte(msgText), srv.crcTable)
}

func (srv *Server) checkMessageTextChecksum(msgText string, checksum uint32) (ok bool) {
	return srv.getMessageTextChecksum(msgText) == checksum
}
