package models

import (
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

// Ping.

type PingParams = cmr.PingParams
type PingResult = cmr.PingResult

// Message.

type SendMessageParams struct {
	Recipient cm.Email `json:"recipient"`
	Subject   cmb.Text `json:"subject"`
	Message   cmb.Text `json:"message"`
}
type SendMessageResult struct {
	cmr.CommonResult
}

// Other.

type ShowDiagnosticDataParams struct{}
type ShowDiagnosticDataResult struct {
	cmr.CommonResult
	cmr.RequestsCount
}
