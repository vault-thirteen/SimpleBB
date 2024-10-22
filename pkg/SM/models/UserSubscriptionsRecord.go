package models

import (
	"database/sql"
	"errors"

	ul "github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

const (
	IdForVirtualUserSubscriptionsRecord = -1
)

type UserSubscriptionsRecord struct {
	Id      cmb.Id      `json:"id"`
	UserId  cmb.Id      `json:"userId"`
	Threads *ul.UidList `json:"threadIds"`
}

func NewUserSubscriptionsRecord() (usr *UserSubscriptionsRecord) {
	return &UserSubscriptionsRecord{}
}

func NewUserSubscriptionsRecordFromScannableSource(src cm.IScannable) (usr *UserSubscriptionsRecord, err error) {
	usr = NewUserSubscriptionsRecord()
	var x = new(ul.UidList)

	err = src.Scan(
		&usr.Id,
		&usr.UserId,
		x, //&us.Threads,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	usr.Threads = x
	return usr, nil
}

func NewUserSubscriptionsListFromScannableSource(rows cm.IScannableSequence) (usrs []UserSubscriptionsRecord, err error) {
	usrs = make([]UserSubscriptionsRecord, 0)

	var usr *UserSubscriptionsRecord
	for rows.Next() {
		usr = NewUserSubscriptionsRecord()
		usr, err = NewUserSubscriptionsRecordFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		if usr != nil {
			usrs = append(usrs, *usr)
		}
	}

	return usrs, nil
}
