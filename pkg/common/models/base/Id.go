package base

import (
	"database/sql/driver"
	"strconv"

	cms "github.com/vault-thirteen/SimpleBB/pkg/common/models/sql"
)

type Id int

func (i *Id) Scan(src any) (err error) {
	var x int
	x, err = cms.ScanSrcAsInt(src)
	if err != nil {
		return err
	}

	*i = Id(x)
	return nil
}

func (i Id) Value() (dv driver.Value, err error) {
	return driver.Value(cms.IntToSql(i.AsInt())), nil
}

func (i Id) ToString() string {
	return strconv.Itoa(i.AsInt())
}

func (i Id) RawValue() int {
	return i.AsInt()
}

func (i Id) AsInt() int {
	return int(i)
}
