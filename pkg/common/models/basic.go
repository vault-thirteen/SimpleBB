package models

import (
	"database/sql"
	"errors"
)

func NewArrayFromScannableSource[T any](src IScannableSequence) (values []T, err error) {
	values = make([]T, 0)
	var value T

	for src.Next() {
		err = src.Scan(&value)
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	if len(values) == 0 {
		values = nil
	}

	return values, nil
}

func NewValueFromScannableSource[T any](src IScannable) (*T, error) {
	var value T

	err := src.Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &value, nil
}
