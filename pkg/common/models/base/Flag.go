package base

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
)

type Flag bool

func (f Flag) Bool() bool {
	return bool(f)
}

func (f *Flag) Scan(src any) (err error) {
	var b bool

	switch src.(type) {
	case int64:
		x := src.(int64)
		if x == 0 {
			b = false
		} else if x == 1 {
			b = true
		} else {
			return fmt.Errorf(ErrF_UnsupportedDataValue, src)
		}

	case bool:
		b = src.(bool)

	case []byte:
		b, err = strconv.ParseBool(string(src.([]byte)))
		if err != nil {
			return err
		}

	case string:
		b, err = strconv.ParseBool(src.(string))
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf(ErrF_UnsupportedDataType, reflect.TypeOf(src).String())
	}

	*f = Flag(b)
	return nil
}

func (f Flag) Value() (dv driver.Value, err error) {
	// https://pkg.go.dev/database/sql/driver#Value
	return driver.Value(bool(f)), nil
}
