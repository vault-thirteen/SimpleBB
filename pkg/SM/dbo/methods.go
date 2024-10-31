package dbo

// Due to the large number of methods, they are sorted alphabetically.

import (
	"database/sql"

	sm "github.com/vault-thirteen/SimpleBB/pkg/SM/models"
	cdbo "github.com/vault-thirteen/SimpleBB/pkg/common/dbo"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	ae "github.com/vault-thirteen/auxie/errors"
)

func (dbo *DatabaseObject) CountUserSubscriptions(userId cmb.Id) (n cmb.Count, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountUserSubscriptions).QueryRow(userId)

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountThreadSubscriptions(threadId cmb.Id) (n cmb.Count, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountThreadSubscriptions).QueryRow(threadId)

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) InitUserSubscriptions(userId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_InitUserSubscriptions).Exec(userId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) InitThreadSubscriptions(threadId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_InitThreadSubscriptions).Exec(threadId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) GetAllThreadSubscriptions() (tsrs []sm.ThreadSubscriptionsRecord, err error) {
	var rows *sql.Rows
	rows, err = dbo.DatabaseObject.PreparedStatement(DbPsid_GetAllThreadSubscriptions).Query()
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return sm.NewThreadSubscriptionsRecordArrayFromRows(rows)
}

func (dbo *DatabaseObject) GetAllUserSubscriptions() (usrs []sm.UserSubscriptionsRecord, err error) {
	var rows *sql.Rows
	rows, err = dbo.DatabaseObject.PreparedStatement(DbPsid_GetAllUserSubscriptions).Query()
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return sm.NewUserSubscriptionsRecordArrayFromRows(rows)
}

func (dbo *DatabaseObject) GetUserSubscriptions(userId cmb.Id) (usr *sm.UserSubscriptionsRecord, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetUserSubscriptions).QueryRow(userId)

	usr, err = sm.NewUserSubscriptionsRecordFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	// If record does not exist, we return a virtual empty record.
	// Virtual record is marked with a negative ID in order to distinguish it
	// from real records.
	if usr == nil {
		usr = &sm.UserSubscriptionsRecord{
			Id:      sm.IdForVirtualUserSubscriptionsRecord,
			UserId:  userId,
			Threads: nil,
		}
	}

	return usr, nil
}

func (dbo *DatabaseObject) GetThreadSubscriptions(threadId cmb.Id) (tsr *sm.ThreadSubscriptionsRecord, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetThreadSubscriptions).QueryRow(threadId)

	tsr, err = sm.NewThreadSubscriptionsRecordFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return tsr, nil
}

func (dbo *DatabaseObject) SaveUserSubscriptions(usr *sm.UserSubscriptionsRecord) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SaveUserSubscriptions).Exec(usr.Threads, usr.UserId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SaveThreadSubscriptions(tsr *sm.ThreadSubscriptionsRecord) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SaveThreadSubscriptions).Exec(tsr.Users, tsr.ThreadId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) ClearThreadSubscriptionRecord(threadId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_ClearThreadSubscriptionRecord).Exec(threadId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}
