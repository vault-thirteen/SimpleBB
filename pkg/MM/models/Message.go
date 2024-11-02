package models

import (
	"database/sql"
	"errors"
	"time"

	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/models/interfaces"
)

type Message struct {
	// Identifier of this message.
	Id cmb.Id `json:"id"`

	// Identifier of a thread containing this message.
	ThreadId cmb.Id `json:"threadId"`

	// Textual information of this message.
	Text cmb.Text `json:"text"`

	// Check sum of the Text field.
	TextChecksum []byte `json:"textChecksum"`

	// Message meta-data.
	EventData
}

func NewMessage() (msg *Message) {
	return &Message{
		EventData: EventData{
			Creator: &EventParameters{},
			Editor:  &OptionalEventParameters{},
		},
	}
}

func NewMessageFromScannableSource(src cmi.IScannable) (msg *Message, err error) {
	msg = NewMessage()

	err = src.Scan(
		&msg.Id,
		&msg.ThreadId,
		&msg.Text,
		&msg.TextChecksum,
		&msg.Creator.UserId,
		&msg.Creator.Time,
		&msg.Editor.UserId,
		&msg.Editor.Time,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return msg, nil
}

func NewMessageArrayFromRows(rows cmi.IScannableSequence) (msgs []Message, err error) {
	msgs = []Message{}
	var msg *Message

	for rows.Next() {
		msg, err = NewMessageFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		msgs = append(msgs, *msg)
	}

	return msgs, nil
}

func (m *Message) GetLastTouchTime() time.Time {
	if m.Editor.Time == nil {
		return m.Creator.Time
	}

	return *m.Editor.Time
}
