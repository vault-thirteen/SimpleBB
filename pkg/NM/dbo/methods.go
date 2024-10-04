package dbo

// Due to the large number of methods, they are sorted alphabetically.

import (
	"database/sql"
	"fmt"
	"net"

	nm "github.com/vault-thirteen/SimpleBB/pkg/NM/models"
	cdbo "github.com/vault-thirteen/SimpleBB/pkg/common/dbo"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

func (dbo *DatabaseObject) DeleteNotificationById(notificationId uint) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_DeleteNotificationById).Exec(notificationId)
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

func (dbo *DatabaseObject) GetNotificationById(notificationId uint) (notification *nm.Notification, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetNotificationById).QueryRow(notificationId)

	notification, err = nm.NewNotificationFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

func (dbo *DatabaseObject) GetAllNotificationsByUserId(userId uint) (notifications []nm.Notification, err error) {
	notifications = make([]nm.Notification, 0, 8)

	var rows *sql.Rows
	rows, err = dbo.DatabaseObject.PreparedStatement(DbPsid_GetAllNotificationsByUserId).Query(userId)
	if err != nil {
		return nil, err
	}

	var notification *nm.Notification
	for rows.Next() {
		notification, err = nm.NewNotificationFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		if notification != nil {
			notifications = append(notifications, *notification)
		}
	}

	return notifications, nil
}

func (dbo *DatabaseObject) GetUnreadNotifications(userId uint) (notifications []nm.Notification, err error) {
	notifications = make([]nm.Notification, 0, 8)

	var rows *sql.Rows
	rows, err = dbo.DatabaseObject.PreparedStatement(DbPsid_GetUnreadNotificationsByUserId).Query(userId)
	if err != nil {
		return nil, err
	}

	var notification *nm.Notification
	for rows.Next() {
		notification, err = nm.NewNotificationFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		if notification != nil {
			notifications = append(notifications, *notification)
		}
	}

	return notifications, nil
}

func (dbo *DatabaseObject) CountUnreadNotificationsByUserId(userId uint) (n int, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountUnreadNotificationsByUserId).QueryRow(userId)

	n, err = cm.NewNonNullValueFromScannableSource[int](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) InsertNewNotification(userId uint, text string) (lastInsertedId int64, err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_InsertNewNotification).Exec(userId, text)
	if err != nil {
		return cdbo.LastInsertedIdOnError, err
	}

	lastInsertedId, err = result.LastInsertId()
	if err != nil {
		return cdbo.LastInsertedIdOnError, err
	}

	return lastInsertedId, nil
}

func (dbo *DatabaseObject) MarkNotificationAsRead(notificationId uint, userId uint) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_MarkNotificationAsRead).Exec(notificationId, userId)
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

func (dbo *DatabaseObject) SaveIncident(module cm.Module, incidentType cm.IncidentType, email string, userIPAB net.IP) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SaveIncident).Exec(module, incidentType, email, userIPAB)
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

func (dbo *DatabaseObject) SaveIncidentWithoutUserIPA(module cm.Module, incidentType cm.IncidentType, email string) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SaveIncidentWithoutUserIPA).Exec(module, incidentType, email)
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
