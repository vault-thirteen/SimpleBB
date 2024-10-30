package models

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type SectionChildType = cmb.Enum

const (
	SectionChildType_Section = cmb.EnumValue(1)
	SectionChildType_Forum   = cmb.EnumValue(2)

	SectionChildTypeMax = SectionChildType_Forum
)

func NewSectionChildType() *SectionChildType {
	return cmb.NewEnumFast(SectionChildTypeMax)
}

func NewSectionChildTypeWithValue(value cmb.EnumValue) SectionChildType {
	sct := NewSectionChildType()
	sct.SetValueFast(value)
	return *sct
}
