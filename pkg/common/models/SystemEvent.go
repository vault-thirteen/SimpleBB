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
	SystemEventData

	// ID and time of the event are automatically set by database and should
	// not be touched manually.
	Id   cmb.Id    `json:"id"`
	Time time.Time `json:"time"`
}

func NewSystemEvent() (se *SystemEvent) {
	return &SystemEvent{}
}

func NewSystemEventWithData(data SystemEventData) (se *SystemEvent, err error) {
	se = &SystemEvent{
		SystemEventData: data,
	}

	_, err = se.CheckParameters()
	if err != nil {
		return nil, err
	}

	return se, nil
}

func NewSystemEventFromScannableSource(src IScannable) (se *SystemEvent, err error) {
	se = NewSystemEvent()

	err = src.Scan(
		&se.Id,
		&se.Type,
		&se.ThreadId,
		&se.MessageId,
		&se.UserId,
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
