package models

import (
	"database/sql"
	"errors"
	"time"

	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type Resource struct {
	Id     cmb.Id
	Type   ResourceType
	Text   *cmb.Text
	Number *cmb.Count
	ToC    time.Time
}

func NewResource() (r *Resource) {
	return &Resource{}
}

func NewResourceFromScannableSource(src IScannable) (r *Resource, err error) {
	r = NewResource()

	err = src.Scan(
		&r.Id,
		&r.Type,
		&r.Text,
		&r.Number,
		&r.ToC,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return r, nil
}

func (r *Resource) IsValid() bool {
	if r == nil {
		return false
	}

	if !r.Type.IsValid() {
		return false
	}

	switch r.Type {
	case ResourceType_Text:
		if r.Text == nil {
			return false
		}

	case ResourceType_Number:
		if r.Number == nil {
			return false
		}

	default:
		return false
	}

	return true
}

func (r *Resource) AsText() cmb.Text {
	return *r.Text
}

func (r *Resource) AsNumber() cmb.Count {
	return *r.Number
}

func (r *Resource) Value() any {
	switch r.Type {
	case ResourceType_Text:
		return *r.Text

	case ResourceType_Number:
		return *r.Number

	default:
		return nil
	}
}
