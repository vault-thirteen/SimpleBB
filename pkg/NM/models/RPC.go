package models

import (
	ul "github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
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

type SendNotificationIfPossibleSParams struct {
	cmr.CommonParams
	cmr.DKeyParams
	UserId cmb.Id   `json:"userId"`
	Text   cmb.Text `json:"text"`
}
type SendNotificationIfPossibleSResult struct {
	cmr.CommonResult

	// ID and status of the created notification when it is available.
	IsSent         cmb.Flag `json:"isSent"`
	NotificationId cmb.Id   `json:"notificationId"`
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

// Resource.

type AddResourceParams struct {
	cmr.CommonParams
	Resource any `json:"resource"`
}
type AddResourceResult struct {
	cmr.CommonResult

	// ID of the created resource.
	ResourceId cmb.Id `json:"resourceId"`
}

type GetResourceParams struct {
	cmr.CommonParams
	ResourceId cmb.Id `json:"resourceId"`
}
type GetResourceResult struct {
	cmr.CommonResult
	Resource *cm.Resource `json:"resource"`
}

type GetResourceValueParams struct {
	cmr.CommonParams
	ResourceId cmb.Id `json:"resourceId"`
}
type GetResourceValueResult struct {
	cmr.CommonResult
	Resource ResourceWithValue `json:"resource"`
}

type GetListOfAllResourcesOnPageParams struct {
	cmr.CommonParams
	Page cmb.Count `json:"page"`
}
type GetListOfAllResourcesOnPageResult struct {
	cmr.CommonResult
	ResourcesOnPage *ResourcesOnPage `json:"rop"`
}

type DeleteResourceParams struct {
	cmr.CommonParams
	ResourceId cmb.Id `json:"resourceId"`
}
type DeleteResourceResult = cmr.CommonResultWithSuccess

type AddFormatStringParams struct {
	cmr.CommonParams
	FormatString cmb.Text `json:"formatString"`
	FSType       cmb.Text `json:"fsType"`
}
type AddFormatStringResult struct {
	cmr.CommonResult

	// ID of the created resource.
	ResourceId cmb.Id `json:"resourceId"`
}

type GetFormatStringParams struct {
	cmr.CommonParams
	FormatStringId cmb.Id `json:"fsId"`
}
type GetFormatStringResult struct {
	cmr.CommonResult
	FormatStringResource *cm.FormatStringResource `json:"fsr"`
}

type DeleteFormatStringParams struct {
	cmr.CommonParams
	FormatStringId cmb.Id `json:"fsId"`
}
type DeleteFormatStringResult = cmr.CommonResultWithSuccess

// Other.

type ProcessSystemEventSParams struct {
	cmr.CommonParams
	cmr.DKeyParams
	SystemEventData cm.SystemEventData `json:"systemEventData"`
}
type ProcessSystemEventSResult = cmr.CommonResultWithSuccess

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
