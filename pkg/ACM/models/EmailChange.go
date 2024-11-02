package models

import (
	"database/sql"
	"errors"
	"net"
	"time"

	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/models/interfaces"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

type EmailChange struct {
	Id             cmb.Id
	UserId         cmb.Id
	TimeOfCreation time.Time
	RequestId      *cm.RequestId

	// IP address of a user. B = Byte array.
	UserIPAB net.IP

	AuthDataBytes     cmr.AuthChallengeData
	IsCaptchaRequired cmb.Flag
	CaptchaId         *cm.CaptchaId
	EmailChangeVerificationFlags

	// Old e-mail.
	VerificationCodeOld *cm.VerificationCode
	IsOldEmailSent      cmb.Flag

	// New e-mail.
	NewEmail            cm.Email
	VerificationCodeNew *cm.VerificationCode
	IsNewEmailSent      cmb.Flag
}

type EmailChangeVerificationFlags struct {
	IsVerifiedByCaptcha  *cmb.Flag
	IsVerifiedByPassword cmb.Flag
	IsVerifiedByOldEmail cmb.Flag
	IsVerifiedByNewEmail cmb.Flag
}

func NewEmailChange() (ec *EmailChange) {
	return &EmailChange{}
}

func NewEmailChangeFromScannableSource(src cmi.IScannable) (ec *EmailChange, err error) {
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return ec, nil
}
