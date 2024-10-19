package models

import (
	ul "github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

// Ping.

type PingParams = cmr.PingParams
type PingResult = cmr.PingResult

// Notification.

type AddNotificationParams struct {
	cmr.CommonParams
	UserId cmb.Id   `json:"userId"`
	Text   cmb.Text `json:"text"`
}
type AddNotificationResult struct {
	cmr.CommonResult

	// ID of the created notification.
	NotificationId cmb.Id `json:"notificationId"`
}

type AddNotificationSParams struct {
	cmr.CommonParams
	cmr.DKeyParams
	UserId cmb.Id   `json:"userId"`
	Text   cmb.Text `json:"text"`
}
type AddNotificationSResult struct {
	cmr.CommonResult

	// ID of the created notification.
	NotificationId cmb.Id `json:"notificationId"`
}

type GetNotificationParams struct {
	cmr.CommonParams
	NotificationId cmb.Id `json:"notificationId"`
}
type GetNotificationResult struct {
	cmr.CommonResult
	Notification *Notification `json:"notification"`
}

type GetNotificationsParams struct {
	cmr.CommonParams
}
type GetNotificationsResult struct {
	cmr.CommonResult
	NotificationIds *ul.UidList    `json:"notificationIds"`
	Notifications   []Notification `json:"notifications"`
}

type GetNotificationsOnPageParams struct {
	cmr.CommonParams
	Page cmb.Count `json:"page"`
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
	NotificationIds *ul.UidList    `json:"notificationIds"`
	Notifications   []Notification `json:"notifications"`
}

type CountUnreadNotificationsParams struct {
	cmr.CommonParams
}
type CountUnreadNotificationsResult struct {
	cmr.CommonResult
	UNC cmb.Count `json:"unc"`
}

type MarkNotificationAsReadParams struct {
	cmr.CommonParams

	// Identifier of a notification.
	NotificationId cmb.Id `json:"notificationId"`
}
type MarkNotificationAsReadResult = cmr.CommonResultWithSuccess

type DeleteNotificationParams struct {
	cmr.CommonParams
	NotificationId cmb.Id `json:"notificationId"`
}
type DeleteNotificationResult = cmr.CommonResultWithSuccess

// Other.

type GetDKeyParams struct {
	cmr.CommonParams
}
type GetDKeyResult struct {
	cmr.CommonResult
	DKey cmb.Text `json:"dKey"`
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
