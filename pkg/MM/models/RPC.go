package models

import (
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

// Ping.

type PingParams = cmr.PingParams
type PingResult = cmr.PingResult

// Section.

type AddSectionParams struct {
	cmr.CommonParams

	// Identifier of a parent section containing this section.
	// Null means that this section is a root section.
	// Only a single root section can exist.
	Parent *cmb.Id `json:"parent"`

	// Name of this section.
	Name cm.Name `json:"name"`
}

type AddSectionResult struct {
	cmr.CommonResult

	// ID of the created section.
	SectionId cmb.Id `json:"sectionId"`
}

type ChangeSectionNameParams struct {
	cmr.CommonParams

	// Identifier of a section.
	SectionId cmb.Id `json:"sectionId"`

	// Name of this section.
	Name cm.Name `json:"name"`
}
type ChangeSectionNameResult = cmr.CommonResultWithSuccess

type ChangeSectionParentParams struct {
	cmr.CommonParams

	// Identifier of a section.
	SectionId cmb.Id `json:"sectionId"`

	// Identifier of a parent section containing this section.
	Parent cmb.Id `json:"parent"`
}
type ChangeSectionParentResult = cmr.CommonResultWithSuccess

type GetSectionParams struct {
	cmr.CommonParams

	SectionId cmb.Id `json:"sectionId"`
}

type GetSectionResult struct {
	cmr.CommonResult

	Section *Section `json:"section"`
}

type MoveSectionUpParams struct {
	cmr.CommonParams

	// Identifier of a section.
	SectionId cmb.Id `json:"sectionId"`
}
type MoveSectionUpResult = cmr.CommonResultWithSuccess

type MoveSectionDownParams struct {
	cmr.CommonParams

	// Identifier of a section.
	SectionId cmb.Id `json:"sectionId"`
}
type MoveSectionDownResult = cmr.CommonResultWithSuccess

type DeleteSectionParams struct {
	cmr.CommonParams

	SectionId cmb.Id `json:"sectionId"`
}
type DeleteSectionResult = cmr.CommonResultWithSuccess

// Forum.

type AddForumParams struct {
	cmr.CommonParams

	// Identifier of a section containing this forum.
	SectionId cmb.Id `json:"sectionId"`

	// Name of this forum.
	Name cm.Name `json:"name"`
}

type AddForumResult struct {
	cmr.CommonResult

	// ID of the created forum.
	ForumId cmb.Id `json:"forumId"`
}

type ChangeForumNameParams struct {
	cmr.CommonParams

	ForumId cmb.Id `json:"forumId"`

	// New name.
	Name cm.Name `json:"name"`
}
type ChangeForumNameResult = cmr.CommonResultWithSuccess

type ChangeForumSectionParams struct {
	cmr.CommonParams

	// Identifier of this forum.
	ForumId cmb.Id `json:"forumId"`

	// Identifier of a section containing this forum.
	SectionId cmb.Id `json:"sectionId"`
}
type ChangeForumSectionResult = cmr.CommonResultWithSuccess

type GetForumParams struct {
	cmr.CommonParams

	ForumId cmb.Id `json:"forumId"`
}

type GetForumResult struct {
	cmr.CommonResult

	Forum *Forum `json:"forum"`
}

type MoveForumUpParams struct {
	cmr.CommonParams

	// Identifier of a forum.
	ForumId cmb.Id `json:"forumId"`
}
type MoveForumUpResult = cmr.CommonResultWithSuccess

type MoveForumDownParams struct {
	cmr.CommonParams

	// Identifier of a forum.
	ForumId cmb.Id `json:"forumId"`
}
type MoveForumDownResult = cmr.CommonResultWithSuccess

type DeleteForumParams struct {
	cmr.CommonParams

	ForumId cmb.Id `json:"forumId"`
}
type DeleteForumResult = cmr.CommonResultWithSuccess

// Thread.

type AddThreadParams struct {
	cmr.CommonParams

	// ID of a forum containing this thread.
	ForumId cmb.Id `json:"forumId"`

	// Thread name.
	Name cm.Name `json:"name"`
}

type AddThreadResult struct {
	cmr.CommonResult

	// ID of the created forum.
	ThreadId cmb.Id `json:"threadId"`
}

type ChangeThreadNameParams struct {
	cmr.CommonParams

	ThreadId cmb.Id `json:"threadId"`

	// New name.
	Name cm.Name `json:"name"`
}
type ChangeThreadNameResult = cmr.CommonResultWithSuccess

type ChangeThreadForumParams struct {
	cmr.CommonParams

	ThreadId cmb.Id `json:"threadId"`

	// ID of a new parent forum.
	ForumId cmb.Id `json:"forumId"`
}
type ChangeThreadForumResult = cmr.CommonResultWithSuccess

type GetThreadParams struct {
	cmr.CommonParams

	ThreadId cmb.Id `json:"threadId"`
}

type GetThreadResult struct {
	cmr.CommonResult

	Thread *Thread `json:"thread"`
}

type GetThreadNamesByIdsParams struct {
	cmr.CommonParams

	ThreadIds []cmb.Id `json:"threadIds"`
}

type GetThreadNamesByIdsResult struct {
	cmr.CommonResult

	ThreadIds   []cmb.Id  `json:"threadIds"`
	ThreadNames []cm.Name `json:"threadNames"`
}

type MoveThreadUpParams struct {
	cmr.CommonParams

	// Identifier of a thread.
	ThreadId cmb.Id `json:"threadId"`
}
type MoveThreadUpResult = cmr.CommonResultWithSuccess

type MoveThreadDownParams struct {
	cmr.CommonParams

	// Identifier of a thread.
	ThreadId cmb.Id `json:"threadId"`
}
type MoveThreadDownResult = cmr.CommonResultWithSuccess

type DeleteThreadParams struct {
	cmr.CommonParams

	ThreadId cmb.Id `json:"threadId"`
}
type DeleteThreadResult = cmr.CommonResultWithSuccess

type ThreadExistsSParams struct {
	cmr.CommonParams
	cmr.DKeyParams

	ThreadId cmb.Id `json:"threadId"`
}

type ThreadExistsSResult struct {
	cmr.CommonResult

	Exists cmb.Flag `json:"exists"`
}

// Message.

type AddMessageParams struct {
	cmr.CommonParams

	// ID of a thread containing this message.
	ThreadId cmb.Id `json:"threadId"`

	// Message text.
	Text cmb.Text `json:"text"`
}

type AddMessageResult struct {
	cmr.CommonResult

	// ID of the created message.
	MessageId cmb.Id `json:"messageId"`
}

type ChangeMessageTextParams struct {
	cmr.CommonParams

	MessageId cmb.Id `json:"messageId"`

	// New text.
	Text cmb.Text `json:"text"`
}
type ChangeMessageTextResult = cmr.CommonResultWithSuccess

type ChangeMessageThreadParams struct {
	cmr.CommonParams

	MessageId cmb.Id `json:"messageId"`

	// ID of a new parent thread.
	ThreadId cmb.Id `json:"threadId"`
}
type ChangeMessageThreadResult = cmr.CommonResultWithSuccess

type GetMessageParams struct {
	cmr.CommonParams

	MessageId cmb.Id `json:"messageId"`
}

type GetMessageResult struct {
	cmr.CommonResult

	Message *Message `json:"message"`
}

type GetLatestMessageOfThreadParams struct {
	cmr.CommonParams

	ThreadId cmb.Id `json:"threadId"`
}

type GetLatestMessageOfThreadResult struct {
	cmr.CommonResult

	Message *Message `json:"message"`
}

type DeleteMessageParams struct {
	cmr.CommonParams

	MessageId cmb.Id `json:"messageId"`
}
type DeleteMessageResult = cmr.CommonResultWithSuccess

// Composite objects.

type ListThreadAndMessagesParams struct {
	cmr.CommonParams

	ThreadId cmb.Id `json:"threadId"`
}

type ListThreadAndMessagesResult struct {
	cmr.CommonResult

	ThreadAndMessages *ThreadAndMessages `json:"tam"`
}

type ListThreadAndMessagesOnPageParams struct {
	cmr.CommonParams

	ThreadId cmb.Id    `json:"threadId"`
	Page     cmb.Count `json:"page"`
}

type ListThreadAndMessagesOnPageResult struct {
	cmr.CommonResult

	ThreadAndMessagesOnPage *ThreadAndMessages `json:"tamop"`
}

type ListForumAndThreadsParams struct {
	cmr.CommonParams

	ForumId cmb.Id `json:"forumId"`
}

type ListForumAndThreadsResult struct {
	cmr.CommonResult

	ForumAndThreads *ForumAndThreads `json:"fat"`
}

type ListForumAndThreadsOnPageParams struct {
	cmr.CommonParams

	ForumId cmb.Id    `json:"forumId"`
	Page    cmb.Count `json:"page"`
}

type ListForumAndThreadsOnPageResult struct {
	cmr.CommonResult

	ForumAndThreadsOnPage *ForumAndThreads `json:"fatop"`
}

type ListSectionsAndForumsParams struct {
	cmr.CommonParams
}

type ListSectionsAndForumsResult struct {
	cmr.CommonResult

	SectionsAndForums *SectionsAndForums `json:"saf"`
}

// Other.

type GetDKeyParams struct {
	cmr.CommonParams
}

type GetDKeyResult struct {
	cmr.CommonResult

	DKey cmb.Text `json:"dKey"`
}

type ShowDiagnosticDataParams struct{}

type ShowDiagnosticDataResult struct {
	cmr.CommonResult
	cmr.RequestsCount
}

type TestParams struct {
	N uint `json:"n"`
}

type TestResult struct {
	cmr.CommonResult
}
