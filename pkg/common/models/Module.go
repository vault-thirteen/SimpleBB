package models

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type Module = cmb.Enum

const (
	Module_ACM  = cmb.EnumValue(1)
	Module_GWM  = cmb.EnumValue(2)
	Module_MM   = cmb.EnumValue(3)
	Module_NM   = cmb.EnumValue(4)
	Module_RCS  = cmb.EnumValue(5)
	Module_SM   = cmb.EnumValue(6)
	Module_SMTP = cmb.EnumValue(7)

	ModuleMax = Module_SMTP
)

func NewModule() *Module {
	return cmb.NewEnumFast(ModuleMax)
}

func NewModuleWithValue(value cmb.EnumValue) Module {
	m := NewModule()
	m.SetValueFast(value)
	return *m
}
