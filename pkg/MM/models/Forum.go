package models

import ul "github.com/vault-thirteen/SimpleBB/pkg/UidList"

type Forum struct {
	// Identifier of this forum.
	Id uint `json:"id"`

	// Identifier of a section containing this forum.
	SectionId uint `json:"sectionId"`

	// Name of this forum.
	Name string `json:"name"`

	// List of identifiers of threads of this forum.
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
