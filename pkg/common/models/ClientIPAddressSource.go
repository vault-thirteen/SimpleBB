package models

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type ClientIPAddressSource = cmb.Enum

const (
	ClientIPAddressSource_Direct       = cmb.EnumValue(1)
	ClientIPAddressSource_CustomHeader = cmb.EnumValue(2)

	ClientIPAddressSourceMax = ClientIPAddressSource_CustomHeader
)

func NewClientIPAddressSource() *ClientIPAddressSource {
	return cmb.NewEnumFast(ClientIPAddressSourceMax)
}

func NewClientIPAddressSourceWithValue(value cmb.EnumValue) ClientIPAddressSource {
	cipas := NewClientIPAddressSource()
	cipas.SetValueFast(value)
	return *cipas
}
