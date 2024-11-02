package models

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/enum"
)

type SectionChildType = enum.Enum

const (
	SectionChildType_Section = enum.EnumValue(1)
	SectionChildType_Forum   = enum.EnumValue(2)

	SectionChildTypeMax = SectionChildType_Forum
)

func NewSectionChildType() *SectionChildType {
	return enum.NewEnumFast(SectionChildTypeMax)
}

func NewSectionChildTypeWithValue(value enum.EnumValue) SectionChildType {
	sct := NewSectionChildType()
	sct.SetValueFast(value)
	return *sct
}
