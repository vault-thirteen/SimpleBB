package models

import (
	"database/sql"
	"net"
	"time"

	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

type EmailChange struct {
	Id             uint
	UserId         uint
	TimeOfCreation time.Time
	RequestId      string

	// IP address of a user. B = Byte array.
	UserIPAB net.IP

	AuthDataBytes     []byte
	IsCaptchaRequired bool
	CaptchaId         sql.NullString
	EmailChangeVerificationFlags

	// Old e-mail.
	VerificationCodeOld string
	IsOldEmailSent      bool

	// New e-mail.
	NewEmail            string
	VerificationCodeNew string
	IsNewEmailSent      bool
}

type EmailChangeVerificationFlags struct {
	IsVerifiedByCaptcha  sql.NullBool
	IsVerifiedByPassword bool
	IsVerifiedByOldEmail bool
	IsVerifiedByNewEmail bool
}

func NewEmailChange() (ec *EmailChange) {
	return &EmailChange{}
}

func NewEmailChangeFromScannableSource(src cm.IScannable) (ec *EmailChange, err error) {
	ec = NewEmailChange()

	err = src.Scan(
		&ec.Id,
		&ec.UserId,
		&ec.TimeOfCreation,
		&ec.RequestId,
		&ec.UserIPAB,
		&ec.AuthDataBytes,
		&ec.IsCaptchaRequired,
		&ec.CaptchaId,
		&ec.IsVerifiedByCaptcha,
		&ec.IsVerifiedByPassword,
		&ec.VerificationCodeOld,
		&ec.IsOldEmailSent,
		&ec.IsVerifiedByOldEmail,
		&ec.NewEmail,
		&ec.VerificationCodeNew,
		&ec.IsNewEmailSent,
		&ec.IsVerifiedByNewEmail,
	)
	if err != nil {
		return nil, err
	}

	return ec, nil
}
