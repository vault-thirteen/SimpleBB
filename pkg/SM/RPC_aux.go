package sm

import (
	"context"
	"fmt"
	"log"
	"sync"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	ac "github.com/vault-thirteen/SimpleBB/pkg/ACM/client"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	mc "github.com/vault-thirteen/SimpleBB/pkg/MM/client"
	mm "github.com/vault-thirteen/SimpleBB/pkg/MM/models"
	sm "github.com/vault-thirteen/SimpleBB/pkg/SM/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	ul "github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
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

// processDatabaseError processes a database error.
func (srv *Server) processDatabaseError(err error) {
	if err == nil {
		return
	}

	if c.IsNetworkError(err) {
		log.Println(fmt.Sprintf(c.ErrFDatabaseNetwork, err.Error()))
		*(srv.dbErrors) <- err
	} else {
		srv.logError(err)
	}

	return
}

// databaseError processes the database error and returns an RPC error.
func (srv *Server) databaseError(err error) (re *jrm1.RpcError) {
	srv.processDatabaseError(err)
	return jrm1.NewRpcErrorByUser(c.RpcErrorCode_Database, c.RpcErrorMsg_Database, err)
}

// Token-related functions.

// mustBeNoAuth ensures that authorisation is not used.
func (srv *Server) mustBeNoAuth(auth *cmr.Auth) (re *jrm1.RpcError) {
	if auth != nil {
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	return nil
}

// mustBeAuthUserIPA ensures that user's IP address is set. If it is not set,
// an error is returned and the caller of this function must stop and return
// this error.
func (srv *Server) mustBeAuthUserIPA(auth *cmr.Auth) (re *jrm1.RpcError) {
	if (auth == nil) ||
		(len(auth.UserIPA) == 0) {
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_Authorisation, c.RpcErrorMsg_Authorisation, nil)
	}

	var err error
	auth.UserIPAB, err = cn.ParseIPA(auth.UserIPA)
	if err != nil {
		srv.logError(err)
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_Authorisation, c.RpcErrorMsg_Authorisation, nil)
	}

	return nil
}

// mustBeNoAuthToken ensures that an authorisation token is not present. If the
// token is present, an error is returned and the caller of this function must
// stop and return this error.
func (srv *Server) mustBeNoAuthToken(auth *cmr.Auth) (re *jrm1.RpcError) {
	re = srv.mustBeAuthUserIPA(auth)
	if re != nil {
		return re
	}

	if len(auth.Token) > 0 {
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
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
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Authorisation, c.RpcErrorMsg_Authorisation, nil)
	}

	userRoles, re = srv.getUserSelfRoles(auth)
	if re != nil {
		return nil, re
	}

	return userRoles, nil
}

// Other functions.

// getUserSelfRoles reads roles of the RPC caller (user) from ACM module.
func (srv *Server) getUserSelfRoles(auth *cmr.Auth) (userRoles *am.GetSelfRolesResult, re *jrm1.RpcError) {
	var params = am.GetSelfRolesParams{
		CommonParams: cmr.CommonParams{
			Auth: auth,
		},
	}

	userRoles = new(am.GetSelfRolesResult)
	var err error
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncGetSelfRoles, params, userRoles)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_RPCCall, c.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Authorisation, c.RpcErrorMsg_Authorisation, nil)
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

// getDKeyForMM receives a DKey from Message module.
func (srv *Server) getDKeyForMM() (dKey *string, re *jrm1.RpcError) {
	params := mm.GetDKeyParams{}
	result := new(mm.GetDKeyResult)
	var err error
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncGetDKey, params, result)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_RPCCall, c.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return nil, re
	}

	// DKey must be non-empty.
	if len(result.DKey) == 0 {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_ModuleSynchronisation, c.RpcErrorMsg_ModuleSynchronisation, nil)
	}

	return &result.DKey, nil
}

