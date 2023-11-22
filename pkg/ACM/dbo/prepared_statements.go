package dbo

import (
	"database/sql"
	"fmt"
)

// Indices of prepared statements.
const (
	DbPsid_CountUsersWithEmailAddr              = 0
	DbPsid_InsertPreRegisteredUser              = 1
	DbPsid_AttachVerificationCodeToPreRegUser   = 2
	DbPsid_CheckVerificationCodeOfPreRegUser    = 3
	DbPsid_DeletePreRegUserIfNotApprovedByEmail = 4
	DbPsid_ApprovePreRegUserEmail               = 5
	DbPsid_SetPreRegUserData                    = 6
	DbPsid_ApprovePreRegUser                    = 7
	DbPsid_RegisterPreRegUserP1                 = 8
	DbPsid_RegisterPreRegUserP2                 = 9
	DbPsid_CountUsersWithName                   = 10
	DbPsid_ClearPreRegUsersTable                = 11
	DbPsid_CountUsersWithEmailAbleToLogIn       = 12
	DbPsid_DeleteAbandonedPreSessions           = 13
	DbPsid_CountStartedUserSessionsByEmail      = 14
	DbPsid_CountCreatedUserPreSessionsByEmail   = 15
	DbPsid_GetUserLastBadLogInTimeByEmail       = 16
	DbPsid_CreatePreliminarySession             = 17
	DbPsid_GetUserIdByEmail                     = 18
	DbPsid_UpdateUserLastBadLogInTimeByEmail    = 19
	DbPsid_GetUserPreSessionByRequestId         = 20
	DbPsid_GetUserPasswordById                  = 21
	DbPsid_DeletePreliminarySessionByRequestId  = 22
	DbPsid_SetUserPreSessionCaptchaFlag         = 23
	DbPsid_SetUserPreSessionPasswordFlag        = 24
	DbPsid_AttachVerificationCodeToPreSession   = 25
	DbPsid_UpdatePreSessionRequestId            = 26
	DbPsid_CheckVerificationCodeOfLogInUser     = 27
	DbPsid_SetUserPreSessionVerificationFlag    = 28
	DbPsid_CreateUserSession                    = 29
	DbPsid_ClearSessions                        = 30
	DbPsid_GetUserById                          = 31
	DbPsid_GetSessionByUserId                   = 32
	DbPsid_DeleteUserSession                    = 33
	DbPsid_SaveIncident                         = 34
	DbPsid_SaveIncidentWithoutUserIPA           = 35
	DbPsid_GetListOfLoggedUsers                 = 36
	DbPsid_CountUserSessionsByUserId            = 37
	DbPsid_GetUserRolesById                     = 38
	DbPsid_GetUserParametersById                = 39
	DbPsid_SetUserRoleAuthor                    = 40
	DbPsid_SetUserRoleWriter                    = 41
	DbPsid_SetUserRoleReader                    = 42
	DbPsid_SetUserRoleCanLogIn                  = 43
	DbPsid_DeleteUserSessionByUserId            = 44
	DbPsid_UpdateUserBanTime                    = 45
)

