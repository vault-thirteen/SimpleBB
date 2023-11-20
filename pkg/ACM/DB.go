package acm

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	m "github.com/vault-thirteen/SQL/mysql"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	"github.com/vault-thirteen/errorz"
)

const (
	SqlScriptFileExt                     = "sql"
	FileExtensionSeparator               = "."
	TableNamePrefixSeparator             = "_"
	SearchPattern_CreateTableIfNotExists = "CREATE TABLE IF NOT EXISTS %s"
	LastInsertedIdOnError                = -1
	CountOnError                         = -1
	IdOnError                            = 0
)

const (
	ErrTableNameIsNotFound = "table name is not found"
	ErrfRowsAffectedCount  = "affected rows count error: %v vs %v"
)

func (srv *Server) initDbM() (err error) {
	fmt.Print(c.MsgConnectingToDatabase)

	err = srv.connectDb()
	if err != nil {
		return err
	}

	err = srv.initDbTables()
	if err != nil {
		return err
	}

	err = srv.prepareDbStatements()
	if err != nil {
		return err
	}

	fmt.Println(c.MsgOK)

	return nil
}

func (srv *Server) connectDb() (err error) {
	mc := mysql.Config{
		Net:                  srv.settings.DbSettings.Net,
		Addr:                 net.JoinHostPort(srv.settings.DbSettings.Host, strconv.FormatUint(uint64(srv.settings.DbSettings.Port), 10)),
		DBName:               srv.settings.DbSettings.DBName,
		User:                 srv.settings.DbSettings.User,
		Passwd:               srv.settings.DbSettings.Password,
		AllowNativePasswords: srv.settings.DbSettings.AllowNativePasswords,
		CheckConnLiveness:    srv.settings.DbSettings.CheckConnLiveness,
		MaxAllowedPacket:     srv.settings.DbSettings.MaxAllowedPacket,
		Params:               srv.settings.DbSettings.Params,
	}

	srv.db, err = sql.Open(srv.settings.DbSettings.DriverName, mc.FormatDSN())
	if err != nil {
		return err
	}

	err = srv.db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) initDbTables() (err error) {
	var tableName string
	var tableExists bool

	// Create those tables which are not found.
	for _, tableNameOriginal := range srv.settings.TablesToInit {
		tableName = srv.prefixTableName(tableNameOriginal)

		tableExists, err = m.TableExists(srv.db, srv.settings.DbSettings.DBName, tableName)
		if err != nil {
			return err
		}

		if !tableExists {
			log.Println(fmt.Sprintf(c.MsgFTableIsNotFound, tableName))

			err = srv.initTable(tableName, tableNameOriginal)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (srv *Server) prefixTableName(tableName string) (tableNameFull string) {
	if len(srv.settings.SystemSettings.TableNamePrefix) > 0 {
		return srv.settings.SystemSettings.TableNamePrefix + TableNamePrefixSeparator + tableName
	}

	return tableName
}

// initTable runs initialisation scripts for tables which require
// initialisation as per configuration.
// tableName is a prefixed table name.
// tableNameOriginal is an original table name.
func (srv *Server) initTable(tableName string, tableNameOriginal string) (err error) {
	log.Println(fmt.Sprintf(c.MsgFInitializingDatabaseTable, tableName))

	sqlScriptFilePath := filepath.Join(srv.settings.SystemSettings.TableInitScriptsFolder, tableNameOriginal+FileExtensionSeparator+SqlScriptFileExt)

	var buf []byte
	buf, err = os.ReadFile(sqlScriptFilePath)
	if err != nil {
		return err
	}

	var cmd string
	cmd, err = srv.replaceTableNameInCreateTableScript(buf, tableName, tableNameOriginal)
	if err != nil {
		return err
	}

	_, err = srv.db.Exec(cmd)
	if err != nil {
		return err
	}

	return nil
}

// replaceTableNameInCreateTableScript replaces a table name in the SQL script
// which creates a table.
// scriptText is contents of an SQL script file.
// tableName is a prefixed table name.
// tableNameOriginal is an original table name.
func (srv *Server) replaceTableNameInCreateTableScript(scriptText []byte, tableName string, tableNameOriginal string) (cmd string, err error) {
	pattern := fmt.Sprintf(SearchPattern_CreateTableIfNotExists, tableNameOriginal)
	replacement := fmt.Sprintf(SearchPattern_CreateTableIfNotExists, tableName)
	scriptTextStr := string(scriptText)

	if strings.Index(scriptTextStr, pattern) < 0 {
		return "", errors.New(ErrTableNameIsNotFound)
	}

	return strings.Replace(scriptTextStr, pattern, replacement, 1), nil
}

func (srv *Server) probeDbM() (err error) {
	err = srv.db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) countUsersWithEmailAddrM(email string) (n int, err error) {
	err = srv.dbPreparedStatements[DbPsid_CountUsersWithEmailAddr].QueryRow(email, email).Scan(&n)
	if err != nil {
		return CountOnError, err
	}

	return n, nil
}

func (srv *Server) insertPreRegisteredUserM(email string) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_InsertPreRegisteredUser].Exec(email)
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

func (srv *Server) attachVerificationCodeToPreRegUserM(email string, code string) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_AttachVerificationCodeToPreRegUser].Exec(code, email)
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

func (srv *Server) checkRegVerificationCodeM(email, code string) (ok bool, err error) {
	var n int
	err = srv.dbPreparedStatements[DbPsid_CheckVerificationCodeOfPreRegUser].QueryRow(email, code).Scan(&n)
	if err != nil {
		return false, err
	}

	if n != 1 {
		return false, nil
	}

	return true, nil
}

func (srv *Server) countUsersWithNameM(name string) (n int, err error) {
	err = srv.dbPreparedStatements[DbPsid_CountUsersWithName].QueryRow(name, name).Scan(&n)
	if err != nil {
		return CountOnError, err
	}

	return n, nil
}

func (srv *Server) deletePreRegUserIfNotApprovedByEmailM(email string) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_DeletePreRegUserIfNotApprovedByEmail].Exec(email)
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

