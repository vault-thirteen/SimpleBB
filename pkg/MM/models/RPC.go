package models

import (
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
	Parent *uint `json:"parent"`

	// Name of this section.
	Name string `json:"name"`
}

type AddSectionResult struct {
	cmr.CommonResult

	// ID of the created section.
	SectionId uint `json:"sectionId"`
}

type ChangeSectionNameParams struct {
	cmr.CommonParams

	// Identifier of a section.
	SectionId uint `json:"sectionId"`

	// Name of this section.
	Name string `json:"name"`
}

type ChangeSectionNameResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

type ChangeSectionParentParams struct {
	cmr.CommonParams

	// Identifier of a section.
	SectionId uint `json:"sectionId"`

	// Identifier of a parent section containing this section.
	Parent uint `json:"parent"`
}

type ChangeSectionParentResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

type GetSectionParams struct {
	cmr.CommonParams

	SectionId uint `json:"sectionId"`
}

type GetSectionResult struct {
	cmr.CommonResult

	Section *Section `json:"section"`
}

type DeleteSectionParams struct {
	cmr.CommonParams

	SectionId uint `json:"sectionId"`
}

type DeleteSectionResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

// Forum.

type AddForumParams struct {
	cmr.CommonParams

	// Identifier of a section containing this forum.
	SectionId uint `json:"sectionId"`

	// Name of this forum.
	Name string `json:"name"`
}

type AddForumResult struct {
	cmr.CommonResult

	// ID of the created forum.
	ForumId uint `json:"forumId"`
}

type ChangeForumNameParams struct {
	cmr.CommonParams

	ForumId uint `json:"forumId"`

	// New name.
	Name string `json:"name"`
}

type ChangeForumNameResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

type ChangeForumSectionParams struct {
	cmr.CommonParams

	// Identifier of this forum.
	ForumId uint `json:"forumId"`

	// Identifier of a section containing this forum.
	SectionId uint `json:"sectionId"`
}

type ChangeForumSectionResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

type GetForumParams struct {
	cmr.CommonParams

	ForumId uint `json:"forumId"`
}

type GetForumResult struct {
	cmr.CommonResult

	Forum *Forum `json:"forum"`
}

type DeleteForumParams struct {
	cmr.CommonParams

	ForumId uint `json:"forumId"`
}

type DeleteForumResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

// Thread.

type AddThreadParams struct {
	cmr.CommonParams

	// ID of a forum containing this thread.
	ForumId uint `json:"forumId"`

	// Thread name.
	Name string `json:"name"`
}

type AddThreadResult struct {
	cmr.CommonResult

	// ID of the created forum.
	ThreadId uint `json:"threadId"`
}

type ChangeThreadNameParams struct {
	cmr.CommonParams

	ThreadId uint `json:"threadId"`

	// New name.
	Name string `json:"name"`
}

type ChangeThreadNameResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

type ChangeThreadForumParams struct {
	cmr.CommonParams

	ThreadId uint `json:"threadId"`

	// ID of a new parent forum.
	ForumId uint `json:"forumId"`
}

type ChangeThreadForumResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

type GetThreadParams struct {
	cmr.CommonParams

	ThreadId uint `json:"threadId"`
}

type GetThreadResult struct {
	cmr.CommonResult

	Thread *Thread `json:"thread"`
}

type DeleteThreadParams struct {
	cmr.CommonParams

	ThreadId uint `json:"threadId"`
}

type DeleteThreadResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

// Message.

type AddMessageParams struct {
	cmr.CommonParams

	// ID of a thread containing this message.
	ThreadId uint `json:"threadId"`

	// Message text.
	Text string `json:"text"`
}

type AddMessageResult struct {
	cmr.CommonResult

	// ID of the created message.
	MessageId uint `json:"messageId"`
}

type ChangeMessageTextParams struct {
	cmr.CommonParams

	MessageId uint `json:"messageId"`

	// New text.
	Text string `json:"text"`
}

type ChangeMessageTextResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

type ChangeMessageThreadParams struct {
	cmr.CommonParams

	MessageId uint `json:"messageId"`

	// ID of a new parent thread.
	ThreadId uint `json:"threadId"`
}

type ChangeMessageThreadResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

type GetMessageParams struct {
	cmr.CommonParams

	MessageId uint `json:"messageId"`
}

type GetMessageResult struct {
	cmr.CommonResult

	Message *Message `json:"message"`
}

type DeleteMessageParams struct {
	cmr.CommonParams

	MessageId uint `json:"messageId"`
}

type DeleteMessageResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

// Composite objects.

type ListThreadAndMessagesParams struct {
	cmr.CommonParams

	ThreadId uint `json:"threadId"`
}

type ListThreadAndMessagesResult struct {
	cmr.CommonResult

	ThreadAndMessages *ThreadAndMessages `json:"tam"`
}

type ListThreadAndMessagesOnPageParams struct {
	cmr.CommonParams

	ThreadId uint `json:"threadId"`
	Page     uint `json:"page"`
}

type ListThreadAndMessagesOnPageResult struct {
	cmr.CommonResult

	ThreadAndMessagesOnPage *ThreadAndMessages `json:"tamop"`
}

type ListForumAndThreadsParams struct {
	cmr.CommonParams

	ForumId uint `json:"forumId"`
}

type ListForumAndThreadsResult struct {
	cmr.CommonResult

	ForumAndThreads *ForumAndThreads `json:"fat"`
}

type ListForumAndThreadsOnPageParams struct {
	cmr.CommonParams

	ForumId uint `json:"forumId"`
	Page    uint `json:"page"`
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
