package mm

import (
	"context"
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
	mm "github.com/vault-thirteen/SimpleBB/pkg/MM/models"
	ul "github.com/vault-thirteen/SimpleBB/pkg/UidList"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
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
	TableSections = "Sections"
	TableForums   = "Forums"
	TableThreads  = "Threads"
	TableMessages = "Messages"
)

const (
	ErrTableNameIsNotFound       = "table name is not found"
	ErrfRowsAffectedCount        = "affected rows count error: %v vs %v"
	ErrNoRootForumsAreFound      = "no root forums are found"
	ErrSeveralRootForumsAreFound = "several root forums are found"
)

func (srv *Server) initDbM() (err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	fmt.Print(c.MsgConnectingToDatabase)

	srv.initDBTableNames()

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

func (srv *Server) initDBTableNames() {
	srv.dbTableNames = &DBTableNames{
		Sections: srv.prefixTableName(TableSections),
		Forums:   srv.prefixTableName(TableForums),
		Threads:  srv.prefixTableName(TableThreads),
		Messages: srv.prefixTableName(TableMessages),
	}
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
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	err = srv.db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) countRootSectionsM() (n int, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return CountOnError, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	err = tx.Stmt(srv.dbPreparedStatements[DbPsid_CountRootSections]).QueryRow().Scan(&n)
	if err != nil {
		return CountOnError, err
	}

	return n, nil
}

func (srv *Server) insertNewSectionM(parent *uint, name string, creatorUserId uint) (lastInsertedId int64, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return LastInsertedIdOnError, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var result sql.Result
	result, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_InsertNewSection]).Exec(parent, name, creatorUserId)
	if err != nil {
		return LastInsertedIdOnError, err
	}

	lastInsertedId, err = result.LastInsertId()
	if err != nil {
		return LastInsertedIdOnError, err
	}

	return lastInsertedId, nil
}

func (srv *Server) insertNewForumM(sectionId uint, name string, creatorUserId uint) (lastInsertedId int64, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return LastInsertedIdOnError, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var result sql.Result
	result, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_InsertNewForum]).Exec(sectionId, name, creatorUserId)
	if err != nil {
		return LastInsertedIdOnError, err
	}

	lastInsertedId, err = result.LastInsertId()
	if err != nil {
		return LastInsertedIdOnError, err
	}

	return lastInsertedId, nil
}

func (srv *Server) countSectionsByIdM(sectionId uint) (n int, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return CountOnError, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	err = tx.Stmt(srv.dbPreparedStatements[DbPsid_CountSectionsById]).QueryRow(sectionId).Scan(&n)
	if err != nil {
		return CountOnError, err
	}

	return n, nil
}

func (srv *Server) countForumsByIdM(forumId uint) (n int, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return CountOnError, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	err = tx.Stmt(srv.dbPreparedStatements[DbPsid_CountForumsById]).QueryRow(forumId).Scan(&n)
	if err != nil {
		return CountOnError, err
	}

	return n, nil
}

func (srv *Server) getSectionChildrenByIdM(sectionId uint) (children *ul.UidList, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	children = ul.New()
	err = tx.Stmt(srv.dbPreparedStatements[DbPsid_GetSectionChildrenById]).QueryRow(sectionId).Scan(children)
	if err != nil {
		return nil, err
	}

	return children, nil
}

func (srv *Server) getSectionChildTypeByIdM(sectionId uint) (childType byte, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return 0, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	err = tx.Stmt(srv.dbPreparedStatements[DbPsid_GetSectionChildTypeById]).QueryRow(sectionId).Scan(&childType)
	if err != nil {
		return 0, err
	}

	return childType, nil
}

func (srv *Server) setSectionChildTypeByIdM(sectionId uint, childType byte) (err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var result sql.Result
	result, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_SetSectionChildTypeById]).Exec(childType, sectionId)
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

func (srv *Server) setSectionChildrenByIdM(sectionId uint, children *ul.UidList) (err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var result sql.Result
	result, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_SetSectionChildrenById]).Exec(children, sectionId)
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

func (srv *Server) setSectionNameByIdM(sectionId uint, name string, editorUserId uint) (err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var result sql.Result
	result, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_SetSectionNameById]).Exec(name, editorUserId, sectionId)
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

func (srv *Server) setForumNameByIdM(forumId uint, name string, editorUserId uint) (err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var result sql.Result
	result, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_SetForumNameById]).Exec(name, editorUserId, forumId)
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

func (srv *Server) setSectionParentByIdM(sectionId uint, parent uint, editorUserId uint) (err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var result sql.Result
	result, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_SetSectionParentById]).Exec(parent, editorUserId, sectionId)
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

func (srv *Server) setForumSectionByIdM(forumId uint, sectionId uint, editorUserId uint) (err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var result sql.Result
	result, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_SetForumSectionById]).Exec(sectionId, editorUserId, forumId)
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

func (srv *Server) getSectionParentByIdM(sectionId uint) (parent *uint, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	parent = new(uint)
	err = tx.Stmt(srv.dbPreparedStatements[DbPsid_GetSectionParentById]).QueryRow(sectionId).Scan(&parent)
	if err != nil {
		return nil, err
	}

	return parent, nil
}

