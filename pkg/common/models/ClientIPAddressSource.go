package models

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/enum"
)

type ClientIPAddressSource = enum.Enum

const (
	ClientIPAddressSource_Direct       = enum.EnumValue(1)
	ClientIPAddressSource_CustomHeader = enum.EnumValue(2)

	ClientIPAddressSourceMax = ClientIPAddressSource_CustomHeader
)

func NewClientIPAddressSource() *ClientIPAddressSource {
	return enum.NewEnumFast(ClientIPAddressSourceMax)
}

func NewClientIPAddressSourceWithValue(value enum.EnumValue) ClientIPAddressSource {
	cipas := NewClientIPAddressSource()
	cipas.SetValueFast(value)
	return *cipas
}
