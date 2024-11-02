package models

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/enum"
)

type SystemEventType = enum.Enum

const (
	SystemEventType_ThreadParentChange    = enum.EnumValue(1) // -> Users subscribed to the thread.
	SystemEventType_ThreadNameChange      = enum.EnumValue(2) // -> Users subscribed to the thread.
	SystemEventType_ThreadDeletion        = enum.EnumValue(3) // -> Users subscribed to the thread.
	SystemEventType_ThreadNewMessage      = enum.EnumValue(4) // -> Users subscribed to the thread.
	SystemEventType_ThreadMessageEdit     = enum.EnumValue(5) // -> Users subscribed to the thread.
	SystemEventType_ThreadMessageDeletion = enum.EnumValue(6) // -> Users subscribed to the thread.
	SystemEventType_MessageTextEdit       = enum.EnumValue(7) // -> Author of the message.
	SystemEventType_MessageParentChange   = enum.EnumValue(8) // -> Author of the message.
	SystemEventType_MessageDeletion       = enum.EnumValue(9) // -> Author of the message.

	SystemEventTypeMax = SystemEventType_MessageDeletion
)

func NewSystemEventType() *SystemEventType {
	return enum.NewEnumFast(SystemEventTypeMax)
}

func NewSystemEventTypeWithValue(value enum.EnumValue) SystemEventType {
	let := NewSystemEventType()
	let.SetValueFast(value)
	return *let
}
