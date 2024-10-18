package dbo

// Due to the large number of methods, they are sorted alphabetically.

import (
	"database/sql"
	"time"

	mm "github.com/vault-thirteen/SimpleBB/pkg/MM/models"
	"github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
	cdbo "github.com/vault-thirteen/SimpleBB/pkg/common/dbo"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

func (dbo *DatabaseObject) CountForumsById(forumId cmb.Id) (n cmb.Count, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountForumsById).QueryRow(forumId)

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountMessagesById(messageId cmb.Id) (n cmb.Count, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountMessagesById).QueryRow(messageId)

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountRootSections() (n cmb.Count, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountRootSections).QueryRow()

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountSectionsById(sectionId cmb.Id) (n cmb.Count, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountSectionsById).QueryRow(sectionId)

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountThreadsById(threadId cmb.Id) (n cmb.Count, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountThreadsById).QueryRow(threadId)

	n, err = cm.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return cdbo.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) DeleteForumById(forumId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_DeleteForumById).Exec(forumId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) DeleteMessageById(messageId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_DeleteMessageById).Exec(messageId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) DeleteSectionById(sectionId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_DeleteSectionById).Exec(sectionId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) DeleteThreadById(threadId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_DeleteThreadById).Exec(threadId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) GetForumById(forumId cmb.Id) (forum *mm.Forum, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetForumById).QueryRow(forumId)

	forum, err = mm.NewForumFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return forum, nil
}

func (dbo *DatabaseObject) GetForumSectionById(forumId cmb.Id) (sectionId cmb.Id, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetForumSectionById).QueryRow(forumId)

	sectionId, err = cm.NewNonNullValueFromScannableSource[cmb.Id](row)
	if err != nil {
		return cdbo.IdOnError, err
	}

	return sectionId, nil
}

func (dbo *DatabaseObject) GetForumThreadsById(forumId cmb.Id) (threads *ul.UidList, err error) {
	threads = ul.New()
	err = dbo.DatabaseObject.PreparedStatement(DbPsid_GetForumThreadsById).QueryRow(forumId).Scan(threads)
	if err != nil {
		return nil, err
	}

	return threads, nil
}

func (dbo *DatabaseObject) GetMessageById(messageId cmb.Id) (message *mm.Message, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetMessageById).QueryRow(messageId)

	message, err = mm.NewMessageFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (dbo *DatabaseObject) GetMessageCreatorAndTimeById(messageId cmb.Id) (creatorUserId cmb.Id, ToC time.Time, ToE *time.Time, err error) {
	// N.B.: ToC can not be null, but ToE can be null !
	err = dbo.DatabaseObject.PreparedStatement(DbPsid_GetMessageCreatorAndTimeById).QueryRow(messageId).Scan(&creatorUserId, &ToC, &ToE)
	if err != nil {
		return cdbo.IdOnError, time.Time{}, nil, err
	}

	return creatorUserId, ToC, ToE, nil
}

func (dbo *DatabaseObject) GetMessageThreadById(messageId cmb.Id) (threadId cmb.Id, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetMessageThreadById).QueryRow(messageId)

	threadId, err = cm.NewNonNullValueFromScannableSource[cmb.Id](row)
	if err != nil {
		return cdbo.IdOnError, err
	}

	return threadId, nil
}

func (dbo *DatabaseObject) GetSectionById(sectionId cmb.Id) (section *mm.Section, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetSectionById).QueryRow(sectionId)

	section, err = mm.NewSectionFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return section, nil
}

func (dbo *DatabaseObject) GetSectionChildTypeById(sectionId cmb.Id) (childType byte, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetSectionChildTypeById).QueryRow(sectionId)

	childType, err = cm.NewNonNullValueFromScannableSource[byte](row)
	if err != nil {
		return 0, err
	}

	return childType, nil
}

func (dbo *DatabaseObject) GetSectionChildrenById(sectionId cmb.Id) (children *ul.UidList, err error) {
	children = ul.New()
	err = dbo.DatabaseObject.PreparedStatement(DbPsid_GetSectionChildrenById).QueryRow(sectionId).Scan(children)
	if err != nil {
		return nil, err
	}

	return children, nil
}

func (dbo *DatabaseObject) GetSectionParentById(sectionId cmb.Id) (parent *cmb.Id, err error) {
	parent = new(cmb.Id)
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetSectionParentById).QueryRow(sectionId)

	parent, err = cm.NewValueFromScannableSource[cmb.Id](row)
	if err != nil {
		return nil, err
	}

	return parent, nil
}

func (dbo *DatabaseObject) GetThreadById(threadId cmb.Id) (thread *mm.Thread, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetThreadByIdM).QueryRow(threadId)

	thread, err = mm.NewThreadFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (dbo *DatabaseObject) GetThreadForumById(threadId cmb.Id) (forumId cmb.Id, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetThreadForumById).QueryRow(threadId)

	forumId, err = cm.NewNonNullValueFromScannableSource[cmb.Id](row)
	if err != nil {
		return cdbo.IdOnError, err
	}

	return forumId, nil
}

