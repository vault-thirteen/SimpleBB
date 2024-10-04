package models

import (
	"database/sql"
	"errors"
	"time"

	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

type Notification struct {
	// Identifier of this notification.
	Id uint `json:"id"`

	// Identifier of a recipient user.
	UserId uint `json:"userId"`

	// Textual information of this notification.
	Text string `json:"text"`

	// Time of creation.
	TimeOfCreation time.Time `json:"toc"`

	// Is the notification read by the recipient ?
	IsRead bool `json:"isRead"`

	// Time of reading.
	TimeOfReading *time.Time `json:"tor"`
}

func NewNotification() (n *Notification) {
	return &Notification{}
}

func NewNotificationFromScannableSource(src cm.IScannable) (n *Notification, err error) {
	n = NewNotification()

	err = src.Scan(
		&n.Id,
		&n.UserId,
		&n.Text,
		&n.TimeOfCreation,
		&n.IsRead,
		&n.TimeOfReading,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return n, nil
}
