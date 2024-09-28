package models

import (
	"database/sql"
	"errors"
)

func NewStringFromScannableSource(src IScannable) (s *string, err error) {
	s = new(string)

	err = src.Scan(&s)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return s, nil
}
