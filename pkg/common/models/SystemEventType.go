package models

import (
	"database/sql/driver"
	"strconv"

	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type SystemEventType byte

const (
	SystemEventType_ThreadParentChange    = SystemEventType(1) // -> Users subscribed to the thread.
	SystemEventType_ThreadNameChange      = SystemEventType(2) // -> Users subscribed to the thread.
	SystemEventType_ThreadDeletion        = SystemEventType(3) // -> Users subscribed to the thread.
	SystemEventType_ThreadNewMessage      = SystemEventType(4) // -> Users subscribed to the thread.
	SystemEventType_ThreadMessageEdit     = SystemEventType(5) // -> Users subscribed to the thread.
	SystemEventType_ThreadMessageDeletion = SystemEventType(6) // -> Users subscribed to the thread.
	SystemEventType_MessageTextEdit       = SystemEventType(7) // -> Author of the message.
	SystemEventType_MessageParentChange   = SystemEventType(8) // -> Author of the message.
	SystemEventType_MessageDeletion       = SystemEventType(9) // -> Author of the message.

	SystemEventTypeLast = SystemEventType_MessageDeletion
)

func (set SystemEventType) IsValid() bool {
	if (set <= 0) || (set > SystemEventTypeLast) {
		return false
	}
	return true
}

func (set SystemEventType) ToString() string {
	return strconv.Itoa(int(set.AsByte()))
}

func (set SystemEventType) AsByte() byte {
	return byte(set)
}

func (set *SystemEventType) Scan(src any) (err error) {
	var x int
	x, err = cmb.ScanSrcAsInt(src)
	if err != nil {
		return err
	}

	*set = SystemEventType(x)
	return nil
}

func (set SystemEventType) Value() (dv driver.Value, err error) {
	// https://pkg.go.dev/database/sql/driver#Value
	return driver.Value(int64(set)), nil
}