func (srv *Server) approveUserByEmailM(email string) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_ApprovePreRegUserEmail].Exec(email)
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

func (srv *Server) setPreRegUserDataM(email string, code string, name string, password []byte) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_SetPreRegUserData].Exec(name, password, email, code)
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

func (srv *Server) approvePreRegUserM(email string) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_ApprovePreRegUser].Exec(email)
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

func (srv *Server) registerPreRegUserM(email string) (err error) {
	// Part 1.
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_RegisterPreRegUserP1].Exec(email)
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
	result, err = srv.dbPreparedStatements[DbPsid_RegisterPreRegUserP2].Exec(email)
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

func (srv *Server) countUsersWithEmailAbleToLogInM(email string) (n int, err error) {
	err = srv.dbPreparedStatements[DbPsid_CountUsersWithEmailAbleToLogIn].QueryRow(email).Scan(&n)
	if err != nil {
		return CountOnError, err
	}

	return n, nil
}

func (srv *Server) deleteAbandonedPreSessionsM() (err error) {
	timeBorder := time.Now().Add(-time.Duration(srv.settings.SystemSettings.PreSessionExpirationTime) * time.Second)

	_, err = srv.dbPreparedStatements[DbPsid_DeleteAbandonedPreSessions].Exec(timeBorder)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) countStartedUserSessionsByEmailM(email string) (n int, err error) {
	err = srv.dbPreparedStatements[DbPsid_CountStartedUserSessionsByEmail].QueryRow(email).Scan(&n)
	if err != nil {
		return CountOnError, err
	}

	return n, nil
}

func (srv *Server) countCreatedUserPreSessionsByEmailM(email string) (n int, err error) {
	err = srv.dbPreparedStatements[DbPsid_CountCreatedUserPreSessionsByEmail].QueryRow(email).Scan(&n)
	if err != nil {
		return CountOnError, err
	}

	return n, nil
}

func (srv *Server) getUserLastBadLogInTimeByEmailM(email string) (lastBadLogInTime *time.Time, err error) {
	var u = &am.User{}
	err = srv.dbPreparedStatements[DbPsid_GetUserLastBadLogInTimeByEmail].QueryRow(email).Scan(&u.LastBadLogInTime)
	if err != nil {
		return nil, err
	}

	return u.LastBadLogInTime, nil
}

func (srv *Server) createPreliminarySessionM(userId uint, requestId string, userIPAB net.IP, pwdSalt []byte, isCaptchaRequired bool, captchaId sql.NullString) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_CreatePreliminarySession].Exec(userId, requestId, userIPAB, pwdSalt, isCaptchaRequired, captchaId)
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

func (srv *Server) getUserIdByEmailM(email string) (id uint, err error) {
	err = srv.dbPreparedStatements[DbPsid_GetUserIdByEmail].QueryRow(email).Scan(&id)
	if err != nil {
		return IdOnError, err
	}

	return id, nil
}

func (srv *Server) updateUserLastBadLogInTimeByEmailM(email string) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_UpdateUserLastBadLogInTimeByEmail].Exec(email)
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

func (srv *Server) getUserPreSessionByRequestIdM(requestId string) (preSession *am.PreSession, err error) {
	preSession = &am.PreSession{}
	err = srv.dbPreparedStatements[DbPsid_GetUserPreSessionByRequestId].QueryRow(requestId).Scan(
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

func (srv *Server) getUserPasswordByIdM(userId uint) (password []byte, err error) {
	err = srv.dbPreparedStatements[DbPsid_GetUserPasswordById].QueryRow(userId).Scan(&password)
	if err != nil {
		return nil, err
	}

	return password, nil
}

func (srv *Server) deletePreliminarySessionByRequestIdM(requestId string) (err error) {
	_, err = srv.dbPreparedStatements[DbPsid_DeletePreliminarySessionByRequestId].Exec(requestId)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) setUserPreSessionCaptchaFlagsM(userId uint, requestId string, isVerifiedByCaptcha bool) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_SetUserPreSessionCaptchaFlag].Exec(isVerifiedByCaptcha, requestId, userId)
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

func (srv *Server) setUserPreSessionPasswordFlagM(userId uint, requestId string, isVerifiedByPassword bool) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_SetUserPreSessionPasswordFlag].Exec(isVerifiedByPassword, requestId, userId)
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

