package models

import (
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

// Ping.

type PingParams = cm.PingParams
type PingResult = cm.PingResult

// Section.

type AddSectionParams struct {
	cm.CommonParams

	// Identifier of a parent section containing this section.
	// Null means that this section is a root section.
	// Only a single root section can exist.
	Parent *uint `json:"parent"`

	// Name of this section.
	Name string `json:"name"`
}

type AddSectionResult struct {
	cm.CommonResult

	// ID of the created section.
	SectionId uint `json:"sectionId"`
}

type ChangeSectionNameParams struct {
	cm.CommonParams

	// Identifier of a section.
	SectionId uint `json:"sectionId"`

	// Name of this section.
	Name string `json:"name"`
}

type ChangeSectionNameResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

type ChangeSectionParentParams struct {
	cm.CommonParams

	// Identifier of a section.
	SectionId uint `json:"sectionId"`

	// Identifier of a parent section containing this section.
	Parent uint `json:"parent"`
}

type ChangeSectionParentResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

type GetSectionParams struct {
	cm.CommonParams

	SectionId uint `json:"sectionId"`
}

type GetSectionResult struct {
	cm.CommonResult

	Section *Section `json:"section"`
}

type DeleteSectionParams struct {
	cm.CommonParams

	SectionId uint `json:"sectionId"`
}

type DeleteSectionResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

// Forum.

type AddForumParams struct {
	cm.CommonParams

	// Identifier of a section containing this forum.
	SectionId uint `json:"sectionId"`

	// Name of this forum.
	Name string `json:"name"`
}

type AddForumResult struct {
	cm.CommonResult

	// ID of the created forum.
	ForumId uint `json:"forumId"`
}

type ChangeForumNameParams struct {
	cm.CommonParams

	ForumId uint `json:"forumId"`

	// New name.
	Name string `json:"name"`
}

type ChangeForumNameResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

type ChangeForumSectionParams struct {
	cm.CommonParams

	// Identifier of this forum.
	ForumId uint `json:"forumId"`

	// Identifier of a section containing this forum.
	SectionId uint `json:"sectionId"`
}

type ChangeForumSectionResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

type GetForumParams struct {
	cm.CommonParams

	ForumId uint `json:"forumId"`
}

type GetForumResult struct {
	cm.CommonResult

	Forum *Forum `json:"forum"`
}

type DeleteForumParams struct {
	cm.CommonParams

	ForumId uint `json:"forumId"`
}

type DeleteForumResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

// Thread.

type AddThreadParams struct {
	cm.CommonParams

	// ID of a forum containing this thread.
	ForumId uint `json:"forumId"`

	// Thread name.
	Name string `json:"name"`
}

type AddThreadResult struct {
	cm.CommonResult

	// ID of the created forum.
	ThreadId uint `json:"threadId"`
}

type ChangeThreadNameParams struct {
	cm.CommonParams

	ThreadId uint `json:"threadId"`

	// New name.
	Name string `json:"name"`
}

type ChangeThreadNameResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

type ChangeThreadForumParams struct {
	cm.CommonParams

	ThreadId uint `json:"threadId"`

	// ID of a new parent forum.
	ForumId uint `json:"forumId"`
}

type ChangeThreadForumResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

type GetThreadParams struct {
	cm.CommonParams

	ThreadId uint `json:"threadId"`
}

type GetThreadResult struct {
	cm.CommonResult

	Thread *Thread `json:"thread"`
}

type DeleteThreadParams struct {
	cm.CommonParams

	ThreadId uint `json:"threadId"`
}

type DeleteThreadResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

// Message.

type AddMessageParams struct {
	cm.CommonParams

	// ID of a thread containing this message.
	ThreadId uint `json:"threadId"`

	// Message text.
	Text string `json:"text"`
}

type AddMessageResult struct {
	cm.CommonResult

	// ID of the created message.
	MessageId uint `json:"messageId"`
}

type ChangeMessageTextParams struct {
	cm.CommonParams

	MessageId uint `json:"messageId"`

	// New text.
	Text string `json:"text"`
}

type ChangeMessageTextResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

type ChangeMessageThreadParams struct {
	cm.CommonParams

	MessageId uint `json:"messageId"`

	// ID of a new parent thread.
	ThreadId uint `json:"threadId"`
}

type ChangeMessageThreadResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

type GetMessageParams struct {
	cm.CommonParams

	MessageId uint `json:"messageId"`
}

type GetMessageResult struct {
	cm.CommonResult

	Message *Message `json:"message"`
}

type DeleteMessageParams struct {
	cm.CommonParams

	MessageId uint `json:"messageId"`
}

type DeleteMessageResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

// Composite objects.

type ListThreadAndMessagesParams struct {
	cm.CommonParams

	ThreadId uint `json:"threadId"`
}

type ListThreadAndMessagesResult struct {
	cm.CommonResult

	ThreadAndMessages *ThreadAndMessages `json:"tam"`
}

type ListThreadAndMessagesOnPageParams struct {
	cm.CommonParams

	ThreadId uint `json:"threadId"`
	Page     uint `json:"page"`
}

type ListThreadAndMessagesOnPageResult struct {
	cm.CommonResult

	ThreadAndMessagesOnPage *ThreadAndMessages `json:"tamop"`
}

type ListForumAndThreadsParams struct {
	cm.CommonParams

	ForumId uint `json:"forumId"`
}

type ListForumAndThreadsResult struct {
	cm.CommonResult

	ForumAndThreads *ForumAndThreads `json:"fat"`
}

type ListForumAndThreadsOnPageParams struct {
	cm.CommonParams

	ForumId uint `json:"forumId"`
	Page    uint `json:"page"`
}

type ListForumAndThreadsOnPageResult struct {
	cm.CommonResult

	ForumAndThreadsOnPage *ForumAndThreads `json:"fatop"`
}

type ListSectionsAndForumsParams struct {
	cm.CommonParams
}

type ListSectionsAndForumsResult struct {
	cm.CommonResult

	SectionsAndForums *SectionsAndForums `json:"saf"`
}

// Other.

type ShowDiagnosticDataParams struct{}

type ShowDiagnosticDataResult struct {
	cm.CommonResult

	TotalRequestsCount      uint64 `json:"totalRequestsCount"`
	SuccessfulRequestsCount uint64 `json:"successfulRequestsCount"`
}

type TestParams struct {
	N uint `json:"n"`
}

type TestResult struct {
	cm.CommonResult
}
