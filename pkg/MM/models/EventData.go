package models

type EventData struct {
	// Parameters of creation.
	Creator *EventParameters `json:"creator"`

	// Parameters of the last edit.
	Editor *OptionalEventParameters `json:"editor"`
}
