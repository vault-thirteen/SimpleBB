package models

import (
	"database/sql"
	"errors"
)

type UserRoles struct {
	IsAdministrator bool `json:"isAdministrator"`
	IsModerator     bool `json:"isModerator"`
	IsAuthor        bool `json:"isAuthor"`
	IsWriter        bool `json:"isWriter"`
	IsReader        bool `json:"isReader"`
	CanLogIn        bool `json:"canLogIn"`
}

func NewUserRoles() (ur *UserRoles) {
	return &UserRoles{}
}

func NewUserRolesFromScannableSource(src IScannable) (ur *UserRoles, err error) {
	ur = NewUserRoles()

	err = src.Scan(
		&ur.IsAuthor,
		&ur.IsWriter,
		&ur.IsReader,
		&ur.CanLogIn,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return ur, nil
}
