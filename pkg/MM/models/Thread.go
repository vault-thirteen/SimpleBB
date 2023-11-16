package models

import ul "github.com/vault-thirteen/SimpleBB/pkg/UidList"

type Thread struct {
	// Identifier of this thread.
	Id uint `json:"id"`

	// Identifier of a forum containing this thread.
	ForumId uint `json:"forumId"`

	// Name of this thread.
	Name string `json:"name"`

	// List of identifiers of messages of this thread.
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
