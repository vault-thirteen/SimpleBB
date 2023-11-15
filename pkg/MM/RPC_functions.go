package mm

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	js "github.com/osamingo/jsonrpc/v2"
	ac "github.com/vault-thirteen/SimpleBB/pkg/ACM/client"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	mm "github.com/vault-thirteen/SimpleBB/pkg/MM/models"
	ul "github.com/vault-thirteen/SimpleBB/pkg/UidList"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
)

// RPC functions.

func (srv *Server) addForum(p *mm.AddForumParams) (result *mm.AddForumResult, jerr *js.Error) {
	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if len(p.Name) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumNameIsNotSet, Message: RpcErrorMsg_ForumNameIsNotSet}
	}

	// If parent is not set, the new forum is a root forum.
	// Only a single root forum may exist.
	var err error
	var n int
	if p.Parent == nil {
		n, err = srv.countRootForumsM()
		if err != nil {
			return nil, srv.databaseError(err)
		}

		if n > 0 {
			return nil, &js.Error{Code: RpcErrorCode_RootForumAlreadyExists, Message: RpcErrorMsg_RootForumAlreadyExists}
		}

		var insertedForumId int64
		insertedForumId, err = srv.insertNewForumM(p.Parent, p.Name, userRoles.UserId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		result = &mm.AddForumResult{
			ForumId: uint(insertedForumId),
		}

		return result, nil
	}

	// Insert a normal forum.
	// Ensure that a parent really exists.
	n, err = srv.countForumsByIdM(*p.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIsNotFound, Message: RpcErrorMsg_ForumIsNotFound}
	}

	var parentChildren *ul.UidList
	parentChildren, err = srv.getForumChildrenByIdM(*p.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var insertedForumId int64
	insertedForumId, err = srv.insertNewForumM(p.Parent, p.Name, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = parentChildren.AddItem(uint(insertedForumId))
	if err != nil {
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.setForumChildrenByIdM(*p.Parent, parentChildren)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.AddForumResult{
		ForumId: uint(insertedForumId),
	}

	return result, nil
}

func (srv *Server) changeForumName(p *mm.ChangeForumNameParams) (result *mm.ChangeForumNameResult, jerr *js.Error) {
	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ForumId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIdIsNotSet, Message: RpcErrorMsg_ForumIdIsNotSet}
	}

	if len(p.Name) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumNameIsNotSet, Message: RpcErrorMsg_ForumNameIsNotSet}
	}

	var err error
	err = srv.setForumNameByIdM(p.ForumId, p.Name, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ChangeForumNameResult{
		OK: true,
	}

	return result, nil
}

