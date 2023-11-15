package rcs

import (
	"encoding/base64"
	"fmt"

	js "github.com/osamingo/jsonrpc/v2"
	rc "github.com/vault-thirteen/RingCaptcha"
	rm "github.com/vault-thirteen/SimpleBB/pkg/RCS/models"
)

// RPC functions.

func (srv *Server) createCaptcha() (result *rm.CreateCaptchaResult, jerr *js.Error) {
	ccResponse, err := srv.captchaManager.CreateCaptcha()
	if err != nil {
		return nil, &js.Error{Code: RpcErrorCode_CreateError, Message: fmt.Sprintf(RpcErrorMsgF_CreateError, err.Error())}
	}

	result = &rm.CreateCaptchaResult{
		TaskId:              ccResponse.TaskId,
		ImageFormat:         ccResponse.ImageFormat,
		IsImageDataReturned: ccResponse.IsImageDataReturned,
	}

	if ccResponse.IsImageDataReturned {
		result.ImageDataB64 = base64.StdEncoding.EncodeToString(ccResponse.ImageData)
	}

	return result, nil
}

func (srv *Server) checkCaptcha(p *rm.CheckCaptchaParams) (result *rm.CheckCaptchaResult, jerr *js.Error) {
	// Check parameters.
	if len(p.TaskId) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_TaskIdIsNotSet, Message: RpcErrorMsg_TaskIdIsNotSet}
	}

	if p.Value == 0 {
		return nil, &js.Error{Code: RpcErrorCode_AnswerIsNotSet, Message: RpcErrorMsg_AnswerIsNotSet}
	}

	resp, err := srv.captchaManager.CheckCaptcha(&rc.CheckCaptchaRequest{TaskId: p.TaskId, Value: p.Value})
	if err != nil {
		return nil, &js.Error{Code: RpcErrorCode_CheckError, Message: fmt.Sprintf(RpcErrorMsgF_CheckError, err.Error())}
	}

	result = &rm.CheckCaptchaResult{
		TaskId:    p.TaskId,
		IsSuccess: resp.IsSuccess,
	}

	return result, nil
}

func (srv *Server) showDiagnosticData() (result *rm.ShowDiagnosticDataResult, jerr *js.Error) {
	result = &rm.ShowDiagnosticDataResult{
		TotalRequestsCount:      srv.diag.GetTotalRequestsCount(),
		SuccessfulRequestsCount: srv.diag.GetSuccessfulRequestsCount(),
	}

	return result, nil
}
