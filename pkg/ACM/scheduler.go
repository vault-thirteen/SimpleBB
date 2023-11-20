package acm

import (
	"time"
)

func (srv *Server) clearPreRegUsersTableM() (err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	timeBorder := time.Now().Add(-time.Duration(srv.settings.SystemSettings.PreRegUserExpirationTime) * time.Second)

	_, err = srv.dbPreparedStatements[DbPsid_ClearPreRegUsersTable].Exec(timeBorder)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) clearSessionsM() (err error) {
	srv.dbGuard.Lock()
	defer srv.dbGuard.Unlock()

	timeBorder := time.Now().Add(-time.Duration(srv.settings.SystemSettings.SessionMaxDuration) * time.Second)

	_, err = srv.dbPreparedStatements[DbPsid_ClearSessions].Exec(timeBorder)
	if err != nil {
		return err
	}

	return nil
}
