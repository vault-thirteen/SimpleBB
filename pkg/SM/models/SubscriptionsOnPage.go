package models

type SubscriptionsOnPage struct {
	// Subscription parameters. If pagination is used, these lists contain
	// information after the application of pagination.

	// Subscriber is a UserId.
	Subscriber uint `json:"subscriber"`

	// Subscriptions is an array of ThreadId's.
	Subscriptions []uint `json:"subscriptions"`

	// Number of the current page of subscriptions.
	Page *uint `json:"page,omitempty"`

	// Total number of available pages of subscriptions.
	TotalPages *uint `json:"totalPages,omitempty"`

	// Total number of available subscriptions.
	TotalSubscriptions *uint `json:"totalSubscriptions,omitempty"`
}

func NewSubscriptionsOnPage() (sop *SubscriptionsOnPage) {
	sop = &SubscriptionsOnPage{}
	return sop
}
