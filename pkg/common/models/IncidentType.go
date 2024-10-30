package models

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type IncidentType = cmb.Enum

const (
	IncidentType_IllegalAccessAttempt            = cmb.EnumValue(1)
	IncidentType_FakeToken                       = cmb.EnumValue(2)
	IncidentType_VerificationCodeMismatch        = cmb.EnumValue(3)
	IncidentType_DoubleLogInAttempt              = cmb.EnumValue(4)
	IncidentType_PreSessionHacking               = cmb.EnumValue(5)
	IncidentType_CaptchaAnswerMismatch           = cmb.EnumValue(6)
	IncidentType_PasswordMismatch                = cmb.EnumValue(7)
	IncidentType_PasswordChangeHacking           = cmb.EnumValue(8)
	IncidentType_EmailChangeHacking              = cmb.EnumValue(9)
	IncidentType_FakeIPA                         = cmb.EnumValue(10)
	IncidentType_ReadingNotificationOfOtherUsers = cmb.EnumValue(11)
	IncidentType_WrongDKey                       = cmb.EnumValue(12)

	IncidentTypeMax = IncidentType_WrongDKey
)

func NewIncidentType() *IncidentType {
	return cmb.NewEnumFast(IncidentTypeMax)
}

func NewIncidentTypeWithValue(value cmb.EnumValue) IncidentType {
	it := NewIncidentType()
	it.SetValueFast(value)
	return *it
}
