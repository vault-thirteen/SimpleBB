package models

import (
	"errors"
	"net"
	"time"

	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

const (
	ErrIncidentIsNotSet = "incident is not set"
)

type Incident struct {
	Id      cmb.Id
	Module  Module
	Type    IncidentType
	Time    time.Time
	Email   Email
	UserIPA net.IP
}

func NewIncident() *Incident {
	return &Incident{
		Module: *NewModule(),
		Type:   *NewIncidentType(),
	}
}

func CheckIncident(inc *Incident) (err error) {
	if inc == nil {
		return errors.New(ErrIncidentIsNotSet)
	}

	return nil
}
