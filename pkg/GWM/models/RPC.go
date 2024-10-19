package models

import (
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

// Ping.

type PingParams = cmr.PingParams
type PingResult = cmr.PingResult

// IP address list.

type BlockIPAddressParams struct {
	cmr.CommonParams

	// IP address of a client to block.
	UserIPA cm.IPAS `json:"userIPA"`

	// Block time in seconds.
	// This is a period for which the specified IP address will be blocked. If
	// for some reason a record with the specified IP address already exists,
	// this time will be added to an already existing value.
	BlockTimeSec cmb.Count `json:"blockTimeSec"`
}
type BlockIPAddressResult = cmr.CommonResultWithSuccess

type IsIPAddressBlockedParams struct {
	cmr.CommonParams

	// IP address of a client to check.
	UserIPA cm.IPAS `json:"userIPA"`
}
type IsIPAddressBlockedResult struct {
	cmr.CommonResult

	IsBlocked cmb.Flag `json:"isBlocked"`
}

// Other.

type ShowDiagnosticDataParams struct{}
type ShowDiagnosticDataResult struct {
	cmr.CommonResult
	cmr.RequestsCount
}
