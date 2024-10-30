package models

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type ResourceType = cmb.Enum

const (
	ResourceType_Text   = cmb.EnumValue(1)
	ResourceType_Number = cmb.EnumValue(2)

	ResourceTypeMax = ResourceType_Number
)

func NewResourceType() *ResourceType {
	return cmb.NewEnumFast(ResourceTypeMax)
}

func NewResourceTypeWithValue(value cmb.EnumValue) ResourceType {
	rt := NewResourceType()
	rt.SetValueFast(value)
	return *rt
}
