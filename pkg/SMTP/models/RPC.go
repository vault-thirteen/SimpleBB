package models

import (
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

type PingParams struct{}

type PingResult struct {
	OK bool `json:"ok"`
}

type SendMessageParams struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Message   string `json:"message"`
}

type SendMessageResult struct {
	cm.CommonResult
}

type ShowDiagnosticDataParams struct{}

type ShowDiagnosticDataResult struct {
	cm.CommonResult

	TotalRequestsCount      uint64 `json:"totalRequestsCount"`
	SuccessfulRequestsCount uint64 `json:"successfulRequestsCount"`
}
