package mm

import (
	"database/sql"
	"fmt"
)

// Indices of prepared statements.
const (
	DbPsid_CountRootForums              = 0
	DbPsid_InsertNewForum               = 1
	DbPsid_CountForumsById              = 2
	DbPsid_GetForumChildrenById         = 3
	DbPsid_SetForumChildrenById         = 4
	DbPsid_SetForumNameById             = 5
	DbPsid_GetRootForumId               = 6
	DbPsid_SetForumParentById           = 7
	DbPsid_GetForumParentById           = 8
	DbPsid_InsertNewThread              = 9
	DbPsid_GetForumThreadsById          = 10
	DbPsid_SetForumThreadsById          = 11
	DbPsid_SetThreadNameById            = 12
	DbPsid_GetThreadForumById           = 13
	DbPsid_SetThreadForumById           = 14
	DbPsid_CountThreadsById             = 15
	DbPsid_GetThreadMessagesById        = 16
	DbPsid_InsertNewMessage             = 17
	DbPsid_SetThreadMessagesById        = 18
	DbPsid_SetMessageTextById           = 19
	DbPsid_GetMessageThreadById         = 20
	DbPsid_SetMessageThreadById         = 21
	DbPsid_GetMessageCreatorAndTimeById = 22
	DbPsid_GetMessageById               = 23
	DbPsid_DeleteMessageById            = 24
	DbPsid_GetThreadByIdM               = 25
	DbPsid_DeleteThreadById             = 26
	DbPsid_GetForumById                 = 27
	DbPsid_DeleteForumById              = 28
	DbPsid_ReadForums                   = 29
)

