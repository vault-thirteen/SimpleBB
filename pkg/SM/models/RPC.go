package models

import (
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

// Ping.

type PingParams = cmr.PingParams
type PingResult = cmr.PingResult

// Subscription.

type AddSubscriptionParams struct {
	cmr.CommonParams

	ThreadId uint `json:"threadId"`
	UserId   uint `json:"userId"`
}

type AddSubscriptionResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

type IsSelfSubscribedParams struct {
	cmr.CommonParams

	ThreadId uint `json:"threadId"`
}

type IsSelfSubscribedResult struct {
	cmr.CommonResult

	UserId       uint `json:"userId"`
	ThreadId     uint `json:"threadId"`
	IsSubscribed bool `json:"isSubscribed"`
}

type IsUserSubscribedParams struct {
	cmr.CommonParams

	UserId   uint `json:"userId"`
	ThreadId uint `json:"threadId"`
}

type IsUserSubscribedResult struct {
	cmr.CommonResult

	UserId       uint `json:"userId"`
	ThreadId     uint `json:"threadId"`
	IsSubscribed bool `json:"isSubscribed"`
}

type GetSelfSubscriptionsParams struct {
	cmr.CommonParams
}

type GetSelfSubscriptionsResult struct {
	cmr.CommonResult

	UserSubscriptions *UserSubscriptions `json:"userSubscriptions"`
}

type GetSelfSubscriptionsOnPageParams struct {
	cmr.CommonParams

	Page uint `json:"page"`
}

type GetSelfSubscriptionsOnPageResult struct {
	cmr.CommonResult

	SubscriptionsOnPage *SubscriptionsOnPage `json:"sop"`
}

type GetUserSubscriptionsParams struct {
	cmr.CommonParams

	UserId uint `json:"userId"`
}

type GetUserSubscriptionsResult struct {
	cmr.CommonResult

	UserSubscriptions *UserSubscriptions `json:"userSubscriptions"`
}

type GetThreadSubscribersSParams struct {
	cmr.CommonParams
	cmr.DKeyParams

	ThreadId uint `json:"threadId"`
}

type GetThreadSubscribersSResult struct {
	cmr.CommonResult

	ThreadSubscriptions *ThreadSubscriptions `json:"threadSubscriptions"`
}

type DeleteSelfSubscriptionParams struct {
	cmr.CommonParams

	ThreadId uint `json:"threadId"`
}

type DeleteSelfSubscriptionResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

type DeleteSubscriptionParams struct {
	cmr.CommonParams

	ThreadId uint `json:"threadId"`
	UserId   uint `json:"userId"`
}

type DeleteSubscriptionResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

type DeleteSubscriptionSParams struct {
	cmr.CommonParams
	cmr.DKeyParams

	ThreadId uint `json:"threadId"`
	UserId   uint `json:"userId"`
}

type DeleteSubscriptionSResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

type ClearThreadSubscriptionsSParams struct {
	cmr.CommonParams
	cmr.DKeyParams

	ThreadId uint `json:"threadId"`
}

type ClearThreadSubscriptionsSResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

// Other.

type GetDKeyParams struct {
	cmr.CommonParams
}

type GetDKeyResult struct {
	cmr.CommonResult

	DKey string `json:"dKey"`
}

type ShowDiagnosticDataParams struct{}

type ShowDiagnosticDataResult struct {
	cmr.CommonResult
	cmr.RequestsCount
}

type TestParams struct {
	N uint `json:"n"`
}

type TestResult struct {
	cmr.CommonResult
}
