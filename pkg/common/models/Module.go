package models

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/enum"
)

type Module = enum.Enum

const (
	Module_ACM  = enum.EnumValue(1)
	Module_GWM  = enum.EnumValue(2)
	Module_MM   = enum.EnumValue(3)
	Module_NM   = enum.EnumValue(4)
	Module_RCS  = enum.EnumValue(5)
	Module_SM   = enum.EnumValue(6)
	Module_SMTP = enum.EnumValue(7)

	ModuleMax = Module_SMTP
)

func NewModule() *Module {
	return enum.NewEnumFast(ModuleMax)
}

func NewModuleWithValue(value enum.EnumValue) Module {
	m := NewModule()
	m.SetValueFast(value)
	return *m
}
