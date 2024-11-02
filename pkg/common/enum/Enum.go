package enum

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

// Enum is a simple enumeration type supporting up to 255 distinct values.
// It is not recommended to use zero as a value. Zero is commonly used as a
// marker of a non-set value. The minimum value should always be equal to 1.

// Limitations for enumeration value.
const (
	// EnumValueMin is a minimal value of an enumeration.
	EnumValueMin = 1

	// Maximum value is set in a child class. Oops. Golang does not know what
	// classes are. So, we emulate some parts of the class behaviour as we can.
)

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
	return e.value.Scan(src)
}

func (e Enum) Value() (dv driver.Value, err error) {
	return e.value.Value()
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

func (e *Enum) UnmarshalJSON(data []byte) (err error) {
	// Unfortunately we can not call the constructor of a base class from the
	// constructor of the child class, while Go language has no classes. We
	// should emulate classes using existing features of Go language ...

	err = e.value.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	//TODO:Here be dragons.
	err = e.checkValue(e.value)
	if err != nil {
		return err
	}

	return nil
}

func (e Enum) MarshalJSON() (ba []byte, err error) {
	return e.value.MarshalJSON()
}
