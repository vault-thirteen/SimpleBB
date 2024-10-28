package dbo

// Due to the large number of methods, they are sorted alphabetically.

import (
	"database/sql"
	"net"

	nm "github.com/vault-thirteen/SimpleBB/pkg/NM/models"
	cdbo "github.com/vault-thirteen/SimpleBB/pkg/common/dbo"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	ae "github.com/vault-thirteen/auxie/errors"
)

func (dbo *DatabaseObject) AddResource(r *cm.Resource) (lastInsertedId cmb.Id, err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_AddResource).Exec(r.Type, r.Text, r.Number)
	if err != nil {
		return cdbo.LastInsertedIdOnError, err
	}

	return cdbo.CheckRowsAffectedAndGetLastInsertedId(result, 1)
}

func (dbo *DatabaseObject) CountAllNotificationsByUserId(userId cmb.Id) (n cmb.Count, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountAllNotificationsByUserId).QueryRow(userId)

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountAllResources() (n cmb.Count, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountAllResources).QueryRow()

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountUnreadNotificationsByUserId(userId cmb.Id) (n cmb.Count, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountUnreadNotificationsByUserId).QueryRow(userId)

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) DeleteNotificationById(notificationId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_DeleteNotificationById).Exec(notificationId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) DeleteResourceById(resourceId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_DeleteResourceById).Exec(resourceId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) GetAllNotificationsByUserId(userId cmb.Id) (notifications []nm.Notification, err error) {
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

func (dbo *DatabaseObject) GetNotificationById(notificationId cmb.Id) (notification *nm.Notification, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetNotificationById).QueryRow(notificationId)

	notification, err = nm.NewNotificationFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

func (dbo *DatabaseObject) GetNotificationsByUserIdOnPage(userId cmb.Id, pageNumber cmb.Count, pageSize cmb.Count) (notifications []nm.Notification, err error) {
	var rows *sql.Rows
	rows, err = dbo.PreparedStatement(DbPsid_GetNotificationsByUserIdOnPage).Query(userId, pageSize, (pageNumber-1)*pageSize)
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return nm.NewNotificationArrayFromRows(rows)
}

func (dbo *DatabaseObject) GetResourceById(resourceId cmb.Id) (r *cm.Resource, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetResourceById).QueryRow(resourceId)

	r, err = cm.NewResourceFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (dbo *DatabaseObject) GetResourceIdsOnPage(pageNumber cmb.Count, pageSize cmb.Count) (resourceIds []cmb.Id, err error) {
	var rows *sql.Rows
	rows, err = dbo.PreparedStatement(DbPsid_ListAllResourceIdsOnPage).Query(pageSize, (pageNumber-1)*pageSize)
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return cm.NewArrayFromScannableSource[cmb.Id](rows)
}

func (dbo *DatabaseObject) GetSystemEventById(systemEventId cmb.Id) (se *cm.SystemEvent, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetSystemEventById).QueryRow(systemEventId)

	se, err = cm.NewSystemEventFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return se, nil
}

func (dbo *DatabaseObject) GetUnreadNotifications(userId cmb.Id) (notifications []nm.Notification, err error) {
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

func (dbo *DatabaseObject) InsertNewNotification(userId cmb.Id, text cmb.Text) (lastInsertedId cmb.Id, err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_InsertNewNotification).Exec(userId, text)
	if err != nil {
		return cdbo.LastInsertedIdOnError, err
	}

	return cdbo.CheckRowsAffectedAndGetLastInsertedId(result, 1)
}

func (dbo *DatabaseObject) MarkNotificationAsRead(notificationId cmb.Id, userId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_MarkNotificationAsRead).Exec(notificationId, userId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SaveIncident(module cm.Module, incidentType cm.IncidentType, email cm.Email, userIPAB net.IP) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SaveIncident).Exec(module, incidentType, email, userIPAB)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SaveIncidentWithoutUserIPA(module cm.Module, incidentType cm.IncidentType, email cm.Email) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SaveIncidentWithoutUserIPA).Exec(module, incidentType, email)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SaveSystemEvent(se cm.SystemEvent) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SaveSystemEvent).Exec(se.Type, se.ThreadId, se.MessageId, se.UserId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}
