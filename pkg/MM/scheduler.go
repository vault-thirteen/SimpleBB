package mm

import (
	"errors"
	"fmt"

	"github.com/vault-thirteen/SimpleBB/pkg/MM/dbo"
	mm "github.com/vault-thirteen/SimpleBB/pkg/MM/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

const (
	Err_TooManyRootSections = "too many root sections"
	ErrF_SectionIsNotFound  = "section is not found, ID=%v"
	ErrF_SectionIsDamaged   = "section is damaged, ID=%v"
	ErrF_ForumIsNotFound    = "forum is not found, ID=%v"
	ErrF_ForumIsDamaged     = "forum is damaged, ID=%v"
	ErrF_ThreadIsNotFound   = "thread is not found, ID=%v"
	ErrF_ThreadIsDamaged    = "thread is damaged, ID=%v"
	ErrF_MessageIsNotFound  = "message is not found, ID=%v"
	ErrF_MessageIsDamaged   = "message is damaged, ID=%v"
)

// checkDatabaseConsistency checks consistency of sections, forums, threads and
// messages. This function is used in the scheduler and is also run once during
// the server's start.
func (srv *Server) checkDatabaseConsistency() (err error) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	fmt.Print(c.MsgDatabaseConsistencyCheck)

	// Sections.
	var sections []mm.Section
	sections, err = srv.dbo.ReadSections()
	if err != nil {
		return err
	}

	sectionsMap := make(map[cmb.Id]mm.Section)
	for _, section := range sections {
		sectionsMap[section.Id] = section
	}

	err = checkSections(sections, sectionsMap)
	if err != nil {
		return err
	}

	// Forums.
	var forums []mm.Forum
	forums, err = srv.dbo.ReadForums()
	if err != nil {
		return err
	}

	forumsMap := make(map[cmb.Id]mm.Forum)
	for _, forum := range forums {
		forumsMap[forum.Id] = forum
	}

	err = checkForums(sections, sectionsMap, forums, forumsMap)
	if err != nil {
		return err
	}

	// Threads.
	var threads []mm.ThreadLink
	threads, err = srv.dbo.ReadThreadLinks()
	if err != nil {
		return err
	}

	threadsMap := make(map[cmb.Id]mm.ThreadLink)
	for _, thread := range threads {
		threadsMap[thread.Id] = thread
	}

	err = checkThreads(forums, forumsMap, threads, threadsMap)
	if err != nil {
		return err
	}

	// Messages.
	err = checkMessages(srv.dbo, threads)
	if err != nil {
		return err
	}

	fmt.Println(c.MsgOK)

	return nil
}

func checkSections(sections []mm.Section, sectionsMap map[cmb.Id]mm.Section) (err error) {
	// Step I. Downward check (parent to child).
	var childSection mm.Section
	var childIds []cmb.Id
	var ok bool
	for _, section := range sections {
		if section.ChildType.GetValue() != mm.SectionChildType_Section {
			continue
		}

		if section.Children.Size() == 0 {
			continue
		}

		childIds = section.Children.AsArray()

		for _, childId := range childIds {
			childSection, ok = sectionsMap[childId]
			if !ok {
				return fmt.Errorf(ErrF_SectionIsNotFound, childId)
			}

			if childSection.Parent == nil {
				return fmt.Errorf(ErrF_SectionIsDamaged, childId)
			}

			if *childSection.Parent != section.Id {
				return fmt.Errorf(ErrF_SectionIsDamaged, childId)
			}
		}
	}

	// Step II. Root section.
	var rootSectionsCount = 0
	for _, section := range sections {
		if section.Parent != nil {
			continue
		}
		rootSectionsCount++
	}
	if rootSectionsCount > 1 {
		return errors.New(Err_TooManyRootSections)
	}

	// Step III. Upward check (child to parent).
	var parentId cmb.Id
	var parentSection mm.Section
	for _, section := range sections {
		if section.Parent == nil {
			continue
		}

		parentId = *section.Parent

		parentSection, ok = sectionsMap[parentId]
		if !ok {
			return fmt.Errorf(ErrF_SectionIsNotFound, parentId)
		}

		if parentSection.ChildType.GetValue() != mm.SectionChildType_Section {
			return fmt.Errorf(ErrF_SectionIsDamaged, parentId)
		}

		if parentSection.Children.Size() == 0 {
			return fmt.Errorf(ErrF_SectionIsDamaged, parentId)
		}

		if !parentSection.Children.HasItem(section.Id) {
			return fmt.Errorf(ErrF_SectionIsDamaged, parentId)
		}
	}

	return nil
}

