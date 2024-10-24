package models

import (
	"database/sql"
	"errors"
	"time"

	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

const (
	ErrUnexpectedNull = "unexpected null"
)

type SystemEvent struct {
	Id        cmb.Id          `json:"id"`
	Type      SystemEventType `json:"type"`
	ThreadId  cmb.Id          `json:"threadId"`
	MessageId cmb.Id          `json:"messageId"`
	Time      time.Time       `json:"time"`

	// Auxiliary fields.

	// ID of a user who initially created the message.
	MessageCreator *cmb.Id `json:"messageCreator,omitempty"`
}

func NewSystemEvent() (se *SystemEvent) {
	return &SystemEvent{}
}

func NewSystemEventFromScannableSource(src IScannable) (se *SystemEvent, err error) {
	se = NewSystemEvent()

	err = src.Scan(
		&se.Id,
		&se.Type,
		&se.ThreadId,
		&se.MessageId,
		&se.Time,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return se, nil
}

func NewSystemEventArrayFromRows(rows *sql.Rows) (ses []SystemEvent, err error) {
	ses = []SystemEvent{}
	var se *SystemEvent
	for rows.Next() {
		se, err = NewSystemEventFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		if se == nil {
			return nil, errors.New(ErrUnexpectedNull)
		}

		ses = append(ses, *se)
	}

	return ses, nil
}

func (se *SystemEvent) CheckType() (ok bool) {
	return se.Type.IsValid()
}

func (se *SystemEvent) isThreadIdSet() (ok bool) {
	return se.ThreadId > 0
}

func (se *SystemEvent) isMessageIdSet() (ok bool) {
	return se.MessageId > 0
}

func (se *SystemEvent) isMessageCreatorSet() (ok bool) {
	return se.MessageCreator != nil
}

func (se *SystemEvent) CheckParameters() (ok bool) {
	// Different event types require different sets of parameters.
	var isThreadIdRequired bool
	var isMessageIdRequired bool
	var isMessageCreatorRequired bool

	switch se.Type {
	case SystemEventType_ThreadParentChange:
		isThreadIdRequired = true

	case SystemEventType_ThreadNameChange:
		isThreadIdRequired = true

	case SystemEventType_ThreadDeletion:
		isThreadIdRequired = true

	case SystemEventType_ThreadNewMessage:
		isThreadIdRequired = true
		isMessageIdRequired = true

	case SystemEventType_ThreadMessageEdit:
		isThreadIdRequired = true
		isMessageIdRequired = true

	case SystemEventType_ThreadMessageDeletion:
		isThreadIdRequired = true
		isMessageIdRequired = true

	case SystemEventType_MessageTextEdit:
		isThreadIdRequired = true
		isMessageIdRequired = true
		isMessageCreatorRequired = true

	case SystemEventType_MessageParentChange:
		isThreadIdRequired = true
		isMessageIdRequired = true
		isMessageCreatorRequired = true

	case SystemEventType_MessageDeletion:
		isThreadIdRequired = true
		isMessageIdRequired = true
		isMessageCreatorRequired = true

	default:
		return false
	}

	if isThreadIdRequired {
		if !se.isThreadIdSet() {
			return false
		}
	}
	if isMessageIdRequired {
		if !se.isMessageIdSet() {
			return false
		}
	}
	if isMessageCreatorRequired {
		if !se.isMessageCreatorSet() {
			return false
		}
	}
	return true
}
