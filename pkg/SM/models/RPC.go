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