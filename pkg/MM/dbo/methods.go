package dbo

// Due to the large number of methods, they are sorted alphabetically.

import (
	"database/sql"
	"fmt"
	"time"

	mm "github.com/vault-thirteen/SimpleBB/pkg/MM/models"
	"github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
	cdbo "github.com/vault-thirteen/SimpleBB/pkg/common/dbo"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

func (dbo *DatabaseObject) CountForumsById(forumId uint) (n int, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountForumsById).QueryRow(forumId)

	n, err = cm.NewNonNullValueFromScannableSource[int](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountMessagesById(messageId uint) (n int, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountMessagesById).QueryRow(messageId)

	n, err = cm.NewNonNullValueFromScannableSource[int](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountRootSections() (n int, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountRootSections).QueryRow()

	n, err = cm.NewNonNullValueFromScannableSource[int](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountSectionsById(sectionId uint) (n int, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountSectionsById).QueryRow(sectionId)

	n, err = cm.NewNonNullValueFromScannableSource[int](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountThreadsById(threadId uint) (n int, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountThreadsById).QueryRow(threadId)

	n, err = cm.NewNonNullValueFromScannableSource[int](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) DeleteForumById(forumId uint) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_DeleteForumById).Exec(forumId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) DeleteMessageById(messageId uint) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_DeleteMessageById).Exec(messageId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) DeleteSectionById(sectionId uint) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_DeleteSectionById).Exec(sectionId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) DeleteThreadById(threadId uint) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_DeleteThreadById).Exec(threadId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) GetForumById(forumId uint) (forum *mm.Forum, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetForumById).QueryRow(forumId)

	forum, err = mm.NewForumFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return forum, nil
}

func (dbo *DatabaseObject) GetForumSectionById(forumId uint) (sectionId uint, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetForumSectionById).QueryRow(forumId)

	sectionId, err = cm.NewNonNullValueFromScannableSource[uint](row)
	if err != nil {
		return cdbo.IdOnError, err
	}

	return sectionId, nil
}

func (dbo *DatabaseObject) GetForumThreadsById(forumId uint) (threads *ul.UidList, err error) {
	threads = ul.New()
	err = dbo.DatabaseObject.PreparedStatement(DbPsid_GetForumThreadsById).QueryRow(forumId).Scan(threads)
	if err != nil {
		return nil, err
	}

	return threads, nil
}

func (dbo *DatabaseObject) GetMessageById(messageId uint) (message *mm.Message, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetMessageById).QueryRow(messageId)

	message, err = mm.NewMessageFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (dbo *DatabaseObject) GetMessageCreatorAndTimeById(messageId uint) (creatorUserId uint, ToC time.Time, ToE *time.Time, err error) {
	// N.B.: ToC can not be null, but ToE can be null !
	err = dbo.DatabaseObject.PreparedStatement(DbPsid_GetMessageCreatorAndTimeById).QueryRow(messageId).Scan(&creatorUserId, &ToC, &ToE)
	if err != nil {
		return cdbo.IdOnError, time.Time{}, nil, err
	}

	return creatorUserId, ToC, ToE, nil
}

func (dbo *DatabaseObject) GetMessageThreadById(messageId uint) (threadId uint, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetMessageThreadById).QueryRow(messageId)

	threadId, err = cm.NewNonNullValueFromScannableSource[uint](row)
	if err != nil {
		return cdbo.IdOnError, err
	}

	return threadId, nil
}

func (dbo *DatabaseObject) GetSectionById(sectionId uint) (section *mm.Section, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetSectionById).QueryRow(sectionId)

	section, err = mm.NewSectionFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return section, nil
}

func (dbo *DatabaseObject) GetSectionChildTypeById(sectionId uint) (childType byte, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetSectionChildTypeById).QueryRow(sectionId)

	childType, err = cm.NewNonNullValueFromScannableSource[byte](row)
	if err != nil {
		return 0, err
	}

	return childType, nil
}

