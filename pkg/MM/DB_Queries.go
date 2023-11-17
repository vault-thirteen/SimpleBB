package mm

import (
	ul "github.com/vault-thirteen/SimpleBB/pkg/UidList"
)

type DBTableNames struct {
	Sections string
	Forums   string
	Threads  string
	Messages string
}

func (srv *Server) dbQuery_ReadMessagesById(messageIds ul.UidList) (query string, err error) {
	var vs string
	vs, err = messageIds.ValuesString()
	if err != nil {
		return "", err
	}

	return `SELECT Id, ThreadId, Text, CreatorUserId, CreatorTime, EditorUserId, EditorTime FROM ` + srv.dbTableNames.Messages + ` WHERE Id IN (` + vs + `);`, nil
}

func (srv *Server) dbQuery_ReadThreadsById(threadIds ul.UidList) (query string, err error) {
	var vs string
	vs, err = threadIds.ValuesString()
	if err != nil {
		return "", err
	}

	return `SELECT Id, ForumId, Name, Messages, CreatorUserId, CreatorTime, EditorUserId, EditorTime FROM ` + srv.dbTableNames.Threads + ` WHERE Id IN (` + vs + `);`, nil
}
