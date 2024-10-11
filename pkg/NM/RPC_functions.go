package nm

import (
	"fmt"
	"sync"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	nm "github.com/vault-thirteen/SimpleBB/pkg/NM/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	ul "github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

// RPC functions.

// Notification.

// addNotification creates a new notification.
// This method is used to send notifications by administrators.
func (srv *Server) addNotification(p *nm.AddNotificationParams) (result *nm.AddNotificationResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	if len(p.Text) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_TextIsNotSet, RpcErrorMsg_TextIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

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
	// Check parameters.
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	if len(p.Text) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_TextIsNotSet, RpcErrorMsg_TextIsNotSet, nil)
	}

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
	// Check parameters.
	if p.NotificationId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_NotificationIdIsNotSet, RpcErrorMsg_NotificationIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.CanLogIn {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

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

// getNotificationsOnPage reads notifications for a user on the selected page.
func (srv *Server) getNotificationsOnPage(p *nm.GetNotificationsOnPageParams) (result *nm.GetNotificationsOnPageResult, re *jrm1.RpcError) {
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
	if !userRoles.CanLogIn {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Get notifications on page.
	notifications, err := srv.dbo.GetNotificationsByUserIdOnPage(userRoles.UserId, p.Page, srv.settings.SystemSettings.PageSize)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Count all notifications.
	var allNotificationsCount int
	allNotificationsCount, err = srv.dbo.CountAllNotificationsByUserId(userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	nop := nm.NewNotificationsOnPage()
	{
		nop.NotificationIds, err = ul.NewFromArray(nm.ListNotificationIds(notifications))
		if err != nil {
			srv.logError(err)
			return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
		}

		nop.Notifications = notifications

		nop.Page = &p.Page

		allNotificationsCountUint := uint(allNotificationsCount)
		tp := c.CalculateTotalPages(allNotificationsCountUint, srv.settings.SystemSettings.PageSize)
		nop.TotalPages = &tp

		nop.TotalNotifications = &allNotificationsCountUint
	}

	result = &nm.GetNotificationsOnPageResult{
		NotificationsOnPage: nop,
	}

	return result, nil
}

// getAllNotifications reads all notifications for a user.
func (srv *Server) getAllNotifications(p *nm.GetAllNotificationsParams) (result *nm.GetAllNotificationsResult, re *jrm1.RpcError) {
	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.CanLogIn {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

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
	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.CanLogIn {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

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
	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.CanLogIn {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

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
	// Check parameters.
	if p.NotificationId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_NotificationIdIsNotSet, RpcErrorMsg_NotificationIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.CanLogIn {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

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
	// Check parameters.
	if p.NotificationId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_NotificationIdIsNotSet, RpcErrorMsg_NotificationIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

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
