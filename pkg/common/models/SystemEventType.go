package models

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type SystemEventType = cmb.Enum

const (
	SystemEventType_ThreadParentChange    = cmb.EnumValue(1) // -> Users subscribed to the thread.
	SystemEventType_ThreadNameChange      = cmb.EnumValue(2) // -> Users subscribed to the thread.
	SystemEventType_ThreadDeletion        = cmb.EnumValue(3) // -> Users subscribed to the thread.
	SystemEventType_ThreadNewMessage      = cmb.EnumValue(4) // -> Users subscribed to the thread.
	SystemEventType_ThreadMessageEdit     = cmb.EnumValue(5) // -> Users subscribed to the thread.
	SystemEventType_ThreadMessageDeletion = cmb.EnumValue(6) // -> Users subscribed to the thread.
	SystemEventType_MessageTextEdit       = cmb.EnumValue(7) // -> Author of the message.
	SystemEventType_MessageParentChange   = cmb.EnumValue(8) // -> Author of the message.
	SystemEventType_MessageDeletion       = cmb.EnumValue(9) // -> Author of the message.

	SystemEventTypeMax = SystemEventType_MessageDeletion
)

func NewSystemEventType() *SystemEventType {
	return cmb.NewEnumFast(SystemEventTypeMax)
}

func NewSystemEventTypeWithValue(value cmb.EnumValue) SystemEventType {
	let := NewSystemEventType()
	let.SetValueFast(value)
	return *let
}
