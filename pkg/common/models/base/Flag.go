package base

import (
	"database/sql/driver"
	"strconv"

	cms "github.com/vault-thirteen/SimpleBB/pkg/common/models/sql"
)

type Flag bool

func (f *Flag) Scan(src any) (err error) {
	var b bool
	b, err = cms.ScanSrcAsBoolean(src)
	if err != nil {
		return err
	}

	*f = Flag(b)
	return nil
}

func (f Flag) Value() (dv driver.Value, err error) {
	return driver.Value(cms.BoolToSql(f.AsBool())), nil
}

func (f Flag) ToString() string {
	return strconv.FormatBool(f.AsBool())
}

func (f Flag) RawValue() bool {
	return f.AsBool()
}

func (f Flag) AsBool() bool {
	return bool(f)
}
