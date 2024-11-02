package enum

import (
	"database/sql/driver"
	"strconv"

	cmj "github.com/vault-thirteen/SimpleBB/pkg/common/models/json"
	cms "github.com/vault-thirteen/SimpleBB/pkg/common/models/sql"
)

// EnumValue is a raw type of the enumeration. It does not restrict the value.
// Validation of the value is performed by the 'Enum' class according to it's
// settings.
type EnumValue byte

func (ev *EnumValue) Scan(src any) (err error) {
	var b byte
	b, err = cms.ScanSrcAsByte(src)
	if err != nil {
		return err
	}

	*ev = EnumValue(b)
	return nil
}

func (ev EnumValue) Value() (dv driver.Value, err error) {
	// https://pkg.go.dev/database/sql/driver#Value
	return driver.Value(cms.ByteToSql(ev.AsByte())), nil
}

func (ev EnumValue) ToString() string {
	return strconv.Itoa(ev.AsInt())
}

func (ev EnumValue) RawValue() byte {
	return ev.AsByte()
}

func (ev EnumValue) AsByte() byte {
	return byte(ev)
}

func (ev EnumValue) AsInt() int {
	return int(ev)
}

func (ev *EnumValue) UnmarshalJSON(data []byte) (err error) {
	var b byte
	b, err = cmj.UnmarshalAsByte(data)
	if err != nil {
		return err
	}

	*ev = EnumValue(b)
	return nil
}

func (ev EnumValue) MarshalJSON() (ba []byte, err error) {
	return cmj.ToJson(ev), nil
}
