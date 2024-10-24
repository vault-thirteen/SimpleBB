package models

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

// Ping.

type PingParams = cmr.PingParams
type PingResult = cmr.PingResult

// Subscription.

type AddSubscriptionParams struct {
	cmr.CommonParams
	ThreadId cmb.Id `json:"threadId"`
	UserId   cmb.Id `json:"userId"`
}
type AddSubscriptionResult = cmr.CommonResultWithSuccess

type IsSelfSubscribedParams struct {
	cmr.CommonParams
	ThreadId cmb.Id `json:"threadId"`
}
type IsSelfSubscribedResult struct {
	cmr.CommonResult
	UserId       cmb.Id   `json:"userId"`
	ThreadId     cmb.Id   `json:"threadId"`
	IsSubscribed cmb.Flag `json:"isSubscribed"`
}

type IsUserSubscribedParams struct {
	cmr.CommonParams
	UserId   cmb.Id `json:"userId"`
	ThreadId cmb.Id `json:"threadId"`
}
type IsUserSubscribedResult struct {
	cmr.CommonResult
	UserId       cmb.Id   `json:"userId"`
	ThreadId     cmb.Id   `json:"threadId"`
	IsSubscribed cmb.Flag `json:"isSubscribed"`
}

type IsUserSubscribedSParams struct {
	cmr.CommonParams
	cmr.DKeyParams
	UserId   cmb.Id `json:"userId"`
	ThreadId cmb.Id `json:"threadId"`
}
type IsUserSubscribedSResult struct {
	cmr.CommonResult
	UserId       cmb.Id   `json:"userId"`
	ThreadId     cmb.Id   `json:"threadId"`
	IsSubscribed cmb.Flag `json:"isSubscribed"`
}

type CountSelfSubscriptionsParams struct {
	cmr.CommonParams
}
type CountSelfSubscriptionsResult struct {
	cmr.CommonResult
	UserSubscriptionsCount cmb.Count `json:"userSubscriptionsCount"`
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
	Page cmb.Count `json:"page"`
}
type GetSelfSubscriptionsOnPageResult struct {
	cmr.CommonResult
	UserSubscriptions *UserSubscriptions `json:"userSubscriptions"`
}

type GetUserSubscriptionsParams struct {
	cmr.CommonParams
	UserId cmb.Id `json:"userId"`
}
type GetUserSubscriptionsResult struct {
	cmr.CommonResult
	UserSubscriptions *UserSubscriptions `json:"userSubscriptions"`
}

type GetUserSubscriptionsOnPageParams struct {
	cmr.CommonParams
	UserId cmb.Id    `json:"userId"`
	Page   cmb.Count `json:"page"`
}
type GetUserSubscriptionsOnPageResult struct {
	cmr.CommonResult
	UserSubscriptions *UserSubscriptions `json:"userSubscriptions"`
}

type GetThreadSubscribersSParams struct {
	cmr.CommonParams
	cmr.DKeyParams
	ThreadId cmb.Id `json:"threadId"`
}
type GetThreadSubscribersSResult struct {
	cmr.CommonResult
	ThreadSubscriptions *ThreadSubscriptionsRecord `json:"threadSubscriptions"`
}

type DeleteSelfSubscriptionParams struct {
	cmr.CommonParams
	ThreadId cmb.Id `json:"threadId"`
}
type DeleteSelfSubscriptionResult = cmr.CommonResultWithSuccess

type DeleteSubscriptionParams struct {
	cmr.CommonParams
	ThreadId cmb.Id `json:"threadId"`
	UserId   cmb.Id `json:"userId"`
}
type DeleteSubscriptionResult = cmr.CommonResultWithSuccess

type DeleteSubscriptionSParams struct {
	cmr.CommonParams
	cmr.DKeyParams
	ThreadId cmb.Id `json:"threadId"`
	UserId   cmb.Id `json:"userId"`
}
type DeleteSubscriptionSResult = cmr.CommonResultWithSuccess

type ClearThreadSubscriptionsSParams struct {
	cmr.CommonParams
	cmr.DKeyParams
	ThreadId cmb.Id `json:"threadId"`
}
type ClearThreadSubscriptionsSResult = cmr.CommonResultWithSuccess

// Other.

type GetDKeyParams struct {
	cmr.CommonParams
}
type GetDKeyResult struct {
	cmr.CommonResult
	DKey cmb.Text `json:"dKey"`
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
