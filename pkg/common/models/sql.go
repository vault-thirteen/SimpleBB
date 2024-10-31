package models

import (
	"database/sql"
	"errors"
)

func NewArrayFromScannableSource[T any](src IScannableSequence) (values []T, err error) {
	values = []T{}
	var value T

	for src.Next() {
		err = src.Scan(&value)
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	return values, nil
}

func NewValueFromScannableSource[T any](src IScannable) (*T, error) {
	var value = new(T)

	err := src.Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return value, nil
}

func NewNonNullValueFromScannableSource[T any](src IScannable) (T, error) {
	var value T

	err := src.Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return value, nil
		} else {
			return value, err
		}
	}

	return value, nil
}
