package models

import (
	"database/sql"
	"errors"

	"github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type Forum struct {
	// Identifier of this forum.
	Id cmb.Id `json:"id"`

	// Identifier of a section containing this forum.
	SectionId cmb.Id `json:"sectionId"`

	// Name of this forum.
	Name cm.Name `json:"name"`

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
	var x = ul.New()

	err = src.Scan(
		&forum.Id,
		&forum.SectionId,
		&forum.Name,
		x, //&forum.Threads,
		&forum.Creator.UserId,
		&forum.Creator.Time,
		&forum.Editor.UserId,
		&forum.Editor.Time,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	forum.Threads = x
	return forum, nil
}
