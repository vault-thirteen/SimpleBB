package models

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
)

type ThreadAndMessages struct {
	// Thread parameters.
	ThreadId      uint   `json:"threadId"`
	ThreadForumId uint   `json:"threadForumId"`
	ThreadName    string `json:"threadName"`
	EventData

	// Message parameters. If pagination is used, these lists contain
	// information after the application of pagination.
	MessageIds *ul.UidList `json:"messageIds"`
	Messages   []Message   `json:"messages"`

	// Number of the current page of messages.
	Page *uint `json:"page,omitempty"`

	// Total number of available pages of messages.
	TotalPages *uint `json:"totalPages,omitempty"`

	// Total number of available messages.
	TotalMessages *uint `json:"totalMessages,omitempty"`
}

func NewThreadAndMessages(thread *Thread) (tam *ThreadAndMessages) {
	if thread == nil {
		return nil
	}

	tam = &ThreadAndMessages{
		ThreadId:      thread.Id,
		ThreadForumId: thread.ForumId,
		ThreadName:    thread.Name,

		EventData: EventData{
			Creator: &EventParameters{
				UserId: thread.Creator.UserId,
				Time:   thread.Creator.Time,
			},
			Editor: &OptionalEventParameters{
				UserId: thread.Editor.UserId,
				Time:   thread.Editor.Time,
			},
		},
	}

	return tam
}
