package c

import (
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/client"
)

// List of supported functions.
const (
	// Ping.
	FuncPing = cc.FuncPing

	// Subscription.
	FuncAddSubscription           = "AddSubscription"
	FuncGetSelfSubscriptions      = "GetSelfSubscriptions"
	FuncGetUserSubscriptions      = "GetUserSubscriptions"
	FuncGetThreadSubscribersS     = "GetThreadSubscribersS"
	FuncDeleteSelfSubscription    = "DeleteSelfSubscription"
	FuncDeleteSubscription        = "DeleteSubscription"
	FuncDeleteSubscriptionS       = "DeleteSubscriptionS"
	FuncClearThreadSubscriptionsS = "ClearThreadSubscriptionsS"

	// Other.
	FuncGetDKey            = "GetDKey"
	FuncShowDiagnosticData = cc.FuncShowDiagnosticData
	FuncTest               = "Test"
)