func (srv *Server) changeForumParent(p *mm.ChangeForumParentParams) (result *mm.ChangeForumParentResult, jerr *js.Error) {
	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ForumId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIdIsNotSet, Message: RpcErrorMsg_ForumIdIsNotSet}
	}

	if p.Parent == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIdIsNotSet, Message: RpcErrorMsg_ForumIdIsNotSet}
	}

	// Ensure that an old parent exists.
	var oldParent *uint
	var err error
	oldParent, err = srv.getForumParentByIdM(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if oldParent == nil {
		return nil, &js.Error{Code: RpcErrorCode_RootForumCanNotBeMoved, Message: RpcErrorMsg_RootForumCanNotBeMoved}
	}

	var n int
	n, err = srv.countForumsByIdM(*oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIsNotFound, Message: RpcErrorMsg_ForumIsNotFound}
	}

	// Ensure that a new parent exists.
	n, err = srv.countForumsByIdM(p.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIsNotFound, Message: RpcErrorMsg_ForumIsNotFound}
	}

	// Update the moved forum.
	err = srv.setForumParentByIdM(p.ForumId, p.Parent, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the new link.
	var childrenR *ul.UidList
	childrenR, err = srv.getForumChildrenByIdM(p.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = childrenR.AddItem(p.ForumId)
	if err != nil {
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.setForumChildrenByIdM(p.Parent, childrenR)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the old link.
	var childrenL *ul.UidList
	childrenL, err = srv.getForumChildrenByIdM(*oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = childrenL.RemoveItem(p.ForumId)
	if err != nil {
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.setForumChildrenByIdM(*oldParent, childrenL)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ChangeForumParentResult{
		OK: true,
	}

	return result, nil
}

func (srv *Server) addThread(p *mm.AddThreadParams) (result *mm.AddThreadResult, jerr *js.Error) {
	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAuthor {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ForumId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIdIsNotSet, Message: RpcErrorMsg_ForumIdIsNotSet}
	}

	if len(p.Name) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadNameIsNotSet, Message: RpcErrorMsg_ThreadNameIsNotSet}
	}

	// Ensure that a parent really exists.
	var err error
	var n int
	n, err = srv.countForumsByIdM(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIsNotFound, Message: RpcErrorMsg_ForumIsNotFound}
	}

	var parentThreads *ul.UidList
	parentThreads, err = srv.getForumThreadsByIdM(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var insertedThreadId int64
	insertedThreadId, err = srv.insertNewThreadM(p.ForumId, p.Name, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = parentThreads.AddItem(uint(insertedThreadId))
	if err != nil {
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.setForumThreadsByIdM(p.ForumId, parentThreads)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.AddThreadResult{
		ThreadId: uint(insertedThreadId),
	}

	return result, nil
}

func (srv *Server) changeThreadName(p *mm.ChangeThreadNameParams) (result *mm.ChangeThreadNameResult, jerr *js.Error) {
	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ThreadId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIdIsNotSet, Message: RpcErrorMsg_ThreadIdIsNotSet}
	}

	if len(p.Name) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadNameIsNotSet, Message: RpcErrorMsg_ThreadNameIsNotSet}
	}

	var err error
	err = srv.setThreadNameByIdM(p.ThreadId, p.Name, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ChangeThreadNameResult{
		OK: true,
	}

	return result, nil
}

