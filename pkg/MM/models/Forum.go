package models

import (
	ul "github.com/vault-thirteen/SimpleBB/pkg/UidList"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

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

func NewForumFromScannableSource(src cm.IScannable) (forum *Forum, err error) {
	forum = NewForum()

	err = src.Scan(
		&forum.Id,
		&forum.SectionId,
		&forum.Name,
		&forum.Threads,
		&forum.Creator.UserId,
		&forum.Creator.Time,
		&forum.Editor.UserId,
		&forum.Editor.Time,
	)
	if err != nil {
		return nil, err
	}

	return forum, nil
}