func (srv *Server) attachVerificationCodeToPreSessionM(userId uint, requestId string, code string) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_AttachVerificationCodeToPreSession].Exec(code, requestId, userId)
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

func (srv *Server) updatePreSessionRequestIdM(userId uint, requestIdOld string, requestIdNew string) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_UpdatePreSessionRequestId].Exec(requestIdNew, requestIdOld, userId)
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

func (srv *Server) checkLogInVerificationCodeM(requestId string, code string) (ok bool, err error) {
	var n int
	err = srv.dbPreparedStatements[DbPsid_CheckVerificationCodeOfLogInUser].QueryRow(requestId, code).Scan(&n)
	if err != nil {
		return false, err
	}

	if n != 1 {
		return false, nil
	}

	return true, nil
}

func (srv *Server) setUserPreSessionVerificationFlagM(userId uint, requestId string, isVerifiedByEmail bool) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_SetUserPreSessionVerificationFlag].Exec(isVerifiedByEmail, requestId, userId)
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

func (srv *Server) createUserSessionM(userId uint, userIPAB net.IP) (lastInsertedId int64, err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_CreateUserSession].Exec(userId, userIPAB)
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

func (srv *Server) getUserByIdM(userId uint) (user *am.User, err error) {
	user = &am.User{}
	err = srv.dbPreparedStatements[DbPsid_GetUserById].QueryRow(userId).Scan(
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

func (srv *Server) getSessionByUserIdM(userId uint) (session *am.Session, err error) {
	session = &am.Session{}
	err = srv.dbPreparedStatements[DbPsid_GetSessionByUserId].QueryRow(userId).Scan(
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

func (srv *Server) deleteUserSessionM(sessionId uint, userId uint, userIPAB net.IP) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_DeleteUserSession].Exec(sessionId, userId, userIPAB)
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

func (srv *Server) saveIncidentM(incidentType IncidentType, email string, userIPAB net.IP) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_SaveIncident].Exec(incidentType, email, userIPAB)
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

func (srv *Server) saveIncidentWithoutUserIPAM(incidentType IncidentType, email string) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_SaveIncidentWithoutUserIPA].Exec(incidentType, email)
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

func (srv *Server) getListOfLoggedUsersM() (userIds []uint, err error) {
	userIds = make([]uint, 0)

	var rows *sql.Rows
	rows, err = srv.dbPreparedStatements[DbPsid_GetListOfLoggedUsers].Query()
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

func (srv *Server) countUserSessionsByUserIdM(userId uint) (n int, err error) {
	err = srv.dbPreparedStatements[DbPsid_CountUserSessionsByUserId].QueryRow(userId).Scan(&n)
	if err != nil {
		return CountOnError, err
	}

	return n, nil
}

func (srv *Server) getUserRolesByIdM(userId uint) (roles *cm.UserRoles, err error) {
	roles = &cm.UserRoles{}
	err = srv.dbPreparedStatements[DbPsid_GetUserRolesById].QueryRow(userId).Scan(
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

func (srv *Server) viewUserParametersByIdM(userId uint) (userParameters *cm.UserParameters, err error) {
	userParameters = &cm.UserParameters{}
	err = srv.dbPreparedStatements[DbPsid_GetUserParametersById].QueryRow(userId).Scan(
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

func (srv *Server) setUserRoleAuthorM(userId uint, isRoleEnabled bool) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_SetUserRoleAuthor].Exec(isRoleEnabled, userId)
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

func (srv *Server) setUserRoleWriterM(userId uint, isRoleEnabled bool) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_SetUserRoleWriter].Exec(isRoleEnabled, userId)
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

func (srv *Server) setUserRoleReaderM(userId uint, isRoleEnabled bool) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_SetUserRoleReader].Exec(isRoleEnabled, userId)
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

func (srv *Server) setUserRoleCanLogInM(userId uint, isRoleEnabled bool) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_SetUserRoleCanLogIn].Exec(isRoleEnabled, userId)
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

func (srv *Server) deleteUserSessionByUserIdM(userId uint) (err error) {
	_, err = srv.dbPreparedStatements[DbPsid_DeleteUserSessionByUserId].Exec(userId)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) updateUserBanTimeM(userId uint) (err error) {
	var result sql.Result
	result, err = srv.dbPreparedStatements[DbPsid_UpdateUserBanTime].Exec(userId)
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
