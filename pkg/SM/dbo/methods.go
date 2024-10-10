package dbo

// Due to the large number of methods, they are sorted alphabetically.

import (
	"database/sql"
	"fmt"

	sm "github.com/vault-thirteen/SimpleBB/pkg/SM/models"
	cdbo "github.com/vault-thirteen/SimpleBB/pkg/common/dbo"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

func (dbo *DatabaseObject) CountUserSubscriptions(userId uint) (n int, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountUserSubscriptions).QueryRow(userId)

	n, err = cm.NewNonNullValueFromScannableSource[int](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountThreadSubscriptions(threadId uint) (n int, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountThreadSubscriptions).QueryRow(threadId)

	n, err = cm.NewNonNullValueFromScannableSource[int](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) InitUserSubscriptions(userId uint) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_InitUserSubscriptions).Exec(userId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) InitThreadSubscriptions(threadId uint) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_InitThreadSubscriptions).Exec(threadId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) GetUserSubscriptions(userId uint) (us *sm.UserSubscriptions, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetUserSubscriptions).QueryRow(userId)

	us, err = sm.NewUserSubscriptionsFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return us, nil
}

func (dbo *DatabaseObject) GetThreadSubscriptions(threadId uint) (ts *sm.ThreadSubscriptions, err error) {
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

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SaveThreadSubscriptions(ts *sm.ThreadSubscriptions) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SaveThreadSubscriptions).Exec(ts.Users, ts.ThreadId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) ClearThreadSubscriptionRecord(threadId uint) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_ClearThreadSubscriptionRecord).Exec(threadId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}
