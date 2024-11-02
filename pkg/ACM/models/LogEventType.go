package models

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/enum"
)

type LogEventType = enum.Enum

const (
	LogEventType_LogIn   = enum.EnumValue(1)
	LogEventType_LogOut  = enum.EnumValue(2) // Self logging out.
	LogEventType_LogOutA = enum.EnumValue(3) // Logging out by an administrator.

	LogEventTypeMax = LogEventType_LogOutA
)

func NewLogEventType() *LogEventType {
	return enum.NewEnumFast(LogEventTypeMax)
}

func NewLogEventTypeWithValue(value enum.EnumValue) LogEventType {
	let := NewLogEventType()
	let.SetValueFast(value)
	return *let
}
