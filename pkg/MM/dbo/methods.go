package dbo

// Due to the large number of methods, they are sorted alphabetically.

import (
	"database/sql"
	"fmt"
	"time"

	mm "github.com/vault-thirteen/SimpleBB/pkg/MM/models"
	ul "github.com/vault-thirteen/SimpleBB/pkg/UidList"
)

func (dbo *DatabaseObject) CountForumsById(forumId uint) (n int, err error) {
	err = dbo.preparedStatements[DbPsid_CountForumsById].QueryRow(forumId).Scan(&n)
	if err != nil {
		return CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountRootSections() (n int, err error) {
	err = dbo.preparedStatements[DbPsid_CountRootSections].QueryRow().Scan(&n)
	if err != nil {
		return CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountSectionsById(sectionId uint) (n int, err error) {
	err = dbo.preparedStatements[DbPsid_CountSectionsById].QueryRow(sectionId).Scan(&n)
	if err != nil {
		return CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountThreadsById(threadId uint) (n int, err error) {
	err = dbo.preparedStatements[DbPsid_CountThreadsById].QueryRow(threadId).Scan(&n)
	if err != nil {
		return CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) DeleteForumById(forumId uint) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_DeleteForumById].Exec(forumId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) DeleteMessageById(messageId uint) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_DeleteMessageById].Exec(messageId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) DeleteSectionById(sectionId uint) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_DeleteSectionById].Exec(sectionId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) DeleteThreadById(threadId uint) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_DeleteThreadById].Exec(threadId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) GetForumById(forumId uint) (forum *mm.Forum, err error) {
	forum = mm.NewForum()
	err = dbo.preparedStatements[DbPsid_GetForumById].QueryRow(forumId).Scan(
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

func (dbo *DatabaseObject) GetForumSectionById(forumId uint) (sectionId uint, err error) {
	err = dbo.preparedStatements[DbPsid_GetForumSectionById].QueryRow(forumId).Scan(&sectionId)
	if err != nil {
		return IdOnError, err
	}

	return sectionId, nil
}

func (dbo *DatabaseObject) GetForumThreadsById(forumId uint) (threads *ul.UidList, err error) {
	threads = ul.New()
	err = dbo.preparedStatements[DbPsid_GetForumThreadsById].QueryRow(forumId).Scan(threads)
	if err != nil {
		return nil, err
	}

	return threads, nil
}

func (dbo *DatabaseObject) GetMessageById(messageId uint) (message *mm.Message, err error) {
	message = mm.NewMessage()
	err = dbo.preparedStatements[DbPsid_GetMessageById].QueryRow(messageId).Scan(
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

func (dbo *DatabaseObject) GetMessageCreatorAndTimeById(messageId uint) (creatorUserId uint, ToC time.Time, ToE *time.Time, err error) {
	// N.B.: ToC can not be null, but ToE can be null !
	err = dbo.preparedStatements[DbPsid_GetMessageCreatorAndTimeById].QueryRow(messageId).Scan(&creatorUserId, &ToC, &ToE)
	if err != nil {
		return IdOnError, time.Time{}, nil, err
	}

	return creatorUserId, ToC, ToE, nil
}

func (dbo *DatabaseObject) GetMessageThreadById(messageId uint) (thread uint, err error) {
	err = dbo.preparedStatements[DbPsid_GetMessageThreadById].QueryRow(messageId).Scan(&thread)
	if err != nil {
		return IdOnError, err
	}

	return thread, nil
}

func (dbo *DatabaseObject) GetSectionById(sectionId uint) (section *mm.Section, err error) {
	section = mm.NewSection()
	err = dbo.preparedStatements[DbPsid_GetSectionById].QueryRow(sectionId).Scan(
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

func (dbo *DatabaseObject) GetSectionChildTypeById(sectionId uint) (childType byte, err error) {
	err = dbo.preparedStatements[DbPsid_GetSectionChildTypeById].QueryRow(sectionId).Scan(&childType)
	if err != nil {
		return 0, err
	}

	return childType, nil
}

func (dbo *DatabaseObject) GetSectionChildrenById(sectionId uint) (children *ul.UidList, err error) {
	children = ul.New()
	err = dbo.preparedStatements[DbPsid_GetSectionChildrenById].QueryRow(sectionId).Scan(children)
	if err != nil {
		return nil, err
	}

	return children, nil
}

func (dbo *DatabaseObject) GetSectionParentById(sectionId uint) (parent *uint, err error) {
	parent = new(uint)
	err = dbo.preparedStatements[DbPsid_GetSectionParentById].QueryRow(sectionId).Scan(&parent)
	if err != nil {
		return nil, err
	}

	return parent, nil
}

func (dbo *DatabaseObject) GetThreadById(threadId uint) (thread *mm.Thread, err error) {
	thread = mm.NewThread()
	err = dbo.preparedStatements[DbPsid_GetThreadByIdM].QueryRow(threadId).Scan(
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

func (dbo *DatabaseObject) GetThreadForumById(threadId uint) (forum uint, err error) {
	err = dbo.preparedStatements[DbPsid_GetThreadForumById].QueryRow(threadId).Scan(&forum)
	if err != nil {
		return IdOnError, err
	}

	return forum, nil
}

func (dbo *DatabaseObject) GetThreadMessagesById(threadId uint) (messages *ul.UidList, err error) {
	messages = ul.New()
	err = dbo.preparedStatements[DbPsid_GetThreadMessagesById].QueryRow(threadId).Scan(messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (dbo *DatabaseObject) InsertNewForum(sectionId uint, name string, creatorUserId uint) (lastInsertedId int64, err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_InsertNewForum].Exec(sectionId, name, creatorUserId)
	if err != nil {
		return LastInsertedIdOnError, err
	}

	lastInsertedId, err = result.LastInsertId()
	if err != nil {
		return LastInsertedIdOnError, err
	}

	return lastInsertedId, nil
}

func (dbo *DatabaseObject) InsertNewMessage(parentThread uint, messageText string, creatorUserId uint) (lastInsertedId int64, err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_InsertNewMessage].Exec(parentThread, messageText, creatorUserId)
	if err != nil {
		return LastInsertedIdOnError, err
	}

	lastInsertedId, err = result.LastInsertId()
	if err != nil {
		return LastInsertedIdOnError, err
	}

	return lastInsertedId, nil
}

func (dbo *DatabaseObject) InsertNewSection(parent *uint, name string, creatorUserId uint) (lastInsertedId int64, err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_InsertNewSection].Exec(parent, name, creatorUserId)
	if err != nil {
		return LastInsertedIdOnError, err
	}

	lastInsertedId, err = result.LastInsertId()
	if err != nil {
		return LastInsertedIdOnError, err
	}

	return lastInsertedId, nil
}

func (dbo *DatabaseObject) InsertNewThread(parentForum uint, threadName string, creatorUserId uint) (lastInsertedId int64, err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_InsertNewThread].Exec(parentForum, threadName, creatorUserId)
	if err != nil {
		return LastInsertedIdOnError, err
	}

	lastInsertedId, err = result.LastInsertId()
	if err != nil {
		return LastInsertedIdOnError, err
	}

	return lastInsertedId, nil
}

func (dbo *DatabaseObject) ReadForums() (forums []mm.Forum, err error) {
	forums = make([]mm.Forum, 0)

	var rows *sql.Rows
	rows, err = dbo.preparedStatements[DbPsid_ReadForums].Query()
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

func (dbo *DatabaseObject) ReadMessagesById(messageIds ul.UidList) (messages []mm.Message, err error) {
	messages = make([]mm.Message, 0, messageIds.Size())

	var query string
	query, err = dbo.dbQuery_ReadMessagesById(messageIds)
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	rows, err = dbo.db.Query(query)
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

func (dbo *DatabaseObject) ReadSections() (sections []mm.Section, err error) {
	sections = make([]mm.Section, 0)

	var rows *sql.Rows
	rows, err = dbo.preparedStatements[DbPsid_ReadSections].Query()
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

func (dbo *DatabaseObject) ReadThreadsById(threadIds ul.UidList) (threads []mm.Thread, err error) {
	threads = make([]mm.Thread, 0, threadIds.Size())

	var query string
	query, err = dbo.dbQuery_ReadThreadsById(threadIds)
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	rows, err = dbo.db.Query(query)
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

func (dbo *DatabaseObject) SetForumNameById(forumId uint, name string, editorUserId uint) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SetForumNameById].Exec(name, editorUserId, forumId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetForumSectionById(forumId uint, sectionId uint, editorUserId uint) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SetForumSectionById].Exec(sectionId, editorUserId, forumId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetForumThreadsById(forumId uint, threads *ul.UidList) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SetForumThreadsById].Exec(threads, forumId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetMessageTextById(messageId uint, text string, editorUserId uint) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SetMessageTextById].Exec(text, editorUserId, messageId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetMessageThreadById(messageId uint, thread uint, editorUserId uint) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SetMessageThreadById].Exec(thread, editorUserId, messageId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetSectionChildTypeById(sectionId uint, childType byte) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SetSectionChildTypeById].Exec(childType, sectionId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetSectionChildrenById(sectionId uint, children *ul.UidList) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SetSectionChildrenById].Exec(children, sectionId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetSectionNameById(sectionId uint, name string, editorUserId uint) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SetSectionNameById].Exec(name, editorUserId, sectionId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetSectionParentById(sectionId uint, parent uint, editorUserId uint) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SetSectionParentById].Exec(parent, editorUserId, sectionId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetThreadForumById(threadId uint, forum uint, editorUserId uint) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SetThreadForumById].Exec(forum, editorUserId, threadId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetThreadMessagesById(threadId uint, messages *ul.UidList) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SetThreadMessagesById].Exec(messages, threadId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}

func (dbo *DatabaseObject) SetThreadNameById(threadId uint, name string, editorUserId uint) (err error) {
	var result sql.Result
	result, err = dbo.preparedStatements[DbPsid_SetThreadNameById].Exec(name, editorUserId, threadId)
	if err != nil {
		return err
	}

	var ra int64
	ra, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return fmt.Errorf(ErrFRowsAffectedCount, 1, ra)
	}

	return nil
}
