package nm

import (
	"fmt"
	"sync"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	nm "github.com/vault-thirteen/SimpleBB/pkg/NM/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

// RPC functions.

// Notification.

// addNotification creates a new notification.
// This method is used to send notifications by administrators.
func (srv *Server) addNotification(p *nm.AddNotificationParams) (result *nm.AddNotificationResult, re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// Check parameters.
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}
	if len(p.Text) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_TextIsNotSet, RpcErrorMsg_TextIsNotSet, nil)
	}

	insertedNotificationId, err := srv.dbo.InsertNewNotification(p.UserId, p.Text)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &nm.AddNotificationResult{
		NotificationId: uint(insertedNotificationId),
	}

	return result, nil
}

// addNotificationS creates a new notification.
// This method is used to send notifications by the system.
func (srv *Server) addNotificationS(p *nm.AddNotificationSParams) (result *nm.AddNotificationSResult, re *jrm1.RpcError) {
	re = srv.mustBeNoAuth(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check the DKey.
	if !srv.dKeyI.CheckString(p.DKey) {
		srv.incidentManager.ReportIncident(cm.IncidentType_WrongDKey, "", nil)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Check parameters.
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}
	if len(p.Text) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_TextIsNotSet, RpcErrorMsg_TextIsNotSet, nil)
	}

	insertedNotificationId, err := srv.dbo.InsertNewNotification(p.UserId, p.Text)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &nm.AddNotificationSResult{
		NotificationId: uint(insertedNotificationId),
	}

	return result, nil
}

// getNotification reads a notification.
func (srv *Server) getNotification(p *nm.GetNotificationParams) (result *nm.GetNotificationResult, re *jrm1.RpcError) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.CanLogIn {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// Check parameters.
	if p.NotificationId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_NotificationIdIsNotSet, RpcErrorMsg_NotificationIdIsNotSet, nil)
	}

	// Read the notification.
	notification, err := srv.dbo.GetNotificationById(p.NotificationId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if notification == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_NotificationIsNotFound, RpcErrorMsg_NotificationIsNotFound, nil)
	}

	// Check the recipient.
	if notification.UserId != userRoles.UserId {
		srv.incidentManager.ReportIncident(cm.IncidentType_ReadingNotificationOfOtherUsers, userRoles.Email, p.Auth.UserIPAB)

		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// All clear.
	result = &nm.GetNotificationResult{
		Notification: notification,
	}

	return result, nil
}

// getAllNotifications reads all notifications for a user.
func (srv *Server) getAllNotifications(p *nm.GetAllNotificationsParams) (result *nm.GetAllNotificationsResult, re *jrm1.RpcError) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.CanLogIn {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// Get notifications.
	notifications, err := srv.dbo.GetAllNotificationsByUserId(userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &nm.GetAllNotificationsResult{
		Notifications: notifications,
	}

	return result, nil
}

// getUnreadNotifications reads all unread notifications for a user.
func (srv *Server) getUnreadNotifications(p *nm.GetUnreadNotificationsParams) (result *nm.GetUnreadNotificationsResult, re *jrm1.RpcError) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.CanLogIn {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// Get notifications.
	notifications, err := srv.dbo.GetUnreadNotifications(userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &nm.GetUnreadNotificationsResult{
		Notifications: notifications,
	}

	return result, nil
}

// countUnreadNotifications counts unread notifications for a user.
func (srv *Server) countUnreadNotifications(p *nm.CountUnreadNotificationsParams) (result *nm.CountUnreadNotificationsResult, re *jrm1.RpcError) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.CanLogIn {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// Count unread notifications.
	n, err := srv.dbo.CountUnreadNotificationsByUserId(userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &nm.CountUnreadNotificationsResult{
		UNC: n,
	}

	return result, nil
}

// markNotificationAsRead marks a notification as read by its recipient.
func (srv *Server) markNotificationAsRead(p *nm.MarkNotificationAsReadParams) (result *nm.MarkNotificationAsReadResult, re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.CanLogIn {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// Check parameters.
	if p.NotificationId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_NotificationIdIsNotSet, RpcErrorMsg_NotificationIdIsNotSet, nil)
	}

	// Get the notification to see its real recipient.
	notification, err := srv.dbo.GetNotificationById(p.NotificationId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if notification == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_NotificationIsNotFound, RpcErrorMsg_NotificationIsNotFound, nil)
	}

	// Check the recipient and status.
	if notification.UserId != userRoles.UserId {
		srv.incidentManager.ReportIncident(cm.IncidentType_ReadingNotificationOfOtherUsers, userRoles.Email, p.Auth.UserIPAB)

		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	if notification.IsRead {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_NotificationIsAlreadyRead, RpcErrorMsg_NotificationIsAlreadyRead, nil)
	}

	// Make the mark.
	err = srv.dbo.MarkNotificationAsRead(p.NotificationId, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &nm.MarkNotificationAsReadResult{
		OK: true,
	}

	return result, nil
}

// deleteNotification removes a notification.
func (srv *Server) deleteNotification(p *nm.DeleteNotificationParams) (result *nm.DeleteNotificationResult, re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// Check parameters.
	if p.NotificationId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_NotificationIdIsNotSet, RpcErrorMsg_NotificationIdIsNotSet, nil)
	}

	// Delete the notification.
	err := srv.dbo.DeleteNotificationById(p.NotificationId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &nm.DeleteNotificationResult{
		OK: true,
	}

	return result, nil
}

// Other.

func (srv *Server) getDKey(p *nm.GetDKeyParams) (result *nm.GetDKeyResult, re *jrm1.RpcError) {
	re = srv.mustBeNoAuth(p.Auth)
	if re != nil {
		return nil, re
	}

	result = &nm.GetDKeyResult{
		DKey: srv.dKeyI.GetString(),
	}

	return result, nil
}

func (srv *Server) showDiagnosticData() (result *nm.ShowDiagnosticDataResult, re *jrm1.RpcError) {
	result = &nm.ShowDiagnosticDataResult{}
	result.TotalRequestsCount, result.SuccessfulRequestsCount = srv.js.GetRequestsCount()

	return result, nil
}

func (srv *Server) test(p *nm.TestParams) (result *nm.TestResult, re *jrm1.RpcError) {
	result = &nm.TestResult{}

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
