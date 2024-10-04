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
	FuncGetNotification          = "GetNotification"
	FuncGetAllNotifications      = "GetAllNotifications"
	FuncGetUnreadNotifications   = "GetUnreadNotifications"
	FuncCountUnreadNotifications = "CountUnreadNotifications"
	FuncMarkNotificationAsRead   = "MarkNotificationAsRead"
	FuncDeleteNotification       = "DeleteNotification"

	// Other.
	FuncShowDiagnosticData = cc.FuncShowDiagnosticData
)
