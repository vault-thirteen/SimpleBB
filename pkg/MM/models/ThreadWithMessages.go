package models

import ul "github.com/vault-thirteen/SimpleBB/pkg/UidList"

type ThreadWithMessages struct {
	// Thread parameters.
	ThreadId      uint   `json:"threadId"`
	ThreadName    string `json:"threadName"`
	ThreadForumId uint   `json:"threadForumId"`
	EventData

	// Number of the current page of messages.
	Page *uint `json:"page,omitempty"`

	// Total number of available pages of messages.
	TotalPages *uint `json:"totalPages,omitempty"`

	// Total number of available messages.
	TotalMessages *uint `json:"totalMessages,omitempty"`

	// Message parameters. If pagination is used, these lists contain
	// information after the application of pagination.
	MessageIds ul.UidList `json:"messageIds"`
	Messages   []Message  `json:"messages"`
}

func NewThreadWithMessages(thread *Thread) (twm *ThreadWithMessages) {
	if thread == nil {
		return nil
	}

	twm = &ThreadWithMessages{
		ThreadId:      thread.Id,
		ThreadName:    thread.Name,
		ThreadForumId: thread.ForumId,

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

	return twm
}
