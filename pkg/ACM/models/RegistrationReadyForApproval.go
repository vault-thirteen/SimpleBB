package models

import (
	"database/sql"
	"errors"

	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

// RegistrationReadyForApproval is a short variant of an object representing a
// request for registration which is ready to be approved.
type RegistrationReadyForApproval struct {
	Id         uint   `json:"id"`
	PreRegTime string `json:"preRegTime"`
	Email      string `json:"email"`
	Name       string `json:"name"`
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
