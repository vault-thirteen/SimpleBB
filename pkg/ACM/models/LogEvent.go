package models

import (
	"net"
	"time"

	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
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

func NewLogEvent() *LogEvent {
	return &LogEvent{
		Type: *NewLogEventType(),
	}
}