// deleteSubscriptionH is a common function to delete subscriptions.
func (srv *Server) deleteSubscriptionH(s *sm.Subscription) (re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Read Subscriptions.
	var us *sm.UserSubscriptions
	var err error
	us, err = srv.dbo.GetUserSubscriptions(s.UserId)
	if err != nil {
		return srv.databaseError(err)
	}
	if us == nil {
		return jrm1.NewRpcErrorByUser(RpcErrorCode_SubscriptionIsNotFound, RpcErrorMsg_SubscriptionIsNotFound, nil)
	}

	var ts *sm.ThreadSubscriptions
	ts, err = srv.dbo.GetThreadSubscriptions(s.ThreadId)
	if err != nil {
		return srv.databaseError(err)
	}
	if ts == nil {
		return jrm1.NewRpcErrorByUser(RpcErrorCode_SubscriptionIsNotFound, RpcErrorMsg_SubscriptionIsNotFound, nil)
	}

	// Remove items.
	err = us.Threads.RemoveItem(s.ThreadId)
	if err != nil {
		srv.logError(err)
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = ts.Users.RemoveItem(s.UserId)
	if err != nil {
		srv.logError(err)
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	// Save changes.
	err = srv.dbo.SaveUserSubscriptions(us)
	if err != nil {
		return srv.databaseError(err)
	}

	err = srv.dbo.SaveThreadSubscriptions(ts)
	if err != nil {
		return srv.databaseError(err)
	}

	return nil
}

// clearThreadSubscriptionsH is a helper function for clearing subscriptions of
// a deleted thread.
func (srv *Server) clearThreadSubscriptionsH(threadId uint) (re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var ts *sm.ThreadSubscriptions
	var err error
	ts, err = srv.dbo.GetThreadSubscriptions(threadId)
	if err != nil {
		return srv.databaseError(err)
	}
	if ts == nil {
		// Nothing to clear.
		return nil
	}

	if ts.Users == nil {
		// No subscribers.
		// Clear the T.S. record only.
		err = srv.dbo.ClearThreadSubscriptionRecord(threadId)
		if err != nil {
			return srv.databaseError(err)
		}

		return nil
	}

	// Delete subscription to the thread from all subscribers.
	var us *sm.UserSubscriptions
	for _, userId := range *ts.Users {
		us, err = srv.dbo.GetUserSubscriptions(userId)
		if err != nil {
			return srv.databaseError(err)
		}

		if us.Threads == nil {
			continue
		}

		err = us.Threads.RemoveItem(threadId)
		if err != nil {
			srv.logError(err)
			return jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
		}

		err = srv.dbo.SaveUserSubscriptions(us)
		if err != nil {
			return srv.databaseError(err)
		}
	}

	err = srv.dbo.ClearThreadSubscriptionRecord(threadId)
	if err != nil {
		return srv.databaseError(err)
	}

	return nil
}

// checkIfThreadExists checks if the thread exists or not.
func (srv *Server) checkIfThreadExists(threadId uint) (exists bool, re *jrm1.RpcError) {
	params := mm.ThreadExistsSParams{
		DKeyParams: cmr.DKeyParams{
			// DKey is set during module start-up, so it is non-null.
			DKey: *srv.dKeyForMM,
		},
		ThreadId: threadId,
	}
	result := new(mm.ThreadExistsSResult)
	var err error
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncThreadExistsS, params, result)
	if err != nil {
		srv.logError(err)
		return false, jrm1.NewRpcErrorByUser(c.RpcErrorCode_RPCCall, c.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return false, re
	}

	return result.Exists, nil
}

// getUserSubscriptionsH is a helper function to read user's subscriptions.
func (srv *Server) getUserSubscriptionsH(userId uint) (us *sm.UserSubscriptions, re *jrm1.RpcError) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read Subscriptions.
	var err error
	us, err = srv.dbo.GetUserSubscriptions(userId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return us, nil
}

// getUserSubscriptionsOnPageH is a helper function to read user's
// subscriptions on the selected page.
func (srv *Server) getUserSubscriptionsOnPageH(userId uint, page uint) (sop *sm.SubscriptionsOnPage, re *jrm1.RpcError) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read all user's subscriptions.
	var us *sm.UserSubscriptions
	var err error
	us, err = srv.dbo.GetUserSubscriptions(userId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if us == nil {
		us = &sm.UserSubscriptions{
			Id:     0,
			UserId: userId,
		}

		us.Threads, err = ul.NewFromArray([]uint{})
		if err != nil {
			srv.logError(err)
			return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
		}
	}

	sop = sm.NewSubscriptionsOnPage()
	{
		sop.Subscriber = userId

		threadIds := us.Threads.OnPage(page, srv.settings.SystemSettings.PageSize)
		if threadIds != nil {
			sop.Subscriptions = *threadIds
		}

		sop.Page = &page

		allSubscriptionsCountUint := uint(us.Threads.Size())
		tp := c.CalculateTotalPages(allSubscriptionsCountUint, srv.settings.SystemSettings.PageSize)
		sop.TotalPages = &tp

		sop.TotalSubscriptions = &allSubscriptionsCountUint
	}

	return sop, nil
}

// isUserSubscribedH is a helper function to check whether the user has a
// subscription to the thread.
func (srv *Server) isUserSubscribedH(userId uint, threadId uint) (isSubscribed bool, re *jrm1.RpcError) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read Subscriptions.
	var us *sm.UserSubscriptions
	var err error
	us, err = srv.dbo.GetUserSubscriptions(userId)
	if err != nil {
		return false, srv.databaseError(err)
	}

	if us == nil {
		return false, nil
	}

	// Search for an item.
	isSubscribed = us.Threads.HasItem(threadId)

	return isSubscribed, nil
}
