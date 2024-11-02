package models

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/enum"
)

type ResourceType = enum.Enum

const (
	ResourceType_Text   = enum.EnumValue(1)
	ResourceType_Number = enum.EnumValue(2)

	ResourceTypeMax = ResourceType_Number
)

func NewResourceType() *ResourceType {
	return enum.NewEnumFast(ResourceTypeMax)
}

func NewResourceTypeWithValue(value enum.EnumValue) ResourceType {
	rt := NewResourceType()
	rt.SetValueFast(value)
	return *rt
}
