package sm

import (
	"fmt"
	"sync"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	sm "github.com/vault-thirteen/SimpleBB/pkg/SM/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

// RPC functions.

// Subscription.

// addSubscription creates a subscription.
func (srv *Server) addSubscription(p *sm.AddSubscriptionParams) (result *sm.AddSubscriptionResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}
	if userRoles.User.Id != p.UserId {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// Check existence of the thread.
	var threadExists cmb.Flag
	threadExists, re = srv.checkIfThreadExists(p.ThreadId)
	if re != nil {
		return nil, re
	}

	// If thread does not exist, we can not subscribe to it.
	if !threadExists {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadDoesNotExist, RpcErrorMsg_ThreadDoesNotExist, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Read Subscriptions. If they are not initialised, initialise them. Other
	// methods reading subscriptions should not initialise un-initialised
	// subscriptions.
	var usr *sm.UserSubscriptionsRecord
	var err error
	usr, err = srv.dbo.GetUserSubscriptions(p.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// If no real record exists, create it.
	if usr.Id == sm.IdForVirtualUserSubscriptionsRecord {
		err = srv.dbo.InitUserSubscriptions(p.UserId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		usr, err = srv.dbo.GetUserSubscriptions(p.UserId)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	var tsr *sm.ThreadSubscriptionsRecord
	tsr, err = srv.dbo.GetThreadSubscriptions(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// If no real record exists, create it.
	if tsr == nil {
		err = srv.dbo.InitThreadSubscriptions(p.ThreadId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		tsr, err = srv.dbo.GetThreadSubscriptions(p.ThreadId)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	// Add items.
	err = usr.Threads.AddItem(p.ThreadId, false)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = tsr.Users.AddItem(p.UserId, false)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	// Save changes.
	err = srv.dbo.SaveUserSubscriptions(usr)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = srv.dbo.SaveThreadSubscriptions(tsr)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &sm.AddSubscriptionResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// isSelfSubscribed checks whether the caller user has a subscription to the thread.
func (srv *Server) isSelfSubscribed(p *sm.IsSelfSubscribedParams) (result *sm.IsSelfSubscribedResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	var isSubscribed cmb.Flag
	isSubscribed, re = srv.isUserSubscribedH(userRoles.User.Id, p.ThreadId)
	if re != nil {
		return nil, re
	}

	result = &sm.IsSelfSubscribedResult{
		UserId:       userRoles.User.Id,
		ThreadId:     p.ThreadId,
		IsSubscribed: isSubscribed,
	}

	return result, nil
}

// isUserSubscribed checks whether the user has a subscription to the thread.
func (srv *Server) isUserSubscribed(p *sm.IsUserSubscribedParams) (result *sm.IsUserSubscribedResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}
	if userRoles.User.Id != p.UserId {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	var isSubscribed cmb.Flag
	isSubscribed, re = srv.isUserSubscribedH(p.UserId, p.ThreadId)
	if re != nil {
		return nil, re
	}

	result = &sm.IsUserSubscribedResult{
		UserId:       p.UserId,
		ThreadId:     p.ThreadId,
		IsSubscribed: isSubscribed,
	}

	return result, nil
}

// countSelfSubscriptions counts subscriptions of the current user.
func (srv *Server) countSelfSubscriptions(p *sm.CountSelfSubscriptionsParams) (result *sm.CountSelfSubscriptionsResult, re *jrm1.RpcError) {
	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	result = &sm.CountSelfSubscriptionsResult{}

	result.UserSubscriptionsCount, re = srv.countSelfSubscriptionsH(userRoles.User.Id)
	if re != nil {
		return nil, re
	}

	return result, nil
}

// getSelfSubscriptions reads subscriptions of the current user.
func (srv *Server) getSelfSubscriptions(p *sm.GetSelfSubscriptionsParams) (result *sm.GetSelfSubscriptionsResult, re *jrm1.RpcError) {
	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	var usr *sm.UserSubscriptionsRecord
	usr, re = srv.getUserSubscriptionsRecordH(userRoles.User.Id)
	if re != nil {
		return nil, re
	}

	result = &sm.GetSelfSubscriptionsResult{
		UserSubscriptions: sm.NewUserSubscriptions(userRoles.User.Id, usr.Threads, 0, 0),
	}

	return result, nil
}

// getSelfSubscriptionsOnPage reads subscriptions of the current user on the selected page.
func (srv *Server) getSelfSubscriptionsOnPage(p *sm.GetSelfSubscriptionsOnPageParams) (result *sm.GetSelfSubscriptionsOnPageResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.Page == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_PageIsNotSet, RpcErrorMsg_PageIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	var usr *sm.UserSubscriptionsRecord
	usr, re = srv.getUserSubscriptionsRecordH(userRoles.User.Id)
	if re != nil {
		return nil, re
	}

	result = &sm.GetSelfSubscriptionsOnPageResult{
		UserSubscriptions: sm.NewUserSubscriptions(userRoles.User.Id, usr.Threads, p.Page, srv.settings.SystemSettings.PageSize),
	}

	return result, nil
}

// getUserSubscriptions reads user subscriptions.
func (srv *Server) getUserSubscriptions(p *sm.GetUserSubscriptionsParams) (result *sm.GetUserSubscriptionsResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}
	if userRoles.User.Id != p.UserId {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	var usr *sm.UserSubscriptionsRecord
	usr, re = srv.getUserSubscriptionsRecordH(p.UserId)
	if re != nil {
		return nil, re
	}

	result = &sm.GetUserSubscriptionsResult{
		UserSubscriptions: sm.NewUserSubscriptions(p.UserId, usr.Threads, 0, 0),
	}

	return result, nil
}

// getUserSubscriptionsOnPage reads user subscriptions on the selected page.
func (srv *Server) getUserSubscriptionsOnPage(p *sm.GetUserSubscriptionsOnPageParams) (result *sm.GetUserSubscriptionsOnPageResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}
	if p.Page == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_PageIsNotSet, RpcErrorMsg_PageIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}
	if userRoles.User.Id != p.UserId {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	var usr *sm.UserSubscriptionsRecord
	usr, re = srv.getUserSubscriptionsRecordH(p.UserId)
	if re != nil {
		return nil, re
	}

	result = &sm.GetUserSubscriptionsOnPageResult{
		UserSubscriptions: sm.NewUserSubscriptions(p.UserId, usr.Threads, p.Page, srv.settings.SystemSettings.PageSize),
	}

	return result, nil
}

// getThreadSubscribersS reads a list of users subscribed to the specified
// thread. This method is used by the system.
func (srv *Server) getThreadSubscribersS(p *sm.GetThreadSubscribersSParams) (result *sm.GetThreadSubscribersSResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	re = srv.mustBeNoAuth(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check the DKey.
	if !srv.dKeyI.CheckString(p.DKey.ToString()) {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// Check existence of the thread.
	var threadExists cmb.Flag
	threadExists, re = srv.checkIfThreadExists(p.ThreadId)
	if re != nil {
		return nil, re
	}

	// If thread does not exist, we do not even set the thread ID in result.
	if !threadExists {
		result = &sm.GetThreadSubscribersSResult{
			ThreadSubscriptions: nil,
		}
		return result, nil
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read Subscriptions.
	var tsr *sm.ThreadSubscriptionsRecord
	var err error
	tsr, err = srv.dbo.GetThreadSubscriptions(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &sm.GetThreadSubscribersSResult{
		ThreadSubscriptions: tsr,
	}

	return result, nil
}

// deleteSelfSubscription deletes a subscription of the caller user.
func (srv *Server) deleteSelfSubscription(p *sm.DeleteSelfSubscriptionParams) (result *sm.DeleteSelfSubscriptionResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// Delete subscription.
	var s = &sm.Subscription{
		ThreadId: p.ThreadId,
		UserId:   userRoles.User.Id,
	}
	re = srv.deleteSubscriptionH(s)
	if re != nil {
		return nil, re
	}

	result = &sm.DeleteSelfSubscriptionResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// deleteSubscription deletes a subscription.
func (srv *Server) deleteSubscription(p *sm.DeleteSubscriptionParams) (result *sm.DeleteSubscriptionResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}
	if userRoles.User.Id != p.UserId {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// Delete subscription.
	var s = &sm.Subscription{
		ThreadId: p.ThreadId,
		UserId:   p.UserId,
	}
	re = srv.deleteSubscriptionH(s)
	if re != nil {
		return nil, re
	}

	result = &sm.DeleteSubscriptionResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// deleteSubscriptionS deletes a subscription. This method is used by the
// system.
func (srv *Server) deleteSubscriptionS(p *sm.DeleteSubscriptionSParams) (result *sm.DeleteSubscriptionSResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	re = srv.mustBeNoAuth(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check the DKey.
	if !srv.dKeyI.CheckString(p.DKey.ToString()) {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// Delete subscription.
	var s = &sm.Subscription{
		ThreadId: p.ThreadId,
		UserId:   p.UserId,
	}
	re = srv.deleteSubscriptionH(s)
	if re != nil {
		return nil, re
	}

	result = &sm.DeleteSubscriptionSResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// clearThreadSubscriptionsS clears remains of subscriptions of a deleted
// thread. This method is used by the system.
func (srv *Server) clearThreadSubscriptionsS(p *sm.ClearThreadSubscriptionsSParams) (result *sm.ClearThreadSubscriptionsSResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	re = srv.mustBeNoAuth(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check the DKey.
	if !srv.dKeyI.CheckString(p.DKey.ToString()) {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// Check existence of the thread.
	var threadExists cmb.Flag
	threadExists, re = srv.checkIfThreadExists(p.ThreadId)
	if re != nil {
		return nil, re
	}

	// If thread exists, we can not clear anything.
	if threadExists {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadExists, RpcErrorMsg_ThreadExists, nil)
	}

	// Clear subscriptions.
	re = srv.clearThreadSubscriptionsH(p.ThreadId)
	if re != nil {
		return nil, re
	}

	result = &sm.ClearThreadSubscriptionsSResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// Other.

func (srv *Server) getDKey(p *sm.GetDKeyParams) (result *sm.GetDKeyResult, re *jrm1.RpcError) {
	re = srv.mustBeNoAuth(p.Auth)
	if re != nil {
		return nil, re
	}

	result = &sm.GetDKeyResult{
		DKey: cmb.Text(srv.dKeyI.GetString()),
	}

	return result, nil
}

func (srv *Server) showDiagnosticData() (result *sm.ShowDiagnosticDataResult, re *jrm1.RpcError) {
	trc, src := srv.js.GetRequestsCount()

	result = &sm.ShowDiagnosticDataResult{
		RequestsCount: cmr.RequestsCount{
			TotalRequestsCount:      cmb.Text(trc),
			SuccessfulRequestsCount: cmb.Text(src),
		},
	}

	return result, nil
}

func (srv *Server) test(p *sm.TestParams) (result *sm.TestResult, re *jrm1.RpcError) {
	result = &sm.TestResult{}

	var wg = new(sync.WaitGroup)
	var errChan = make(chan error, p.N)

	for i := uint(1); i <= p.N; i++ {
		wg.Add(1)
		go srv.doTestA(wg, errChan)
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			srv.logError(err)
			return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_TestError, fmt.Sprintf(RpcErrorMsgF_TestError, err.Error()), nil)
		}
	}

	return result, nil
}
