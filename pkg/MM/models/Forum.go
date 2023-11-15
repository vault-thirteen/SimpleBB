package models

import ul "github.com/vault-thirteen/SimpleBB/pkg/UidList"

type Forum struct {
	Id uint `json:"id"`

	// ID of a parent forum containing this forum.
	// Null means that this forum is a root forum.
	// Only a single forum can be a root forum.
	Parent *uint `json:"parent"`

	// List of IDs of child forums of this forum.
	// Null means that this forum has no sub-forums.
	Children *ul.UidList `json:"children"`

	// Forum name.
	Name string `json:"name"`

	// List of IDs of threads of this forum.
	Threads *ul.UidList `json:"threads"`

	// Forum meta-data.
	EventData
}

func NewForum() (frm *Forum) {
	return &Forum{
		EventData: EventData{
			Creator: &EventParameters{},
			Editor:  &OptionalEventParameters{},
		},
	}
}
