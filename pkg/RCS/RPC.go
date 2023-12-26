package rcs

// RPC handlers.

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	js "github.com/osamingo/jsonrpc/v2"
	rm "github.com/vault-thirteen/SimpleBB/pkg/RCS/models"
)

// Ping.

type PingHandler struct {
	Server *Server
}

func (h PingHandler) ServeJSONRPC(_ context.Context, _ *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	result := rm.PingResult{OK: true}
	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

// Captcha.

type CreateCaptchaHandler struct {
	Server *Server
}

func (h CreateCaptchaHandler) ServeJSONRPC(_ context.Context, _ *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	result, jerr := h.Server.createCaptcha()
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type CheckCaptchaHandler struct {
	Server *Server
}

func (h CheckCaptchaHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p rm.CheckCaptchaParams
	jerr = js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *rm.CheckCaptchaResult
	result, jerr = h.Server.checkCaptcha(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

// Other.

type ShowDiagnosticDataHandler struct {
	Server *Server
}

func (h ShowDiagnosticDataHandler) ServeJSONRPC(_ context.Context, _ *json.RawMessage) (resp any, jerr *js.Error) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(fmt.Sprintf("%v, %s", x, string(debug.Stack())))
			jerr = &js.Error{Code: RpcErrorCode_Exception, Message: RpcErrorMsg_Exception}
		}
	}()

	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	result, jerr := h.Server.showDiagnosticData()
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}
