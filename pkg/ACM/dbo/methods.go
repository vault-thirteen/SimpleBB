package dbo

// Due to the large number of methods, they are sorted alphabetically.

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"time"

	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	"github.com/vault-thirteen/errorz"
)

func (dbo *DatabaseObject) ApprovePreRegUser(email string) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_ApprovePreRegUser].Exec(email)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) ApproveUserByEmail(email string) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_ApprovePreRegUserEmail].Exec(email)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) AttachVerificationCodeToPreRegUser(email string, code string) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_AttachVerificationCodeToPreRegUser].Exec(code, email)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) AttachVerificationCodeToPreSession(userId uint, requestId string, code string) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_AttachVerificationCodeToPreSession].Exec(code, requestId, userId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) CheckLogInVerificationCode(requestId string, code string) (ok bool, err error) {
	var n int
	err = dbo.preparedStatements[DbPsid_CheckVerificationCodeOfLogInUser].QueryRow(requestId, code).Scan(&n)
	if err != nil {
		return false, err
	}

	if n != 1 {
		return false, nil
	}

	return true, nil
}

func (dbo *DatabaseObject) CheckRegVerificationCode(email, code string) (ok bool, err error) {
	var n int
	err = dbo.preparedStatements[DbPsid_CheckVerificationCodeOfPreRegUser].QueryRow(email, code).Scan(&n)
	if err != nil {
		return false, err
	}

	if n != 1 {
		return false, nil
	}

	return true, nil
}