func (srv *Server) prepareDbStatements() (err error) {
	var q string
	qs := make([]string, 0)

	// 0.
	q = fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE Parent IS NULL;`, srv.dbTableNames.Forums)
	qs = append(qs, q)

	// 1.
	q = fmt.Sprintf(`INSERT INTO %s (Parent, NAME, CreatorUserId, CreatorTime) VALUES (?, ?, ?, Now());`, srv.dbTableNames.Forums)
	qs = append(qs, q)

	// 2.
	q = fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE Id = ?;`, srv.dbTableNames.Forums)
	qs = append(qs, q)

	// 3.
	q = fmt.Sprintf(`SELECT Children FROM %s WHERE Id = ?;`, srv.dbTableNames.Forums)
	qs = append(qs, q)

	// 4.
	q = fmt.Sprintf(`UPDATE %s SET Children = ? WHERE Id = ?;`, srv.dbTableNames.Forums)
	qs = append(qs, q)

	// 5.
	q = fmt.Sprintf(`UPDATE %s SET NAME = ?, EditorUserId = ?, EditorTime = Now() WHERE Id = ?;`, srv.dbTableNames.Forums)
	qs = append(qs, q)

	// 6.
	q = fmt.Sprintf(`SELECT Id FROM %s WHERE Parent IS NULL;`, srv.dbTableNames.Forums)
	qs = append(qs, q)

	// 7.
	q = fmt.Sprintf(`UPDATE %s SET Parent = ?, EditorUserId = ?, EditorTime = Now() WHERE Id = ?;`, srv.dbTableNames.Forums)
	qs = append(qs, q)

	// 8.
	q = fmt.Sprintf(`SELECT Parent FROM %s WHERE Id = ?;`, srv.dbTableNames.Forums)
	qs = append(qs, q)

	// 9.
	q = fmt.Sprintf(`INSERT INTO %s (ForumId, NAME, CreatorUserId, CreatorTime) VALUES (?, ?, ?, Now());`, srv.dbTableNames.Threads)
	qs = append(qs, q)

	// 10.
	q = fmt.Sprintf(`SELECT Threads FROM %s WHERE Id = ?;`, srv.dbTableNames.Forums)
	qs = append(qs, q)

	// 11.
	q = fmt.Sprintf(`UPDATE %s SET Threads = ? WHERE Id = ?;`, srv.dbTableNames.Forums)
	qs = append(qs, q)

	// 12.
	q = fmt.Sprintf(`UPDATE %s SET NAME = ?, EditorUserId = ?, EditorTime = Now() WHERE Id = ?;`, srv.dbTableNames.Threads)
	qs = append(qs, q)

	// 13.
	q = fmt.Sprintf(`SELECT ForumId FROM %s WHERE Id = ?;`, srv.dbTableNames.Threads)
	qs = append(qs, q)

	// 14.
	q = fmt.Sprintf(`UPDATE %s SET ForumId = ?, EditorUserId = ?, EditorTime = Now() WHERE Id = ?;`, srv.dbTableNames.Threads)
	qs = append(qs, q)

	// 15.
	q = fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE Id = ?;`, srv.dbTableNames.Threads)
	qs = append(qs, q)

	// 16.
	q = fmt.Sprintf(`SELECT Messages FROM %s WHERE Id = ?;`, srv.dbTableNames.Threads)
	qs = append(qs, q)

	// 17.
	q = fmt.Sprintf(`INSERT INTO %s (ThreadId, TEXT, CreatorUserId, CreatorTime) VALUES (?, ?, ?, Now());`, srv.dbTableNames.Messages)
	qs = append(qs, q)

	// 18.
	q = fmt.Sprintf(`UPDATE %s SET Messages = ? WHERE Id = ?;`, srv.dbTableNames.Threads)
	qs = append(qs, q)

	// 19.
	q = fmt.Sprintf(`UPDATE %s SET TEXT = ?, EditorUserId = ?, EditorTime = Now() WHERE Id = ?;`, srv.dbTableNames.Messages)
	qs = append(qs, q)

	// 20.
	q = fmt.Sprintf(`SELECT ThreadId FROM %s WHERE Id = ?;`, srv.dbTableNames.Messages)
	qs = append(qs, q)

	// 21.
	q = fmt.Sprintf(`UPDATE %s SET ThreadId = ?, EditorUserId = ?, EditorTime = Now() WHERE Id = ?;`, srv.dbTableNames.Messages)
	qs = append(qs, q)

	// 22.
	q = fmt.Sprintf(`SELECT CreatorUserId, CreatorTime, EditorTime FROM %s WHERE Id = ?;`, srv.dbTableNames.Messages)
	qs = append(qs, q)

	// 23.
	q = fmt.Sprintf(`SELECT Id, ThreadId, Text, CreatorUserId, CreatorTime, EditorUserId, EditorTime FROM %s WHERE Id = ?;`, srv.dbTableNames.Messages)
	qs = append(qs, q)

	// 24.
	q = fmt.Sprintf(`DELETE FROM %s WHERE Id = ?;`, srv.dbTableNames.Messages)
	qs = append(qs, q)

	// 25.
	q = fmt.Sprintf(`SELECT Id, ForumId, Name, Messages, CreatorUserId, CreatorTime, EditorUserId, EditorTime FROM %s WHERE Id = ?;`, srv.dbTableNames.Threads)
	qs = append(qs, q)

	// 26.
	q = fmt.Sprintf(`DELETE FROM %s WHERE Id = ?;`, srv.dbTableNames.Threads)
	qs = append(qs, q)

	// 27.
	q = fmt.Sprintf(`SELECT Id, Parent, Children, Name, Threads, CreatorUserId, CreatorTime, EditorUserId, EditorTime FROM %s WHERE Id = ?;`, srv.dbTableNames.Forums)
	qs = append(qs, q)

	// 28.
	q = fmt.Sprintf(`DELETE FROM %s WHERE Id = ?;`, srv.dbTableNames.Forums)
	qs = append(qs, q)

	// 29.
	q = fmt.Sprintf(`SELECT Id, Parent, Children, Name, Threads, CreatorUserId, CreatorTime, EditorUserId, EditorTime FROM %s;`, srv.dbTableNames.Forums)
	qs = append(qs, q)

	srv.dbPreparedStatementQueries = make([]string, 0, len(qs))
	for _, q = range qs {
		srv.dbPreparedStatementQueries = append(srv.dbPreparedStatementQueries, q)
	}
	qs = nil

	srv.dbPreparedStatements = make([]*sql.Stmt, 0, len(srv.dbPreparedStatementQueries))
	var st *sql.Stmt
	for _, psq := range srv.dbPreparedStatementQueries {
		st, err = srv.db.Prepare(psq)
		if err != nil {
			return err
		}

		srv.dbPreparedStatements = append(srv.dbPreparedStatements, st)
	}

	return nil
}
