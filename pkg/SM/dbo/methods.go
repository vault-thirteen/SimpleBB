package dbo

// Due to the large number of methods, they are sorted alphabetically.

import (
	"database/sql"

	sm "github.com/vault-thirteen/SimpleBB/pkg/SM/models"
	cdbo "github.com/vault-thirteen/SimpleBB/pkg/common/dbo"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
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

func (dbo *DatabaseObject) GetUserSubscriptions(userId cmb.Id) (us *sm.UserSubscriptions, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetUserSubscriptions).QueryRow(userId)

	us, err = sm.NewUserSubscriptionsFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return us, nil
}

func (dbo *DatabaseObject) GetThreadSubscriptions(threadId cmb.Id) (ts *sm.ThreadSubscriptions, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetThreadSubscriptions).QueryRow(threadId)

	ts, err = sm.NewThreadSubscriptionsFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func (dbo *DatabaseObject) SaveUserSubscriptions(us *sm.UserSubscriptions) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SaveUserSubscriptions).Exec(us.Threads, us.UserId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SaveThreadSubscriptions(ts *sm.ThreadSubscriptions) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SaveThreadSubscriptions).Exec(ts.Users, ts.ThreadId)
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
