package models

import (
	"database/sql"
	"errors"

	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/models/interfaces"
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

func NewMessageLinkFromScannableSource(src cmi.IScannable) (ml *MessageLink, err error) {
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

func NewMessageLinkArrayFromRows(rows cmi.IScannableSequence) (mls []MessageLink, err error) {
	mls = []MessageLink{}
	var ml *MessageLink

	for rows.Next() {
		ml, err = NewMessageLinkFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		mls = append(mls, *ml)
	}

	return mls, nil
}
