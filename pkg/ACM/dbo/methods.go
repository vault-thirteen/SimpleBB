package dbo

// Due to the large number of methods, they are sorted alphabetically.

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"time"

	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	cdbo "github.com/vault-thirteen/SimpleBB/pkg/common/dbo"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	"github.com/vault-thirteen/errorz"
)

func (dbo *DatabaseObject) ApprovePreRegUser(email string) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_ApprovePreRegUser).Exec(email)
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

func (dbo *DatabaseObject) ApproveUserByEmail(email string) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_ApprovePreRegUserEmail).Exec(email)
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

func (dbo *DatabaseObject) AttachVerificationCodeToPreRegUser(email string, code string) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_AttachVerificationCodeToPreRegUser).Exec(code, email)
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

func (dbo *DatabaseObject) AttachVerificationCodeToPreSession(userId uint, requestId string, code string) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_AttachVerificationCodeToPreSession).Exec(code, requestId, userId)
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

func (dbo *DatabaseObject) CheckVerificationCodeForLogIn(requestId string, code string) (ok bool, err error) {
	var n int
	err = dbo.PreparedStatement(DbPsid_CheckVerificationCodeForLogIn).QueryRow(requestId, code).Scan(&n)
	if err != nil {
		return false, err
	}

	if n != 1 {
		return false, nil
	}

	return true, nil
}

func (dbo *DatabaseObject) CheckVerificationCodeForPwdChange(requestId string, code string) (ok bool, err error) {
	var n int
	err = dbo.PreparedStatement(DbPsid_CheckVerificationCodeForPwdChange).QueryRow(requestId, code).Scan(&n)
	if err != nil {
		return false, err
	}

	if n != 1 {
		return false, nil
	}

	return true, nil
}

func (dbo *DatabaseObject) CheckVerificationCodeForPreReg(email, code string) (ok bool, err error) {
	var n int
	err = dbo.PreparedStatement(DbPsid_CheckVerificationCodeForPreReg).QueryRow(email, code).Scan(&n)
	if err != nil {
		return false, err
	}

	if n != 1 {
		return false, nil
	}

	return true, nil
}

func (dbo *DatabaseObject) CheckVerificationCodesForEmailChange(requestId string, codeOld string, codeNew string) (ok bool, err error) {
	var n int
	err = dbo.PreparedStatement(DbPsid_CheckVerificationCodesForEmailChange).QueryRow(requestId, codeOld, codeNew).Scan(&n)
	if err != nil {
		return false, err
	}

	if n != 1 {
		return false, nil
	}

	return true, nil
}

func (dbo *DatabaseObject) CountEmailChangesByUserId(userId uint) (n int, err error) {
	err = dbo.PreparedStatement(DbPsid_CountEmailChangesByUserId).QueryRow(userId).Scan(&n)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountPasswordChangesByUserId(userId uint) (n int, err error) {
	err = dbo.PreparedStatement(DbPsid_CountPasswordChangesByUserId).QueryRow(userId).Scan(&n)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountPreSessionsByUserEmail(email string) (n int, err error) {
	err = dbo.PreparedStatement(DbPsid_CountPreSessionsByUserEmail).QueryRow(email).Scan(&n)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountSessionsByUserEmail(email string) (n int, err error) {
	err = dbo.PreparedStatement(DbPsid_CountSessionsByUserEmail).QueryRow(email).Scan(&n)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountSessionsByUserId(userId uint) (n int, err error) {
	err = dbo.PreparedStatement(DbPsid_CountSessionsByUserId).QueryRow(userId).Scan(&n)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountUsersWithEmailAbleToLogIn(email string) (n int, err error) {
	err = dbo.PreparedStatement(DbPsid_CountUsersWithEmailAbleToLogIn).QueryRow(email).Scan(&n)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountUsersWithEmail(email string) (n int, err error) {
	err = dbo.PreparedStatement(DbPsid_CountUsersWithEmail).QueryRow(email, email).Scan(&n)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountUsersWithName(name string) (n int, err error) {
	err = dbo.PreparedStatement(DbPsid_CountUsersWithName).QueryRow(name, name).Scan(&n)
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
		pcr.NewPassword,
	)
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

func (dbo *DatabaseObject) CreatePreSession(userId uint, requestId string, userIPAB net.IP, pwdSalt []byte, isCaptchaRequired bool, captchaId sql.NullString) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_CreatePreSession).Exec(userId, requestId, userIPAB, pwdSalt, isCaptchaRequired, captchaId)
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

func (dbo *DatabaseObject) CreateSession(userId uint, userIPAB net.IP) (lastInsertedId int64, err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_CreateSession).Exec(userId, userIPAB)
	if err != nil {
		return cdbo.LastInsertedIdOnError, err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return cdbo.LastInsertedIdOnError, err
	}

	if ra != 1 {
		return cdbo.LastInsertedIdOnError, fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	lastInsertedId, err = result.LastInsertId()
	if err != nil {
		return cdbo.LastInsertedIdOnError, err
	}

	return lastInsertedId, nil
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

func (dbo *DatabaseObject) DeleteEmailChangeByRequestId(requestId string) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_DeleteEmailChangeByRequestId).Exec(requestId)
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

func (dbo *DatabaseObject) DeletePasswordChangeByRequestId(requestId string) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_DeletePasswordChangeByRequestId).Exec(requestId)
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

func (dbo *DatabaseObject) DeletePreRegUserIfNotApprovedByEmail(email string) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_DeletePreRegUserIfNotApprovedByEmail).Exec(email)
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

func (dbo *DatabaseObject) DeletePreSessionByRequestId(requestId string) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_DeletePreSessionByRequestId).Exec(requestId)
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

func (dbo *DatabaseObject) DeleteSession(sessionId uint, userId uint, userIPAB net.IP) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_DeleteSession).Exec(sessionId, userId, userIPAB)
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

func (dbo *DatabaseObject) DeleteSessionByUserId(userId uint) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_DeleteSessionByUserId).Exec(userId)
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

func (dbo *DatabaseObject) GetEmailChangeByRequestId(requestId string) (ecr *am.EmailChange, err error) {
	row := dbo.PreparedStatement(DbPsid_GetEmailChangeByRequestId).QueryRow(requestId)

	ecr, err = am.NewEmailChangeFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return ecr, nil
}

func (dbo *DatabaseObject) GetListOfLoggedUsers() (userIds []uint, err error) {
	userIds = make([]uint, 0)

	var rows *sql.Rows
	rows, err = dbo.PreparedStatement(DbPsid_GetListOfLoggedUsers).Query()
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = errorz.Combine(err, derr)
		}
	}()

	var userId uint
	for rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return userIds, nil // No rows.
			} else {
				return nil, err
			}
		}

		userIds = append(userIds, userId)
	}

	return userIds, nil
}