func (srv *Server) getForumSectionByIdM(forumId uint) (sectionId uint, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return IdOnError, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	err = tx.Stmt(srv.dbPreparedStatements[DbPsid_GetForumSectionById]).QueryRow(forumId).Scan(&sectionId)
	if err != nil {
		return IdOnError, err
	}

	return sectionId, nil
}

func (srv *Server) insertNewThreadM(parentForum uint, threadName string, creatorUserId uint) (lastInsertedId int64, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return LastInsertedIdOnError, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var result sql.Result
	result, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_InsertNewThread]).Exec(parentForum, threadName, creatorUserId)
	if err != nil {
		return LastInsertedIdOnError, err
	}

	lastInsertedId, err = result.LastInsertId()
	if err != nil {
		return LastInsertedIdOnError, err
	}

	return lastInsertedId, nil
}

func (srv *Server) getForumThreadsByIdM(forumId uint) (threads *ul.UidList, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	threads = ul.New()
	err = tx.Stmt(srv.dbPreparedStatements[DbPsid_GetForumThreadsById]).QueryRow(forumId).Scan(threads)
	if err != nil {
		return nil, err
	}

	return threads, nil
}

func (srv *Server) setForumThreadsByIdM(forumId uint, threads *ul.UidList) (err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var result sql.Result
	result, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_SetForumThreadsById]).Exec(threads, forumId)
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

func (srv *Server) setThreadNameByIdM(threadId uint, name string, editorUserId uint) (err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var result sql.Result
	result, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_SetThreadNameById]).Exec(name, editorUserId, threadId)
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

func (srv *Server) getThreadForumByIdM(threadId uint) (forum uint, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return IdOnError, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	err = tx.Stmt(srv.dbPreparedStatements[DbPsid_GetThreadForumById]).QueryRow(threadId).Scan(&forum)
	if err != nil {
		return IdOnError, err
	}

	return forum, nil
}

func (srv *Server) setThreadForumByIdM(threadId uint, forum uint, editorUserId uint) (err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var result sql.Result
	result, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_SetThreadForumById]).Exec(forum, editorUserId, threadId)
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

func (srv *Server) countThreadsByIdM(threadId uint) (n int, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return CountOnError, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	err = tx.Stmt(srv.dbPreparedStatements[DbPsid_CountThreadsById]).QueryRow(threadId).Scan(&n)
	if err != nil {
		return CountOnError, err
	}

	return n, nil
}

func (srv *Server) getThreadMessagesByIdM(threadId uint) (messages *ul.UidList, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	messages = ul.New()
	err = tx.Stmt(srv.dbPreparedStatements[DbPsid_GetThreadMessagesById]).QueryRow(threadId).Scan(messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (srv *Server) insertNewMessageM(parentThread uint, messageText string, creatorUserId uint) (lastInsertedId int64, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return LastInsertedIdOnError, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var result sql.Result
	result, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_InsertNewMessage]).Exec(parentThread, messageText, creatorUserId)
	if err != nil {
		return LastInsertedIdOnError, err
	}

	lastInsertedId, err = result.LastInsertId()
	if err != nil {
		return LastInsertedIdOnError, err
	}

	return lastInsertedId, nil
}

func (srv *Server) setThreadMessagesByIdM(threadId uint, messages *ul.UidList) (err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var result sql.Result
	result, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_SetThreadMessagesById]).Exec(messages, threadId)
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

func (srv *Server) setMessageTextByIdM(messageId uint, text string, editorUserId uint) (err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var result sql.Result
	result, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_SetMessageTextById]).Exec(text, editorUserId, messageId)
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

func (srv *Server) getMessageThreadByIdM(messageId uint) (thread uint, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return IdOnError, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	err = tx.Stmt(srv.dbPreparedStatements[DbPsid_GetMessageThreadById]).QueryRow(messageId).Scan(&thread)
	if err != nil {
		return IdOnError, err
	}

	return thread, nil
}

func (srv *Server) setMessageThreadByIdM(messageId uint, thread uint, editorUserId uint) (err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var result sql.Result
	result, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_SetMessageThreadById]).Exec(thread, editorUserId, messageId)
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

func (srv *Server) getMessageCreatorAndTimeByIdM(messageId uint) (creatorUserId uint, ToC time.Time, ToE *time.Time, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return IdOnError, time.Time{}, nil, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	// N.B.: ToC can not be null, but ToE can be null !
	err = tx.Stmt(srv.dbPreparedStatements[DbPsid_GetMessageCreatorAndTimeById]).QueryRow(messageId).Scan(&creatorUserId, &ToC, &ToE)
	if err != nil {
		return IdOnError, time.Time{}, nil, err
	}

	return creatorUserId, ToC, ToE, nil
}

func (srv *Server) getMessageByIdM(messageId uint) (message *mm.Message, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	message = mm.NewMessage()
	err = tx.Stmt(srv.dbPreparedStatements[DbPsid_GetMessageById]).QueryRow(messageId).Scan(
		&message.Id,
		&message.ThreadId,
		&message.Text,
		&message.Creator.UserId,
		&message.Creator.Time,
		&message.Editor.UserId,
		&message.Editor.Time,
	)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (srv *Server) deleteMessageByIdM(messageId uint) (err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var result sql.Result
	result, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_DeleteMessageById]).Exec(messageId)
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

