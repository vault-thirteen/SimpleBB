package models

import (
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

// Ping.

type PingParams = cm.PingParams
type PingResult = cm.PingResult

// Captcha.

type CreateCaptchaParams struct{}

type CreateCaptchaResult struct {
	cm.CommonResult

	TaskId              string `json:"taskId"`
	ImageFormat         string `json:"imageFormat"`
	IsImageDataReturned bool   `json:"isImageDataReturned"`
	ImageDataB64        string `json:"imageDataB64,omitempty"`
}

type CheckCaptchaParams struct {
	TaskId string `json:"taskId"`
	Value  uint   `json:"value"`
}

type CheckCaptchaResult struct {
	cm.CommonResult

	TaskId    string `json:"taskId"`
	IsSuccess bool   `json:"isSuccess"`
}

// Other.

type ShowDiagnosticDataParams struct{}

type ShowDiagnosticDataResult struct {
	cm.CommonResult

	TotalRequestsCount      uint64 `json:"totalRequestsCount"`
	SuccessfulRequestsCount uint64 `json:"successfulRequestsCount"`
}
