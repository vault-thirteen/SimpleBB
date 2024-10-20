package c

import (
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/client"
)

// List of supported functions.
const (
	// Ping.
	FuncPing = cc.FuncPing

	// Notification.
	FuncAddNotification          = "AddNotification"
	FuncAddNotificationS         = "AddNotificationS"
	FuncGetNotification          = "GetNotification"
	FuncGetNotifications         = "GetNotifications"
	FuncGetNotificationsOnPage   = "GetNotificationsOnPage"
	FuncGetUnreadNotifications   = "GetUnreadNotifications"
	FuncCountUnreadNotifications = "CountUnreadNotifications"
	FuncMarkNotificationAsRead   = "MarkNotificationAsRead"
	FuncDeleteNotification       = "DeleteNotification"

	// Other.
	FuncGetDKey            = "GetDKey"
	FuncShowDiagnosticData = cc.FuncShowDiagnosticData
	FuncTest               = "Test"
)
