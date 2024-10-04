package models

import (
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

// Ping.

type PingParams = cmr.PingParams
type PingResult = cmr.PingResult

// Notification.

type AddNotificationParams struct {
	cmr.CommonParams

	UserId uint   `json:"userId"`
	Text   string `json:"text"`
}

type AddNotificationResult struct {
	cmr.CommonResult

	// ID of the created notification.
	NotificationId uint `json:"notificationId"`
}

type GetNotificationParams struct {
	cmr.CommonParams

	NotificationId uint `json:"notificationId"`
}

type GetNotificationResult struct {
	cmr.CommonResult

	Notification *Notification `json:"notification"`
}

type GetAllNotificationsParams struct {
	cmr.CommonParams
}

type GetAllNotificationsResult struct {
	cmr.CommonResult

	Notifications []Notification `json:"notifications"`
}

type GetUnreadNotificationsParams struct {
	cmr.CommonParams
}

type GetUnreadNotificationsResult struct {
	cmr.CommonResult

	Notifications []Notification `json:"notifications"`
}

type CountUnreadNotificationsParams struct {
	cmr.CommonParams
}

type CountUnreadNotificationsResult struct {
	cmr.CommonResult

	UNC int `json:"unc"`
}

type MarkNotificationAsReadParams struct {
	cmr.CommonParams

	// Identifier of a notification.
	NotificationId uint `json:"notificationId"`
}

type MarkNotificationAsReadResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

type DeleteNotificationParams struct {
	cmr.CommonParams

	NotificationId uint `json:"notificationId"`
}

type DeleteNotificationResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

// Other.

type ShowDiagnosticDataParams struct{}

type ShowDiagnosticDataResult struct {
	cmr.CommonResult
	cmr.RequestsCount
}

type TestParams struct {
	N uint `json:"n"`
}

type TestResult struct {
	cmr.CommonResult
}
