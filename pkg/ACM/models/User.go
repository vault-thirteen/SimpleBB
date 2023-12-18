package models

import (
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

type User struct {
	Password string `json:"-"`
	cm.UserParameters
}

func NewUser() (u *User) {
	return &User{}
}

func NewUserFromScannableSource(src cm.IScannable) (u *User, err error) {
	u = NewUser()

	var up *cm.UserParameters
	up, err = cm.NewUserParametersFromScannableSource(src)
	if err != nil {
		return nil, err
	}

	//u.Password = "" // Password is never shown.
	u.UserParameters = *up

	return u, nil
}
