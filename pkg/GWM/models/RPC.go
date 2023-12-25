package models

import (
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

type RequestId = uint

// Ping.

type PingParams = cm.PingParams
type PingResult = cm.PingResult

// IP address list.

type BlockIPAddressParams struct {
	cm.CommonParams

	// IP address of a client to block.
	UserIPA string `json:"userIPA"`

	// Block time in seconds.
	// This is a period for which the specified IP address will be blocked. If
	// for some reason a record with the specified IP address already exists,
	// this time will be added to an already existing value.
	BlockTimeSec uint `json:"blockTimeSec"`
}

type BlockIPAddressResult struct {
	cm.CommonResult

	OK bool `json:"ok"`
}

type IsIPAddressBlockedParams struct {
	cm.CommonParams

	// IP address of a client to check.
	UserIPA string `json:"userIPA"`
}

type IsIPAddressBlockedResult struct {
	cm.CommonResult

	IsBlocked bool `json:"isBlocked"`
}

// Other.

type ShowDiagnosticDataParams struct{}

type ShowDiagnosticDataResult struct {
	cm.CommonResult

	TotalRequestsCount      uint64 `json:"totalRequestsCount"`
	SuccessfulRequestsCount uint64 `json:"successfulRequestsCount"`
}
