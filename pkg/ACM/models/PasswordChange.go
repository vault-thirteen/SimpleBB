package models

import (
	"database/sql"
	"net"
	"time"

	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

type PasswordChange struct {
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
	VerificationCode     string
	IsEmailSent          bool
	IsVerifiedByEmail    bool
	NewPassword          []byte
}

type PasswordChangeVerificationFlags struct {
	IsVerifiedByCaptcha  sql.NullBool
	IsVerifiedByPassword bool
	IsVerifiedByEmail    bool
}

func NewPasswordChange() (pc *PasswordChange) {
	return &PasswordChange{}
}

func NewPasswordChangeFromScannableSource(src cm.IScannable) (pc *PasswordChange, err error) {
	pc = NewPasswordChange()

	err = src.Scan(
		&pc.Id,
		&pc.UserId,
		&pc.TimeOfCreation,
		&pc.RequestId,
		&pc.UserIPAB,
		&pc.AuthDataBytes,
		&pc.IsCaptchaRequired,
		&pc.CaptchaId,
		&pc.IsVerifiedByCaptcha,
		&pc.IsVerifiedByPassword,
		&pc.VerificationCode,
		&pc.IsEmailSent,
		&pc.IsVerifiedByEmail,
		&pc.NewPassword,
	)
	if err != nil {
		return nil, err
	}

	return pc, nil
}
