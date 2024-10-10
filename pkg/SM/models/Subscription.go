package models

type Subscription struct {
	// ID of a thread to which a user is subscribed.
	ThreadId uint `json:"threadId"`

	// ID of a subscribed user.
	UserId uint `json:"userId"`
}