func (dbo *DatabaseObject) GetPasswordChangeByRequestId(requestId string) (pcr *am.PasswordChange, err error) {
	row := dbo.PreparedStatement(DbPsid_GetPasswordChangeByRequestId).QueryRow(requestId)

	pcr, err = am.NewPasswordChangeFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return pcr, nil
}

func (dbo *DatabaseObject) GetPreSessionByRequestId(requestId string) (preSession *am.PreSession, err error) {
	row := dbo.PreparedStatement(DbPsid_GetPreSessionByRequestId).QueryRow(requestId)

	preSession, err = am.NewPreSessionFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return preSession, nil
}

func (dbo *DatabaseObject) GetSessionByUserId(userId uint) (session *am.Session, err error) {
	row := dbo.PreparedStatement(DbPsid_GetSessionByUserId).QueryRow(userId)

	session, err = am.NewSessionFromScannableSource(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No rows.
		} else {
			return nil, err
		}
	}

	return session, nil
}

func (dbo *DatabaseObject) GetUserById(userId uint) (user *am.User, err error) {
	row := dbo.PreparedStatement(DbPsid_GetUserById).QueryRow(userId)

	user, err = am.NewUserFromScannableSource(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No rows.
		} else {
			return nil, err
		}
	}

	return user, nil
}

func (dbo *DatabaseObject) GetUserIdByEmail(email string) (id uint, err error) {
	err = dbo.PreparedStatement(DbPsid_GetUserIdByEmail).QueryRow(email).Scan(&id)
	if err != nil {
		return cdbo.IdOnError, err
	}

	return id, nil
}

func (dbo *DatabaseObject) GetUserLastBadActionTimeById(userId uint) (lastBadActionTime *time.Time, err error) {
	var u = &am.User{}
	err = dbo.PreparedStatement(DbPsid_GetUserLastBadActionTimeById).QueryRow(userId).Scan(&u.LastBadActionTime)
	if err != nil {
		return nil, err
	}

	return u.LastBadActionTime, nil
}

func (dbo *DatabaseObject) GetUserLastBadLogInTimeByEmail(email string) (lastBadLogInTime *time.Time, err error) {
	var u = &am.User{}
	err = dbo.PreparedStatement(DbPsid_GetUserLastBadLogInTimeByEmail).QueryRow(email).Scan(&u.LastBadLogInTime)
	if err != nil {
		return nil, err
	}

	return u.LastBadLogInTime, nil
}

func (dbo *DatabaseObject) GetUserPasswordById(userId uint) (password []byte, err error) {
	err = dbo.PreparedStatement(DbPsid_GetUserPasswordById).QueryRow(userId).Scan(&password)
	if err != nil {
		return nil, err
	}

	return password, nil
}

