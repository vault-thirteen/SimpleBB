package models

import (
	"errors"
	"net"
	"time"
)

const (
	ErrLogEventIsNotSet  = "log event is not set"
	ErrLogEventTypeError = "log event type error"
)

type LogEvent struct {
	Id       uint
	Time     time.Time
	Type     LogEventType
	UserId   uint
	Email    string
	UserIPAB net.IP

	// ID of administrator for those events which were started by an
	// administrator.
	AdminId *uint
}

func CheckLogEvent(le *LogEvent) (err error) {
	if le == nil {
		return errors.New(ErrLogEventIsNotSet)
	}

	if !le.Type.IsValid() {
		return errors.New(ErrLogEventTypeError)
	}

	return nil
}