func (dbo *DatabaseObject) CountCreatedUserPreSessionsByEmail(email string) (n int, err error) {
	err = dbo.preparedStatements[DbPsid_CountCreatedUserPreSessionsByEmail].QueryRow(email).Scan(&n)
	if err != nil {
		return CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountStartedUserSessionsByEmail(email string) (n int, err error) {
	err = dbo.preparedStatements[DbPsid_CountStartedUserSessionsByEmail].QueryRow(email).Scan(&n)
	if err != nil {
		return CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountUserSessionsByUserId(userId uint) (n int, err error) {
	err = dbo.preparedStatements[DbPsid_CountUserSessionsByUserId].QueryRow(userId).Scan(&n)
	if err != nil {
		return CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountUsersWithEmailAbleToLogIn(email string) (n int, err error) {
	err = dbo.preparedStatements[DbPsid_CountUsersWithEmailAbleToLogIn].QueryRow(email).Scan(&n)
	if err != nil {
		return CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountUsersWithEmailAddr(email string) (n int, err error) {
	err = dbo.preparedStatements[DbPsid_CountUsersWithEmailAddr].QueryRow(email, email).Scan(&n)
	if err != nil {
		return CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountUsersWithName(name string) (n int, err error) {
	err = dbo.preparedStatements[DbPsid_CountUsersWithName].QueryRow(name, name).Scan(&n)
	if err != nil {
		return CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CreatePreliminarySession(userId uint, requestId string, userIPAB net.IP, pwdSalt []byte, isCaptchaRequired bool, captchaId sql.NullString) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_CreatePreliminarySession].Exec(userId, requestId, userIPAB, pwdSalt, isCaptchaRequired, captchaId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) CreateUserSession(userId uint, userIPAB net.IP) (lastInsertedId int64, err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_CreateUserSession].Exec(userId, userIPAB)
	if err != nil {
		return LastInsertedIdOnError, err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return LastInsertedIdOnError, err
	}

	if ra != 1 {
		return LastInsertedIdOnError, fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	lastInsertedId, err = result.LastInsertId()
	if err != nil {
		return LastInsertedIdOnError, err
	}

	return lastInsertedId, nil
}

func (dbo *DatabaseObject) DeleteAbandonedPreSessions() (err error) {
	timeBorder := time.Now().Add(-time.Duration(dbo.sp.PreSessionExpirationTime) * time.Second)

	_, err = dbo.preparedStatements[DbPsid_DeleteAbandonedPreSessions].Exec(timeBorder)
	if err != nil {
		return err
	}

	return nil
}

func (dbo *DatabaseObject) DeletePreRegUserIfNotApprovedByEmail(email string) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_DeletePreRegUserIfNotApprovedByEmail].Exec(email)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) DeletePreliminarySessionByRequestId(requestId string) (err error) {
	_, err = dbo.preparedStatements[DbPsid_DeletePreliminarySessionByRequestId].Exec(requestId)
	if err != nil {
		return err
	}

	return nil
}

func (dbo *DatabaseObject) DeleteUserSession(sessionId uint, userId uint, userIPAB net.IP) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_DeleteUserSession].Exec(sessionId, userId, userIPAB)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) DeleteUserSessionByUserId(userId uint) (err error) {
	_, err = dbo.preparedStatements[DbPsid_DeleteUserSessionByUserId].Exec(userId)
	if err != nil {
		return err
	}

	return nil
}

func (dbo *DatabaseObject) GetListOfLoggedUsers() (userIds []uint, err error) {
	userIds = make([]uint, 0)

	var rows *sql.Rows
	rows, err = dbo.preparedStatements[DbPsid_GetListOfLoggedUsers].Query()
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

func (dbo *DatabaseObject) GetSessionByUserId(userId uint) (session *am.Session, err error) {
	session = &am.Session{}
	err = dbo.preparedStatements[DbPsid_GetSessionByUserId].QueryRow(userId).Scan(
		&session.Id,
		&session.UserId,
		&session.StartTime,
		&session.UserIPAB,
	)
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
	user = &am.User{}
	err = dbo.preparedStatements[DbPsid_GetUserById].QueryRow(userId).Scan(
		&user.Id,
		&user.PreRegTime,
		&user.Email,
		&user.Name,
		&user.ApprovalTime,
		&user.RegTime,
		&user.IsAuthor,
		&user.IsWriter,
		&user.IsReader,
		&user.CanLogIn,
		&user.LastBadLogInTime,
		&user.BanTime,
	)
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
	err = dbo.preparedStatements[DbPsid_GetUserIdByEmail].QueryRow(email).Scan(&id)
	if err != nil {
		return IdOnError, err
	}

	return id, nil
}

func (dbo *DatabaseObject) GetUserLastBadLogInTimeByEmail(email string) (lastBadLogInTime *time.Time, err error) {
	var u = &am.User{}
	err = dbo.preparedStatements[DbPsid_GetUserLastBadLogInTimeByEmail].QueryRow(email).Scan(&u.LastBadLogInTime)
	if err != nil {
		return nil, err
	}

	return u.LastBadLogInTime, nil
}

func (dbo *DatabaseObject) GetUserPasswordById(userId uint) (password []byte, err error) {
	err = dbo.preparedStatements[DbPsid_GetUserPasswordById].QueryRow(userId).Scan(&password)
	if err != nil {
		return nil, err
	}

	return password, nil
}

func (dbo *DatabaseObject) GetUserPreSessionByRequestId(requestId string) (preSession *am.PreSession, err error) {
	preSession = &am.PreSession{}
	err = dbo.preparedStatements[DbPsid_GetUserPreSessionByRequestId].QueryRow(requestId).Scan(
		&preSession.Id,
		&preSession.UserId,
		&preSession.RequestId,
		&preSession.UserIPAB,
		&preSession.AuthDataBytes,
		&preSession.IsCaptchaRequired,
		&preSession.CaptchaId,
	)
	if err != nil {
		return nil, err
	}

	return preSession, nil
}

func (dbo *DatabaseObject) GetUserRolesById(userId uint) (roles *cm.UserRoles, err error) {
	roles = &cm.UserRoles{}
	err = dbo.preparedStatements[DbPsid_GetUserRolesById].QueryRow(userId).Scan(
		&roles.IsAuthor,
		&roles.IsWriter,
		&roles.IsReader,
		&roles.CanLogIn,
	)
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
	result, err = dbo.preparedStatements[DbPsid_InsertPreRegisteredUser].Exec(email)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) RegisterPreRegUser(email string) (err error) {
	// Part 1.
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_RegisterPreRegUserP1].Exec(email)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	// Part 2.
	result, err = dbo.preparedStatements[DbPsid_RegisterPreRegUserP2].Exec(email)
	if err != nil {
		return err
	}

	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SaveIncident(incidentType am.IncidentType, email string, userIPAB net.IP) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SaveIncident].Exec(incidentType, email, userIPAB)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SaveIncidentWithoutUserIPA(incidentType am.IncidentType, email string) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SaveIncidentWithoutUserIPA].Exec(incidentType, email)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetPreRegUserData(email string, code string, name string, password []byte) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SetPreRegUserData].Exec(name, password, email, code)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetUserPreSessionCaptchaFlags(userId uint, requestId string, isVerifiedByCaptcha bool) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SetUserPreSessionCaptchaFlag].Exec(isVerifiedByCaptcha, requestId, userId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetUserPreSessionPasswordFlag(userId uint, requestId string, isVerifiedByPassword bool) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SetUserPreSessionPasswordFlag].Exec(isVerifiedByPassword, requestId, userId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetUserPreSessionVerificationFlag(userId uint, requestId string, isVerifiedByEmail bool) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SetUserPreSessionVerificationFlag].Exec(isVerifiedByEmail, requestId, userId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetUserRoleAuthor(userId uint, isRoleEnabled bool) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SetUserRoleAuthor].Exec(isRoleEnabled, userId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetUserRoleCanLogIn(userId uint, isRoleEnabled bool) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SetUserRoleCanLogIn].Exec(isRoleEnabled, userId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetUserRoleReader(userId uint, isRoleEnabled bool) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SetUserRoleReader].Exec(isRoleEnabled, userId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetUserRoleWriter(userId uint, isRoleEnabled bool) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SetUserRoleWriter].Exec(isRoleEnabled, userId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) UpdatePreSessionRequestId(userId uint, requestIdOld string, requestIdNew string) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_UpdatePreSessionRequestId].Exec(requestIdNew, requestIdOld, userId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) UpdateUserBanTime(userId uint) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_UpdateUserBanTime].Exec(userId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) UpdateUserLastBadLogInTimeByEmail(email string) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_UpdateUserLastBadLogInTimeByEmail].Exec(email)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrfRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) ViewUserParametersById(userId uint) (userParameters *cm.UserParameters, err error) {
	userParameters = &cm.UserParameters{}
	err = dbo.preparedStatements[DbPsid_GetUserParametersById].QueryRow(userId).Scan(
		&userParameters.Id,
		&userParameters.PreRegTime,
		&userParameters.Email,
		&userParameters.Name,
		&userParameters.ApprovalTime,
		&userParameters.RegTime,
		&userParameters.IsAuthor,
		&userParameters.IsWriter,
		&userParameters.IsReader,
		&userParameters.CanLogIn,
		&userParameters.LastBadLogInTime,
		&userParameters.BanTime,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return userParameters, nil
}
