package models

import (
	"database/sql"
	"errors"

	"github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

const (
	ChildTypeNone    = 0
	ChildTypeSection = 1
	ChildTypeForum   = 2
)

type Section struct {
	// Identifier of a section.
	Id cmb.Id `json:"id"`

	// Identifier of a parent section containing this section.
	// Null means that this section is a root section.
	// Only a single root section can exist.
	Parent *cmb.Id `json:"parent"`

	// Type of child elements: either sub-sections or forums.
	ChildType byte `json:"childType"`

	// List of IDs of child elements (either sub-sections or forums).
	// Null means that this section has no children.
	Children *ul.UidList `json:"children"`

	// Name of this section.
	Name cm.Name `json:"name"`

	// Section meta-data.
	EventData
}

func NewSection() (sec *Section) {
	return &Section{
		EventData: EventData{
			Creator: &EventParameters{},
			Editor:  &OptionalEventParameters{},
		},
	}
}

func NewSectionFromScannableSource(src cm.IScannable) (sec *Section, err error) {
	sec = NewSection()
	var x = ul.New()

	err = src.Scan(
		&sec.Id,
		&sec.Parent,
		&sec.ChildType,
		x, //&sec.Children,
		&sec.Name,
		&sec.Creator.UserId,
		&sec.Creator.Time,
		&sec.Editor.UserId,
		&sec.Editor.Time,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	sec.Children = x
	return sec, nil
}
