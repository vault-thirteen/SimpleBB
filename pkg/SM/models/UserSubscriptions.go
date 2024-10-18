package models

import (
	"database/sql"
	"errors"

	ul "github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type UserSubscriptions struct {
	Id      cmb.Id      `json:"id"`
	UserId  cmb.Id      `json:"userId"`
	Threads *ul.UidList `json:"threadIds"`
}

func NewUserSubscriptions() (us *UserSubscriptions) {
	return &UserSubscriptions{}
}

func NewUserSubscriptionsFromScannableSource(src cm.IScannable) (us *UserSubscriptions, err error) {
	us = NewUserSubscriptions()
	var x = new(ul.UidList)

	err = src.Scan(
		&us.Id,
		&us.UserId,
		x, //&us.Threads,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	us.Threads = x
	return us, nil
}