func (dbo *DatabaseObject) GetThreadMessagesById(threadId cmb.Id) (messages *ul.UidList, err error) {
	messages = ul.New()
	err = dbo.DatabaseObject.PreparedStatement(DbPsid_GetThreadMessagesById).QueryRow(threadId).Scan(messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (dbo *DatabaseObject) InsertNewForum(sectionId cmb.Id, name cm.Name, creatorUserId cmb.Id) (lastInsertedId cmb.Id, err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_InsertNewForum).Exec(sectionId, name, creatorUserId)
	if err != nil {
		return cdbo.LastInsertedIdOnError, err
	}

	return cdbo.CheckRowsAffectedAndGetLastInsertedId(result, 1)
}

func (dbo *DatabaseObject) InsertNewMessage(parentThread cmb.Id, messageText cmb.Text, textChecksum []byte, creatorUserId cmb.Id) (lastInsertedId cmb.Id, err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_InsertNewMessage).Exec(parentThread, messageText, textChecksum, creatorUserId)
	if err != nil {
		return cdbo.LastInsertedIdOnError, err
	}

	return cdbo.CheckRowsAffectedAndGetLastInsertedId(result, 1)
}

func (dbo *DatabaseObject) InsertNewSection(parent *cmb.Id, name cm.Name, creatorUserId cmb.Id) (lastInsertedId cmb.Id, err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_InsertNewSection).Exec(parent, name, creatorUserId)
	if err != nil {
		return cdbo.LastInsertedIdOnError, err
	}

	return cdbo.CheckRowsAffectedAndGetLastInsertedId(result, 1)
}

func (dbo *DatabaseObject) InsertNewThread(parentForum cmb.Id, threadName cm.Name, creatorUserId cmb.Id) (lastInsertedId cmb.Id, err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_InsertNewThread).Exec(parentForum, threadName, creatorUserId)
	if err != nil {
		return cdbo.LastInsertedIdOnError, err
	}

	return cdbo.CheckRowsAffectedAndGetLastInsertedId(result, 1)
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

func (dbo *DatabaseObject) ReadMessagesById(messageIds *ul.UidList) (messages []mm.Message, err error) {
	if messageIds == nil {
		return []mm.Message{}, nil
	}

	messages = make([]mm.Message, 0, messageIds.Size())

	var query string
	query, err = dbo.dbQuery_ReadMessagesById(*messageIds)
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

func (dbo *DatabaseObject) ReadThreadNamesByIds(threadIds ul.UidList) (threadNames []cm.Name, err error) {
	threadNames = make([]cm.Name, 0, threadIds.Size())

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

	var threadName cm.Name
	for rows.Next() {
		threadName, err = cm.NewNonNullValueFromScannableSource[cm.Name](rows)
		if err != nil {
			return nil, err
		}

		threadNames = append(threadNames, threadName)
	}

	return threadNames, nil
}

func (dbo *DatabaseObject) ReadThreadsById(threadIds *ul.UidList) (threads []mm.Thread, err error) {
	if threadIds == nil {
		return []mm.Thread{}, nil
	}

	threads = make([]mm.Thread, 0, threadIds.Size())

	var query string
	query, err = dbo.dbQuery_ReadThreadsById(*threadIds)
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

func (dbo *DatabaseObject) SetForumNameById(forumId cmb.Id, name cm.Name, editorUserId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetForumNameById).Exec(name, editorUserId, forumId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetForumSectionById(forumId cmb.Id, sectionId cmb.Id, editorUserId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetForumSectionById).Exec(sectionId, editorUserId, forumId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetForumThreadsById(forumId cmb.Id, threads *ul.UidList) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetForumThreadsById).Exec(threads, forumId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetMessageTextById(messageId cmb.Id, text cmb.Text, textChecksum []byte, editorUserId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetMessageTextById).Exec(text, textChecksum, editorUserId, messageId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetMessageThreadById(messageId cmb.Id, threadId cmb.Id, editorUserId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetMessageThreadById).Exec(threadId, editorUserId, messageId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetSectionChildTypeById(sectionId cmb.Id, childType byte) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetSectionChildTypeById).Exec(childType, sectionId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetSectionChildrenById(sectionId cmb.Id, children *ul.UidList) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetSectionChildrenById).Exec(children, sectionId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetSectionNameById(sectionId cmb.Id, name cm.Name, editorUserId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetSectionNameById).Exec(name, editorUserId, sectionId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetSectionParentById(sectionId cmb.Id, parent cmb.Id, editorUserId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetSectionParentById).Exec(parent, editorUserId, sectionId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetThreadForumById(threadId cmb.Id, forumId cmb.Id, editorUserId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetThreadForumById).Exec(forumId, editorUserId, threadId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetThreadMessagesById(threadId cmb.Id, messages *ul.UidList) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetThreadMessagesById).Exec(messages, threadId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetThreadNameById(threadId cmb.Id, name cm.Name, editorUserId cmb.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetThreadNameById).Exec(name, editorUserId, threadId)
	if err != nil {
		return err
	}

	return cdbo.CheckRowsAffected(result, 1)
}
