package models

import (
	"errors"
	"net"
	"time"

	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

const (
	ErrLogEventIsNotSet  = "log event is not set"
	ErrLogEventTypeError = "log event type error"
)

type LogEvent struct {
	Id       cmb.Id
	Time     time.Time
	Type     LogEventType
	UserId   cmb.Id
	Email    cm.Email
	UserIPAB net.IP

	// ID of administrator for those events which were started by an
	// administrator.
	AdminId *cmb.Id
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