func (dbo *DatabaseObject) GetSectionChildrenById(sectionId uint) (children *ul.UidList, err error) {
	children = ul.New()
	err = dbo.DatabaseObject.PreparedStatement(DbPsid_GetSectionChildrenById).QueryRow(sectionId).Scan(children)
	if err != nil {
		return nil, err
	}

	return children, nil
}

func (dbo *DatabaseObject) GetSectionParentById(sectionId uint) (parent *uint, err error) {
	parent = new(uint)
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetSectionParentById).QueryRow(sectionId)

	parent, err = cm.NewValueFromScannableSource[uint](row)
	if err != nil {
		return nil, err
	}

	return parent, nil
}

func (dbo *DatabaseObject) GetThreadById(threadId uint) (thread *mm.Thread, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetThreadByIdM).QueryRow(threadId)

	thread, err = mm.NewThreadFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (dbo *DatabaseObject) GetThreadForumById(threadId uint) (forumId uint, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetThreadForumById).QueryRow(threadId)

	forumId, err = cm.NewNonNullValueFromScannableSource[uint](row)
	if err != nil {
		return cdbo.IdOnError, err
	}

	return forumId, nil
}

func (dbo *DatabaseObject) GetThreadMessagesById(threadId uint) (messages *ul.UidList, err error) {
	messages = ul.New()
	err = dbo.DatabaseObject.PreparedStatement(DbPsid_GetThreadMessagesById).QueryRow(threadId).Scan(messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (dbo *DatabaseObject) InsertNewForum(sectionId uint, name string, creatorUserId uint) (lastInsertedId int64, err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_InsertNewForum).Exec(sectionId, name, creatorUserId)
	if err != nil {
		return cdbo.LastInsertedIdOnError, err
	}

	lastInsertedId, err = result.LastInsertId()
	if err != nil {
		return cdbo.LastInsertedIdOnError, err
	}

	return lastInsertedId, nil
}

func (dbo *DatabaseObject) InsertNewMessage(parentThread uint, messageText string, textChecksum uint32, creatorUserId uint) (lastInsertedId int64, err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_InsertNewMessage).Exec(parentThread, messageText, textChecksum, creatorUserId)
	if err != nil {
		return cdbo.LastInsertedIdOnError, err
	}

	lastInsertedId, err = result.LastInsertId()
	if err != nil {
		return cdbo.LastInsertedIdOnError, err
	}

	return lastInsertedId, nil
}

func (dbo *DatabaseObject) InsertNewSection(parent *uint, name string, creatorUserId uint) (lastInsertedId int64, err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_InsertNewSection).Exec(parent, name, creatorUserId)
	if err != nil {
		return cdbo.LastInsertedIdOnError, err
	}

	lastInsertedId, err = result.LastInsertId()
	if err != nil {
		return cdbo.LastInsertedIdOnError, err
	}

	return lastInsertedId, nil
}

func (dbo *DatabaseObject) InsertNewThread(parentForum uint, threadName string, creatorUserId uint) (lastInsertedId int64, err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_InsertNewThread).Exec(parentForum, threadName, creatorUserId)
	if err != nil {
		return cdbo.LastInsertedIdOnError, err
	}

	lastInsertedId, err = result.LastInsertId()
	if err != nil {
		return cdbo.LastInsertedIdOnError, err
	}

	return lastInsertedId, nil
}

func (dbo *DatabaseObject) ReadForums() (forums []mm.Forum, err error) {
	forums = make([]mm.Forum, 0)

	var rows *sql.Rows
	rows, err = dbo.DatabaseObject.PreparedStatement(DbPsid_ReadForums).Query()
	if err != nil {
		return nil, err
	}

	var forum *mm.Forum
	for rows.Next() {
		forum, err = mm.NewForumFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		if forum != nil {
			forums = append(forums, *forum)
		}
	}

	return forums, nil
}

