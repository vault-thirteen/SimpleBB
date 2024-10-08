package models

import (
	"database/sql"
	"errors"
	"time"

	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

type Message struct {
	// Identifier of this message.
	Id uint `json:"id"`

	// Identifier of a thread containing this message.
	ThreadId uint `json:"threadId"`

	// Textual information of this message.
	Text string `json:"text"`

	// Check sum of the Text field.
	TextChecksum uint32 `json:"textChecksum"`

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

func NewMessageFromScannableSource(src cm.IScannable) (msg *Message, err error) {
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

func (m *Message) GetLastTouchTime() time.Time {
	if m.Editor.Time == nil {
		return m.Creator.Time
	}

	return *m.Editor.Time
}
