package dbo

// Due to the large number of methods, they are sorted alphabetically.

import (
	"database/sql"
	"net"
	"time"

	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	cdbo "github.com/vault-thirteen/SimpleBB/pkg/common/dbo"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	ae "github.com/vault-thirteen/auxie/errors"
)

func (dbo *DatabaseObject) ApprovePreRegUser(email cm.Email) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_ApprovePreRegUser).Exec(email)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) ApproveUserByEmail(email cm.Email) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_ApprovePreRegUserEmail).Exec(email)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) AttachVerificationCodeToPreRegUser(email cm.Email, code cm.VerificationCode) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_AttachVerificationCodeToPreRegUser).Exec(code, email)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) AttachVerificationCodeToPreSession(userId cmb.Id, requestId cm.RequestId, code cm.VerificationCode) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_AttachVerificationCodeToPreSession).Exec(code, requestId, userId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) CheckVerificationCodeForLogIn(requestId cm.RequestId, code cm.VerificationCode) (ok bool, err error) {
	row := dbo.PreparedStatement(DbPsid_CheckVerificationCodeForLogIn).QueryRow(requestId, code)

	var n int
	n, err = cm.NewNonNullValueFromScannableSource[int](row)
	if err != nil {
		return false, err
	}

	if n != 1 {
		return false, nil
	}

	return true, nil
}

func (dbo *DatabaseObject) CheckVerificationCodeForPwdChange(requestId cm.RequestId, code cm.VerificationCode) (ok bool, err error) {
	row := dbo.PreparedStatement(DbPsid_CheckVerificationCodeForPwdChange).QueryRow(requestId, code)

	var n int
	n, err = cm.NewNonNullValueFromScannableSource[int](row)
	if err != nil {
		return false, err
	}

	if n != 1 {
		return false, nil
	}

	return true, nil
}

func (dbo *DatabaseObject) CheckVerificationCodeForPreReg(email cm.Email, code cm.VerificationCode) (ok bool, err error) {
	row := dbo.PreparedStatement(DbPsid_CheckVerificationCodeForPreReg).QueryRow(email, code)

	var n int
	n, err = cm.NewNonNullValueFromScannableSource[int](row)
	if err != nil {
		return false, err
	}

	if n != 1 {
		return false, nil
	}

	return true, nil
}

func (dbo *DatabaseObject) CheckVerificationCodesForEmailChange(requestId cm.RequestId, codeOld cm.VerificationCode, codeNew cm.VerificationCode) (ok bool, err error) {
	row := dbo.PreparedStatement(DbPsid_CheckVerificationCodesForEmailChange).QueryRow(requestId, codeOld, codeNew)

	var n int
	n, err = cm.NewNonNullValueFromScannableSource[int](row)
	if err != nil {
		return false, err
	}

	if n != 1 {
		return false, nil
	}

	return true, nil
}

func (dbo *DatabaseObject) CountAllUsers() (n cmb.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountAllUsers).QueryRow()

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountRegistrationsReadyForApproval() (n cmb.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountRegistrationsReadyForApproval).QueryRow()

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountEmailChangesByUserId(userId cmb.Id) (n cmb.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountEmailChangesByUserId).QueryRow(userId)

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountLoggedUsers() (n cmb.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountLoggedUsers).QueryRow()

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountPasswordChangesByUserId(userId cmb.Id) (n cmb.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountPasswordChangesByUserId).QueryRow(userId)

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountPreSessionsByUserEmail(email cm.Email) (n cmb.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountPreSessionsByUserEmail).QueryRow(email)

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountSessionsByUserEmail(email cm.Email) (n cmb.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountSessionsByUserEmail).QueryRow(email)

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountSessionsByUserId(userId cmb.Id) (n cmb.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountSessionsByUserId).QueryRow(userId)

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountUsersWithEmailAbleToLogIn(email cm.Email) (n cmb.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountUsersWithEmailAbleToLogIn).QueryRow(email)

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountUsersWithEmail(email cm.Email) (n cmb.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountUsersWithEmail).QueryRow(email, email)

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountUsersWithName(name cm.Name) (n cmb.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountUsersWithName).QueryRow(name, name)

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CreateEmailChangeRequest(ecr *am.EmailChange) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_CreateEmailChangeRequest).Exec(
		ecr.UserId,
		ecr.RequestId,
		ecr.UserIPAB,
		ecr.AuthDataBytes,
		ecr.IsCaptchaRequired,
		ecr.CaptchaId,
		ecr.VerificationCodeOld,
		true,
		ecr.NewEmail,
		ecr.VerificationCodeNew,
		true,
	)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) CreatePasswordChangeRequest(pcr *am.PasswordChange) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_CreatePasswordChangeRequest).Exec(
		pcr.UserId,
		pcr.RequestId,
		pcr.UserIPAB,
		pcr.AuthDataBytes,
		pcr.IsCaptchaRequired,
		pcr.CaptchaId,
		pcr.VerificationCode,
		true,
		pcr.NewPasswordBytes,
	)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) CreatePreSession(userId cmb.Id, requestId cm.RequestId, userIPAB net.IP, pwdSalt []byte, isCaptchaRequired cmb.Flag, captchaId *cm.CaptchaId) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_CreatePreSession).Exec(userId, requestId, userIPAB, pwdSalt, isCaptchaRequired, captchaId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) CreateSession(userId cmb.Id, userIPAB net.IP) (lastInsertedId cmb.Id, err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_CreateSession).Exec(userId, userIPAB)
	if err != nil {
		return cdbo.LastInsertedIdOnError, err
	}

	return cdbo.CheckRowsAffectedAndGetLastInsertedId(result, 1)
}

