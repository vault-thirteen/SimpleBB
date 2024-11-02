package models

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/enum"
)

type IncidentType = enum.Enum

const (
	IncidentType_IllegalAccessAttempt            = enum.EnumValue(1)
	IncidentType_FakeToken                       = enum.EnumValue(2)
	IncidentType_VerificationCodeMismatch        = enum.EnumValue(3)
	IncidentType_DoubleLogInAttempt              = enum.EnumValue(4)
	IncidentType_PreSessionHacking               = enum.EnumValue(5)
	IncidentType_CaptchaAnswerMismatch           = enum.EnumValue(6)
	IncidentType_PasswordMismatch                = enum.EnumValue(7)
	IncidentType_PasswordChangeHacking           = enum.EnumValue(8)
	IncidentType_EmailChangeHacking              = enum.EnumValue(9)
	IncidentType_FakeIPA                         = enum.EnumValue(10)
	IncidentType_ReadingNotificationOfOtherUsers = enum.EnumValue(11)
	IncidentType_WrongDKey                       = enum.EnumValue(12)

	IncidentTypeMax = IncidentType_WrongDKey
)

func NewIncidentType() *IncidentType {
	return enum.NewEnumFast(IncidentTypeMax)
}

func NewIncidentTypeWithValue(value enum.EnumValue) IncidentType {
	it := NewIncidentType()
	it.SetValueFast(value)
	return *it
}
