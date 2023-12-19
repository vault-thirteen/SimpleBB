package models

import (
	"database/sql"
	"errors"
	"net"
	"time"

	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

type PreSession struct {
	Id             uint
	UserId         uint
	TimeOfCreation time.Time
	RequestId      string

	// IP address of a user. B = Byte array.
	UserIPAB net.IP

	AuthDataBytes        []byte
	IsCaptchaRequired    bool
	CaptchaId            sql.NullString
	IsVerifiedByCaptcha  sql.NullBool
	IsVerifiedByPassword bool

	// Verification code is set on Step 2, so it is NULL on Step 1.
	VerificationCode sql.NullString

	IsEmailSent       bool
	IsVerifiedByEmail bool
}

func NewPreSession() (ps *PreSession) {
	return &PreSession{}
}

func NewPreSessionFromScannableSource(src cm.IScannable) (ps *PreSession, err error) {
	ps = NewPreSession()

	err = src.Scan(
		&ps.Id,
		&ps.UserId,
		&ps.TimeOfCreation,
		&ps.RequestId,
		&ps.UserIPAB,
		&ps.AuthDataBytes,
		&ps.IsCaptchaRequired,
		&ps.CaptchaId,
		&ps.IsVerifiedByCaptcha,
		&ps.IsVerifiedByPassword,
		&ps.VerificationCode,
		&ps.IsEmailSent,
		&ps.IsVerifiedByEmail,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return ps, nil
}
