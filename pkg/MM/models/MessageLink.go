package models

import (
	"database/sql"
	"errors"

	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

// MessageLink is a short variant of a message which stores only IDs.
type MessageLink struct {
	// Identifier of this message.
	Id cmb.Id `json:"id"`

	// Identifier of a thread containing this message.
	ThreadId cmb.Id `json:"threadId"`
}

func NewMessageLink() (ml *MessageLink) {
	return &MessageLink{}
}

func NewMessageLinkFromScannableSource(src cm.IScannable) (ml *MessageLink, err error) {
	ml = NewMessageLink()

	err = src.Scan(
		&ml.Id,
		&ml.ThreadId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return ml, nil
}
