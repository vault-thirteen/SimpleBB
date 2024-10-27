package base

import (
	"database/sql/driver"
	"strconv"
)

type Count int

func (c Count) ToString() string {
	return strconv.Itoa(c.AsInt())
}

func (c Count) AsInt() int {
	return int(c)
}

func (c *Count) Scan(src any) (err error) {
	var x int
	x, err = ScanSrcAsInt(src)
	if err != nil {
		return err
	}

	*c = Count(x)
	return nil
}

func (c Count) Value() (dv driver.Value, err error) {
	// https://pkg.go.dev/database/sql/driver#Value
	return driver.Value(int64(c)), nil
}
