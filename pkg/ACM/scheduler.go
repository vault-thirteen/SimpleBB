package acm

import (
	"time"

	"github.com/vault-thirteen/SimpleBB/pkg/ACM/dbo"
)

func (srv *Server) clearPreRegUsersTableM() (err error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	timeBorder := time.Now().Add(-time.Duration(srv.settings.SystemSettings.PreRegUserExpirationTime) * time.Second)

	_, err = srv.dbo.GetPreparedStatementByIndex(dbo.DbPsid_ClearPreRegUsersTable).Exec(timeBorder)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) clearSessionsM() (err error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	timeBorder := time.Now().Add(-time.Duration(srv.settings.SystemSettings.SessionMaxDuration) * time.Second)

	_, err = srv.dbo.GetPreparedStatementByIndex(dbo.DbPsid_ClearSessions).Exec(timeBorder)
	if err != nil {
		return err
	}

	return nil
}
