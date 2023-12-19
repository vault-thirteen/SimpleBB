package models

import (
	"database/sql"
	"errors"
	"time"
)

type UserParameters struct {
	Id           uint      `json:"id"`
	PreRegTime   time.Time `json:"preRegTime"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	ApprovalTime time.Time `json:"approvalTime"`
	RegTime      time.Time `json:"regTime"`
	UserRoles
	LastBadLogInTime  *time.Time `json:"lastBadLogInTime"`
	BanTime           *time.Time `json:"banTime"`
	LastBadActionTime *time.Time `json:"lastBadActionTime"`
}

func NewUserParameters() (up *UserParameters) {
	return &UserParameters{}
}

func NewUserParametersFromScannableSource(src IScannable) (up *UserParameters, err error) {
	up = NewUserParameters()

	err = src.Scan(
		&up.Id,
		&up.PreRegTime,
		&up.Email,
		&up.Name,
		&up.ApprovalTime,
		&up.RegTime,
		&up.IsAuthor,
		&up.IsWriter,
		&up.IsReader,
		&up.CanLogIn,
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
