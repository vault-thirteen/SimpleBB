package models

import (
	"database/sql"
	"errors"
	"time"

	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

// RegistrationReadyForApproval is a short variant of an object representing a
// request for registration which is ready to be approved.
type RegistrationReadyForApproval struct {
	Id         cmb.Id    `json:"id"`
	PreRegTime time.Time `json:"preRegTime"`
	Email      cm.Email  `json:"email"`
	Name       *cm.Name  `json:"name"`
}

func NewRegistrationReadyForApproval() (r *RegistrationReadyForApproval) {
	return &RegistrationReadyForApproval{}
}

func NewRegistrationReadyForApprovalFromScannableSource(src cm.IScannable) (r *RegistrationReadyForApproval, err error) {
	r = NewRegistrationReadyForApproval()

	err = src.Scan(
		&r.Id,
		&r.PreRegTime,
		&r.Email,
		&r.Name,
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

func NewRegistrationReadyForApprovalArrayFromRows(rows cm.IScannableSequence) (rs []RegistrationReadyForApproval, err error) {
	rs = []RegistrationReadyForApproval{}
	var r *RegistrationReadyForApproval

	for rows.Next() {
		r, err = NewRegistrationReadyForApprovalFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		rs = append(rs, *r)
	}

	return rs, nil
}
