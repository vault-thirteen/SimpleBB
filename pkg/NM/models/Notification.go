package models

import (
	"database/sql"
	"errors"
	"time"

	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

const (
	ErrUnexpectedNull = "unexpected null"
)

type Notification struct {
	// Identifier of this notification.
	Id cmb.Id `json:"id"`

	// Identifier of a recipient user.
	UserId cmb.Id `json:"userId"`

	// Textual information of this notification.
	Text cmb.Text `json:"text"`

	// Time of creation.
	TimeOfCreation time.Time `json:"toc"`

	// Is the notification read by the recipient ?
	IsRead cmb.Flag `json:"isRead"`

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

func NewNotificationArrayFromRows(rows *sql.Rows) (ns []Notification, err error) {
	ns = []Notification{}
	var n *Notification
	for rows.Next() {
		n, err = NewNotificationFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		if n == nil {
			return nil, errors.New(ErrUnexpectedNull)
		}

		ns = append(ns, *n)
	}

	return ns, nil
}

func ListNotificationIds(notifications []Notification) (ids []cmb.Id) {
	ids = make([]cmb.Id, 0, len(notifications))

	for _, n := range notifications {
		ids = append(ids, n.Id)
	}

	return ids
}
