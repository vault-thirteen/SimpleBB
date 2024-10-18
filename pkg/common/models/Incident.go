package models

import (
	"errors"
	"net"
	"time"

	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

const (
	ErrIncidentIsNotSet  = "incident is not set"
	ErrIncidentTypeError = "incident type error"
)

type Incident struct {
	Id      cmb.Id
	Module  Module
	Type    IncidentType
	Time    time.Time
	Email   Email
	UserIPA net.IP
}

func CheckIncident(inc *Incident) (err error) {
	if inc == nil {
		return errors.New(ErrIncidentIsNotSet)
	}

	if !inc.Type.IsValid() {
		return errors.New(ErrIncidentTypeError)
	}

	return nil
}
