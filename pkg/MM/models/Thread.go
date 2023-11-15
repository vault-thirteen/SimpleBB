package models

import ul "github.com/vault-thirteen/SimpleBB/pkg/UidList"

type Thread struct {
	Id uint `json:"id"`

	// ID of a forum containing this thread.
	ForumId uint `json:"forumId"`

	Name string `json:"name"`

	// List of IDs of messages of this thread.
	Messages *ul.UidList `json:"messages"`

	// Thread meta-data.
	EventData
}

func NewThread() (thr *Thread) {
	return &Thread{
		EventData: EventData{
			Creator: &EventParameters{},
			Editor:  &OptionalEventParameters{},
		},
	}
}
