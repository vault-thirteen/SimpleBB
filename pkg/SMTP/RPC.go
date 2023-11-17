package smtp

// RPC handlers.

import (
	"context"
	"encoding/json"
	"time"

	js "github.com/osamingo/jsonrpc/v2"
	sm "github.com/vault-thirteen/SimpleBB/pkg/SMTP/models"
)

// Ping.

type PingHandler struct {
	Server *Server
}

func (h PingHandler) ServeJSONRPC(_ context.Context, _ *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	result := sm.PingResult{OK: true}
	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

// Message.

type SendMessageHandler struct {
	Server *Server
}

func (h SendMessageHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p sm.SendMessageParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	result, jerr := h.Server.sendEmailMessage(p.Recipient, p.Subject, p.Message)
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