func (dbo *DatabaseObject) GetUserRolesById(userId uint) (roles *cm.UserRoles, err error) {
	row := dbo.PreparedStatement(DbPsid_GetUserRolesById).QueryRow(userId)

	roles, err = cm.NewUserRolesFromScannableSource(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return roles, nil
}

func (dbo *DatabaseObject) InsertPreRegisteredUser(email string) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_InsertPreRegisteredUser).Exec(email)
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

func (dbo *DatabaseObject) RegisterPreRegUser(email string) (err error) {
	// Part 1.
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_RegisterPreRegUserP1).Exec(email)
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

	// Part 2.
	result, err = dbo.PreparedStatement(DbPsid_RegisterPreRegUserP2).Exec(email)
	if err != nil {
		return err
	}

	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SaveIncident(incidentType am.IncidentType, email string, userIPAB net.IP) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SaveIncident).Exec(incidentType, email, userIPAB)
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

func (dbo *DatabaseObject) SaveIncidentWithoutUserIPA(incidentType am.IncidentType, email string) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SaveIncidentWithoutUserIPA).Exec(incidentType, email)
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

func (dbo *DatabaseObject) SaveLogEvent(logEvent *am.LogEvent) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SaveLogEvent).Exec(
		logEvent.Type,
		logEvent.UserId,
		logEvent.Email,
		logEvent.UserIPAB,
	)
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

func (dbo *DatabaseObject) SetEmailChangeVFlags(userId uint, requestId string, ecvf *am.EmailChangeVerificationFlags) (err error) {
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

func (dbo *DatabaseObject) SetPasswordChangeVFlags(userId uint, requestId string, pcvf *am.PasswordChangeVerificationFlags) (err error) {
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

func (dbo *DatabaseObject) SetPreRegUserData(email string, code string, name string, password []byte) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetPreRegUserData).Exec(name, password, email, code)
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

func (dbo *DatabaseObject) SetPreRegUserEmailSendStatus(emailSendStatus bool, email string) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetPreRegUserEmailSendStatus).Exec(emailSendStatus, email)
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

func (dbo *DatabaseObject) SetPreSessionCaptchaFlags(userId uint, requestId string, isVerifiedByCaptcha bool) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetPreSessionCaptchaFlag).Exec(isVerifiedByCaptcha, requestId, userId)
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

func (dbo *DatabaseObject) SetPreSessionEmailSendStatus(userId uint, requestId string, emailSendStatus bool) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetPreSessionEmailSendStatus).Exec(emailSendStatus, requestId, userId)
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

func (dbo *DatabaseObject) SetPreSessionPasswordFlag(userId uint, requestId string, isVerifiedByPassword bool) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetPreSessionPasswordFlag).Exec(isVerifiedByPassword, requestId, userId)
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

func (dbo *DatabaseObject) SetPreSessionVerificationFlag(userId uint, requestId string, isVerifiedByEmail bool) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetPreSessionVerificationFlag).Exec(isVerifiedByEmail, requestId, userId)
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

func (dbo *DatabaseObject) SetUserEmail(userId uint, email string, newEmail string) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetUserEmail).Exec(newEmail, userId, email)
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

func (dbo *DatabaseObject) SetUserPassword(userId uint, email string, newPassword []byte) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetUserPassword).Exec(newPassword, userId, email)
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

func (dbo *DatabaseObject) SetUserRoleAuthor(userId uint, isRoleEnabled bool) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetUserRoleAuthor).Exec(isRoleEnabled, userId)
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

func (dbo *DatabaseObject) SetUserRoleCanLogIn(userId uint, isRoleEnabled bool) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetUserRoleCanLogIn).Exec(isRoleEnabled, userId)
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

func (dbo *DatabaseObject) SetUserRoleReader(userId uint, isRoleEnabled bool) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetUserRoleReader).Exec(isRoleEnabled, userId)
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

func (dbo *DatabaseObject) SetUserRoleWriter(userId uint, isRoleEnabled bool) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetUserRoleWriter).Exec(isRoleEnabled, userId)
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

func (dbo *DatabaseObject) UpdatePreSessionRequestId(userId uint, requestIdOld string, requestIdNew string) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_UpdatePreSessionRequestId).Exec(requestIdNew, requestIdOld, userId)
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

func (dbo *DatabaseObject) UpdateUserBanTime(userId uint) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_UpdateUserBanTime).Exec(userId)
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

func (dbo *DatabaseObject) UpdateUserLastBadActionTimeById(userId uint) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_UpdateUserLastBadActionTimeById).Exec(userId)
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

func (dbo *DatabaseObject) UpdateUserLastBadLogInTimeByEmail(email string) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_UpdateUserLastBadLogInTimeByEmail).Exec(email)
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

func (dbo *DatabaseObject) ViewUserParametersById(userId uint) (userParameters *cm.UserParameters, err error) {
	row := dbo.PreparedStatement(DbPsid_GetUserParametersById).QueryRow(userId)

	userParameters, err = cm.NewUserParametersFromScannableSource(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return userParameters, nil
}
