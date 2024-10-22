package models

import (
	"database/sql"
	"errors"

	ul "github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type ThreadSubscriptionsRecord struct {
	Id       cmb.Id      `json:"id"`
	ThreadId cmb.Id      `json:"threadId"`
	Users    *ul.UidList `json:"userIds"`
}

func NewThreadSubscriptionsRecord() (tsr *ThreadSubscriptionsRecord) {
	return &ThreadSubscriptionsRecord{}
}

func NewThreadSubscriptionsRecordFromScannableSource(src cm.IScannable) (tsr *ThreadSubscriptionsRecord, err error) {
	tsr = NewThreadSubscriptionsRecord()
	var x = ul.New()

	err = src.Scan(
		&tsr.Id,
		&tsr.ThreadId,
		x, //&tsr.Users,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	tsr.Users = x
	return tsr, nil
}

func NewThreadSubscriptionsListFromScannableSource(rows cm.IScannableSequence) (tsrs []ThreadSubscriptionsRecord, err error) {
	tsrs = make([]ThreadSubscriptionsRecord, 0)

	var tsr *ThreadSubscriptionsRecord
	for rows.Next() {
		tsr = NewThreadSubscriptionsRecord()
		tsr, err = NewThreadSubscriptionsRecordFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		if tsr != nil {
			tsrs = append(tsrs, *tsr)
		}
	}

	return tsrs, nil
}
