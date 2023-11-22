package dbo

import (
	ul "github.com/vault-thirteen/SimpleBB/pkg/UidList"
)

func (dbo *DatabaseObject) dbQuery_ReadMessagesById(messageIds ul.UidList) (query string, err error) {
	var vs string
	vs, err = messageIds.ValuesString()
	if err != nil {
		return "", err
	}

	return `SELECT Id, ThreadId, Text, CreatorUserId, CreatorTime, EditorUserId, EditorTime FROM ` + dbo.tableNames.Messages + ` WHERE Id IN (` + vs + `);`, nil
}

func (dbo *DatabaseObject) dbQuery_ReadThreadsById(threadIds ul.UidList) (query string, err error) {
	var vs string
	vs, err = threadIds.ValuesString()
	if err != nil {
		return "", err
	}

	return `SELECT Id, ForumId, Name, Messages, CreatorUserId, CreatorTime, EditorUserId, EditorTime FROM ` + dbo.tableNames.Threads + ` WHERE Id IN (` + vs + `);`, nil
}
