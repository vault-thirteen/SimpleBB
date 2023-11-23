package models

import (
	"net"
	"time"
)

const (
	IncidentType_IllegalAccessAttempt     = 1
	IncidentType_FakeToken                = 2
	IncidentType_VerificationCodeMismatch = 3
	IncidentType_DoubleLogInAttempt       = 4
	IncidentType_PreSessionHacking        = 5
	IncidentType_CaptchaAnswerMismatch    = 6
	IncidentType_PasswordMismatch         = 7
)

type Incident struct {
	Id      uint
	Time    time.Time
	Type    IncidentType
	Email   string
	UserIPA net.IP
}

type IncidentType byte