func (dbo *DatabaseObject) DeleteAbandonedPreSessions() (err error) {
	timeBorder := time.Now().Add(-time.Duration(dbo.sp.PreSessionExpirationTime) * time.Second)

	_, err = dbo.PreparedStatement(DbPsid_DeleteAbandonedPreSessions).Exec(timeBorder)
	if err != nil {
		return err
	}

	// Affected rows are not checked while they may be none or multiple.

	return nil
}

func (dbo *DatabaseObject) DeleteEmailChangeByRequestId(requestId cm.RequestId) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_DeleteEmailChangeByRequestId).Exec(requestId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) DeletePasswordChangeByRequestId(requestId cm.RequestId) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_DeletePasswordChangeByRequestId).Exec(requestId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) DeletePreRegUserIfNotApprovedByEmail(email cm.Email) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_DeletePreRegUserIfNotApprovedByEmail).Exec(email)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) DeletePreSessionByRequestId(requestId cm.RequestId) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_DeletePreSessionByRequestId).Exec(requestId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) DeleteSession(sessionId cmb.Id, userId cmb.Id, userIPAB net.IP) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_DeleteSession).Exec(sessionId, userId, userIPAB)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) DeleteSessionByUserId(userId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_DeleteSessionByUserId).Exec(userId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) GetEmailChangeByRequestId(requestId cm.RequestId) (ecr *am.EmailChange, err error) {
	row := dbo.PreparedStatement(DbPsid_GetEmailChangeByRequestId).QueryRow(requestId)
	return am.NewEmailChangeFromScannableSource(row)
}

