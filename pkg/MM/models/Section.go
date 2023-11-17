package models

import ul "github.com/vault-thirteen/SimpleBB/pkg/UidList"

const (
	ChildTypeNone    = 0
	ChildTypeSection = 1
	ChildTypeForum   = 2
)

type Section struct {
	// Identifier of a section.
	Id uint `json:"id"`

	// Identifier of a parent section containing this section.
	// Null means that this section is a root section.
	// Only a single root section can exist.
	Parent *uint `json:"parent"`

	// Type of child elements: either sub-sections or forums.
	ChildType byte `json:"childType"`

	// List of IDs of child elements (either sub-sections or forums).
	// Null means that this section has no children.
	Children *ul.UidList `json:"children"`

	// Name of this section.
	Name string `json:"name"`

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