func (srv *Server) changeThreadForum(p *mm.ChangeThreadForumParams) (result *mm.ChangeThreadForumResult, jerr *js.Error) {
	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ThreadId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIdIsNotSet, Message: RpcErrorMsg_ThreadIdIsNotSet}
	}

	if p.ForumId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIdIsNotSet, Message: RpcErrorMsg_ForumIdIsNotSet}
	}

	// Ensure that an old parent exists.
	var oldParent uint
	var err error
	oldParent, err = srv.getThreadForumByIdM(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var n int
	n, err = srv.countForumsByIdM(oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIsNotFound, Message: RpcErrorMsg_ForumIsNotFound}
	}

	// Ensure that a new parent exists.
	n, err = srv.countForumsByIdM(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIsNotFound, Message: RpcErrorMsg_ForumIsNotFound}
	}

	// Update the moved thread.
	err = srv.setThreadForumByIdM(p.ThreadId, p.ForumId, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the new link.
	var threadsR *ul.UidList
	threadsR, err = srv.getForumThreadsByIdM(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = threadsR.AddItem(p.ThreadId)
	if err != nil {
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.setForumThreadsByIdM(p.ForumId, threadsR)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the old link.
	var threadsL *ul.UidList
	threadsL, err = srv.getForumThreadsByIdM(oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = threadsL.RemoveItem(p.ThreadId)
	if err != nil {
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.setForumThreadsByIdM(oldParent, threadsL)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ChangeThreadForumResult{
		OK: true,
	}

	return result, nil
}

func (srv *Server) addMessage(p *mm.AddMessageParams) (result *mm.AddMessageResult, jerr *js.Error) {
	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsWriter {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ThreadId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIdIsNotSet, Message: RpcErrorMsg_ThreadIdIsNotSet}
	}

	if len(p.Text) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_MessageTextIsNotSet, Message: RpcErrorMsg_MessageTextIsNotSet}
	}

	// Ensure that a parent really exists.
	var err error
	var n int
	n, err = srv.countThreadsByIdM(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIsNotFound, Message: RpcErrorMsg_ThreadIsNotFound}
	}

	var parentMessages *ul.UidList
	parentMessages, err = srv.getThreadMessagesByIdM(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var insertedThreadId int64
	insertedThreadId, err = srv.insertNewMessageM(p.ThreadId, p.Text, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = parentMessages.AddItem(uint(insertedThreadId))
	if err != nil {
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.setThreadMessagesByIdM(p.ThreadId, parentMessages)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.AddMessageResult{
		MessageId: uint(insertedThreadId),
	}

	return result, nil
}

func (srv *Server) changeMessageText(p *mm.ChangeMessageTextParams) (result *mm.ChangeMessageTextResult, jerr *js.Error) {
	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	var ok = false
	if userRoles.IsModerator {
		ok = true
	} else {
		// User can edit its own messages if they are not too old.
		if userRoles.IsWriter {
			var creatorUserId uint
			var ToC time.Time
			var ToE *time.Time
			var err error
			creatorUserId, ToC, ToE, err = srv.getMessageCreatorAndTimeByIdM(p.MessageId)
			if err != nil {
				return nil, srv.databaseError(err)
			}

			if userRoles.UserId == creatorUserId {
				var lastTouchTime time.Time
				if ToE != nil {
					lastTouchTime = *ToE
				} else {
					lastTouchTime = ToC
				}

				if time.Now().Before(lastTouchTime.Add(time.Second * time.Duration(srv.settings.SystemSettings.MessageEditTime))) {
					ok = true
				}
			}
		}
	}

	if !ok {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.MessageId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_MessageIdIsNotSet, Message: RpcErrorMsg_MessageIdIsNotSet}
	}

	if len(p.Text) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_MessageTextIsNotSet, Message: RpcErrorMsg_MessageTextIsNotSet}
	}

	var err error
	err = srv.setMessageTextByIdM(p.MessageId, p.Text, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ChangeMessageTextResult{
		OK: true,
	}

	return result, nil
}

func (srv *Server) changeMessageThread(p *mm.ChangeMessageThreadParams) (result *mm.ChangeMessageThreadResult, jerr *js.Error) {
	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.MessageId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_MessageIdIsNotSet, Message: RpcErrorMsg_MessageIdIsNotSet}
	}

	if p.ThreadId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIdIsNotSet, Message: RpcErrorMsg_ThreadIdIsNotSet}
	}

	// Ensure that an old parent exists.
	var oldParent uint
	var err error
	oldParent, err = srv.getMessageThreadByIdM(p.MessageId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var n int
	n, err = srv.countThreadsByIdM(oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIsNotFound, Message: RpcErrorMsg_ThreadIsNotFound}
	}

	// Ensure that a new parent exists.
	n, err = srv.countThreadsByIdM(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIsNotFound, Message: RpcErrorMsg_ThreadIsNotFound}
	}

	// Update the moved message.
	err = srv.setMessageThreadByIdM(p.MessageId, p.ThreadId, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the new link.
	var messagesR *ul.UidList
	messagesR, err = srv.getThreadMessagesByIdM(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = messagesR.AddItem(p.MessageId)
	if err != nil {
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.setThreadMessagesByIdM(p.ThreadId, messagesR)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the old link.
	var messagesL *ul.UidList
	messagesL, err = srv.getThreadMessagesByIdM(oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = messagesL.RemoveItem(p.MessageId)
	if err != nil {
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.setThreadMessagesByIdM(oldParent, messagesL)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ChangeMessageThreadResult{
		OK: true,
	}

	return result, nil
}

func (srv *Server) deleteMessage(p *mm.DeleteMessageParams) (result *mm.DeleteMessageResult, jerr *js.Error) {
	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.MessageId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_MessageIdIsNotSet, Message: RpcErrorMsg_MessageIdIsNotSet}
	}

	// Read the message.
	var message *mm.Message
	var err error
	message, err = srv.getMessageByIdM(p.MessageId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the link.
	var linkMessages *ul.UidList
	linkMessages, err = srv.getThreadMessagesByIdM(message.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = linkMessages.RemoveItem(p.MessageId)
	if err != nil {
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.setThreadMessagesByIdM(message.ThreadId, linkMessages)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Delete the message.
	err = srv.deleteMessageByIdM(p.MessageId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.DeleteMessageResult{
		OK: true,
	}

	return result, nil
}

func (srv *Server) getMessage(p *mm.GetMessageParams) (result *mm.GetMessageResult, jerr *js.Error) {
	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsReader {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.MessageId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_MessageIdIsNotSet, Message: RpcErrorMsg_MessageIdIsNotSet}
	}

	// Read the message.
	var message *mm.Message
	var err error
	message, err = srv.getMessageByIdM(p.MessageId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.GetMessageResult{
		Message: message,
	}

	return result, nil
}

func (srv *Server) getThread(p *mm.GetThreadParams) (result *mm.GetThreadResult, jerr *js.Error) {
	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsReader {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ThreadId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIdIsNotSet, Message: RpcErrorMsg_ThreadIdIsNotSet}
	}

	// Read the thread.
	var thread *mm.Thread
	var err error
	thread, err = srv.getThreadByIdM(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.GetThreadResult{
		Thread: thread,
	}

	return result, nil
}

func (srv *Server) deleteThread(p *mm.DeleteThreadParams) (result *mm.DeleteThreadResult, jerr *js.Error) {
	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ThreadId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIdIsNotSet, Message: RpcErrorMsg_ThreadIdIsNotSet}
	}

	// Read the thread.
	var thread *mm.Thread
	var err error
	thread, err = srv.getThreadByIdM(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Check for children.
	if thread.Messages.Size() > 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIsNotEmpty, Message: RpcErrorMsg_ThreadIsNotEmpty}
	}

	// Update the link.
	var linkThreads *ul.UidList
	linkThreads, err = srv.getForumThreadsByIdM(thread.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = linkThreads.RemoveItem(p.ThreadId)
	if err != nil {
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.setForumThreadsByIdM(thread.ForumId, linkThreads)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Delete the thread.
	err = srv.deleteThreadByIdM(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.DeleteThreadResult{
		OK: true,
	}

	return result, nil
}

func (srv *Server) getForum(p *mm.GetForumParams) (result *mm.GetForumResult, jerr *js.Error) {
	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsReader {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ForumId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIdIsNotSet, Message: RpcErrorMsg_ForumIdIsNotSet}
	}

	// Read the forum.
	var forum *mm.Forum
	var err error
	forum, err = srv.getForumByIdM(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.GetForumResult{
		Forum: forum,
	}

	return result, nil
}

func (srv *Server) deleteForum(p *mm.DeleteForumParams) (result *mm.DeleteForumResult, jerr *js.Error) {
	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ForumId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIdIsNotSet, Message: RpcErrorMsg_ForumIdIsNotSet}
	}

	// Read the forum.
	var forum *mm.Forum
	var err error
	forum, err = srv.getForumByIdM(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var isRootForum = false
	if forum.Parent == nil {
		isRootForum = true
	}

	// Check for children.
	if forum.Children.Size() > 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumHasChildren, Message: RpcErrorMsg_ForumHasChildren}
	}

	if forum.Threads.Size() > 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumHasThreads, Message: RpcErrorMsg_ForumHasThreads}
	}

	// Update the link.
	if !isRootForum {
		var linkForums *ul.UidList
		linkForums, err = srv.getForumChildrenByIdM(*forum.Parent)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		err = linkForums.RemoveItem(p.ForumId)
		if err != nil {
			return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
		}

		err = srv.setForumChildrenByIdM(*forum.Parent, linkForums)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	// Delete the forum.
	err = srv.deleteForumByIdM(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.DeleteForumResult{
		OK: true,
	}

	return result, nil
}

func (srv *Server) listThreadMessages(p *mm.ListThreadMessagesParams) (result *mm.ListThreadMessagesResult, jerr *js.Error) {
	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsReader {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ThreadId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIdIsNotSet, Message: RpcErrorMsg_ThreadIdIsNotSet}
	}

	// Read the thread.
	var thread *mm.Thread
	var err error
	thread, err = srv.getThreadByIdM(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	twm := mm.NewThreadWithMessages(thread)

	// Read all the messages.
	twm.MessageIds = *thread.Messages
	twm.Messages, err = srv.readMessagesByIdM(*thread.Messages)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ListThreadMessagesResult{
		ThreadWithMessages: twm,
	}

	return result, nil
}

func (srv *Server) listThreadMessagesOnPage(p *mm.ListThreadMessagesOnPageParams) (result *mm.ListThreadMessagesOnPageResult, jerr *js.Error) {
	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsReader {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ThreadId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIdIsNotSet, Message: RpcErrorMsg_ThreadIdIsNotSet}
	}

	if p.Page == 0 {
		return nil, &js.Error{Code: RpcErrorCode_PageIsNotSet, Message: RpcErrorMsg_PageIsNotSet}
	}

	// Read the thread.
	var thread *mm.Thread
	var err error
	thread, err = srv.getThreadByIdM(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	twmop := mm.NewThreadWithMessages(thread)

	// Total numbers before pagination.
	tm := uint(thread.Messages.Size())
	twmop.TotalMessages = &tm
	tp := uint(math.Ceil(float64(tm) / float64(srv.settings.SystemSettings.PageSize)))
	twmop.TotalPages = &tp

	// Read messages of a specified page.
	twmop.Page = &p.Page
	twmop.MessageIds = thread.Messages.OnPage(p.Page, srv.settings.SystemSettings.PageSize)
	if len(twmop.MessageIds) > 0 {
		twmop.Messages, err = srv.readMessagesByIdM(twmop.MessageIds)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	result = &mm.ListThreadMessagesOnPageResult{
		ThreadWithMessages: twmop,
	}

	return result, nil
}

func (srv *Server) listForumThreads(p *mm.ListForumThreadsParams) (result *mm.ListForumThreadsResult, jerr *js.Error) {
	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsReader {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ForumId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIdIsNotSet, Message: RpcErrorMsg_ForumIdIsNotSet}
	}

	// Read the forum.
	var forum *mm.Forum
	var err error
	forum, err = srv.getForumByIdM(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	fwt := mm.NewForumWithThreads(forum)

	// Read all the threads.
	fwt.ThreadIds = *forum.Threads
	fwt.Threads, err = srv.readThreadsByIdM(*forum.Threads)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ListForumThreadsResult{
		ForumWithThreads: fwt,
	}

	return result, nil
}

func (srv *Server) listForumThreadsOnPage(p *mm.ListForumThreadsOnPageParams) (result *mm.ListForumThreadsOnPageResult, jerr *js.Error) {
	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsReader {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ForumId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIdIsNotSet, Message: RpcErrorMsg_ForumIdIsNotSet}
	}

	if p.Page == 0 {
		return nil, &js.Error{Code: RpcErrorCode_PageIsNotSet, Message: RpcErrorMsg_PageIsNotSet}
	}

	// Read the forum.
	var forum *mm.Forum
	var err error
	forum, err = srv.getForumByIdM(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	fwtop := mm.NewForumWithThreads(forum)

	// Total numbers before pagination.
	tt := uint(forum.Threads.Size())
	fwtop.TotalThreads = &tt
	tp := uint(math.Ceil(float64(tt) / float64(srv.settings.SystemSettings.PageSize)))
	fwtop.TotalPages = &tp

	// Read threads of a specified page.
	fwtop.Page = &p.Page
	fwtop.ThreadIds = forum.Threads.OnPage(p.Page, srv.settings.SystemSettings.PageSize)
	if len(fwtop.ThreadIds) > 0 {
		fwtop.Threads, err = srv.readThreadsByIdM(fwtop.ThreadIds)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	result = &mm.ListForumThreadsOnPageResult{
		ForumWithThreads: fwtop,
	}

	return result, nil
}

func (srv *Server) listForums(p *mm.ListForumsParams) (result *mm.ListForumsResult, jerr *js.Error) {
	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsReader {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Read all the forums.
	var forums []mm.Forum
	var err error
	forums, err = srv.readForums()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ListForumsResult{
		Forums: forums,
	}

	return result, nil
}

func (srv *Server) showDiagnosticData() (result *mm.ShowDiagnosticDataResult, jerr *js.Error) {
	result = &mm.ShowDiagnosticDataResult{
		TotalRequestsCount:      srv.diag.GetTotalRequestsCount(),
		SuccessfulRequestsCount: srv.diag.GetSuccessfulRequestsCount(),
	}

	return result, nil
}

func (srv *Server) doTest(p *mm.TestParams) (result *mm.TestResult, jerr *js.Error) {
	result = &mm.TestResult{}

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
			return nil, &js.Error{Code: RpcErrorCode_TestError, Message: fmt.Sprintf(RpcErrorMsgF_TestError, err.Error())}
		}
	}

	return result, nil
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
