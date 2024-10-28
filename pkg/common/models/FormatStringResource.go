package models

import (
	"database/sql"
	"errors"
	"time"

	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

// FormatStringResource is a format string as a resource.
// This model is used for interaction with database.
type FormatStringResource struct {
	Id     cmb.Id       `json:"id"`
	Type   ResourceType `json:"type"`
	FSType *cmb.Text    `json:"fsType"`
	Text   *cmb.Text    `json:"text"`
	ToC    time.Time    `json:"toc"`
}

func NewFormatStringResource() (fs *FormatStringResource) {
	return &FormatStringResource{}
}

func NewFormatStringResourceFromScannableSource(src IScannable) (fsr *FormatStringResource, err error) {
	fsr = NewFormatStringResource()

	err = src.Scan(
		&fsr.Id,
		&fsr.Type,
		&fsr.FSType,
		&fsr.Text,
		&fsr.ToC,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return fsr, nil
}
