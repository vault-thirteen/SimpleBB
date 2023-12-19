package models

import (
	"database/sql"
	"errors"
	"net"
	"time"

	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

type Session struct {
	Id        uint
	UserId    uint
	StartTime time.Time

	// IP address of a user. B = Byte array.
	UserIPAB net.IP
}

func NewSession() (s *Session) {
	return &Session{}
}

func NewSessionFromScannableSource(src cm.IScannable) (s *Session, err error) {
	s = NewSession()

	err = src.Scan(
		&s.Id,
		&s.UserId,
		&s.StartTime,
		&s.UserIPAB,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return s, nil
}
