package nm

import (
	"context"
	"fmt"
	"log"
	"sync"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	ac "github.com/vault-thirteen/SimpleBB/pkg/ACM/client"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	nm "github.com/vault-thirteen/SimpleBB/pkg/NM/models"
	sc "github.com/vault-thirteen/SimpleBB/pkg/SM/client"
	sm "github.com/vault-thirteen/SimpleBB/pkg/SM/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
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
		srv.incidentManager.ReportIncident(cm.IncidentType_IllegalAccessAttempt, "", auth.UserIPAB)
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
		srv.incidentManager.ReportIncident(cm.IncidentType_IllegalAccessAttempt, "", nil)
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

// getDKeyForSM receives a DKey from Subscription module.
func (srv *Server) getDKeyForSM() (dKey *cmb.Text, re *jrm1.RpcError) {
	params := sm.GetDKeyParams{}
	result := new(sm.GetDKeyResult)
	var err error
	re, err = srv.smServiceClient.MakeRequest(context.Background(), sc.FuncGetDKey, params, result)
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

func (srv *Server) processSystemEvent_ThreadParentChange(se *cm.SystemEvent) (re *jrm1.RpcError) {
	return srv.sendNotificationsToThreadSubscribers(se)
}

func (srv *Server) processSystemEvent_ThreadNameChange(se *cm.SystemEvent) (re *jrm1.RpcError) {
	return srv.sendNotificationsToThreadSubscribers(se)
}

func (srv *Server) processSystemEvent_ThreadDeletion(se *cm.SystemEvent) (re *jrm1.RpcError) {
	re = srv.sendNotificationsToThreadSubscribers(se)
	if re != nil {
		return re
	}

	// Ask the SM module to clear the subscriptions.
	re = srv.clearSubscriptionsOfDeletedThread(*se.ThreadId)
	if re != nil {
		return re
	}

	return nil
}

func (srv *Server) processSystemEvent_ThreadNewMessage(se *cm.SystemEvent) (re *jrm1.RpcError) {
	return srv.sendNotificationsToThreadSubscribers(se)
}

func (srv *Server) processSystemEvent_ThreadMessageEdit(se *cm.SystemEvent) (re *jrm1.RpcError) {
	return srv.sendNotificationsToThreadSubscribers(se)
}

func (srv *Server) processSystemEvent_ThreadMessageDeletion(se *cm.SystemEvent) (re *jrm1.RpcError) {
	return srv.sendNotificationsToThreadSubscribers(se)
}

func (srv *Server) processSystemEvent_MessageTextEdit(se *cm.SystemEvent) (re *jrm1.RpcError) {
	var isSubscribed cmb.Flag
	isSubscribed, re = srv.isUserSubscribed(*se.ThreadId, *se.Creator)
	if re != nil {
		return re
	}

	// If user is subscribed to the thread, it does not need the second
	// notification about this message.
	if isSubscribed {
		return nil
	}

	// The actor does not need a notification about self action.
	if *se.UserId == *se.Creator {
		return nil
	}

	return srv.sendNotificationToCreator(se)
}

func (srv *Server) processSystemEvent_MessageParentChange(se *cm.SystemEvent) (re *jrm1.RpcError) {
	// The actor does not need a notification about self action.
	if *se.UserId == *se.Creator {
		return nil
	}

	return srv.sendNotificationToCreator(se)
}

func (srv *Server) processSystemEvent_MessageDeletion(se *cm.SystemEvent) (re *jrm1.RpcError) {
	var isSubscribed cmb.Flag
	isSubscribed, re = srv.isUserSubscribed(*se.ThreadId, *se.Creator)
	if re != nil {
		return re
	}

	// If user is subscribed to the thread, it does not need the second
	// notification about this message.
	if isSubscribed {
		return nil
	}

	// The actor does not need a notification about self action.
	if *se.UserId == *se.Creator {
		return nil
	}

	return srv.sendNotificationToCreator(se)
}

// sendNotificationsToThreadSubscribers sends notifications to thread
// subscribers.
func (srv *Server) sendNotificationsToThreadSubscribers(se *cm.SystemEvent) (re *jrm1.RpcError) {
	var tsr *sm.ThreadSubscriptionsRecord
	tsr, re = srv.getThreadSubscribers(*se.ThreadId)
	if re != nil {
		return re
	}

	var notificationText cmb.Text
	notificationText, re = srv.composeNotificationText(se)
	if re != nil {
		return re
	}

	if tsr != nil {
		for _, userId := range tsr.Users.AsArray() {
			// The performer of the action does not need the notification.
			if se.SystemEventData.UserId != nil {
				if userId == *se.SystemEventData.UserId {
					continue
				}
			}

			_, re = srv.sendNotificationIfPossibleH(userId, notificationText)
			if re != nil {
				return re
			}
		}
	}

	return nil
}

// sendNotificationToCreator sends a notification to the initial creator of the
// object.
func (srv *Server) sendNotificationToCreator(se *cm.SystemEvent) (re *jrm1.RpcError) {
	var notificationText cmb.Text
	notificationText, re = srv.composeNotificationText(se)
	if re != nil {
		return re
	}

	_, re = srv.sendNotificationIfPossibleH(*se.Creator, notificationText)
	if re != nil {
		return re
	}

	return nil
}

// composeNotificationText creates a text for notification about the system
// event.
func (srv *Server) composeNotificationText(se *cm.SystemEvent) (text cmb.Text, re *jrm1.RpcError) {
	switch se.Type {
	case cm.SystemEventType_ThreadParentChange:
		// Template: FUT.
		text = cmb.Text(fmt.Sprintf("A user (%d) has moved the thread (%d) into another forum.", *se.UserId, *se.ThreadId))

	case cm.SystemEventType_ThreadNameChange:
		// Template: FUT.
		text = cmb.Text(fmt.Sprintf("A user (%d) has renamed the thread (%d).", *se.UserId, *se.ThreadId))

	case cm.SystemEventType_ThreadDeletion:
		// Template: FUT.
		text = cmb.Text(fmt.Sprintf("A user (%d) has deleted the thread (%d).", *se.UserId, *se.ThreadId))

	case cm.SystemEventType_ThreadNewMessage:
		// Template: FUMT.
		text = cmb.Text(fmt.Sprintf("A user (%d) has added a new message (%d) into the thread (%d).", *se.UserId, *se.MessageId, *se.ThreadId))

	case cm.SystemEventType_ThreadMessageEdit:
		// Template: FUMT.
		text = cmb.Text(fmt.Sprintf("A user (%d) has edited a message (%d) in the thread (%d).", *se.UserId, *se.MessageId, *se.ThreadId))

	case cm.SystemEventType_ThreadMessageDeletion:
		// Template: FUMT.
		text = cmb.Text(fmt.Sprintf("A user (%d) has deleted a message (%d) from the thread (%d).", *se.UserId, *se.MessageId, *se.ThreadId))

	case cm.SystemEventType_MessageTextEdit:
		// Template: FMTU.
		text = cmb.Text(fmt.Sprintf("Your message (%d) in the thread (%d) was edited by a user (%d).", *se.MessageId, *se.ThreadId, *se.UserId))

	case cm.SystemEventType_MessageParentChange:
		// Template: FMTU.
		text = cmb.Text(fmt.Sprintf("Your message (%d) in the thread (%d) was moved into another thread by a user (%d).", *se.MessageId, *se.ThreadId, *se.UserId))

	case cm.SystemEventType_MessageDeletion:
		// Template: FMTU.
		text = cmb.Text(fmt.Sprintf("Your message (%d) in the thread (%d) was deleted by a user (%d).", *se.MessageId, *se.ThreadId, *se.UserId))

	default:
		return "", jrm1.NewRpcErrorByUser(RpcErrorCode_SystemEvent, RpcErrorMsg_SystemEvent, nil)
	}

	return text, nil
}

// clearSubscriptionsOfDeletedThread clears remains of subscriptions of a
// deleted thread.
func (srv *Server) clearSubscriptionsOfDeletedThread(threadId cmb.Id) (re *jrm1.RpcError) {
	params := sm.ClearThreadSubscriptionsSParams{
		DKeyParams: cmr.DKeyParams{
			// DKey is set during module start-up, so it is non-null.
			DKey: *srv.dKeyForSM,
		},
		ThreadId: threadId,
	}
	result := new(sm.ClearThreadSubscriptionsSResult)
	var err error
	re, err = srv.smServiceClient.MakeRequest(context.Background(), sc.FuncClearThreadSubscriptionsS, params, result)
	if err != nil {
		srv.logError(err)
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_RPCCall, c.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return re
	}

	return nil
}

// getThreadSubscribers gets a list of users subscribed to the thread.
func (srv *Server) getThreadSubscribers(threadId cmb.Id) (tsr *sm.ThreadSubscriptionsRecord, re *jrm1.RpcError) {
	params := sm.GetThreadSubscribersSParams{
		DKeyParams: cmr.DKeyParams{
			// DKey is set during module start-up, so it is non-null.
			DKey: *srv.dKeyForSM,
		},
		ThreadId: threadId,
	}
	result := new(sm.GetThreadSubscribersSResult)
	var err error
	re, err = srv.smServiceClient.MakeRequest(context.Background(), sc.FuncGetThreadSubscribersS, params, result)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_RPCCall, c.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return nil, re
	}

	return result.ThreadSubscriptions, nil
}

// isUserSubscribed checks whether the user is subscribed to the thread.
func (srv *Server) isUserSubscribed(threadId cmb.Id, userId cmb.Id) (isSubscribed cmb.Flag, re *jrm1.RpcError) {
	params := sm.IsUserSubscribedSParams{
		DKeyParams: cmr.DKeyParams{
			// DKey is set during module start-up, so it is non-null.
			DKey: *srv.dKeyForSM,
		},
		ThreadId: threadId,
		UserId:   userId,
	}
	result := new(sm.IsUserSubscribedSResult)
	var err error
	re, err = srv.smServiceClient.MakeRequest(context.Background(), sc.FuncIsUserSubscribedS, params, result)
	if err != nil {
		srv.logError(err)
		return false, jrm1.NewRpcErrorByUser(c.RpcErrorCode_RPCCall, c.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return false, re
	}

	return result.IsSubscribed, nil
}

// sendNotificationIfPossibleH is a helper function which tries to send a
// notification to a user when it is possible.
func (srv *Server) sendNotificationIfPossibleH(userId cmb.Id, text cmb.Text) (result *nm.SendNotificationIfPossibleSResult, re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var err error
	var unc cmb.Count
	unc, err = srv.dbo.CountUnreadNotificationsByUserId(userId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Notification box is full.
	if unc >= srv.settings.SystemSettings.NotificationCountLimit {
		result = &nm.SendNotificationIfPossibleSResult{
			IsSent: false,
		}
		return result, nil
	}

	var insertedNotificationId cmb.Id
	insertedNotificationId, err = srv.dbo.InsertNewNotification(userId, text)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &nm.SendNotificationIfPossibleSResult{
		IsSent:         true,
		NotificationId: insertedNotificationId,
	}

	return result, nil
}

// saveSystemEventH saves the system event into database.
func (srv *Server) saveSystemEventH(se *cm.SystemEvent) (re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var err error
	err = srv.dbo.SaveSystemEvent(*se)
	if err != nil {
		return srv.databaseError(err)
	}

	return nil
}
