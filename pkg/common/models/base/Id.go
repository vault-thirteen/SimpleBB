package base

import (
	"database/sql/driver"
	"strconv"
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
	x, err = ScanSrcAsInt(src)
	if err != nil {
		return err
	}

	*i = Id(x)
	return nil
}

func (i Id) Value() (dv driver.Value, err error) {
	// https://pkg.go.dev/database/sql/driver#Value
	return driver.Value(int64(i)), nil
}