func (dbo *DatabaseObject) prepareStatements() (err error) {
	// Add prepared statements into an array (slice) for future usage.
	tblPreRegisteredUsers := dbo.prefixTableName("PreRegisteredUsers")
	tblUsers := dbo.prefixTableName("Users")
	tblPreSessions := dbo.prefixTableName("PreSessions")
	tblSessions := dbo.prefixTableName("Sessions")
	tblIncidents := dbo.prefixTableName("Incidents")

	var q string
	qs := make([]string, 0)

	// 0.
	q = fmt.Sprintf(`SELECT (SELECT COUNT(Email) FROM %s WHERE Email = ?) + (SELECT COUNT(Email) FROM %s WHERE Email = ?) AS n;`, tblUsers, tblPreRegisteredUsers)
	qs = append(qs, q)

	// 1.
	q = fmt.Sprintf(`INSERT INTO %s (Email) VALUES (?);`, tblPreRegisteredUsers)
	qs = append(qs, q)

	// 2.
	q = fmt.Sprintf(`UPDATE %s SET VerificationCode = ? WHERE Email = ?;`, tblPreRegisteredUsers)
	qs = append(qs, q)

	// 3.
	q = fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE Email = ? AND VerificationCode = ?;`, tblPreRegisteredUsers)
	qs = append(qs, q)

	// 4.
	q = fmt.Sprintf(`DELETE FROM %s WHERE Email = ? AND IsEmailApproved = FALSE;`, tblPreRegisteredUsers)
	qs = append(qs, q)

	// 5.
	q = fmt.Sprintf(`UPDATE %s SET IsEmailApproved = TRUE WHERE Email = ?;`, tblPreRegisteredUsers)
	qs = append(qs, q)

	// 6.
	q = fmt.Sprintf(`UPDATE %s SET NAME = ?, PASSWORD = ?, IsReadyForApproval = TRUE WHERE Email = ? AND VerificationCode = ? AND IsReadyForApproval = FALSE;`, tblPreRegisteredUsers)
	qs = append(qs, q)

	// 7.
	q = fmt.Sprintf(`UPDATE %s SET IsApproved = TRUE, ApprovalTime = Now() WHERE Email = ? AND IsEmailApproved = TRUE AND IsReadyForApproval = TRUE;`, tblPreRegisteredUsers)
	qs = append(qs, q)

	// 8.
	q = fmt.Sprintf(`INSERT INTO %s (PreRegTime, Email, NAME, PASSWORD, ApprovalTime, RegTime, IsReader, CanLogIn) SELECT PreRegTime, Email, NAME, PASSWORD, ApprovalTime, Now(), TRUE, TRUE FROM %s AS pru WHERE pru.Email = ? AND pru.IsApproved = TRUE;`, tblUsers, tblPreRegisteredUsers)
	qs = append(qs, q)

	// 9.
	q = fmt.Sprintf(`DELETE FROM %s WHERE Email = ? AND IsApproved = TRUE;`, tblPreRegisteredUsers)
	qs = append(qs, q)

	// 10.
	q = fmt.Sprintf(`SELECT (SELECT COUNT(Name) FROM %s WHERE NAME = ?) + (SELECT COUNT(Name) FROM %s WHERE NAME = ?) AS n;`, tblUsers, tblPreRegisteredUsers)
	qs = append(qs, q)

	// 11.
	q = fmt.Sprintf(`DELETE FROM %s WHERE IsReadyForApproval = FALSE AND PreRegTime < ?;`, tblPreRegisteredUsers)
	qs = append(qs, q)

	// 12.
	q = fmt.Sprintf(`SELECT COUNT(Email) AS n FROM %s WHERE Email = ? AND CanLogIn = TRUE;`, tblUsers)
	qs = append(qs, q)

	// 13.
	q = fmt.Sprintf(`DELETE FROM %s WHERE TimeOfCreation < ?;`, tblPreSessions)
	qs = append(qs, q)

	// 14.
	q = fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE UserId = (SELECT Id FROM %s WHERE Email = ?);`, tblSessions, tblUsers)
	qs = append(qs, q)

	// 15.
	q = fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE UserId = (SELECT Id FROM %s WHERE Email = ?);`, tblPreSessions, tblUsers)
	qs = append(qs, q)

	// 16.
	q = fmt.Sprintf(`SELECT LastBadLogInTime FROM %s WHERE Email = ?;`, tblUsers)
	qs = append(qs, q)

	// 17.
	q = fmt.Sprintf(`INSERT INTO %s (UserId, RequestId, UserIPAB, AuthDataBytes, IsCaptchaRequired, CaptchaId) VALUES (?, ?, ?, ?, ?, ?);`, tblPreSessions)
	qs = append(qs, q)

	// 18.
	q = fmt.Sprintf(`SELECT Id FROM %s WHERE Email = ?;`, tblUsers)
	qs = append(qs, q)

	// 19.
	q = fmt.Sprintf(`UPDATE %s SET LastBadLogInTime = Now() WHERE Email = ?;`, tblUsers)
	qs = append(qs, q)

	// 20.
	q = fmt.Sprintf(`SELECT Id, UserId, RequestId, UserIPAB, AuthDataBytes, IsCaptchaRequired, CaptchaId FROM %s WHERE RequestId = ?;`, tblPreSessions)
	qs = append(qs, q)

	// 21.
	q = fmt.Sprintf(`SELECT Password FROM %s WHERE Id = ?;`, tblUsers)
	qs = append(qs, q)

	// 22.
	q = fmt.Sprintf(`DELETE FROM %s WHERE RequestId = ?;`, tblPreSessions)
	qs = append(qs, q)

	// 23.
	q = fmt.Sprintf(`UPDATE %s SET IsVerifiedByCaptcha = ? WHERE RequestId = ? AND UserId = ?;`, tblPreSessions)
	qs = append(qs, q)

	// 24.
	q = fmt.Sprintf(`UPDATE %s SET IsVerifiedByPassword = ? WHERE RequestId = ? AND UserId = ?;`, tblPreSessions)
	qs = append(qs, q)

	// 25.
	q = fmt.Sprintf(`UPDATE %s SET VerificationCode = ? WHERE RequestId = ? AND UserId = ?;`, tblPreSessions)
	qs = append(qs, q)

	// 26.
	q = fmt.Sprintf(`UPDATE %s SET RequestId = ? WHERE RequestId = ? AND UserId = ?;`, tblPreSessions)
	qs = append(qs, q)

	// 27.
	q = fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE RequestId = ? AND VerificationCode = ?;`, tblPreSessions)
	qs = append(qs, q)

	// 28.
	q = fmt.Sprintf(`UPDATE %s SET IsVerifiedByEmail = ? WHERE RequestId = ? AND UserId = ?;`, tblPreSessions)
	qs = append(qs, q)

	// 29.
	q = fmt.Sprintf(`INSERT INTO %s (UserId, UserIPAB) VALUES (?, ?);`, tblSessions)
	qs = append(qs, q)

	// 30.
	q = fmt.Sprintf(`DELETE FROM %s WHERE StartTime < ?;`, tblSessions)
	qs = append(qs, q)

	// 31.
	q = fmt.Sprintf(`SELECT Id, PreRegTime, Email, Name, ApprovalTime, RegTime, IsAuthor, IsWriter, IsReader, CanLogIn, LastBadLogInTime, BanTime FROM %s WHERE Id = ?;`, tblUsers)
	qs = append(qs, q)

	// 32.
	q = fmt.Sprintf(`SELECT Id, UserId, StartTime, UserIPAB FROM %s WHERE UserId = ?;`, tblSessions)
	qs = append(qs, q)

	// 33.
	q = fmt.Sprintf(`DELETE FROM %s WHERE Id = ? AND UserId = ? AND UserIPAB = ?;`, tblSessions)
	qs = append(qs, q)

	// 34.
	q = fmt.Sprintf(`INSERT INTO %s (Type, Email, UserIPAB) VALUES (?, ?, ?);`, tblIncidents)
	qs = append(qs, q)

	// 35.
	q = fmt.Sprintf(`INSERT INTO %s (Type, Email) VALUES (?, ?);`, tblIncidents)
	qs = append(qs, q)

	// 36.
	q = fmt.Sprintf(`SELECT UserId FROM %s;`, tblSessions)
	qs = append(qs, q)

	// 37.
	q = fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE UserId = ?;`, tblSessions)
	qs = append(qs, q)

	// 38.
	q = fmt.Sprintf(`SELECT IsAuthor, IsWriter, IsReader, CanLogIn FROM %s WHERE Id = ?;`, tblUsers)
	qs = append(qs, q)

	// 39.
	q = fmt.Sprintf(`SELECT Id, PreRegTime, Email, Name, ApprovalTime, RegTime, IsAuthor, IsWriter, IsReader, CanLogIn, LastBadLogInTime, BanTime FROM %s WHERE Id = ?;`, tblUsers)
	qs = append(qs, q)

	// 40.
	q = fmt.Sprintf(`UPDATE %s SET IsAuthor = ? WHERE Id = ?;`, tblUsers)
	qs = append(qs, q)

	// 41.
	q = fmt.Sprintf(`UPDATE %s SET IsWriter = ? WHERE Id = ?;`, tblUsers)
	qs = append(qs, q)

	// 42.
	q = fmt.Sprintf(`UPDATE %s SET IsReader = ? WHERE Id = ?;`, tblUsers)
	qs = append(qs, q)

	// 43.
	q = fmt.Sprintf(`UPDATE %s SET CanLogIn = ? WHERE Id = ?;`, tblUsers)
	qs = append(qs, q)

	// 44.
	q = fmt.Sprintf(`DELETE FROM %s WHERE UserId = ?;`, tblSessions)
	qs = append(qs, q)

	// 45.
	q = fmt.Sprintf(`UPDATE %s SET BanTime = Now() WHERE Id = ?;`, tblUsers)
	qs = append(qs, q)

	dbo.preparedStatementQueries = make([]string, 0, len(qs))
	for _, q = range qs {
		dbo.preparedStatementQueries = append(dbo.preparedStatementQueries, q)
	}
	qs = nil

	dbo.preparedStatements = make([]*sql.Stmt, 0, len(dbo.preparedStatementQueries))
	var st *sql.Stmt
	for _, psq := range dbo.preparedStatementQueries {
		st, err = dbo.db.Prepare(psq)
		if err != nil {
			return err
		}

		dbo.preparedStatements = append(dbo.preparedStatements, st)
	}

	return nil
}

func (dbo *DatabaseObject) GetPreparedStatementByIndex(i int) (ps *sql.Stmt) {
	return dbo.preparedStatements[i]
}
