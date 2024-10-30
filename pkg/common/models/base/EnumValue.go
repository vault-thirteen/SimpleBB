package base

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

const EnumValueMin = 1

type EnumValue byte

func (ev EnumValue) ToString() string {
	return strconv.Itoa(int(ev.AsByte()))
}

func (ev EnumValue) AsByte() byte {
	return byte(ev)
}

func (ev EnumValue) AsInt() int {
	return int(ev)
}

func (ev *EnumValue) Scan(src any) (err error) {
	var x int
	x, err = ScanSrcAsInt(src)
	if err != nil {
		return err
	}

	// Type casting (int to byte) must be safe.
	if (x < 0) || (x > 255) {
		return fmt.Errorf(ErrF_Overflow, x)
	}
	var b = byte(x)

	*ev = EnumValue(b)
	return nil
}

func (ev EnumValue) Value() (dv driver.Value, err error) {
	// https://pkg.go.dev/database/sql/driver#Value
	return driver.Value(int64(ev)), nil
}

func (ev *EnumValue) UnmarshalJSON(data []byte) (err error) {
	var x int
	err = json.Unmarshal(data, &x)
	if err != nil {
		return err
	}

	// Type casting (int to byte) must be safe.
	if (x < 0) || (x > 255) {
		return fmt.Errorf(ErrF_Overflow, x)
	}
	var b = byte(x)

	*ev = EnumValue(b)
	return nil
}

func (ev EnumValue) MarshalJSON() (ba []byte, err error) {
	return []byte(ev.ToString()), nil
}
