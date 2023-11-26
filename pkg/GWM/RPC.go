package acm

// RPC handlers.

import (
	"context"
	"encoding/json"
	"time"

	js "github.com/osamingo/jsonrpc/v2"
	gm "github.com/vault-thirteen/SimpleBB/pkg/GWM/models"
)

// Ping.

type PingHandler struct {
	Server *Server
}

func (h PingHandler) ServeJSONRPC(_ context.Context, _ *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	result := gm.PingResult{OK: true}
	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

// IP address list.

type BlockIPAddressHandler struct {
	Server *Server
}

func (h BlockIPAddressHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p gm.BlockIPAddressParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *gm.BlockIPAddressResult
	result, jerr = h.Server.blockIPAddress(&p)
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

type IsIPAddressBlockedHandler struct {
	Server *Server
}

func (h IsIPAddressBlockedHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p gm.IsIPAddressBlockedParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *gm.IsIPAddressBlockedResult
	result, jerr = h.Server.isIPAddressBlocked(&p)
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

func (h ShowDiagnosticDataHandler) ServeJSONRPC(_ context.Context, _ *json.RawMessage) (any, *js.Error) {
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
