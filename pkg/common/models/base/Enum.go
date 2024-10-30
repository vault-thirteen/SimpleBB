package base

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	Err_MaxValue           = "max value error"
	ErrF_EnumValueOverflow = "enumeration value overflow: %v"
)

// Enum is a simple enumeration type supporting up to 255 distinct values.
// It is not recommended to use zero as a value. Zero is commonly used as a
// marker of a non-set value. The minimum value should always be equal to 1.
type Enum struct {
	value    EnumValue
	maxValue EnumValue
}

func NewEnum(maxValue EnumValue) (e *Enum, err error) {
	e = &Enum{maxValue: maxValue}

	err = e.checkMaxValue(maxValue)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func NewEnumFast(maxValue EnumValue) (e *Enum) {
	e = &Enum{maxValue: maxValue}

	err := e.checkMaxValue(maxValue)
	if err != nil {
		panic(err)
	}

	return e
}

func (e *Enum) checkMaxValue(maxValue EnumValue) (err error) {
	if maxValue < EnumValueMin {
		return errors.New(Err_MaxValue)
	}
	return nil
}

func (e *Enum) SetValue(value EnumValue) (err error) {
	err = e.checkValue(value)
	if err != nil {
		return err
	}

	e.value = value
	return nil
}

func (e *Enum) SetValueFast(value EnumValue) {
	err := e.checkValue(value)
	if err != nil {
		panic(err)
	}

	e.value = value
}

func (e *Enum) setValueJson(value EnumValue) (err error) {
	e.value = value
	return nil
}

func (e *Enum) checkValue(v EnumValue) (err error) {
	if (v < EnumValueMin) || (v > e.maxValue) {
		return fmt.Errorf(ErrF_EnumValueOverflow, v)
	}
	return nil
}

func (e *Enum) GetValue() EnumValue {
	return e.value
}

func (e *Enum) Scan(src any) (err error) {
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

	err = e.SetValue(EnumValue(b))
	if err != nil {
		return err
	}

	return nil
}

func (e Enum) ToString() string {
	return e.value.ToString()
}

func (e Enum) AsByte() byte {
	return e.value.AsByte()
}

func (e Enum) AsInt() int {
	return e.value.AsInt()
}

func (e Enum) Value() (dv driver.Value, err error) {
	return e.value.Value()
}

func (e *Enum) UnmarshalJSON(data []byte) (err error) {
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

	err = e.setValueJson(EnumValue(b))
	if err != nil {
		return err
	}

	return nil
}

func (e Enum) MarshalJSON() (ba []byte, err error) {
	return []byte(e.ToString()), nil
}
