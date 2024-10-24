package models

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"

	num "github.com/vault-thirteen/auxie/number"
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

const (
	ErrF_UnsupportedDataType = "unsupported data type: %s"
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

	switch src.(type) {
	case int64:
		x = int(src.(int64))

	case []byte:
		x, err = num.ParseInt(string(src.([]byte)))
		if err != nil {
			return err
		}

	case string:
		x, err = num.ParseInt(src.(string))
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf(ErrF_UnsupportedDataType, reflect.TypeOf(src).String())
	}

	*set = SystemEventType(x)
	return nil
}

func (set SystemEventType) Value() (dv driver.Value, err error) {
	// https://pkg.go.dev/database/sql/driver#Value
	return driver.Value(int64(set)), nil
}
