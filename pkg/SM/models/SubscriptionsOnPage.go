package models

import (
	ul "github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

type SubscriptionsOnPage struct {
	// Subscription parameters. If pagination is used, these lists contain
	// information after the application of pagination.

	// Subscriber is a UserId.
	Subscriber cmb.Id `json:"subscriber"`

	// Subscriptions is an array of ThreadId's.
	Subscriptions *ul.UidList `json:"subscriptions"`

	PageData *cmr.PageData `json:"pageData,omitempty"`
}

func NewSubscriptionsOnPage() (sop *SubscriptionsOnPage) {
	sop = &SubscriptionsOnPage{}
	return sop
}
