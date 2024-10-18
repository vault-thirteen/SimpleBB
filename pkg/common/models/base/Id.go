package base

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"

	num "github.com/vault-thirteen/auxie/number"
)

type Id int

func (i Id) ToString() string {
	return strconv.Itoa(i.AsInt())
}

func (i Id) AsInt() int {
	return int(i)
}

func (i *Id) Scan(src any) (err error) {
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

	*i = Id(x)
	return nil
}

func (i Id) Value() (dv driver.Value, err error) {
	// https://pkg.go.dev/database/sql/driver#Value
	return driver.Value(int64(i)), nil
}
