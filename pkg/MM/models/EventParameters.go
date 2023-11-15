package models

import "time"

type EventParameters struct {
	// ID of a user who initiated the event.
	UserId uint `json:"userId"`

	// Time of the event.
	Time time.Time `json:"time"`
}

type OptionalEventParameters struct {
	// ID of a user who initiated the event.
	UserId *uint `json:"userId"`

	// Time of the event.
	Time *time.Time `json:"time"`
}
