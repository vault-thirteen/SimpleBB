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

type PasswordChange struct {
	Id             cmb.Id
	UserId         cmb.Id
	TimeOfCreation time.Time
	RequestId      *cm.RequestId

	// IP address of a user. B = Byte array.
	UserIPAB net.IP

	AuthDataBytes        cmr.AuthChallengeData
	IsCaptchaRequired    cmb.Flag
	CaptchaId            *cm.CaptchaId
	IsVerifiedByCaptcha  *cmb.Flag
	IsVerifiedByPassword cmb.Flag
	VerificationCode     *cm.VerificationCode
	IsEmailSent          cmb.Flag
	IsVerifiedByEmail    cmb.Flag
	NewPasswordBytes     []byte
}

type PasswordChangeVerificationFlags struct {
	IsVerifiedByCaptcha  *cmb.Flag
	IsVerifiedByPassword cmb.Flag
	IsVerifiedByEmail    cmb.Flag
}

func NewPasswordChange() (pc *PasswordChange) {
	return &PasswordChange{}
}

func NewPasswordChangeFromScannableSource(src cmi.IScannable) (pc *PasswordChange, err error) {
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
		&pc.NewPasswordBytes,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return pc, nil
}
