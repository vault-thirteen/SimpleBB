package smtp

// RPC handlers.

import (
	"encoding/json"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	sm "github.com/vault-thirteen/SimpleBB/pkg/SMTP/models"
	cs "github.com/vault-thirteen/SimpleBB/pkg/common/settings"
)

func (srv *Server) initRpc() (err error) {
	rpcDurationFieldName := cs.RpcDurationFieldName
	rpcRequestIdFieldName := cs.RpcRequestIdFieldName

	ps := &jrm1.ProcessorSettings{
		CatchExceptions:    true,
		LogExceptions:      true,
		CountRequests:      true,
		DurationFieldName:  &rpcDurationFieldName,
		RequestIdFieldName: &rpcRequestIdFieldName,
	}

	srv.js, err = jrm1.NewProcessor(ps)
	if err != nil {
		return err
	}

	fns := []jrm1.RpcFunction{
		srv.Ping,
		srv.SendMessage,
		srv.ShowDiagnosticData,
	}

	for _, fn := range fns {
		err = srv.js.AddFunc(fn)
		if err != nil {
			return err
		}
	}

	return nil
}

// Ping.

func (srv *Server) Ping(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	return sm.PingResult{OK: true}, nil
}

// Message.

func (srv *Server) SendMessage(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *sm.SendMessageParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *sm.SendMessageResult
	r, re = srv.sendMessage(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

// Other.

func (srv *Server) ShowDiagnosticData(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *sm.ShowDiagnosticDataParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *sm.ShowDiagnosticDataResult
	r, re = srv.showDiagnosticData()
	if re != nil {
		return nil, re
	}

	return r, nil
}
