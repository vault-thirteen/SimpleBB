package models

import (
	"database/sql"
	"net"
	"time"
)

type PreSession struct {
	Id             uint
	UserId         uint
	TimeOfCreation time.Time
	RequestId      string

	// IP address of a user. B = Byte array.
	UserIPAB net.IP

	AuthDataBytes     []byte
	IsCaptchaRequired bool
	CaptchaId         sql.NullString
}
