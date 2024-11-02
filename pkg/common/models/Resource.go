package models

import (
	"database/sql"
	"errors"
	"time"

	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/models/interfaces"
)

type Resource struct {
	Id     cmb.Id       `json:"id"`
	Type   ResourceType `json:"type"`
	Text   *cmb.Text    `json:"text"`
	Number *cmb.Count   `json:"number"`
	ToC    time.Time    `json:"toc"`
}

func NewResource() (r *Resource) {
	return &Resource{
		Type: *NewResourceType(),
	}
}

func NewResourceFromValue(value any) (r *Resource) {
	switch value.(type) {
	case float64:
		return NewResourceFromNumber(cmb.Count(value.(float64)))

	case int64:
		return NewResourceFromNumber(cmb.Count(value.(int64)))

	case int:
		return NewResourceFromNumber(cmb.Count(value.(int)))

	case []byte:
		return NewResourceFromText(cmb.Text(value.([]byte)))

	case string:
		return NewResourceFromText(cmb.Text(value.(string)))
	}

	return nil
}

func NewResourceFromText(t cmb.Text) (r *Resource) {
	return &Resource{
		Type: NewResourceTypeWithValue(ResourceType_Text),
		Text: &t,
	}
}

func NewResourceFromNumber(n cmb.Count) (r *Resource) {
	return &Resource{
		Type:   NewResourceTypeWithValue(ResourceType_Number),
		Number: &n,
	}
}

func NewResourceFromScannableSource(src cmi.IScannable) (r *Resource, err error) {
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

	switch r.Type.GetValue() {
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

func (r *Resource) GetValue() any {
	switch r.Type.GetValue() {
	case ResourceType_Text:
		return *r.Text

	case ResourceType_Number:
		return *r.Number

	default:
		return nil
	}
}
