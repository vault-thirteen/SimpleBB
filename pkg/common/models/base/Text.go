package base

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

type Text string

func (t Text) ToString() string {
	return string(t)
}

func (t Text) ToBytes() []byte {
	return []byte(t)
}

func (t *Text) Scan(src any) (err error) {
	var s string

	switch src.(type) {
	case string:
		s = src.(string)

	case []byte:
		s = string(src.([]byte))

	default:
		return fmt.Errorf(ErrF_UnsupportedDataType, reflect.TypeOf(src).String())
	}

	*t = Text(s)
	return nil
}

func (t Text) Value() (dv driver.Value, err error) {
	// https://pkg.go.dev/database/sql/driver#Value
	return driver.Value(t.ToBytes()), nil
}
