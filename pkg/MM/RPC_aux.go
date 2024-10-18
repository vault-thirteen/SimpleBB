package mm

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	ac "github.com/vault-thirteen/SimpleBB/pkg/ACM/client"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	mm "github.com/vault-thirteen/SimpleBB/pkg/MM/models"
	nc "github.com/vault-thirteen/SimpleBB/pkg/NM/client"
	nm "github.com/vault-thirteen/SimpleBB/pkg/NM/models"
	sc "github.com/vault-thirteen/SimpleBB/pkg/SM/client"
	sm "github.com/vault-thirteen/SimpleBB/pkg/SM/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	ul "github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	cn "github.com/vault-thirteen/SimpleBB/pkg/common/net"
	ah "github.com/vault-thirteen/auxie/hash"
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

// getDKeyForNM receives a DKey from Notification module.
func (srv *Server) getDKeyForNM() (dKey *cmb.Text, re *jrm1.RpcError) {
	params := nm.GetDKeyParams{}
	result := new(nm.GetDKeyResult)
	var err error
	re, err = srv.nmServiceClient.MakeRequest(context.Background(), nc.FuncGetDKey, params, result)
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

// sendNotificationToUser sends a system notification to a user.
func (srv *Server) sendNotificationToUser(userId cmb.Id, text cmb.Text) (re *jrm1.RpcError) {
	params := nm.AddNotificationSParams{
		DKeyParams: cmr.DKeyParams{
			// DKey is set during module start-up, so it is non-null.
			DKey: *srv.dKeyForNM,
		},
		UserId: userId,
		Text:   text,
	}
	result := new(nm.AddNotificationSResult)
	var err error
	re, err = srv.nmServiceClient.MakeRequest(context.Background(), nc.FuncAddNotificationS, params, result)
	if err != nil {
		srv.logError(err)
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_RPCCall, c.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return re
	}

	return nil
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

func (srv *Server) getMessageTextChecksum(msgText cmb.Text) (checksum []byte) {
	x, _ := ah.CalculateCrc32(msgText.ToBytes())
	return x[:]
}

func (srv *Server) checkMessageTextChecksum(msgText cmb.Text, checksum []byte) (ok bool) {
	return bytes.Compare(srv.getMessageTextChecksum(msgText), checksum) == 0
}

// canUserEditMessage checks whether a user (specified by the 'userRoles'
// argument) can edit a message (specified as an 'message' argument).
func (srv *Server) canUserEditMessage(userRoles *am.GetSelfRolesResult, message *mm.Message) (ok bool) {
	// Are you stupid, kido ?
	if (userRoles == nil) || (message == nil) {
		return false
	}

	// Moderators have extended rights to edit messages of any users.
	if userRoles.User.Roles.IsModerator {
		return true
	}

	// Writers can edit their own messages.
	if !userRoles.User.Roles.IsWriter {
		return false
	}

	// User can not edit messages created by other users.
	if userRoles.User.Id != message.Creator.UserId {
		return false
	}

	// User can edit its own messages if they are not too old.
	if time.Now().Before(srv.getMessageMaxEditTime(message)) {
		return true
	}

	return false
}

// canUserAddMessage checks whether a user (specified by the 'userRoles'
// argument) can add a new message into a thread in case when there is a
// [latest] message in the thread (specified as an 'latestMessageInThread'
// argument). If the thread is empty, i.e. no latest message is available, it
// must be set as null.
func (srv *Server) canUserAddMessage(userRoles *am.GetSelfRolesResult, latestMessageInThread *mm.Message) (ok bool) {
	// Are you stupid, kido ?
	if userRoles == nil {
		return false
	}

	// Only writers can add new messages.
	if !userRoles.User.Roles.IsWriter {
		return false
	}

	// If there is no latest message in the thread, the thread is empty.
	if latestMessageInThread == nil {
		return true
	}

	// If the latest message in the thread was written by another [that] user,
	// this user can add a new message.
	if latestMessageInThread.Creator.UserId != userRoles.User.Id {
		return true
	}

	// If the user is trying to add its another message into a thread, we must
	// check for collision. If the latest message can be edited, it should be
	// edited, and another new message is not allowed.
	if time.Now().Before(srv.getMessageMaxEditTime(latestMessageInThread)) {
		return false
	}

	return true
}

// getLatestMessageOfThreadH is a helper function used by other functions to
// get the latest message of a thread.
func (srv *Server) getLatestMessageOfThreadH(threadId cmb.Id) (message *mm.Message, re *jrm1.RpcError) {
	// Get the thread.
	var thread *mm.Thread
	var err error
	thread, err = srv.dbo.GetThreadById(threadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if thread == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	if thread.Messages == nil {
		return nil, nil
	}

	latestMessageId := thread.Messages.LastElement()
	if latestMessageId == nil {
		return nil, nil
	}

	// Read the message.
	message, err = srv.dbo.GetMessageById(*latestMessageId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if message == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MessageIsNotFound, RpcErrorMsg_MessageIsNotFound, nil)
	}

	return message, nil
}

// getMessageMaxEditTime returns the time border after which message editing is
// not allowed for an ordinary (non-moderator) user.
func (srv *Server) getMessageMaxEditTime(message *mm.Message) time.Time {
	return message.GetLastTouchTime().Add(time.Second * time.Duration(srv.settings.SystemSettings.MessageEditTime))
}

// deleteThreadH is a helper function used by other functions to delete a
// thread.
func (srv *Server) deleteThreadH(threadId cmb.Id) (re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Read the thread.
	var thread *mm.Thread
	var err error
	thread, err = srv.dbo.GetThreadById(threadId)
	if err != nil {
		return srv.databaseError(err)
	}

	if thread == nil {
		return jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	// Check for children.
	if thread.Messages.Size() > 0 {
		return jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotEmpty, RpcErrorMsg_ThreadIsNotEmpty, nil)
	}

	// Update the link.
	var linkThreads *ul.UidList
	linkThreads, err = srv.dbo.GetForumThreadsById(thread.ForumId)
	if err != nil {
		return srv.databaseError(err)
	}

	err = linkThreads.RemoveItem(threadId)
	if err != nil {
		srv.logError(err)
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetForumThreadsById(thread.ForumId, linkThreads)
	if err != nil {
		return srv.databaseError(err)
	}

	// Delete the thread.
	err = srv.dbo.DeleteThreadById(threadId)
	if err != nil {
		return srv.databaseError(err)
	}

	return nil
}
