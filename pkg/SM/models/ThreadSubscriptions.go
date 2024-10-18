package models

import (
	"database/sql"
	"errors"

	ul "github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type ThreadSubscriptions struct {
	Id       cmb.Id      `json:"id"`
	ThreadId cmb.Id      `json:"threadId"`
	Users    *ul.UidList `json:"userIds"`
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
