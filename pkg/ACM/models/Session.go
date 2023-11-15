package models

import (
	"net"
	"time"
)

type Session struct {
	Id        uint
	UserId    uint
	StartTime time.Time

	// IP address of a user. B = Byte array.
	UserIPAB net.IP
}
