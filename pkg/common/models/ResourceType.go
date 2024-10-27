package models

import (
	"database/sql/driver"
	"strconv"

	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type ResourceType byte

const (
	ResourceType_Text   = ResourceType(1)
	ResourceType_Number = ResourceType(2)

	ResourceTypeLast = ResourceType_Number
)

func (rt ResourceType) IsValid() bool {
	if (rt <= 0) || (rt > ResourceTypeLast) {
		return false
	}
	return true
}

func (rt ResourceType) ToString() string {
	return strconv.Itoa(int(rt.AsByte()))
}

func (rt ResourceType) AsByte() byte {
	return byte(rt)
}

func (rt *ResourceType) Scan(src any) (err error) {
	var x int
	x, err = cmb.ScanSrcAsInt(src)
	if err != nil {
		return err
	}

	*rt = ResourceType(x)
	return nil
}

func (rt ResourceType) Value() (dv driver.Value, err error) {
	// https://pkg.go.dev/database/sql/driver#Value
	return driver.Value(int64(rt)), nil
}