func checkForums(sections []mm.Section, sectionsMap map[cmb.Id]mm.Section, forums []mm.Forum, forumsMap map[cmb.Id]mm.Forum) (err error) {
	// Step I. Downward check (parent to child).
	var childIds []cmb.Id
	var ok bool
	for _, section := range sections {
		if section.ChildType.GetValue() != mm.SectionChildType_Forum {
			continue
		}

		if section.Children.Size() == 0 {
			continue
		}

		childIds = section.Children.AsArray()

		var forum mm.Forum
		for _, childId := range childIds {
			forum, ok = forumsMap[childId]
			if !ok {
				return fmt.Errorf(ErrF_ForumIsNotFound, childId)
			}

			if forum.SectionId != section.Id {
				return fmt.Errorf(ErrF_ForumIsDamaged, childId)
			}
		}
	}

	// Step II. Upward check (child to parent).
	var parentId cmb.Id
	var parentSection mm.Section
	for _, forum := range forums {
		parentId = forum.SectionId

		parentSection, ok = sectionsMap[parentId]
		if !ok {
			return fmt.Errorf(ErrF_SectionIsNotFound, parentId)
		}

		if parentSection.ChildType.GetValue() != mm.SectionChildType_Forum {
			return fmt.Errorf(ErrF_SectionIsDamaged, parentId)
		}

		if parentSection.Children.Size() == 0 {
			return fmt.Errorf(ErrF_SectionIsDamaged, parentId)
		}

		if !parentSection.Children.HasItem(forum.Id) {
			return fmt.Errorf(ErrF_SectionIsDamaged, parentId)
		}
	}

	return nil
}

func checkThreads(forums []mm.Forum, forumsMap map[cmb.Id]mm.Forum, threads []mm.ThreadLink, threadsMap map[cmb.Id]mm.ThreadLink) (err error) {
	// Step I. Downward check (parent to child).
	var childIds []cmb.Id
	var ok bool
	for _, forum := range forums {
		if forum.Threads.Size() == 0 {
			continue
		}

		childIds = forum.Threads.AsArray()

		var thread mm.ThreadLink
		for _, childId := range childIds {
			thread, ok = threadsMap[childId]
			if !ok {
				return fmt.Errorf(ErrF_ThreadIsNotFound, childId)
			}

			if thread.ForumId != forum.Id {
				return fmt.Errorf(ErrF_ThreadIsDamaged, childId)
			}
		}
	}

	// Step II. Upward check (child to parent).
	var parentId cmb.Id
	var parentForum mm.Forum
	for _, thread := range threads {
		parentId = thread.ForumId

		parentForum, ok = forumsMap[parentId]
		if !ok {
			return fmt.Errorf(ErrF_ForumIsNotFound, parentId)
		}

		if parentForum.Threads.Size() == 0 {
			return fmt.Errorf(ErrF_ForumIsDamaged, parentId)
		}

		if !parentForum.Threads.HasItem(thread.Id) {
			return fmt.Errorf(ErrF_ForumIsDamaged, parentId)
		}
	}

	return nil
}

func checkMessages(dbo *dbo.DatabaseObject, threads []mm.ThreadLink) (err error) {
	// Step I. Downward check (parent to child).
	var messages []mm.MessageLink
	for _, thread := range threads {
		if thread.Messages.Size() == 0 {
			continue
		}

		messages, err = dbo.ReadMessageLinksById(thread.Messages)
		if err != nil {
			return err
		}

		messagesMap := make(map[cmb.Id]mm.MessageLink)
		for _, message := range messages {
			messagesMap[message.Id] = message
		}

		var ok bool
		var message mm.MessageLink
		for _, messageId := range thread.Messages.AsArray() {
			message, ok = messagesMap[messageId]
			if !ok {
				return fmt.Errorf(ErrF_MessageIsNotFound, messageId)
			}

			if message.ThreadId != thread.Id {
				return fmt.Errorf(ErrF_MessageIsDamaged, message.Id)
			}
		}
	}

	// Step II. Upward check (child to parent).
	// This kind of check requires huge amount of time.
	// It is not implemented.

	return nil
}
