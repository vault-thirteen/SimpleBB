package models

import (
	"database/sql"
	"errors"

	ul "github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

type ThreadSubscriptions struct {
	Id       uint        `json:"id"`
	ThreadId uint        `json:"threadId"`
	Users    *ul.UidList `json:"users"`
}

func NewThreadSubscriptions() (ts *ThreadSubscriptions) {
	return &ThreadSubscriptions{}
}

func NewThreadSubscriptionsFromScannableSource(src cm.IScannable) (ts *ThreadSubscriptions, err error) {
	ts = NewThreadSubscriptions()
	var x = ul.New()

	err = src.Scan(
		&ts.Id,
		&ts.ThreadId,
		x, //&ts.Users,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	ts.Users = x
	return ts, nil
}
