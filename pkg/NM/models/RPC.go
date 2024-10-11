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

type AddNotificationSParams struct {
	cmr.CommonParams
	cmr.DKeyParams

	UserId uint   `json:"userId"`
	Text   string `json:"text"`
}

type AddNotificationSResult struct {
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

type GetNotificationsOnPageParams struct {
	cmr.CommonParams

	Page uint `json:"page"`
}

type GetNotificationsOnPageResult struct {
	cmr.CommonResult

	NotificationsOnPage *NotificationsOnPage `json:"nop"`
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

type GetDKeyParams struct {
	cmr.CommonParams
}

type GetDKeyResult struct {
	cmr.CommonResult

	DKey string `json:"dKey"`
}

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
