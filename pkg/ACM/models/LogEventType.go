package models

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type LogEventType = cmb.Enum

const (
	LogEventType_LogIn   = cmb.EnumValue(1)
	LogEventType_LogOut  = cmb.EnumValue(2) // Self logging out.
	LogEventType_LogOutA = cmb.EnumValue(3) // Logging out by an administrator.

	LogEventTypeMax = LogEventType_LogOutA
)

func NewLogEventType() *LogEventType {
	return cmb.NewEnumFast(LogEventTypeMax)
}

func NewLogEventTypeWithValue(value cmb.EnumValue) LogEventType {
	let := NewLogEventType()
	let.SetValueFast(value)
	return *let
}
