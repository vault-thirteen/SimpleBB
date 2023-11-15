package models

import (
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

type PingParams struct{}

type PingResult struct {
	OK bool `json:"ok"`
}

type AddForumParams struct {
	cm.CommonParams

	// ID of a parent forum containing this forum.
	// Null means that this forum is a root forum.
	// Only a single forum can be a root forum.
	Parent *uint `json:"parent"`

	// Forum name.
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

type ChangeForumParentParams struct {
	cm.CommonParams

	ForumId uint `json:"forumId"`

	// ID of a new parent forum.
	Parent uint `json:"parent"`
}

type ChangeForumParentResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

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

type DeleteMessageParams struct {
	cm.CommonParams

	MessageId uint `json:"messageId"`
}

type DeleteMessageResult struct {
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

type ListThreadMessagesParams struct {
	cm.CommonParams

	ThreadId uint `json:"threadId"`
}

type ListThreadMessagesResult struct {
	cm.CommonResult

	ThreadWithMessages *ThreadWithMessages `json:"twm"`
}

type ListThreadMessagesOnPageParams struct {
	cm.CommonParams

	ThreadId uint `json:"threadId"`
	Page     uint `json:"page"`
}

type ListThreadMessagesOnPageResult struct {
	cm.CommonResult

	ThreadWithMessages *ThreadWithMessages `json:"twmop"`
}

type ListForumThreadsParams struct {
	cm.CommonParams

	ForumId uint `json:"forumId"`
}

type ListForumThreadsResult struct {
	cm.CommonResult

	ForumWithThreads *ForumWithThreads `json:"fwt"`
}

type ListForumThreadsOnPageParams struct {
	cm.CommonParams

	ForumId uint `json:"forumId"`
	Page    uint `json:"page"`
}

type ListForumThreadsOnPageResult struct {
	cm.CommonResult

	ForumWithThreads *ForumWithThreads `json:"fwtop"`
}

type ListForumsParams struct {
	cm.CommonParams
}

type ListForumsResult struct {
	cm.CommonResult

	Forums []Forum `json:"forums"`
}

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
