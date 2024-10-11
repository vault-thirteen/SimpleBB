package models

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
)

type NotificationsOnPage struct {
	// Notification parameters. If pagination is used, these lists contain
	// information after the application of pagination.
	NotificationIds *ul.UidList    `json:"notificationIds"`
	Notifications   []Notification `json:"notifications"`

	// Number of the current page of notifications.
	Page *uint `json:"page,omitempty"`

	// Total number of available pages of notifications.
	TotalPages *uint `json:"totalPages,omitempty"`

	// Total number of available notifications.
	TotalNotifications *uint `json:"totalNotifications,omitempty"`
}

func NewNotificationsOnPage() (nop *NotificationsOnPage) {
	nop = &NotificationsOnPage{}
	return nop
}
