package models

import (
	"database/sql"
	"errors"
	"time"

	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/models/interfaces"
)

type UserParameters struct {
	Id                cmb.Id     `json:"id"`
	PreRegTime        *time.Time `json:"preRegTime,omitempty"`
	Email             Email      `json:"email,omitempty"`
	Name              Name       `json:"name,omitempty"`
	ApprovalTime      *time.Time `json:"approvalTime,omitempty"`
	RegTime           *time.Time `json:"regTime,omitempty"`
	Roles             *UserRoles `json:"roles,omitempty"`
	LastBadLogInTime  *time.Time `json:"lastBadLogInTime,omitempty"`
	BanTime           *time.Time `json:"banTime,omitempty"`
	LastBadActionTime *time.Time `json:"lastBadActionTime,omitempty"`
}

func NewUserParameters() (up *UserParameters) {
	return &UserParameters{
		Roles: NewUserRoles(),
	}
}

func NewUserParametersFromScannableSource(src cmi.IScannable) (up *UserParameters, err error) {
	up = NewUserParameters()

	err = src.Scan(
		&up.Id,
		&up.PreRegTime,
		&up.Email,
		&up.Name,
		&up.ApprovalTime,
		&up.RegTime,
		&up.Roles.IsAuthor,
		&up.Roles.IsWriter,
		&up.Roles.IsReader,
		&up.Roles.CanLogIn,
		&up.LastBadLogInTime,
		&up.BanTime,
		&up.LastBadActionTime,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return up, nil
}
