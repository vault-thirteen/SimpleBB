package models

import (
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

// Ping.

type PingParams = cmr.PingParams
type PingResult = cmr.PingResult

// Captcha.

type CreateCaptchaParams struct{}

type CreateCaptchaResult struct {
	cmr.CommonResult
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
	cmr.CommonResult
	TaskId    string `json:"taskId"`
	IsSuccess bool   `json:"isSuccess"`
}

// Other.

type ShowDiagnosticDataParams struct{}

type ShowDiagnosticDataResult struct {
	cmr.CommonResult
	cmr.RequestsCount
}
