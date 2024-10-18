package base

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"

	num "github.com/vault-thirteen/auxie/number"
)

type Count int

func (c Count) ToString() string {
	return strconv.Itoa(c.AsInt())
}

func (c Count) AsInt() int {
	return int(c)
}

func (c *Count) Scan(src any) (err error) {
	var i int

	switch src.(type) {
	case int64:
		i = int(src.(int64))

	case []byte:
		i, err = num.ParseInt(string(src.([]byte)))
		if err != nil {
			return err
		}

	case string:
		i, err = num.ParseInt(src.(string))
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf(ErrF_UnsupportedDataType, reflect.TypeOf(src).String())
	}

	*c = Count(i)
	return nil
}

func (c Count) Value() (dv driver.Value, err error) {
	// https://pkg.go.dev/database/sql/driver#Value
	return driver.Value(int64(c)), nil
}
