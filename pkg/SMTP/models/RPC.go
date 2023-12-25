package models

import (
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

// Ping.

type PingParams = cm.PingParams
type PingResult = cm.PingResult

// Message.

type SendMessageParams struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Message   string `json:"message"`
}

type SendMessageResult struct {
	cm.CommonResult
}

// Other.

type ShowDiagnosticDataParams struct{}

type ShowDiagnosticDataResult struct {
	cm.CommonResult

	TotalRequestsCount      uint64 `json:"totalRequestsCount"`
	SuccessfulRequestsCount uint64 `json:"successfulRequestsCount"`
}