func (dbo *DatabaseObject) ReadMessagesById(messageIds ul.UidList) (messages []mm.Message, err error) {
	messages = make([]mm.Message, 0, messageIds.Size())

	var query string
	query, err = dbo.dbQuery_ReadMessagesById(messageIds)
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	rows, err = dbo.DatabaseObject.DB().Query(query)
	if err != nil {
		return nil, err
	}

	var msg *mm.Message
	for rows.Next() {
		msg, err = mm.NewMessageFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		if msg != nil {
			messages = append(messages, *msg)
		}
	}

	return messages, nil
}

func (dbo *DatabaseObject) ReadSections() (sections []mm.Section, err error) {
	sections = make([]mm.Section, 0)

	var rows *sql.Rows
	rows, err = dbo.DatabaseObject.PreparedStatement(DbPsid_ReadSections).Query()
	if err != nil {
		return nil, err
	}

	var section *mm.Section
	for rows.Next() {
		section, err = mm.NewSectionFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		if section != nil {
			sections = append(sections, *section)
		}
	}

	return sections, nil
}

func (dbo *DatabaseObject) ReadThreadNamesByIds(threadIds ul.UidList) (threadNames []string, err error) {
	threadNames = make([]string, 0, threadIds.Size())

	var query string
	query, err = dbo.dbQuery_ReadThreadNamesByIds(threadIds)
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	rows, err = dbo.DB().Query(query)
	if err != nil {
		return nil, err
	}

	var threadName string
	for rows.Next() {
		threadName, err = cm.NewNonNullValueFromScannableSource[string](rows)
		if err != nil {
			return nil, err
		}

		threadNames = append(threadNames, threadName)
	}

	return threadNames, nil
}

func (dbo *DatabaseObject) ReadThreadsById(threadIds ul.UidList) (threads []mm.Thread, err error) {
	threads = make([]mm.Thread, 0, threadIds.Size())

	var query string
	query, err = dbo.dbQuery_ReadThreadsById(threadIds)
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	rows, err = dbo.DB().Query(query)
	if err != nil {
		return nil, err
	}

	var thread *mm.Thread
	for rows.Next() {
		thread, err = mm.NewThreadFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		if thread != nil {
			threads = append(threads, *thread)
		}
	}

	return threads, nil
}

func (dbo *DatabaseObject) SetForumNameById(forumId uint, name string, editorUserId uint) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetForumNameById).Exec(name, editorUserId, forumId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetForumSectionById(forumId uint, sectionId uint, editorUserId uint) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetForumSectionById).Exec(sectionId, editorUserId, forumId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetForumThreadsById(forumId uint, threads *ul.UidList) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetForumThreadsById).Exec(threads, forumId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetMessageTextById(messageId uint, text string, textChecksum uint32, editorUserId uint) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetMessageTextById).Exec(text, textChecksum, editorUserId, messageId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetMessageThreadById(messageId uint, thread uint, editorUserId uint) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetMessageThreadById).Exec(thread, editorUserId, messageId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetSectionChildTypeById(sectionId uint, childType byte) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetSectionChildTypeById).Exec(childType, sectionId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetSectionChildrenById(sectionId uint, children *ul.UidList) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetSectionChildrenById).Exec(children, sectionId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetSectionNameById(sectionId uint, name string, editorUserId uint) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetSectionNameById).Exec(name, editorUserId, sectionId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetSectionParentById(sectionId uint, parent uint, editorUserId uint) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetSectionParentById).Exec(parent, editorUserId, sectionId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetThreadForumById(threadId uint, forum uint, editorUserId uint) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetThreadForumById).Exec(forum, editorUserId, threadId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetThreadMessagesById(threadId uint, messages *ul.UidList) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetThreadMessagesById).Exec(messages, threadId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetThreadNameById(threadId uint, name string, editorUserId uint) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetThreadNameById).Exec(name, editorUserId, threadId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(cdbo.ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}