func (dbo *DatabaseObject) GetListOfAllUsers() (userIds []cmb.Id, err error) {
	var rows *sql.Rows
	rows, err = dbo.PreparedStatement(DbPsid_GetListOfAllUsers).Query()
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

func (dbo *DatabaseObject) GetListOfAllUsersOnPage(pageNumber cmb.Count, pageSize cmb.Count) (userIds []cmb.Id, err error) {
	var rows *sql.Rows
	rows, err = dbo.PreparedStatement(DbPsid_GetListOfAllUsersOnPage).Query(pageSize, (pageNumber-1)*pageSize)
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

func (dbo *DatabaseObject) GetListOfLoggedUsers() (userIds []cmb.Id, err error) {
	var rows *sql.Rows
	rows, err = dbo.PreparedStatement(DbPsid_GetListOfLoggedUsers).Query()
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

func (dbo *DatabaseObject) GetListOfLoggedUsersOnPage(pageNumber cmb.Count, pageSize cmb.Count) (userIds []cmb.Id, err error) {
	var rows *sql.Rows
	rows, err = dbo.PreparedStatement(DbPsid_GetListOfLoggedUsersOnPage).Query(pageSize, (pageNumber-1)*pageSize)
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

func (dbo *DatabaseObject) GetListOfRegistrationsReadyForApproval(pageNumber cmb.Count, pageSize cmb.Count) (rrfas []am.RegistrationReadyForApproval, err error) {
	var rows *sql.Rows
	rows, err = dbo.PreparedStatement(DbPsid_GetListOfRegistrationsReadyForApproval).Query(pageSize, (pageNumber-1)*pageSize)
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return am.NewRegistrationReadyForApprovalArrayFromRows(rows)
}

func (dbo *DatabaseObject) GetPasswordChangeByRequestId(requestId cm.RequestId) (pcr *am.PasswordChange, err error) {
	row := dbo.PreparedStatement(DbPsid_GetPasswordChangeByRequestId).QueryRow(requestId)
	return am.NewPasswordChangeFromScannableSource(row)
}

func (dbo *DatabaseObject) GetPreSessionByRequestId(requestId cm.RequestId) (preSession *am.PreSession, err error) {
	row := dbo.PreparedStatement(DbPsid_GetPreSessionByRequestId).QueryRow(requestId)
	return am.NewPreSessionFromScannableSource(row)
}

func (dbo *DatabaseObject) GetSessionByUserId(userId cmb.Id) (session *am.Session, err error) {
	row := dbo.PreparedStatement(DbPsid_GetSessionByUserId).QueryRow(userId)
	return am.NewSessionFromScannableSource(row)
}

func (dbo *DatabaseObject) GetUserNameById(userId cmb.Id) (userName *cm.Name, err error) {
	row := dbo.PreparedStatement(DbPsid_GetUserNameById).QueryRow(userId)
	return cm.NewValueFromScannableSource[cm.Name](row)
}

func (dbo *DatabaseObject) GetUserById(userId cmb.Id) (user *am.User, err error) {
	row := dbo.PreparedStatement(DbPsid_GetUserById).QueryRow(userId)
	return am.NewUserFromScannableSource(row)
}

func (dbo *DatabaseObject) GetUserIdByEmail(email cm.Email) (userId cmb.Id, err error) {
	row := dbo.PreparedStatement(DbPsid_GetUserIdByEmail).QueryRow(email)

	userId, err = cm.NewNonNullValueFromScannableSource[cmb.Id](row)
	if err != nil {
		return cdbo.IdOnError, err
	}

	return userId, nil
}

func (dbo *DatabaseObject) GetUserLastBadActionTimeById(userId cmb.Id) (lastBadActionTime *time.Time, err error) {
	var u = &am.User{}
	err = dbo.PreparedStatement(DbPsid_GetUserLastBadActionTimeById).QueryRow(userId).Scan(&u.LastBadActionTime)
	if err != nil {
		return nil, err
	}

	return u.LastBadActionTime, nil
}

func (dbo *DatabaseObject) GetUserLastBadLogInTimeByEmail(email cm.Email) (lastBadLogInTime *time.Time, err error) {
	var u = &am.User{}
	err = dbo.PreparedStatement(DbPsid_GetUserLastBadLogInTimeByEmail).QueryRow(email).Scan(&u.LastBadLogInTime)
	if err != nil {
		return nil, err
	}

	return u.LastBadLogInTime, nil
}

func (dbo *DatabaseObject) GetUserPasswordById(userId cmb.Id) (password *[]byte, err error) {
	row := dbo.PreparedStatement(DbPsid_GetUserPasswordById).QueryRow(userId)
	return cm.NewValueFromScannableSource[[]byte](row)
}

func (dbo *DatabaseObject) GetUserRolesById(userId cmb.Id) (roles *cm.UserRoles, err error) {
	row := dbo.PreparedStatement(DbPsid_GetUserRolesById).QueryRow(userId)
	return cm.NewUserRolesFromScannableSource(row)
}

func (dbo *DatabaseObject) InsertPreRegisteredUser(email cm.Email) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_InsertPreRegisteredUser).Exec(email)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) RegisterPreRegUser(email cm.Email) (err error) {
	// Part 1.
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_RegisterPreRegUserP1).Exec(email)
	if err != nil {
		return err
	}

	err = cdbo.CheckRowsAffected(result, 1)
	if err != nil {
		return err
	}

	// Part 2.
	result, err = dbo.PreparedStatement(DbPsid_RegisterPreRegUserP2).Exec(email)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) RejectRegistrationRequest(id cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_RejectRegistrationRequest).Exec(id)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SaveIncident(module cmb.EnumValue, incidentType cm.IncidentType, email cm.Email, userIPAB net.IP) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SaveIncident).Exec(module, incidentType, email, userIPAB)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SaveIncidentWithoutUserIPA(module cmb.EnumValue, incidentType cm.IncidentType, email cm.Email) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SaveIncidentWithoutUserIPA).Exec(module, incidentType, email)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SaveLogEvent(logEvent *am.LogEvent) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SaveLogEvent).Exec(
		logEvent.Type,
		logEvent.UserId,
		logEvent.Email,
		logEvent.UserIPAB,
		logEvent.AdminId,
	)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetEmailChangeVFlags(userId cmb.Id, requestId cm.RequestId, ecvf *am.EmailChangeVerificationFlags) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetEmailChangeVFlags).Exec(
		ecvf.IsVerifiedByCaptcha,
		ecvf.IsVerifiedByPassword,
		ecvf.IsVerifiedByOldEmail,
		ecvf.IsVerifiedByNewEmail,
		requestId,
		userId,
	)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetPasswordChangeVFlags(userId cmb.Id, requestId cm.RequestId, pcvf *am.PasswordChangeVerificationFlags) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetPasswordChangeVFlags).Exec(
		pcvf.IsVerifiedByCaptcha,
		pcvf.IsVerifiedByPassword,
		pcvf.IsVerifiedByEmail,
		requestId,
		userId,
	)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetPreRegUserData(email cm.Email, code cm.VerificationCode, name cm.Name, password []byte) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetPreRegUserData).Exec(name, password, email, code)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetPreRegUserEmailSendStatus(emailSendStatus cmb.Flag, email cm.Email) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetPreRegUserEmailSendStatus).Exec(emailSendStatus, email)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetPreSessionCaptchaFlags(userId cmb.Id, requestId cm.RequestId, isVerifiedByCaptcha cmb.Flag) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetPreSessionCaptchaFlag).Exec(isVerifiedByCaptcha, requestId, userId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetPreSessionEmailSendStatus(userId cmb.Id, requestId cm.RequestId, emailSendStatus cmb.Flag) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetPreSessionEmailSendStatus).Exec(emailSendStatus, requestId, userId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetPreSessionPasswordFlag(userId cmb.Id, requestId cm.RequestId, isVerifiedByPassword cmb.Flag) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetPreSessionPasswordFlag).Exec(isVerifiedByPassword, requestId, userId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetPreSessionVerificationFlag(userId cmb.Id, requestId cm.RequestId, isVerifiedByEmail cmb.Flag) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetPreSessionVerificationFlag).Exec(isVerifiedByEmail, requestId, userId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetUserEmail(userId cmb.Id, email cm.Email, newEmail cm.Email) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetUserEmail).Exec(newEmail, userId, email)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetUserPassword(userId cmb.Id, email cm.Email, newPassword []byte) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetUserPassword).Exec(newPassword, userId, email)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetUserRoleAuthor(userId cmb.Id, isRoleEnabled cmb.Flag) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetUserRoleAuthor).Exec(isRoleEnabled, userId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetUserRoleCanLogIn(userId cmb.Id, isRoleEnabled cmb.Flag) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetUserRoleCanLogIn).Exec(isRoleEnabled, userId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetUserRoleReader(userId cmb.Id, isRoleEnabled cmb.Flag) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetUserRoleReader).Exec(isRoleEnabled, userId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetUserRoleWriter(userId cmb.Id, isRoleEnabled cmb.Flag) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetUserRoleWriter).Exec(isRoleEnabled, userId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) UpdatePreSessionRequestId(userId cmb.Id, requestIdOld cm.RequestId, requestIdNew cm.RequestId) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_UpdatePreSessionRequestId).Exec(requestIdNew, requestIdOld, userId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) UpdateUserBanTime(userId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_UpdateUserBanTime).Exec(userId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) UpdateUserLastBadActionTimeById(userId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_UpdateUserLastBadActionTimeById).Exec(userId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) UpdateUserLastBadLogInTimeByEmail(email cm.Email) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_UpdateUserLastBadLogInTimeByEmail).Exec(email)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) ViewUserParametersById(userId cmb.Id) (userParameters *cm.UserParameters, err error) {
	row := dbo.PreparedStatement(DbPsid_GetUserParametersById).QueryRow(userId)
	return cm.NewUserParametersFromScannableSource(row)
}
