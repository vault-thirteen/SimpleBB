package acm

import (
	"context"
	"database/sql"
	"time"

	"github.com/vault-thirteen/errorz"
)

func (srv *Server) clearPreRegUsersTableM() (err error) {
	timeBorder := time.Now().Add(-time.Duration(srv.settings.SystemSettings.PreRegUserExpirationTime) * time.Second)

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

	_, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_ClearPreRegUsersTable]).Exec(timeBorder)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) clearSessionsM() (err error) {
	timeBorder := time.Now().Add(-time.Duration(srv.settings.SystemSettings.SessionMaxDuration) * time.Second)

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

	_, err = tx.Stmt(srv.dbPreparedStatements[DbPsid_ClearSessions]).Exec(timeBorder)
	if err != nil {
		return err
	}

	return nil
}
