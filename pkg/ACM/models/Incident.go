package models

import (
	"errors"
	"net"
	"time"
)

const (
	ErrIncidentIsNotSet  = "incident is not set"
	ErrIncidentTypeError = "incident type error"
)

type Incident struct {
	Id      uint
	Time    time.Time
	Type    IncidentType
	Email   string
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
