package models

import (
	ul "github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

type UserSubscriptions struct {
	// Subscription parameters. If pagination is used, these lists contain
	// information after the application of pagination.

	// Subscriber is an ID of a user. This is equivalent to UserId.
	Subscriber cmb.Id `json:"subscriber"`

	// Subscriptions is an array of IDs of threads to which a user is
	// subscribed. This is equivalent to ThreadIds.
	Subscriptions []cmb.Id `json:"subscriptions"`

	PageData *cmr.PageData `json:"pageData,omitempty"`
}

// NewUserSubscriptions is a constructor of UserSubscriptions. When pagination
// is not needed, pag number must be zero. When page number is positive,
// pagination is enabled.
func NewUserSubscriptions(userId cmb.Id, allThreadIds *ul.UidList, pageNumber cmb.Count, pageSize cmb.Count) (us *UserSubscriptions) {
	if pageNumber <= 0 {
		us = &UserSubscriptions{
			Subscriber:    userId,
			Subscriptions: allThreadIds.AsArray(),
			PageData:      nil,
		}
		return us
	}

	var threadIdsOnPage = allThreadIds.OnPage(pageNumber, pageSize)

	us = &UserSubscriptions{
		Subscriber:    userId,
		Subscriptions: threadIdsOnPage.AsArray(),
		PageData: &cmr.PageData{
			PageNumber:  pageNumber,
			TotalPages:  cmb.CalculateTotalPages(allThreadIds.Size(), pageSize),
			PageSize:    pageSize,
			ItemsOnPage: threadIdsOnPage.Size(),
			TotalItems:  allThreadIds.Size(),
		},
	}
	return us
}
