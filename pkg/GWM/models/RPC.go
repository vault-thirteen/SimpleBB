package models

import (
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

type RequestId = uint

// Ping.

type PingParams = cmr.PingParams
type PingResult = cmr.PingResult

// IP address list.

type BlockIPAddressParams struct {
	cmr.CommonParams

	// IP address of a client to block.
	UserIPA string `json:"userIPA"`

	// Block time in seconds.
	// This is a period for which the specified IP address will be blocked. If
	// for some reason a record with the specified IP address already exists,
	// this time will be added to an already existing value.
	BlockTimeSec uint `json:"blockTimeSec"`
}

type BlockIPAddressResult struct {
	cmr.CommonResult

	OK bool `json:"ok"`
}

type IsIPAddressBlockedParams struct {
	cmr.CommonParams

	// IP address of a client to check.
	UserIPA string `json:"userIPA"`
}

type IsIPAddressBlockedResult struct {
	cmr.CommonResult

	IsBlocked bool `json:"isBlocked"`
}

// Other.

type ShowDiagnosticDataParams struct{}

type ShowDiagnosticDataResult struct {
	cmr.CommonResult
	cmr.RequestsCount
}