func (srv *Server) getThreadByIdM(threadId uint) (thread *mm.Thread, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	thread = mm.NewThread()
	err = tx.Stmt(srv.dbPreparedStatements[DbPsid_GetThreadByIdM]).QueryRow(threadId).Scan(
		&thread.Id,
		&thread.ForumId,
		&thread.Name,
		&thread.Messages,
		&thread.Creator.UserId,
		&thread.Creator.Time,
		&thread.Editor.UserId,
		&thread.Editor.Time,
	)
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (srv *Server) deleteThreadByIdM(threadId uint) (err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var result sql.Result
	result, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_DeleteThreadById]).Exec(threadId)
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

func (srv *Server) getForumByIdM(forumId uint) (forum *mm.Forum, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	forum = mm.NewForum()
	err = tx.Stmt(srv.dbPreparedStatements[DbPsid_GetForumById]).QueryRow(forumId).Scan(
		&forum.Id,
		&forum.SectionId,
		&forum.Name,
		&forum.Threads,
		&forum.Creator.UserId,
		&forum.Creator.Time,
		&forum.Editor.UserId,
		&forum.Editor.Time,
	)
	if err != nil {
		return nil, err
	}

	return forum, nil
}

func (srv *Server) deleteForumByIdM(forumId uint) (err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var result sql.Result
	result, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_DeleteForumById]).Exec(forumId)
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

func (srv *Server) readMessagesByIdM(messageIds ul.UidList) (messages []mm.Message, err error) {
	messages = make([]mm.Message, 0, messageIds.Size())

	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var query string
	query, err = srv.dbQuery_ReadMessagesById(messageIds)
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	rows, err = tx.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		msg := mm.NewMessage()

		err = rows.Scan(
			&msg.Id,
			&msg.ThreadId,
			&msg.Text,
			&msg.Creator.UserId,
			&msg.Creator.Time,
			&msg.Editor.UserId,
			&msg.Editor.Time,
		)
		if err != nil {
			return nil, err
		}

		messages = append(messages, *msg)
	}

	return messages, nil
}

func (srv *Server) readThreadsByIdM(threadIds ul.UidList) (threads []mm.Thread, err error) {
	threads = make([]mm.Thread, 0, threadIds.Size())

	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var query string
	query, err = srv.dbQuery_ReadThreadsById(threadIds)
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	rows, err = tx.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		thr := mm.NewThread()

		err = rows.Scan(
			&thr.Id,
			&thr.ForumId,
			&thr.Name,
			&thr.Messages,
			&thr.Creator.UserId,
			&thr.Creator.Time,
			&thr.Editor.UserId,
			&thr.Editor.Time,
		)
		if err != nil {
			return nil, err
		}

		threads = append(threads, *thr)
	}

	return threads, nil
}

func (srv *Server) readForumsM() (forums []mm.Forum, err error) {
	forums = make([]mm.Forum, 0)

	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var rows *sql.Rows
	rows, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_ReadForums]).Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		forum := mm.NewForum()

		err = rows.Scan(
			&forum.Id,
			&forum.SectionId,
			&forum.Name,
			&forum.Threads,
			&forum.Creator.UserId,
			&forum.Creator.Time,
			&forum.Editor.UserId,
			&forum.Editor.Time,
		)
		if err != nil {
			return nil, err
		}

		forums = append(forums, *forum)
	}

	return forums, nil
}

func (srv *Server) getSectionByIdM(sectionId uint) (section *mm.Section, err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	section = mm.NewSection()
	err = tx.Stmt(srv.dbPreparedStatements[DbPsid_GetSectionById]).QueryRow(sectionId).Scan(
		&section.Id,
		&section.Parent,
		&section.ChildType,
		&section.Children,
		&section.Name,
		&section.Creator.UserId,
		&section.Creator.Time,
		&section.Editor.UserId,
		&section.Editor.Time,
	)
	if err != nil {
		return nil, err
	}

	return section, nil
}

func (srv *Server) deleteSectionByIdM(sectionId uint) (err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var result sql.Result
	result, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_DeleteSectionById]).Exec(sectionId)
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

func (srv *Server) readSectionsM() (sections []mm.Section, err error) {
	sections = make([]mm.Section, 0)

	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	var tx *sql.Tx
	tx, err = srv.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()
		} else {
			txErr = tx.Commit()
		}

		if txErr != nil {
			err = errorz.Combine(err, txErr)
		}
	}()

	var rows *sql.Rows
	rows, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_ReadSections]).Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		section := mm.NewSection()

		err = rows.Scan(
			&section.Id,
			&section.Parent,
			&section.ChildType,
			&section.Children,
			&section.Name,
			&section.Creator.UserId,
			&section.Creator.Time,
			&section.Editor.UserId,
			&section.Editor.Time,
		)
		if err != nil {
			return nil, err
		}

		sections = append(sections, *section)
	}

	return sections, nil
}
