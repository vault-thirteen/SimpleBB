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

type PreSession struct {
	Id             cmb.Id
	UserId         cmb.Id
	TimeOfCreation time.Time
	RequestId      cm.RequestId

	// IP address of a user. B = Byte array.
	UserIPAB net.IP

	AuthDataBytes        cmr.AuthChallengeData
	IsCaptchaRequired    cmb.Flag
	CaptchaId            *cm.CaptchaId
	IsVerifiedByCaptcha  *cmb.Flag
	IsVerifiedByPassword cmb.Flag

	// Verification code is set on Step 2, so it is NULL on Step 1.
	VerificationCode *cm.VerificationCode

	IsEmailSent       cmb.Flag
	IsVerifiedByEmail cmb.Flag
}

func NewPreSession() (ps *PreSession) {
	return &PreSession{}
}

func NewPreSessionFromScannableSource(src cmi.IScannable) (ps *PreSession, err error) {
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
